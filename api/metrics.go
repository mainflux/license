// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/license"
)

var _ license.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     license.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc license.Service, counter metrics.Counter, latency metrics.Histogram) license.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateLicense(ctx context.Context, token string, l license.License) (id string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_license").Add(1)
		ms.latency.With("method", "create_license").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateLicense(ctx, token, l)
}

func (ms *metricsMiddleware) RetrieveLicense(ctx context.Context, owner, id string) (l license.License, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_license").Add(1)
		ms.latency.With("method", "retrieve_license").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveLicense(ctx, owner, id)
}

func (ms *metricsMiddleware) UpdateLicense(ctx context.Context, token string, l license.License) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_license").Add(1)
		ms.latency.With("method", "update_license").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateLicense(ctx, token, l)
}

func (ms *metricsMiddleware) RemoveLicense(ctx context.Context, owner, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_license").Add(1)
		ms.latency.With("method", "remove_license").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveLicense(ctx, owner, id)
}
