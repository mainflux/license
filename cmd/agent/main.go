// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mainflux/license/agent"
	"github.com/mainflux/license/agent/api"
	"github.com/mainflux/license/pkg/crypto"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
)

const (
	defLogLevel    = "error"
	defSvcURL      = "http://localhost:8180/licenses"
	defLicenseFile = "./license"
	defClientTLS   = "false"
	defServerCert  = ""
	defServerKey   = ""
	defPort        = "3000"

	envLogLevel    = "MF_LICENSE_LOG_LEVEL"
	envSvcURL      = "MF_LICENSE_SERVICE_URL"
	envLicenseFile = "LICENSE_FILE"
	envClientTLS   = "MF_LICENSE_CLIENT_TLS"
	envServerCert  = "MF_LICENSE_SERVER_CERT"
	envServerKey   = "MF_AGENT_SERVER_KEY"
	envPort        = "MF_AGENT_PORT"
)

type config struct {
	logLevel    string
	svcURL      string
	licenseFile string
	licenseID   string
	licenseKey  string
	tls         bool
	serverCert  string
	serverKey   string
	port        string
}

func main() {
	cfg := loadConfig()

	logger, err := mflog.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	cfg.licenseID = addresses(logger)
	c := crypto.New()
	a := agent.New(cfg.svcURL, cfg.licenseFile, cfg.licenseID, cfg.licenseKey, c, nil)
	a = api.NewLoggingMiddleware(a, logger)
	go a.Do()
	for {
		logger.Info("Loading the license...")
		if err := a.Load(); err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}

	errs := make(chan error, 2)

	go startHTTPServer(api.MakeHandler(logger, a), cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("License agent terminated: %s", err))
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
		serverKey:   mainflux.Env(envServerKey, defServerKey),
		port:        mainflux.Env(envPort, defPort),
	}
}

func startHTTPServer(h http.Handler, cfg config, logger mflog.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.port)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("License agent started using https on port %s with cert %s key %s",
			cfg.port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, h)
		return
	}
	logger.Info(fmt.Sprintf("License agent started using http on port %s", cfg.port))
	errs <- http.ListenAndServe(p, h)
}

func addresses(logger mflog.Logger) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn(fmt.Sprintf("Unable to read id: %s", err))
		return ""
	}
	sort(ifaces)

	// Add the rest of the interfaces with HardwareAddr in Index order.
	// Iterate over the list to preserve order.
	for _, ifc := range ifaces {
		addr := ifc.HardwareAddr.String()
		if addr != "" && strings.HasPrefix(ifc.Name, "e") {
			return addr
		}
	}

	logger.Warn(fmt.Sprintf("License id not found"))
	return ""
}

// Sort sorts network interfaces by index.
func sort(ifaces []net.Interface) {
	n := len(ifaces)
	for i := range ifaces {
		swap := i
		for j := i; j < n; j++ {
			if ifaces[j].Name < ifaces[swap].Name {
				swap = j
			}
		}
		if i != swap {
			ifaces[i], ifaces[swap] = ifaces[swap], ifaces[i]
		}
	}
}
