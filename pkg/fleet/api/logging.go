/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/fleet"
	"go.uber.org/zap"
)

var _ fleet.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    fleet.Service
}

func (l loggingMiddleware) Add() (fleet.Agent, error) {
	panic("implement me")
}

func NewLoggingMiddleware(svc fleet.Service, logger *zap.Logger) fleet.Service {
	return &loggingMiddleware{logger, svc}
}
