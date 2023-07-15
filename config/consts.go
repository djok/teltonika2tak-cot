package config

type MyKey struct {
	KeyName string
}

var (
	ContextConfigKey = MyKey{
		KeyName: "config",
	}
)

const (
	AppName                         = "haltonika"
	ViperEnvPrefix                  = AppName
	Verbose                         = "verbose"
	Debug                           = "debug"
	AllowedIMEIs                    = "imeilist"
	InfluxConfigUrl                 = "url"
	InfluxConfigUsername            = "username"
	InfluxConfigPassword            = "password"
	InfluxConfigDatabase            = "database"
	InfluxConfigMeasurement         = "measurement"
	TeltonikaListeningIp            = "listenip"
	TeltonikaListeningPort          = "listenport"
	MetricsListeningIp              = "metricsip"
	MetricsListeningPort            = "metricsport"
	MetricsTeltonikaMetricsFileName = "mp"
	DefaultDebug                    = true
	DefaultVerbose                  = true
	DefaultAllowedIMEIs             = "352016701836447" // list, separated by comma
	DefaultTeltonikaListeningIP     = "0.0.0.0"
	DefaultTeltonikaListeningPort   = 7809
	DefaultMetricsListeningIP       = "0.0.0.0"
	// DefaultMetricsListeningPort            = 9161
	// DefaultMetricsTeltonikaMetricsFileName = AppName + ".met"
)
