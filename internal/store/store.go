package store

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type CertMeta struct {
	Domain    string    `json:"domain"`
	Email     string    `json:"email"`
	Challenge string    `json:"challenge"`
	NotAfter  time.Time `json:"not_after"`
	CertPath  string    `json:"cert_path"`
	KeyPath   string    `json:"key_path"`
	CertPEM   []byte    `json:"-"`
}

func baseDir() string {
	if dir := os.Getenv("LECERT_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".lecert", "certs")
}

type SaveOptions struct {
	Domain    string
	Email     string
	Challenge string
}

func Save(opts SaveOptions, certPEM, keyPEM []byte) (string, time.Time, error) {
	domain := opts.Domain
	dir := filepath.Join(baseDir(), domain)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", time.Time{}, fmt.Errorf("mkdir: %w", err)
	}

	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")

	if err := writeSecure(certPath, certPEM, 0644); err != nil {
		return "", time.Time{}, err
	}
	if err := writeSecure(keyPath, keyPEM, 0600); err != nil {
		return "", time.Time{}, err
	}

	notAfter, err := parseCertExpiry(certPEM)
	if err != nil {
		return "", time.Time{}, err
	}

	meta := CertMeta{
		Domain:    domain,
		Email:     opts.Email,
		Challenge: opts.Challenge,
		NotAfter:  notAfter,
		CertPath:  certPath,
		KeyPath:   keyPath,
	}

	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
	metaPath := filepath.Join(dir, "meta.json")
	if err := writeSecure(metaPath, metaBytes, 0644); err != nil {
		return "", time.Time{}, err
	}

	return dir, notAfter, nil
}

func Load(domain string) (*CertMeta, error) {
	dir := filepath.Join(baseDir(), domain)
	metaPath := filepath.Join(dir, "meta.json")

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read meta: %w", err)
	}

	var meta CertMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	certPEM, err := os.ReadFile(meta.CertPath)
	if err != nil {
		return nil, fmt.Errorf("read cert: %w", err)
	}
	meta.CertPEM = certPEM

	return &meta, nil
}

func List() ([]CertMeta, error) {
	dir := baseDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var certs []CertMeta
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		meta, err := Load(entry.Name())
		if err != nil {
			continue
		}
		certs = append(certs, *meta)
	}
	return certs, nil
}

func writeSecure(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	if runtime.GOOS != "windows" {
		return os.Chmod(path, perm)
	}
	return nil
}

func parseCertExpiry(certPEM []byte) (time.Time, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return time.Time{}, fmt.Errorf("failed to decode PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse cert: %w", err)
	}
	return cert.NotAfter, nil
}
