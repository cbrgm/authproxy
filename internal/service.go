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
	"github.com/cbrgm/authproxy/provider"
)

// Service represents the middleware used by authproxy
type Service interface {
	Login(username, password string) (*models.TokenReviewRequest, error)
	Authenticate(bearerToken string) (*models.TokenReviewRequest, error)
}

// service represents the middleware implementation
type service struct {
	provider provider.Provider
}

// NewService returns a new middleware service configured with the given provider
func NewService(prv *provider.Provider) *service {
	return &service{
		provider: *prv,
	}
}

// Login wraps the provider specific login implementation
func (s *service) Login(username, password string) (*models.TokenReviewRequest, error) {
	return s.provider.Login(username, password)
}

// Authenticate wraps the provider specific authentication implementation
func (s *service) Authenticate(bearerToken string) (*models.TokenReviewRequest, error) {
	return s.provider.Authenticate(bearerToken)
}
