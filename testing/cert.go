package testing

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"launchpad.net/juju-core/cert"
	"time"
)

func init() {
	if err := verifyCertificates(); err != nil {
		panic(err)
	}
}

// CACert and CAKey make up a CA key pair.
// CACertX509 and CAKeyRSA hold their parsed equivalents.
var (
	CACert, CAKey = mustNewCA()

	CACertX509, CAKeyRSA = mustParseCertAndKey([]byte(CACert), []byte(CAKey))

	serverCert, serverKey = mustNewServer()
)

func verifyCertificates() error {
	_, err := tls.X509KeyPair([]byte(CACert), []byte(CAKey))
	if err != nil {
		return fmt.Errorf("bad CA cert key pair: %v", err)
	}
	_, err = tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
	if err != nil {
		return fmt.Errorf("bad server cert key pair: %v", err)
	}
	return cert.Verify([]byte(serverCert), []byte(CACert), time.Now())
}

func mustNewCA() (string, string) {
	cert.KeyBits = 512
	caCert, caKey, err := cert.NewCA("juju testing", time.Now().AddDate(10, 0, 0))
	if err != nil {
		panic(err)
	}
	return string(caCert), string(caKey)
}

func mustNewServer() (string, string) {
	cert.KeyBits = 512
	srvCert, srvKey, err := cert.NewServer("testing-env", []byte(CACert), []byte(CAKey), time.Now().AddDate(10, 0, 0))
	if err != nil {
		panic(err)
	}
	return string(srvCert), string(srvKey)
}

// CACertPEM and CAKeyPEM make up a CA key pair.
// CACertX509 and CAKeyRSA hold their parsed equivalents.
var (
	CACertPEM, CAKeyPEM = mustNewCA()

	CACertX509, CAKeyRSA = mustParseCertAndKey([]byte(CACertPEM), []byte(CAKeyPEM))

	serverCertPEM, serverKeyPEM = mustNewServer()
)

func mustParseCert(pemData string) *x509.Certificate {
	cert, err := cert.ParseCertificate([]byte(pemData))
	if err != nil {
		panic(err)
	}
	return cert
}

func mustParseCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, *rsa.PrivateKey) {
	cert, key, err := cert.ParseCertAndKey(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return cert, key
}
