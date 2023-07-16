package tak

import (
	"context"
	"crypto"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"teltonika2tak-cot/config"

	"github.com/spf13/viper"
	"software.sslmate.com/src/go-pkcs12"
)

func connect(ctx context.Context, takHost string, takPort string, takTls bool) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%s", takHost, takPort)
	log := config.GetLogger(ctx)

	if takTls {
		log.Infof("connecting with SSL to %s...", addr)
		// tlsCert := loadCerts(ctx)

		log.Infof("load cert from %s", viper.GetString("ssl.cert"))
		p12Data, err := os.ReadFile(viper.GetString("ssl.cert"))
		if err != nil {
			log.Fatal(err)
		}

		key, cert, _, err := pkcs12.DecodeChain(p12Data, viper.GetString("ssl.password"))
		if err != nil {
			log.Fatal(err)
		}

		tlsCert := &tls.Certificate{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key.(crypto.PrivateKey),
			Leaf:        cert,
		}

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{*tlsCert}, InsecureSkipVerify: true}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return nil, err
		}
		log.Debugf("handshake...")

		if err := conn.Handshake(); err != nil {
			return conn, err
		}
		cs := conn.ConnectionState()

		log.Infof("Handshake complete: %t", cs.HandshakeComplete)
		log.Infof("version: %d", cs.Version)
		for i, cert := range cs.PeerCertificates {
			log.Infof("cert #%d subject: %s", i, cert.Subject.String())
			log.Infof("cert #%d issuer: %s", i, cert.Issuer.String())
			log.Infof("cert #%d dns_names: %s", i, strings.Join(cert.DNSNames, ","))
		}
		return conn, nil
	} else {
		log.Infof("connecting to %s...", addr)
		return net.DialTimeout("tcp", addr, time.Second*5)
	}
}

// func getTlsConfig(ctx context.Context) *tls.Config {
// 	tlsCert := loadCerts(ctx)
// 	return &tls.Config{Certificates: []tls.Certificate{*tlsCert}, InsecureSkipVerify: true}
// }

// func loadCerts(ctx context.Context) *tls.Certificate {
// 	log := config.GetLogger(ctx)
// 	log.Infof("load cert from %s", viper.GetString("ssl.cert"))
// 	p12Data, err := os.ReadFile(viper.GetString("ssl.cert"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	key, cert, _, err := pkcs12.DecodeChain(p12Data, viper.GetString("ssl.password"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tlsCert := &tls.Certificate{
// 		Certificate: [][]byte{cert.Raw},
// 		PrivateKey:  key.(crypto.PrivateKey),
// 		Leaf:        cert,
// 	}
// 	return tlsCert
// }
