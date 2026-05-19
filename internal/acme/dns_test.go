package acme

import (
	"os"
	"testing"
)

func TestDetectDNSProviderCloudflare(t *testing.T) {
	t.Setenv("CF_DNS_API_TOKEN", "test-token")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("DO_AUTH_TOKEN")

	provider, ok := DetectDNSProvider()
	if !ok {
		t.Fatal("Expected detection")
	}
	if provider != "cloudflare" {
		t.Errorf("Got %q, want cloudflare", provider)
	}
}

func TestDetectDNSProviderRoute53(t *testing.T) {
	os.Unsetenv("CF_DNS_API_TOKEN")
	os.Unsetenv("CF_API_KEY")
	t.Setenv("AWS_ACCESS_KEY_ID", "test-key")
	os.Unsetenv("DO_AUTH_TOKEN")

	provider, ok := DetectDNSProvider()
	if !ok {
		t.Fatal("Expected detection")
	}
	if provider != "route53" {
		t.Errorf("Got %q, want route53", provider)
	}
}

func TestDetectDNSProviderNone(t *testing.T) {
	os.Unsetenv("CF_DNS_API_TOKEN")
	os.Unsetenv("CF_API_KEY")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("DO_AUTH_TOKEN")

	_, ok := DetectDNSProvider()
	if ok {
		t.Fatal("Expected no detection")
	}
}

func TestGetDNSProviderUnsupported(t *testing.T) {
	_, err := GetDNSProvider("notreal")
	if err == nil {
		t.Fatal("Expected error for unsupported provider")
	}
}

func TestListDNSProviders(t *testing.T) {
	providers := ListDNSProviders()
	if len(providers) < 3 {
		t.Errorf("Expected at least 3 providers, got %d", len(providers))
	}
	if _, ok := providers["cloudflare"]; !ok {
		t.Error("Missing cloudflare")
	}
}
