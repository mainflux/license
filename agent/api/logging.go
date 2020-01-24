// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"fmt"
	"github.com/mainflux/license"
	log "github.com/mainflux/mainflux/logger"
	"time"
)

var _ license.Agent = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	agent  license.Agent
}

// NewLoggingMiddleware adds logging facilities to the Licensing agent.
func NewLoggingMiddleware(a license.Agent, logger log.Logger) license.Agent {
	return &loggingMiddleware{logger, a}
}

func (lm *loggingMiddleware) Load() (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method load took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.agent.Load()
}

func (lm *loggingMiddleware) Save() (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method save took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.agent.Save()
}

func (lm *loggingMiddleware) Validate(services []string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method validate for service %v took %s to complete", services, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.agent.Validate(services)
}

func (lm *loggingMiddleware) Do() {
	lm.agent.Do()
}
