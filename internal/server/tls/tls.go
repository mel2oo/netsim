package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"
)

// load certificate
func Load(certfile, keyfile, cafile string) (*tls.Config, error) {
	if len(certfile) != 0 && len(keyfile) != 0 {
		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			return nil, err
		}

		cp, err := loadCA(cafile)
		if err != nil {
			return nil, err
		}

		return &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    cp,
		}, nil
	} else {
		cert, err := GenCertificate()
		if err != nil {
			return nil, err
		}

		return &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}, nil
	}
}

func loadCA(cafile string) (*x509.CertPool, error) {
	cp := x509.NewCertPool()

	data, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}

	if !cp.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("AppendCertsFromPEM")
	}

	return cp, nil
}

// generate certificate
func GenCertificate() (*tls.Certificate, error) {
	rawCert, rawKey, err := generateKeyPair()
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(rawCert, rawKey)
	return &cert, err
}

func generateKeyPair() (rawCert, rawKey []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	validFor := time.Hour * 24 * 365 * 10
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"netsim"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return
	}

	rawCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	rawKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return
}
