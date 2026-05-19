package acme

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/providers/dns/digitalocean"
	"github.com/go-acme/lego/v4/providers/dns/route53"
)

var supportedDNSProviders = map[string]string{
	"cloudflare":   "CF_DNS_API_TOKEN or CF_API_EMAIL+CF_API_KEY",
	"route53":      "AWS_ACCESS_KEY_ID + AWS_SECRET_ACCESS_KEY + AWS_REGION",
	"digitalocean": "DO_AUTH_TOKEN",
}

func GetDNSProvider(name string) (challenge.Provider, error) {
	switch strings.ToLower(name) {
	case "cloudflare":
		return cloudflare.NewDNSProvider()
	case "route53":
		return route53.NewDNSProvider()
	case "digitalocean":
		return digitalocean.NewDNSProvider()
	default:
		return nil, fmt.Errorf("unsupported DNS provider: %s\nSupported: %s", name, listProviders())
	}
}

func DetectDNSProvider() (string, bool) {
	if os.Getenv("CF_DNS_API_TOKEN") != "" || os.Getenv("CF_API_KEY") != "" {
		return "cloudflare", true
	}
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		return "route53", true
	}
	if os.Getenv("DO_AUTH_TOKEN") != "" {
		return "digitalocean", true
	}
	return "", false
}

func ListDNSProviders() map[string]string {
	return supportedDNSProviders
}

func listProviders() string {
	names := make([]string, 0, len(supportedDNSProviders))
	for k := range supportedDNSProviders {
		names = append(names, k)
	}
	return strings.Join(names, ", ")
}
