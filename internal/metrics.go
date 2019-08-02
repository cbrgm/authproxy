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
	"github.com/go-kit/kit/metrics"
)

type metricsService struct {
	loginAttempts metrics.Counter
	service       Service
}

func NewMetricsService(loginAttempts metrics.Counter, service Service) Service {
	// Initialize counters with 0
	loginAttempts.With("status", "failure").Add(0)
	loginAttempts.With("status", "success").Add(0)

	return &metricsService{loginAttempts: loginAttempts, service: service}
}

func (s *metricsService) Login(username, password string) (*models.TokenReviewRequest, error) {
	trr, err := s.service.Login(username, password)

	if err != nil || !trr.Status.Authenticated {
		s.loginAttempts.With("status", "failure").Add(1)
	} else {
		s.loginAttempts.With("status", "success").Add(1)
	}

	return trr, err
}

func (s *metricsService) Authenticate(bearerToken string) (*models.TokenReviewRequest, error) {
	// Don't do anything here
	return s.service.Authenticate(bearerToken)
}
