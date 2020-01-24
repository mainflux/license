// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/mainflux/license/agent"
	"github.com/mainflux/license/agent/api"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
)

const (
	defLogLevel    = "error"
	defSvcURL      = "http://localhost:8180"
	defLicenseFile = "./license"
	defClientTLS   = "false"
	defServerCert  = ""

	envLogLevel    = "MF_LICENSE_LOG_LEVEL"
	envSvcURL      = "MF_LICENSE_SERVICE_URL"
	envLicenseFile = "LICENSE_FILE"
	envClientTLS   = "MF_LICENSE_CLIENT_TLS"
	envServerCert  = "MF_LICENSE_SERVER_CERT"
)

type config struct {
	logLevel    string
	svcURL      string
	licenseFile string
	tls         bool
	serverCert  string
}

func main() {
	cfg := loadConfig()

	logger, err := mflog.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	a := agent.New(cfg.svcURL, cfg.licenseFile)
	a = api.NewLoggingMiddleware(a, logger)
	go a.Do()
}

func loadConfig() config {
	tls, err := strconv.ParseBool(mainflux.Env(envClientTLS, defClientTLS))
	if err != nil {
		tls = false
	}

	return config{
		svcURL:      mainflux.Env(envSvcURL, defSvcURL),
		logLevel:    mainflux.Env(envLogLevel, defLogLevel),
		tls:         tls,
		licenseFile: mainflux.Env(envLicenseFile, defLicenseFile),
		serverCert:  mainflux.Env(envServerCert, defServerCert),
	}
}
