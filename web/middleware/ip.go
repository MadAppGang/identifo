package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/model"
)

// IP tries to extract real IP from user.
func IP(errorPath string, appStorage model.AppStorage, logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-REAL-IP")
			netIP := net.ParseIP(ip)
			if netIP == nil {
				// Get IP from X-FORWARDED-FOR header
				ips := r.Header.Get("X-FORWARDED-FOR")
				splitIps := strings.Split(ips, ",")
				for _, ip := range splitIps {
					netIP = net.ParseIP(ip)
					if netIP != nil {
						break
					}
				}
			}

			if netIP == nil {
				// Get IP from RemoteAddr
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err == nil {
					netIP = net.ParseIP(ip)
				}
			}

			// we managed to found the IP
			if netIP != nil {
				ctx := context.WithValue(r.Context(), model.IPContextKey, netIP.String())
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IPFromContext returns IP from request context.
func IPFromContext(ctx context.Context) string {
	value := ctx.Value(model.IPContextKey)

	if value == nil {
		return ""
	}

	return value.(string)
}
