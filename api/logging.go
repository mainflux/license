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

func (lm *loggingMiddleware) CreateLicense(ctx context.Context, token string, l license.License) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_license with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateLicense(ctx, token, l)
}

func (lm *loggingMiddleware) RetrieveLicense(ctx context.Context, owner, id string) (l license.License, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method retrieve_license with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveLicense(ctx, owner, id)
}

func (lm *loggingMiddleware) UpdateLicense(ctx context.Context, token string, l license.License) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_license with ID %s took %s to complete", l.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateLicense(ctx, token, l)
}

func (lm *loggingMiddleware) RemoveLicense(ctx context.Context, owner, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_license with ID %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveLicense(ctx, owner, id)
}
