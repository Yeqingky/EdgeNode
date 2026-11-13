# Changelog

## v1.4.2-rc.1 (2026-11-13)

- 升级全部依赖链并升级 Go 工具链到 1.27.4, govulncheck 符号级扫描可达漏洞 18 -> 0
- andybalholm/brotli v1.0.5 -> v1.3.0, 修复解压恶意流内存耗尽 (GO-2026-6888)
- klauspost/compress v1.17.8 -> v1.20.4, 修复两块帧 panic (GO-2026-6865)
- golang.org/x/image v0.16.0 -> v0.48.0, 修复 webp/vp8/bmp 解码 panic 与内存放大 (GO-2026-6222, GO-2026-5061, GO-2026-5031, GO-2026-4961)
- google/nftables v0.2.0 -> v0.3.0, google.golang.org/grpc v1.63.2 -> v1.85.0, golang.org/x/net v0.25.0 -> v0.62.0
- 依赖 v8go 与 quic-go 无代码引用, 已随 go mod tidy 移除
