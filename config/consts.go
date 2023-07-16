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
	AppName                       = "teltonika2tak-cot"
	ViperEnvPrefix                = AppName
	Verbose                       = "verbose"
	Debug                         = "debug"
	AllowedIMEIs                  = "imeilist"
	TeltonikaListeningIP          = "teltonikalistenip"
	TeltonikaListeningPort        = "teltonikalistenport"
	TakHostIP                     = "takserverip"
	TakHostPort                   = "takserverport"
	TakHostProtocol               = "takserverprotocol"
	DefaultDebug                  = true
	DefaultVerbose                = true
	DefaultAllowedIMEIs           = "352016701836447" // list, separated by comma
	DefaultTeltonikaListeningIP   = "0.0.0.0"
	DefaultTeltonikaListeningPort = 7809
	DefaultTakHostIP              = "127.0.0.1"
	DefaultTakHostPort            = 8087
	DefaultTakHostProtocol        = "tcp"
)
