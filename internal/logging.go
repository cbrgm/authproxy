/*
 * Copyright 2019, authproxy authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package internal

import (
	"github.com/cbrgm/authproxy/api/v1/models"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type loggingService struct {
	logger  log.Logger
	service Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger: logger, service: s}
}

func (s *loggingService) Login(username, password string) (*models.TokenReviewRequest, error) {
	start := time.Now()

	tkn, err := s.service.Login(username, password)

	logger := log.With(s.logger,
		"method", "Login",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log("msg", "failed to login user", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return tkn, err
}

func (s *loggingService) Authenticate(bearerToken string) (*models.TokenReviewRequest, error) {
	start := time.Now()

	trr, err := s.service.Authenticate(bearerToken)

	logger := log.With(s.logger,
		"method", "Authenticate",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log("msg", "failed to authenticate token", "err", err)
	} else {
		level.Debug(logger).Log()
	}

	return trr, err
}
