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

package fake

import (
	"github.com/cbrgm/authproxy/api/v1/models"
)

// FakeProvider represents a fake identity provider
type FakeProvider struct {
	Name string
}

// NewFakeProvider returns a new fake identity provider
func NewFakeProvider() *FakeProvider {
	return &FakeProvider{
		Name: "fake-authenticator",
	}
}

// Login implements login functionality for user foo and password bar
func (provider *FakeProvider) Login(username, password string) (*models.TokenReviewRequest, error) {
	var isAuthenticated = false
	var bearerToken = ""

	if username == "foo" && password == "bar" {
		isAuthenticated = true

		// generate bearerToken to be used for authentication
		bearerToken = "AbCdEf123456"
	}

	return &models.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Status: &models.TokenReviewStatus{
			// Required: let the client know if the user has successfully authenticated
			Authenticated: isAuthenticated,

			// optional: add user information for the client
			User: &models.UserInfo{
				Username: username,
				UID:      "1",
				Groups:   []string{"developers"},
			},
		},
		// Required: return the token for the client
		Spec: &models.TokenReviewSpec{
			Token: bearerToken,
		},
	}, nil
}

// Authenticate implements bearer token validation functionalities
func (provider *FakeProvider) Authenticate(bearerToken string) (*models.TokenReviewRequest, error) {
	var isTokenValid = false

	if bearerToken == "AbCdEf123456" {
		isTokenValid = true
	}

	return &models.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		// Required: let the client know that the token is valid or not
		Status: &models.TokenReviewStatus{
			Authenticated: isTokenValid,
		},
	}, nil
}
