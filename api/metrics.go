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

func (ms *metricsMiddleware) Create(ctx context.Context, token string, l license.License) (id string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create").Add(1)
		ms.latency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Create(ctx, token, l)
}

func (ms *metricsMiddleware) Retrieve(ctx context.Context, owner, id string) (l license.License, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve").Add(1)
		ms.latency.With("method", "retrieve").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Retrieve(ctx, owner, id)
}

func (ms *metricsMiddleware) Update(ctx context.Context, token string, l license.License) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update").Add(1)
		ms.latency.With("method", "update").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Update(ctx, token, l)
}

func (ms *metricsMiddleware) Remove(ctx context.Context, owner, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove").Add(1)
		ms.latency.With("method", "remove").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Remove(ctx, owner, id)
}

func (ms *metricsMiddleware) ChangeActive(ctx context.Context, token, id string, active bool) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "change_active").Add(1)
		ms.latency.With("method", "change_active").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ChangeActive(ctx, token, id, active)
}

func (ms *metricsMiddleware) Validate(ctx context.Context, svc, id string, payload []byte) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "validate").Add(1)
		ms.latency.With("method", "validate").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Validate(ctx, svc, id, payload)
}
