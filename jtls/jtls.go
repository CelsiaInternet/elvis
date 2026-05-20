package jtls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/file"
	"github.com/celsiainternet/elvis/logs"
)

/**
* Create
* @param fileCrt string
* @param fileKey string
* @param hosts []string
* @param expire time.Duration
* @return error
**/
func Create(fileCrt, fileKey string, hosts []string, expire time.Duration) error {
	logs.Logf("pipe", "generate certificates TLS...")

	file.RemoveFiles(fileCrt, fileKey)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	dnsNames := []string{"localhost"}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1")}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			ipAddresses = append(ipAddresses, ip)
		} else if h != "" && h != "localhost" {
			dnsNames = append(dnsNames, h)
		}
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),

		Subject: pkix.Name{
			CommonName: dnsNames[0],
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(expire),

		IsCA:                  true,
		BasicConstraintsValid: true,

		KeyUsage: x509.KeyUsageCertSign |
			x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment,

		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},

		DNSNames:    dnsNames,
		IPAddresses: ipAddresses,
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&priv.PublicKey,
		priv,
	)
	if err != nil {
		return err
	}

	certOut, err := os.Create(fileCrt)
	if err != nil {
		return err
	}

	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	certOut.Close()

	keyOut, err := os.Create(fileKey)
	if err != nil {
		return err
	}

	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	keyOut.Close()

	return nil
}

/**
* Load
* @param path string
* @param hosts []string
* @param expire time.Duration
* @return (tls.Certificate, error)
**/
func Load(path string, hosts []string, expire time.Duration) (tls.Certificate, error) {
	if !file.ExistPath(path) {
		_, err := file.MakeFolder(path)
		if err != nil {
			return tls.Certificate{}, err
		}
	}

	fileCrt := filepath.Join(path, "server.crt")
	fileKey := filepath.Join(path, "server.key")
	if file.ExistPath(fileCrt) && file.ExistPath(fileKey) {
		cert, err := tls.LoadX509KeyPair(fileCrt, fileKey)
		if err != nil {
			return tls.Certificate{}, err
		}
		return cert, nil
	}

	err := Create(fileCrt, fileKey, hosts, expire)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert, err := tls.LoadX509KeyPair(fileCrt, fileKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}

/**
* Pool
* @param path string
* @param hosts []string
* @param expire time.Duration
* @return (*x509.CertPool, error)
**/
func Pool(path string, hosts []string, expire time.Duration) (*x509.CertPool, error) {
	if !file.ExistPath(path) {
		_, err := file.MakeFolder(path)
		if err != nil {
			return nil, err
		}
	}

	fileCrt := filepath.Join(path, "server.crt")
	fileKey := filepath.Join(path, "server.key")
	if file.ExistPath(fileCrt) && file.ExistPath(fileKey) {
		certPool := x509.NewCertPool()
		certData, err := os.ReadFile(fileCrt)
		if err != nil {
			return nil, err
		}
		ok := certPool.AppendCertsFromPEM(certData)
		if !ok {
			return nil, fmt.Errorf("failed to append certificate")
		}
		return certPool, nil
	}

	err := Create(fileCrt, fileKey, hosts, expire)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certData, err := os.ReadFile(fileCrt)
	if err != nil {
		return nil, err
	}
	ok := certPool.AppendCertsFromPEM(certData)
	if !ok {
		return nil, fmt.Errorf("failed to append certificate")
	}
	return certPool, nil
}

/**
* Deal
* @param path string
* @param host string
* @param port int
* @param expire time.Duration
* @return (*tls.Conn, error)
**/
func Deal(path, host string, port int, expire time.Duration) (*tls.Conn, error) {
	cert, err := Pool(path, []string{host}, expire)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:            cert,
		InsecureSkipVerify: envar.GetBool(false, "PIPE_INSECURE_SKIP_VERIFY"),
	}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

/**
* Wrapper
* @param path string
* @param hosts []string
* @param owner net.Listener
* @param expire time.Duration
* @return net.Listener
**/
func Wrapper(path string, hosts []string, owner net.Listener, expire time.Duration) net.Listener {
	cert, err := Load(path, hosts, expire)
	if err != nil {
		logs.Panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: envar.GetBool(false, "PIPE_INSECURE_SKIP_VERIFY"),
	}

	return tls.NewListener(owner, tlsConfig)
}
