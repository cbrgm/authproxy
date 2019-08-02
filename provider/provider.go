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

package provider

import "github.com/cbrgm/authproxy/api/v1/models"

// Provider is an interface representing an identity provider.
// The identity provider is responsible for issuing bearer tokens for clients (login) and offers authentication of the tokens (authenticate)
type Provider interface {

	// Login issues bearer tokens for a client.
	// A client has to identity itself by his correct username and password pair.
	// In response, a TokenReviewRequest will be send containing the issued bearer token for the client and user information
	Login(username, password string) (*models.TokenReviewRequest, error)

	// Authenticate authenticates a client by a bearer token.
	// In response, a TokenReviewRequest is sent to the client.
	Authenticate(bearerToken string) (*models.TokenReviewRequest, error)
}
