// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/license"
	api "github.com/mainflux/license/api"
	"github.com/mainflux/license/postgres"
	uuid "github.com/mainflux/license/uuid"
	mainflux "github.com/mainflux/mainflux"
	authapi "github.com/mainflux/mainflux/authn/api/grpc"
	mflog "github.com/mainflux/mainflux/logger"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	defLogLevel      = "error"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defDBName        = "licenses"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defClientTLS     = "false"
	defCACerts       = ""
	defPort          = "8180"
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""
	defAuthURL       = "localhost:8181"
	defAuthTimeout   = "1" // in seconds

	envLogLevel      = "MF_LICENSE_LOG_LEVEL"
	envDBHost        = "MF_LICENSE_DB_HOST"
	envDBPort        = "MF_LICENSE_DB_PORT"
	envDBUser        = "MF_LICENSE_DB_USER"
	envDBPass        = "MF_LICENSE_DB_PASS"
	envDBName        = "MF_LICENSE_DB"
	envDBSSLMode     = "MF_LICENSE_DB_SSL_MODE"
	envDBSSLCert     = "MF_LICENSE_DB_SSL_CERT"
	envDBSSLKey      = "MF_LICENSE_DB_SSL_KEY"
	envDBSSLRootCert = "MF_LICENSE_DB_SSL_ROOT_CERT"
	envEncryptKey    = "MF_LICENSE_ENCRYPT_KEY"
	envClientTLS     = "MF_LICENSE_CLIENT_TLS"
	envCACerts       = "MF_LICENSE_CA_CERTS"
	envPort          = "MF_LICENSE_PORT"
	envServerCert    = "MF_LICENSE_SERVER_CERT"
	envServerKey     = "MF_LICENSE_SERVER_KEY"
	envJaegerURL     = "MF_JAEGER_URL"
	envAuthURL       = "MF_AUTH_URL"
	envAuthTimeout   = "MF_AUTH_TIMEOUT"
)

type config struct {
	logLevel    string
	dbConfig    postgres.Config
	clientTLS   bool
	encKey      []byte
	caCerts     string
	httpPort    string
	serverCert  string
	serverKey   string
	jaegerURL   string
	authURL     string
	authTimeout time.Duration
}

func main() {
	cfg := loadConfig()

	logger, err := mflog.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)

	authTracer, authCloser := initJaeger("auth", cfg.jaegerURL, logger)
	defer authCloser.Close()

	authConn := connectToAuth(cfg, logger)
	defer authConn.Close()

	auth := authapi.NewClient(authTracer, authConn, cfg.authTimeout)

	svc := newService(auth, db, logger, cfg)
	errs := make(chan error, 2)

	licenseTracer, licenseCloser := initJaeger("license", cfg.jaegerURL, logger)
	defer licenseCloser.Close()

	go startHTTPServer(api.MakeHandler(licenseTracer, logger, svc), cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("License service terminated: %s", err))
}

func loadConfig() config {
	tls, err := strconv.ParseBool(mainflux.Env(envClientTLS, defClientTLS))
	if err != nil {
		tls = false
	}
	dbConfig := postgres.Config{
		Host:        mainflux.Env(envDBHost, defDBHost),
		Port:        mainflux.Env(envDBPort, defDBPort),
		User:        mainflux.Env(envDBUser, defDBUser),
		Pass:        mainflux.Env(envDBPass, defDBPass),
		Name:        mainflux.Env(envDBName, defDBName),
		SSLMode:     mainflux.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     mainflux.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      mainflux.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: mainflux.Env(envDBSSLRootCert, defDBSSLRootCert),
	}

	timeout, err := strconv.ParseInt(mainflux.Env(envAuthTimeout, defAuthTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envAuthTimeout, err.Error())
	}

	return config{
		logLevel:    mainflux.Env(envLogLevel, defLogLevel),
		dbConfig:    dbConfig,
		clientTLS:   tls,
		caCerts:     mainflux.Env(envCACerts, defCACerts),
		httpPort:    mainflux.Env(envPort, defPort),
		serverCert:  mainflux.Env(envServerCert, defServerCert),
		serverKey:   mainflux.Env(envServerKey, defServerKey),
		jaegerURL:   mainflux.Env(envJaegerURL, defJaegerURL),
		authURL:     mainflux.Env(envAuthURL, defAuthURL),
		authTimeout: time.Duration(timeout) * time.Second,
	}
}

func connectToDB(cfg postgres.Config, logger mflog.Logger) *sqlx.DB {
	db, err := postgres.Connect(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func initJaeger(svcName, url string, logger mflog.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to init Jaeger client: %s", err))
		os.Exit(1)
	}

	return tracer, closer
}

func newService(auth mainflux.AuthNServiceClient, db *sqlx.DB, logger mflog.Logger, cfg config) license.Service {
	licenseRepo := postgres.New(db)
	idp := uuid.New()

	svc := license.New(licenseRepo, idp, auth)
	svc = api.NewLoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "license",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "license",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func connectToAuth(cfg config, logger mflog.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	if cfg.clientTLS {
		if cfg.caCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create tls credentials: %s", err))
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
		logger.Info("gRPC communication is not encrypted")
	}

	conn, err := grpc.Dial(cfg.authURL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to authn service: %s", err))
		os.Exit(1)
	}

	return conn
}

func startHTTPServer(h http.Handler, cfg config, logger mflog.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.httpPort)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("License service started using https on port %s with cert %s key %s",
			cfg.httpPort, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, h)
		return
	}
	logger.Info(fmt.Sprintf("License service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, h)
}
