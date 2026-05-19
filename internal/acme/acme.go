package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/valorisa/lecert/internal/store"
)

const (
	prodURL    = lego.LEDirectoryProduction
	stagingURL = lego.LEDirectoryStaging
)

type ObtainRequest struct {
	Domain      string
	Email       string
	Challenge   string
	DNSProvider string
	Staging     bool
}

type CertResult struct {
	Domain   string
	NotAfter time.Time
	Path     string
}

type user struct {
	email        string
	registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *user) GetEmail() string                        { return u.email }
func (u *user) GetRegistration() *registration.Resource { return u.registration }
func (u *user) GetPrivateKey() crypto.PrivateKey        { return u.key }

func Obtain(req ObtainRequest) (*CertResult, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	u := &user{email: req.Email, key: privateKey}

	config := lego.NewConfig(u)
	config.Certificate.KeyType = certcrypto.EC256
	if req.Staging {
		config.CADirURL = stagingURL
	} else {
		config.CADirURL = prodURL
	}

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	switch req.Challenge {
	case "http-01", "":
		err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
		if err != nil {
			return nil, fmt.Errorf("set http01 provider: %w", err)
		}
	case "dns-01":
		dnsProvider, err := GetDNSProvider(req.DNSProvider)
		if err != nil {
			return nil, fmt.Errorf("dns provider: %w", err)
		}
		err = client.Challenge.SetDNS01Provider(dnsProvider)
		if err != nil {
			return nil, fmt.Errorf("set dns01 provider: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported challenge type: %s (use http-01 or dns-01)", req.Challenge)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	u.registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{req.Domain},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("obtain certificate: %w", err)
	}

	path, notAfter, err := store.Save(store.SaveOptions{
		Domain:    req.Domain,
		Email:     req.Email,
		Challenge: req.Challenge,
	}, certificates.Certificate, certificates.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("store certificate: %w", err)
	}

	return &CertResult{
		Domain:   req.Domain,
		NotAfter: notAfter,
		Path:     path,
	}, nil
}

func Renew(domain string, force bool, staging bool) (*CertResult, error) {
	meta, err := store.Load(domain)
	if err != nil {
		return nil, fmt.Errorf("load cert: %w", err)
	}

	if !force && time.Until(meta.NotAfter) > 30*24*time.Hour {
		return nil, fmt.Errorf("certificate expires %s — not due for renewal (use --force to override)", meta.NotAfter.Format("2006-01-02"))
	}

	return Obtain(ObtainRequest{
		Domain:    domain,
		Email:     meta.Email,
		Challenge: meta.Challenge,
		Staging:   staging,
	})
}

func Revoke(domain string, staging bool) error {
	meta, err := store.Load(domain)
	if err != nil {
		return fmt.Errorf("load cert: %w", err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	u := &user{email: meta.Email, key: privateKey}

	config := lego.NewConfig(u)
	if staging {
		config.CADirURL = stagingURL
	} else {
		config.CADirURL = prodURL
	}

	client, err := lego.NewClient(config)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}
	u.registration = reg

	return client.Certificate.Revoke(meta.CertPEM)
}
