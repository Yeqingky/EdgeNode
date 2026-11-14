package nodes

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"strconv"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/configutils"
	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/remotelogs"
	"github.com/iwind/TeaGo/types"
)

// OriginConnect 连接源站
func OriginConnect(origin *serverconfigs.OriginConfig, serverPort int, remoteAddr string, tlsHost string) (originConn net.Conn, originAddr string, err error) {
	if origin.Addr == nil {
		return nil, "", errors.New("origin server address should not be empty")
	}

	if origin.Addr.TunnelNodeId > 0 {
		originAddr = origin.Addr.PickAddress()
		if origin.FollowPort && serverPort > 0 {
			originAddr = configutils.QuoteIP(origin.Addr.Host) + ":" + types.String(serverPort)
		}
		timeout := origin.ConnTimeoutDuration()
		if timeout <= 0 {
			timeout = 15 * time.Second
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err := dialOriginContext(ctx, origin, "tcp", originAddr, timeout)
		if err != nil {
			return nil, originAddr, err
		}
		if origin.Addr.Protocol != serverconfigs.ProtocolHTTPS && origin.Addr.Protocol != serverconfigs.ProtocolTLS {
			return conn, originAddr, nil
		}
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		if origin.Cert != nil {
			if cert := origin.Cert.CertObject(); cert != nil {
				tlsConfig.InsecureSkipVerify = false
				tlsConfig.Certificates = []tls.Certificate{*cert}
				tlsConfig.ServerName = origin.Cert.ServerName
			}
		}
		if tlsHost != "" {
			tlsConfig.ServerName = tlsHost
		}
		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName, _, _ = net.SplitHostPort(originAddr)
		}
		tlsConn := tls.Client(conn, tlsConfig)
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			conn.Close()
			return nil, originAddr, err
		}
		return tlsConn, originAddr, nil
	}

	// 支持TOA的连接
	// 这个条件很重要,如果没有传递remoteAddr,表示不使用TOA
	if len(remoteAddr) > 0 {
		var toaConfig = sharedTOAManager.Config()
		if toaConfig != nil && toaConfig.IsOn {
			var retries = 3
			for i := 1; i <= retries; i++ {
				var port = int(toaConfig.RandLocalPort())
				err = sharedTOAManager.SendMsg("add:" + strconv.Itoa(port) + ":" + remoteAddr)
				if err != nil {
					remotelogs.Error("TOA", "add failed: "+err.Error())
				} else {
					var dialer = net.Dialer{
						Timeout: origin.ConnTimeoutDuration(),
						LocalAddr: &net.TCPAddr{
							Port: port,
						},
					}
					originAddr = origin.Addr.PickAddress()

					// 端口跟随
					if origin.FollowPort && serverPort > 0 {
						originAddr = configutils.QuoteIP(origin.Addr.Host) + ":" + types.String(serverPort)
					}

					var conn net.Conn
					switch origin.Addr.Protocol {
					case "", serverconfigs.ProtocolTCP, serverconfigs.ProtocolHTTP:
						// TODO 支持TCP4/TCP6
						// TODO 支持指定特定网卡
						conn, err = dialer.Dial("tcp", originAddr)
					case serverconfigs.ProtocolTLS, serverconfigs.ProtocolHTTPS:
						// TODO 支持TCP4/TCP6
						// TODO 支持指定特定网卡

						var tlsConfig = &tls.Config{
							InsecureSkipVerify: true,
						}
						if origin.Cert != nil {
							var obj = origin.Cert.CertObject()
							if obj != nil {
								tlsConfig.InsecureSkipVerify = false
								tlsConfig.Certificates = []tls.Certificate{*obj}
								if len(origin.Cert.ServerName) > 0 {
									tlsConfig.ServerName = origin.Cert.ServerName
								}
							}
						}
						if len(tlsHost) > 0 {
							tlsConfig.ServerName = tlsHost
						}

						conn, err = tls.DialWithDialer(&dialer, "tcp", originAddr, tlsConfig)
					}

					// TODO 需要在合适的时机删除TOA记录
					if err == nil || i == retries {
						return conn, originAddr, err
					}
				}
			}
		}
	}

	originAddr = origin.Addr.PickAddress()

	// 端口跟随
	if origin.FollowPort && serverPort > 0 {
		originAddr = configutils.QuoteIP(origin.Addr.Host) + ":" + types.String(serverPort)
	}

	switch origin.Addr.Protocol {
	case "", serverconfigs.ProtocolTCP, serverconfigs.ProtocolHTTP:
		// TODO 支持TCP4/TCP6
		// TODO 支持指定特定网卡
		originConn, err = net.DialTimeout("tcp", originAddr, origin.ConnTimeoutDuration())
		return originConn, originAddr, err
	case serverconfigs.ProtocolTLS, serverconfigs.ProtocolHTTPS:
		// TODO 支持TCP4/TCP6
		// TODO 支持指定特定网卡

		var tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		if origin.Cert != nil {
			var obj = origin.Cert.CertObject()
			if obj != nil {
				tlsConfig.InsecureSkipVerify = false
				tlsConfig.Certificates = []tls.Certificate{*obj}
				if len(origin.Cert.ServerName) > 0 {
					tlsConfig.ServerName = origin.Cert.ServerName
				}
			}
		}
		if len(tlsHost) > 0 {
			tlsConfig.ServerName = tlsHost
		}

		originConn, err = tls.Dial("tcp", originAddr, tlsConfig)
		return originConn, originAddr, err
	case serverconfigs.ProtocolUDP:
		addr, err := net.ResolveUDPAddr("udp", originAddr)
		if err != nil {
			return nil, originAddr, err
		}
		originConn, err = net.DialUDP("udp", nil, addr)
		return originConn, originAddr, err
	}

	// TODO 支持从Unix,Pipe,HTTP,HTTPS中读取数据

	return nil, originAddr, errors.New("invalid origin scheme '" + origin.Addr.Protocol.String() + "'")
}
