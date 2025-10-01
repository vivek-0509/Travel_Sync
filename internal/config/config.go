package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	Port           string
	PostgresURI    string
	FrontendURL    string
	AllowedOrigins []string
	CookieSecure   bool
	CookieDomain   string
	GinMode        string
	TrustedProxies []string
}

func LoadConfig() *AppConfig {
	allowed := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if allowed != "" {
		for _, o := range strings.Split(allowed, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
	}
	return &AppConfig{
		Port:           os.Getenv("PORT"),
		PostgresURI:    os.Getenv("POSTGRES_URI"),
		FrontendURL:    os.Getenv("FRONTEND_URL"),
		AllowedOrigins: origins,
		CookieSecure:   strings.ToLower(os.Getenv("COOKIE_SECURE")) == "true",
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
		GinMode:        os.Getenv("GIN_MODE"),
		TrustedProxies: splitAndTrim(os.Getenv("TRUSTED_PROXIES")),
	}

}

func splitAndTrim(v string) []string {
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
