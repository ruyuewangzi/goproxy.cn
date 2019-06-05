package cfg

import (
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/aofei/air"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
)

var (
	// a is the `air.Default`.
	a = air.Default

	// Zerolog is the Zerolog configuration items.
	Zerolog struct {
		// LoggerLevel is the logger level of the Zerolog.
		LoggerLevel string `mapstructure:"logger_level"`
	}

	// Qiniu is the Qiniu configuration items.
	Qiniu struct {
		// AccessKey is the access key of the Qiniu.
		AccessKey string `mapstructure:"access_key"`

		// SecretKey is the secret key of the Qiniu.
		SecretKey string `mapstructure:"secret_key"`

		// BucketName is the bucket name of the Qiniu.
		BucketName string `mapstructure:"bucket_name"`

		// BucketEndpoint is the bucket endpoint of the Qiniu.
		BucketEndpoint string `mapstructure:"bucket_endpoint"`
	}

	// Goproxy is the Goproxy configuration items.
	Goproxy struct {
		// GoBinName is the name of the Go binary of the Goproxy.
		GoBinName string `mapstructure:"go_bin_name"`

		// MaxGoBinWorkers is the maximum number of the Go binary
		// commands of the Goproxy that are allowed to execute at the
		// same time.
		MaxGoBinWorkers int `mapstructure:"max_go_bin_workers"`

		// SupportedSUMDBHosts is the supported checksum database host
		// of the Goproxy.
		SupportedSUMDBHosts []string `mapstructure:"supported_sumdb_hosts"`
	}
)

func init() {
	cf := flag.String("config", "config.toml", "configuration file")
	flag.Parse()

	m := map[string]interface{}{}
	if _, err := toml.DecodeFile(*cf, &m); err != nil {
		panic(fmt.Errorf(
			"failed to decode configuration file: %v",
			err,
		))
	}

	if err := mapstructure.Decode(m["air"], a); err != nil {
		panic(fmt.Errorf(
			"failed to decode air configuration items: %v",
			err,
		))
	}

	if err := mapstructure.Decode(m["zerolog"], &Zerolog); err != nil {
		panic(fmt.Errorf(
			"failed to decode zerolog configuration items: %v",
			err,
		))
	}

	if err := mapstructure.Decode(m["goproxy"], &Goproxy); err != nil {
		panic(fmt.Errorf(
			"failed to decode goproxy configuration items: %v",
			err,
		))
	}

	zerolog.TimeFieldFormat = ""
	switch Zerolog.LoggerLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "no":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	if a.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
