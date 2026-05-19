package store

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func generateTestCert(domain string, notAfter time.Time) ([]byte, []byte) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: domain},
		NotBefore:    time.Now(),
		NotAfter:     notAfter,
		DNSNames:     []string{domain},
	}

	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyDER, _ := x509.MarshalECPrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return certPEM, keyPEM
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("LECERT_DIR", tmpDir)

	domain := "test.example.com"
	expiry := time.Now().Add(90 * 24 * time.Hour)
	certPEM, keyPEM := generateTestCert(domain, expiry)

	opts := SaveOptions{
		Domain:    domain,
		Email:     "test@example.com",
		Challenge: "http-01",
	}

	dir, notAfter, err := Save(opts, certPEM, keyPEM)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if dir == "" {
		t.Fatal("Save returned empty directory")
	}

	if notAfter.Before(time.Now()) {
		t.Fatal("NotAfter is in the past")
	}

	// Verify files exist
	if _, err := os.Stat(filepath.Join(dir, "cert.pem")); err != nil {
		t.Fatalf("cert.pem not found: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "key.pem")); err != nil {
		t.Fatalf("key.pem not found: %v", err)
	}

	// Load and verify
	meta, err := Load(domain)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if meta.Domain != domain {
		t.Errorf("Domain mismatch: got %s, want %s", meta.Domain, domain)
	}
	if meta.Email != "test@example.com" {
		t.Errorf("Email mismatch: got %s", meta.Email)
	}
	if meta.Challenge != "http-01" {
		t.Errorf("Challenge mismatch: got %s", meta.Challenge)
	}
	if len(meta.CertPEM) == 0 {
		t.Error("CertPEM is empty after Load")
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("LECERT_DIR", tmpDir)

	// Empty list
	certs, err := List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(certs) != 0 {
		t.Fatalf("Expected 0 certs, got %d", len(certs))
	}

	// Save two certs
	for _, domain := range []string{"a.example.com", "b.example.com"} {
		certPEM, keyPEM := generateTestCert(domain, time.Now().Add(90*24*time.Hour))
		opts := SaveOptions{Domain: domain, Email: "t@t.com", Challenge: "http-01"}
		if _, _, err := Save(opts, certPEM, keyPEM); err != nil {
			t.Fatalf("Save %s: %v", domain, err)
		}
	}

	certs, err = List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(certs) != 2 {
		t.Fatalf("Expected 2 certs, got %d", len(certs))
	}
}

func TestKeyPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod not supported on Windows — ACL-based security used instead")
	}

	tmpDir := t.TempDir()
	t.Setenv("LECERT_DIR", tmpDir)

	domain := "perm.example.com"
	certPEM, keyPEM := generateTestCert(domain, time.Now().Add(90*24*time.Hour))
	opts := SaveOptions{Domain: domain, Email: "t@t.com", Challenge: "http-01"}

	dir, _, err := Save(opts, certPEM, keyPEM)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	keyPath := filepath.Join(dir, "key.pem")
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("Stat key: %v", err)
	}

	if info.Mode().Perm()&0077 != 0 {
		t.Errorf("Key file permissions too open: %o (want 0600)", info.Mode().Perm())
	}
}
