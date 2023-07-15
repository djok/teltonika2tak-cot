package main

import (
	"context"
	"flag"
	"fmt"

	"teltonika2tak-cot/config"
	"teltonika2tak-cot/tmt250"

	// influxdb2 "github.com/halacs/haltonika/influxdb"
	// "github.com/halacs/haltonika/messaging"
	// m "github.com/halacs/haltonika/metrics"
	// mi "github.com/halacs/haltonika/metrics/impl"
	"os"
	"os/signal"
	"strings"
	"sync"

	// "teltonika2tak-cot/tmt250"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	// InfluxDB client configs
	// Teltonika server configs
	flag.String(config.TeltonikaListeningIp, config.DefaultTeltonikaListeningIP, "Teltonika server listening IP address (IPv4 or IPv6)")
	flag.Int(config.TeltonikaListeningPort, config.DefaultTeltonikaListeningPort, "Teltonika server listening UDP port")
	// Metrics server configs
	flag.String(config.MetricsListeningIp, config.DefaultMetricsListeningIP, "Metrics server listening IP address (IPv4 or IPv6)")
	// flag.Int(config.MetricsListeningPort, config.DefaultMetricsListeningPort, "Metrics server listening port")
	// flag.String(config.MetricsTeltonikaMetricsFileName, config.DefaultMetricsTeltonikaMetricsFileName, "File where metrics are written")

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
		Host:         viper.GetString(config.TeltonikaListeningIp),
		Port:         viper.GetInt(config.TeltonikaListeningPort),
		AllowedIMEIs: allowedIMEIs,
	}

	// metricsConfig := &config.MetricsConfig{
	// 	Host:                     viper.GetString(config.MetricsListeningIp),
	// 	Port:                     viper.GetInt(config.MetricsListeningPort),
	// 	TeltonikaMetricsFileName: viper.GetString(config.MetricsTeltonikaMetricsFileName),
	// }

	cfg := config.NewConfig(log, teltonikaConfig)
	return cfg
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
