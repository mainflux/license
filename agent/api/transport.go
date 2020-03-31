// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/license/errors"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/license"
	"github.com/mainflux/mainflux"
	log "github.com/mainflux/mainflux/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger log.Logger

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(l log.Logger, agent license.Agent) http.Handler {
	logger = l
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()

	r.Post("/licenses/validate", kithttp.NewServer(
		validateEndpoint(agent),
		decodeValidate,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", mainflux.Version("license-agent"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeValidate(_ context.Context, r *http.Request) (interface{}, error) {
	return ioutil.ReadAll(r.Body)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err := w.Write(response.([]byte))
	return err
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case io.ErrUnexpectedEOF, io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	}
	switch e := err.(type) {
	case errors.Error:
		switch {
		case errors.Contains(e, license.ErrMalformedEntity), errors.Contains(e, license.ErrExpired):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(e, license.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusForbidden)
		case errors.Contains(e, license.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		}
		if err := json.NewEncoder(w).Encode(errorRes{Err: e.Msg()}); err != nil {
			logger.Warn(fmt.Sprintf("failed to send error response %s", err))
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
