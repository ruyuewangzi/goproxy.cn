package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/air-gases/limiter"
	"github.com/aofei/air"
	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy.cn/cfg"
	"github.com/goproxy/goproxy/cachers"
	"github.com/rs/zerolog/log"
)

// a is the `air.Default`.
var a = air.Default

func init() {
	if err := os.Setenv("GO111MODULE", "on"); err != nil {
		log.Fatal().Err(err).
			Str("app_name", a.AppName).
			Msg("failed to set $GO111MODULE")
	}

	if err := os.Setenv("GOPROXY", "direct"); err != nil {
		log.Fatal().Err(err).
			Str("app_name", a.AppName).
			Msg("failed to set $GOPROXY")
	}

	if err := os.Setenv("GOSUMDB", "off"); err != nil {
		log.Fatal().Err(err).
			Str("app_name", a.AppName).
			Msg("failed to set $GOSUMDB")
	}

	g := goproxy.New()
	g.GoBinName = cfg.Goproxy.GoBinName
	g.MaxGoBinWorkers = cfg.Goproxy.MaxGoBinWorkers
	g.Cacher = &cachers.Kodo{
		AccessKey:      cfg.Qiniu.AccessKey,
		SecretKey:      cfg.Qiniu.SecretKey,
		BucketName:     cfg.Qiniu.BucketName,
		BucketEndpoint: cfg.Qiniu.BucketEndpoint,
	}

	g.SupportedSUMDBHosts = cfg.Goproxy.SupportedSUMDBHosts
	g.ErrorLogger = a.ErrorLogger

	a.BATCH(
		[]string{http.MethodGet, http.MethodHead},
		"/",
		indexPageHandler,
	)
	a.BATCH(
		nil,
		"/*",
		air.WrapHTTPHandler(g),
		limiter.RateGas(limiter.RateGasConfig{
			MaxRequests:   1200,
			ResetInterval: time.Hour,
		}),
	)
}

// indexPageHandler handles requests to get index page.
func indexPageHandler(req *air.Request, res *air.Response) error {
	return res.Redirect("https://github.com/goproxy/goproxy.cn")
}
