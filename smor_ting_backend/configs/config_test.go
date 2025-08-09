package configs_test

import (
	"encoding/base64"
	"os"
	"testing"

	cfg "github.com/smorting/backend/configs"
)

// helper to set and restore env
func withEnv(t *testing.T, kv map[string]string, fn func()) {
	t.Helper()
	old := make(map[string]string)
	for k, v := range kv {
		old[k] = os.Getenv(k)
		if v == "__UNSET__" {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}
	defer func() {
		for k, v := range old {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()
	fn()
}

func TestLoadConfig_FailsClosedInProduction_WhenSecretsMissingOrDefault(t *testing.T) {
	withEnv(t, map[string]string{
		"ENV":                    "production",
		"JWT_ACCESS_SECRET":      "__UNSET__",
		"JWT_REFRESH_SECRET":     "__UNSET__",
		"ENCRYPTION_KEY":         "__UNSET__",
		"PAYMENT_ENCRYPTION_KEY": "__UNSET__",
	}, func() {
		if _, err := cfg.LoadConfig(); err == nil {
			t.Fatalf("expected error when secrets missing in production, got nil")
		}
	})
}

func TestLoadConfig_SucceedsInProduction_WithValidBase64Secrets(t *testing.T) {
	// 32-byte secrets
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	enc := make([]byte, 32)
	pay := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	for i := range enc {
		enc[i] = 3
	}
	for i := range pay {
		pay[i] = 4
	}

	withEnv(t, map[string]string{
		"ENV":                    "production",
		"JWT_ACCESS_SECRET":      base64.StdEncoding.EncodeToString(access),
		"JWT_REFRESH_SECRET":     base64.StdEncoding.EncodeToString(refresh),
		"ENCRYPTION_KEY":         base64.StdEncoding.EncodeToString(enc),
		"PAYMENT_ENCRYPTION_KEY": base64.StdEncoding.EncodeToString(pay),
		// minimal MoMo production requirements
		"MOMO_BASE_URL":           "https://example.com",
		"MOMO_API_USER":           "user",
		"MOMO_API_KEY":            "key",
		"MOMO_SUB_KEY_COLLECTION": "sub-col",
		// minimal SmileID production requirements
		"SMILEID_BASE_URL":   "https://kyc.example.com",
		"SMILEID_PARTNER_ID": "pid",
		"SMILEID_API_KEY":    "skey",
	}, func() {
		if _, err := cfg.LoadConfig(); err != nil {
			t.Fatalf("expected success with valid base64 secrets, got error: %v", err)
		}
	})
}

func TestLoadConfig_FailsInProduction_WithInvalidBase64(t *testing.T) {
	withEnv(t, map[string]string{
		"ENV":                    "production",
		"JWT_ACCESS_SECRET":      "not-base64",
		"JWT_REFRESH_SECRET":     "also-not-base64",
		"ENCRYPTION_KEY":         "bad",
		"PAYMENT_ENCRYPTION_KEY": "bad",
	}, func() {
		if _, err := cfg.LoadConfig(); err == nil {
			t.Fatalf("expected error with invalid base64 secrets")
		}
	})
}

func TestLoadConfig_FailsClosedInStaging_WhenSecretsMissingOrDefault(t *testing.T) {
	withEnv(t, map[string]string{
		"ENV":                    "staging",
		"JWT_ACCESS_SECRET":      "__UNSET__",
		"JWT_REFRESH_SECRET":     "__UNSET__",
		"ENCRYPTION_KEY":         "__UNSET__",
		"PAYMENT_ENCRYPTION_KEY": "__UNSET__",
	}, func() {
		if _, err := cfg.LoadConfig(); err == nil {
			t.Fatalf("expected error when secrets missing in staging, got nil")
		}
	})
}

func TestLoadConfig_SucceedsInStaging_WithValidBase64Secrets(t *testing.T) {
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	enc := make([]byte, 32)
	pay := make([]byte, 32)
	for i := range access {
		access[i] = 1
	}
	for i := range refresh {
		refresh[i] = 2
	}
	for i := range enc {
		enc[i] = 3
	}
	for i := range pay {
		pay[i] = 4
	}

	withEnv(t, map[string]string{
		"ENV":                    "staging",
		"JWT_ACCESS_SECRET":      base64.StdEncoding.EncodeToString(access),
		"JWT_REFRESH_SECRET":     base64.StdEncoding.EncodeToString(refresh),
		"ENCRYPTION_KEY":         base64.StdEncoding.EncodeToString(enc),
		"PAYMENT_ENCRYPTION_KEY": base64.StdEncoding.EncodeToString(pay),
		// minimal MoMo staging requirements
		"MOMO_BASE_URL":           "https://staging.example.com",
		"MOMO_API_USER":           "user",
		"MOMO_API_KEY":            "key",
		"MOMO_SUB_KEY_COLLECTION": "sub-col",
		// minimal SmileID staging requirements
		"SMILEID_BASE_URL":   "https://kyc-staging.example.com",
		"SMILEID_PARTNER_ID": "pid",
		"SMILEID_API_KEY":    "skey",
	}, func() {
		if _, err := cfg.LoadConfig(); err != nil {
			t.Fatalf("expected success with valid base64 secrets in staging, got error: %v", err)
		}
	})
}
