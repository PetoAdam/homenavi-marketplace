package config

import (
	"os"
	"strings"
)

type Config struct {
	BindAddress        string
	DatabaseURL        string
	AllowedOrigin      []string
	OIDCIssuer         string
	OIDCAudience       string
	OIDCVerifyWorkflow string
	OIDCTagPrefix      string
	GitHubAPIToken     string
}

func Load() Config {
	bind := getEnv("BIND_ADDRESS", ":8098")
	dbURL := os.Getenv("DATABASE_URL")
	origins := splitCSV(os.Getenv("ALLOWED_ORIGINS"))
	issuer := getEnv("OIDC_ISSUER", "https://token.actions.githubusercontent.com")
	audience := getEnv("OIDC_AUDIENCE", "homenavi-marketplace")
	verifyWorkflow := getEnv("OIDC_VERIFY_WORKFLOW", "verify.yml")
	tagPrefix := getEnv("OIDC_TAG_PREFIX", "v")
	githubToken := os.Getenv("GITHUB_API_TOKEN")

	return Config{
		BindAddress:        bind,
		DatabaseURL:        dbURL,
		AllowedOrigin:      origins,
		OIDCIssuer:         issuer,
		OIDCAudience:       audience,
		OIDCVerifyWorkflow: verifyWorkflow,
		OIDCTagPrefix:      tagPrefix,
		GitHubAPIToken:     githubToken,
	}
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
