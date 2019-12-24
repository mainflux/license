// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/mainflux/license"
	log "github.com/mainflux/mainflux/logger"
)

var _ license.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    license.Service
}

// NewLoggingMiddleware adds logging facilities to the core service.
func NewLoggingMiddleware(svc license.Service, logger log.Logger) license.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Create(ctx context.Context, token string, l license.License) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Create(ctx, token, l)
}

func (lm *loggingMiddleware) Retrieve(ctx context.Context, owner, id string) (l license.License, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method retrieve with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Retrieve(ctx, owner, id)
}

func (lm *loggingMiddleware) Fetch(ctx context.Context, key, id string) (l license.License, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method fetch with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Fetch(ctx, key, id)
}

func (lm *loggingMiddleware) Update(ctx context.Context, token string, l license.License) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update with ID %s took %s to complete", l.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Update(ctx, token, l)
}

func (lm *loggingMiddleware) Remove(ctx context.Context, owner, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Remove(ctx, owner, id)
}

func (lm *loggingMiddleware) ChangeActive(ctx context.Context, token, id string, active bool) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method change_active to %t with ID %s took %s to complete", active, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ChangeActive(ctx, token, id, active)
}

func (lm *loggingMiddleware) Validate(ctx context.Context, svc, id string, payload []byte) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method validate for license %s and service %s took %s to complete", id, svc, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.Validate(ctx, svc, id, payload)
}
