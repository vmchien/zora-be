package utils

import (
	"net"
	"net/http"
	"strings"
)

// NormalizeIP normalizes loopback addresses for consistency.
// - ::1        -> 127.0.0.1
// - 127.0.0.1  -> 127.0.0.1
func NormalizeIP(ip net.IP) net.IP {
	if ip == nil {
		return nil
	}

	if ip.IsLoopback() {
		return net.IPv4(127, 0, 0, 1)
	}

	return ip
}

// GetClientIPFromNginx returns the real client IP.
//
// Assumptions (IMPORTANT):
//   - All production traffic goes through trusted nginx.
//   - nginx sets and sanitizes:
//     proxy_set_header X-Real-IP $remote_addr;
//   - The application is NOT directly exposed to the internet.
//
// Behavior:
//   - Prefer X-Real-IP (authoritative).
//   - Fallback to RemoteAddr for local dev / debugging.
//   - No CIDR checks, no X-Forwarded-For parsing.
func GetClientIPFromNginx(req *http.Request) string {
	// fmt.Printf(
	// 	"client-ip resolve: remote=%s xreal=%s xff=%s\n",
	// 	req.RemoteAddr,
	// 	req.Header.Get("X-Real-Ip"),
	// 	req.Header.Get("X-Forwarded-For"),
	// )

	// 1. Authoritative source: X-Real-IP (set by nginx)
	if xr := strings.TrimSpace(req.Header.Get("X-Real-Ip")); xr != "" {
		if ip := NormalizeIP(net.ParseIP(xr)); ip != nil {
			return ip.String()
		}
	}

	// 2. Fallback: direct connection (local dev / tests)
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		if ip := NormalizeIP(net.ParseIP(req.RemoteAddr)); ip != nil {
			return ip.String()
		}
		return req.RemoteAddr
	}

	if ip := NormalizeIP(net.ParseIP(host)); ip != nil {
		return ip.String()
	}

	return host
}

// GetUserAgent returns the raw User-Agent header.
func GetUserAgent(req *http.Request) string {
	return req.UserAgent()
}

// IsPrivateIP reports whether ip is a private or link-local address.
// NOTE: This MUST NOT be used for proxy trust or security decisions.
// Intended for logging, analytics, or soft anti-abuse checks only.
func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16", // link-local
	}

	for _, cidr := range privateCIDRs {
		if _, n, err := net.ParseCIDR(cidr); err == nil && n.Contains(ip) {
			return true
		}
	}

	return false
}
