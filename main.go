package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"net"
	"time"

	"os"
	"os/signal"
	"strings"
	"sync"
	"teltonika2tak-cot/config"
	"teltonika2tak-cot/tmt250"

	"github.com/kdudkov/goatak/cot"
	"github.com/kdudkov/goatak/cotproto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"crypto"
	"crypto/tls"

	"software.sslmate.com/src/go-pkcs12"
)

func parseConfig() *config.Config {
	// Initialize logger
	log := config.NewLogger()

	// Read configuration
	viper.SetConfigName("cfg")                                     // Name of cfg file (without extension)
	viper.SetConfigType("yaml")                                    // REQUIRED if the cfg file does not have the extension in the name
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", config.AppName))   // path to look for the cfg file in
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", config.AppName)) // call multiple times to add many search paths
	viper.AddConfigPath(".")                                       // Optionally look for cfg in the working directory
	viper.SetEnvPrefix(config.ViperEnvPrefix)
	viper.AutomaticEnv() // Use environment variables if defined

	viper.SetDefault("cot.proto", "tcp")
	viper.SetDefault("cot.type", "a-n-G")
	viper.SetDefault("cot.stale", time.Minute*10)

	err := viper.ReadInConfig() // Find and read the cfg file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Infof("Config file was not found. Using defaults.")
	} else {
		log.Fatalf("Failed to parse cfg file. %v", err)
	}

	// General configs
	flag.Bool(config.Debug, config.DefaultDebug, "Set log level to debug")
	flag.Bool(config.Verbose, config.DefaultVerbose, "Set log level to verbose")
	flag.String(config.AllowedIMEIs, config.DefaultAllowedIMEIs, "IMEI identifiers needs to be processed. Separated by comma. Example: 123456789012345,123456789012345,123456789012345")
	// Teltonika server configs
	flag.String(config.TeltonikaListeningIP, config.DefaultTeltonikaListeningIP, "Teltonika server listening IP address (IPv4 or IPv6)")
	flag.Int(config.TeltonikaListeningPort, config.DefaultTeltonikaListeningPort, "Teltonika server listening UDP port")
	// Tak server configs
	flag.String(config.TakHostIP, config.DefaultTakHostIP, "Tak server IP address (IPv4 or IPv6)")
	flag.Int(config.TakHostPort, config.DefaultTakHostPort, "Tak server TCP port")
	flag.String(config.TakHostProtocol, config.DefaultTakHostProtocol, "Tak server Protocol (tcp or udp)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Errorf("Failed to bindPFlags. %v", err)
	}

	verbose := viper.GetBool(config.Verbose)
	debug := viper.GetBool(config.Debug)
	if verbose {
		log.SetLevel(logrus.TraceLevel)
		log.Warningf("Active log level: %s", log.GetLevel())
	} else if debug {
		log.SetLevel(logrus.DebugLevel)
		log.Warningf("Active log level: %s", log.GetLevel())
	}

	allowedIMEIs := strings.Split(viper.GetString(config.AllowedIMEIs), ",")

	teltonikaConfig := &config.TeltonikaConfig{
		Host:         viper.GetString(config.TeltonikaListeningIP),
		Port:         viper.GetInt(config.TeltonikaListeningPort),
		AllowedIMEIs: allowedIMEIs,
	}

	takConfig := &config.TakConfig{
		Host:     viper.GetString(config.TakHostIP),
		Port:     viper.GetInt(config.TakHostPort),
		Protocol: viper.GetString(config.TakHostProtocol),
	}

	cfg := config.NewConfig(log, teltonikaConfig, takConfig)
	return cfg
}

func makeEvent(id, callSign string, lat, lon float64) *cot.Event {
	evt := cot.BasicMsg(viper.GetString("cot.type"), id, viper.GetDuration("cot.stale"))
	evt.CotEvent.How = "a-g"
	evt.CotEvent.Detail = &cotproto.Detail{
		Contact: &cotproto.Contact{Callsign: callSign},
	}
	evt.CotEvent.Lon = lon
	evt.CotEvent.Lat = lat

	return cot.ProtoToEvent(evt)
}

func sendCotMessage(ctx context.Context, evt *cot.Event, cotServer string, cotPort int, cotProtocol string) {
	// uri := cotServer + ":" + fmt.Sprint(cotPort)
	log := config.NewLogger()
	msg, err := xml.Marshal(evt)
	if err != nil {
		log.Errorf("marshal error: %v", err)
		return
	}

	conn, err := connect2tak(cotServer, cotPort, true)
	// conn, err := net.Dial(cotProtocol, uri)
	if err != nil {
		log.Errorf("connection error: %v", err)
		return
	}

	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	if _, err := conn.Write(msg); err != nil {
		log.Errorf("write error: %v", err)
	}
	_ = conn.Close()
}

func connect2tak(takHost string, takPort int, takTls bool) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", takHost, takPort)
	log := config.NewLogger()

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

func main() {
	var wg sync.WaitGroup

	cfg := parseConfig()

	log := cfg.GetLogger()
	log.Tracef("Used Teltonika server configuration: %+v", cfg.GetTeltonikaConfig())

	// Initialize context
	ctxSignals, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx := context.WithValue(context.Background(), config.ContextConfigKey, cfg)

	server := tmt250.NewServer(ctx, &wg, cfg.GetTeltonikaConfig().Host, cfg.GetTeltonikaConfig().Port, cfg.GetTeltonikaConfig().AllowedIMEIs, func(ctx context.Context, message tmt250.TeltonikaMessage) {
		log := cfg.GetLogger()

		log.Debugf("PACKET ARRIVED: %+v", message)
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.OpenFile("output.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if _, err := file.Write(jsonMessage); err != nil {
			log.Fatal(err)
		}
		if err := file.Sync(); err != nil {
			log.Fatal(err)
		}

		// evt := makeEvent(
		// 					fmt.Sprintf("tg-%d", message.From.ID),
		// 					fmt.Sprintf("tg-%s", message.From.UserName),
		// 					loc.Latitude,
		// 					loc.Longitude)
		// Forward data internally for further processing
		// messenger.Publish(message)
	})

	// Start Teltonika server
	err := server.Start()
	if err != nil {
		log.Errorf("Failed to start Teltonika server. %v", err)
	}

	<-ctxSignals.Done()
	log.Infof("Exiting")
	wg.Wait()
}
