// Helper para extraer y normalizar el dominio a partir de una URL (para filtrado).
package usecase

import (
	"net/url"
	"strings"
)

func extractDomain(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		if !strings.Contains(raw, "://") {
			raw = "https://" + raw
			if u2, err2 := url.Parse(raw); err2 == nil {
				host := u2.Hostname()
				return strings.TrimPrefix(strings.ToLower(host), "www.")
			}
		}
		return ""
	}
	host := u.Hostname()
	return strings.TrimPrefix(strings.ToLower(host), "www.")
}
