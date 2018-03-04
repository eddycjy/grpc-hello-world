package util

import (
	"crypto/tls"
    "io/ioutil"
    "log"
    "net"

    "golang.org/x/net/http2"
)

func GetTLSConfig(certPemPath, certKeyPath string) *tls.Config {
    var certKeyPair *tls.Certificate
    cert, _ := ioutil.ReadFile(certPemPath)
    key, _ := ioutil.ReadFile(certKeyPath)
    
    pair, err := tls.X509KeyPair(cert, key)
    if err != nil {
        log.Println("TLS KeyPair err: %v\n", err)
    }
    
    certKeyPair = &pair

    return &tls.Config{
        Certificates: []tls.Certificate{*certKeyPair},
        NextProtos:   []string{http2.NextProtoTLS},
    }
}

func NewTLSListener(inner net.Listener, config *tls.Config) net.Listener {
    return tls.NewListener(inner, config)
}