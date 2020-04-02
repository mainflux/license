// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mainflux/license/errors"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/license"
	"github.com/mainflux/mainflux"
	log "github.com/mainflux/mainflux/logger"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const contentType = "application/json"

var (
	errUnsupportedContentType = errors.New("unsupported content type")

	logger log.Logger
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(tracer opentracing.Tracer, l log.Logger, svc license.Service) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/licenses", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_license")(createEndpoint(svc)),
		decodeCreation,
		encodeResponse,
		opts...,
	))

	r.Get("/licenses/view/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_license")(viewEndpoint(svc)),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.Get("/licenses/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "fetch_license")(fetchEndpoint(svc)),
		decodeView,
		encodeFetch,
		opts...,
	))

	r.Get("/licenses/devices", kithttp.NewServer(
		kitot.TraceServer(tracer, "fetch_by_device_id")(viewByDeviceIDEndpoint(svc)),
		decodeView,
		encodeFetch,
		opts...,
	))

	r.Patch("/licenses/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_license")(updateEndpoint(svc)),
		decodeUpdate,
		encodeResponse,
		opts...,
	))

	r.Delete("/licenses/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_license")(removeEndpoint(svc)),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.Post("/licenses/validate/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "validate_license")(validateEndpoint(svc)),
		decodeValidate,
		encodeResponse,
		opts...,
	))

	r.Patch("/licenses/activation/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "activation_license")(activationEndpoint(svc, true)),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.Delete("/licenses/activation/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "activation_license")(activationEndpoint(svc, false)),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", mainflux.Version("license"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := createReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeView(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := licenseReq{
		token: r.Header.Get("Authorization"),
	}

	return req, nil
}

func decodeValidate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := validateReq{
		service: r.URL.Query().Get("service"),
		id:      bone.GetValue(r, "id"),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeFetch(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err := w.Write(response.([]byte))
	return err
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
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
		case errors.Contains(e, license.ErrConflict):
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		if err := json.NewEncoder(w).Encode(errorRes{Err: e.Msg()}); err != nil {
			logger.Warn(fmt.Sprintf("failed to send error response %s", err))
		}
	case *json.SyntaxError, *json.UnmarshalTypeError:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
