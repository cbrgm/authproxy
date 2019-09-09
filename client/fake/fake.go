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
	"encoding/base64"
	"errors"
	"github.com/cbrgm/authproxy/client"
	v1 "github.com/cbrgm/authproxy/client/v1"
	"strings"
)

// fakeClient represents the authproxy fake client implementation
type fakeClient struct {
	tokens map[string]string
}

// NewForConfig returns a new client for a given config
func NewFakeClient() (client.ClientSet, error) {
	var res client.ClientSet = &fakeClient{
		tokens: map[string]string{},
	}
	return res, nil
}

// Login logs is a user or an application by username and password.
// It will return a bearer token that can be used to authorize actions performed by he client
func (c *fakeClient) Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("invalid arguments: username or password is empty")
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(username + "," + password))
	c.tokens[username] = encoded
	return encoded, nil
}

// Authenticate authenticates actions performed by the client
// It will return true, if the client in authenticated, false if not
func (c *fakeClient) Authenticate(bearerToken string) (*v1.TokenReviewRequest, error) {

	decoded, err := base64.StdEncoding.DecodeString(bearerToken)
	if err != nil {
		return nil, err
	}
	username := strings.Split(string(decoded), ",")[0]

	result := &v1.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Status: &v1.TokenReviewStatus{
			Authenticated: false,
			User: &v1.UserInfo{
				Username: username,
				UID:      "1",
				Groups:   []string{"developers"},
			},
		},
	}

	if bearerToken == "" {
		result.Status.Authenticated = false
		return result, errors.New("invalid arguments: token is missing")
	}

	for _, v := range c.tokens {
		if v == bearerToken {
			result.Status.Authenticated = true
			break
		}
	}
	return result, nil
}
