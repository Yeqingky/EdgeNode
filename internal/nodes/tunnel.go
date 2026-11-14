package nodes

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/frp"
)

type tunnelConn struct {
	net.Conn
	reader *bufio.Reader
}

func (c *tunnelConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func dialOriginContext(ctx context.Context, origin *serverconfigs.OriginConfig, network, addr string, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout, KeepAlive: time.Minute}
	if origin.Addr.TunnelNodeId <= 0 {
		return dialer.DialContext(ctx, network, addr)
	}
	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, fmt.Errorf("tunnel does not support %s", network)
	}
	if strings.ContainsAny(addr, "\r\n") {
		return nil, fmt.Errorf("invalid tunnel origin address")
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	paths := frp.SharedManager.TunnelSocketPaths(origin.Addr.TunnelNodeId)
	if len(paths) == 0 {
		return nil, fmt.Errorf("tunnel %d is unavailable", origin.Addr.TunnelNodeId)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	attemptTimeout := timeout / time.Duration(len(paths))
	if attemptTimeout < time.Second {
		attemptTimeout = time.Second
	}
	var lastErr error
	for _, path := range paths {
		conn, err := dialTunnelPath(ctx, path, addr, attemptTimeout)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		if ctx.Err() != nil {
			break
		}
	}
	return nil, fmt.Errorf("tunnel %d: %w", origin.Addr.TunnelNodeId, lastErr)
}

func dialTunnelPath(ctx context.Context, path, addr string, timeout time.Duration) (net.Conn, error) {
	conn, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "unix", path)
	if err != nil {
		return nil, err
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.SetDeadline(time.Now()) })
	defer stop()
	deadline := time.Now().Add(timeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetDeadline(deadline); err != nil {
		conn.Close()
		return nil, err
	}
	request := &http.Request{Method: http.MethodConnect, URL: &url.URL{Opaque: addr}, Host: addr, Header: make(http.Header)}
	if err := request.Write(conn); err != nil {
		conn.Close()
		return nil, err
	}
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("tunnel connection failed: %s", response.Status)
	}
	if err := ctx.Err(); err != nil {
		conn.Close()
		return nil, err
	}
	if err := conn.SetDeadline(time.Time{}); err != nil {
		conn.Close()
		return nil, err
	}
	return &tunnelConn{Conn: conn, reader: reader}, nil
}
