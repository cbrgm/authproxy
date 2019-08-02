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

package client

import (
	"context"
	"errors"
	swagger "github.com/cbrgm/authproxy/client/v1"
	httpclient "github.com/go-openapi/runtime/client"
)

// AuthClientV1 represents a v1 authproxy client
type AuthClientV1 interface {
	Login(username, password string) (string, error)
	Authenticate(bearerToken string) (bool, error)
}

// authClientV1 represents the v1 authproxy client implementation
type authClientV1 struct {
	v1 swagger.APIClient
}

// newClientV1ForConfig returns a new v1 client for a given config
func newClientV1ForConfig(c *AuthClientConfig) (AuthClientV1, error) {

	// build a tls client from config
	tlsConfig := httpclient.TLSClientOptions{
		CA:          c.TLSClientCA,
		Certificate: c.TLSCert,
		Key:         c.TLSKey,
	}

	tlsClient, err := httpclient.TLSClient(tlsConfig)
	if err != nil {
		return nil, err
	}

	// create a new clientset using tls and basepath
	config := swagger.NewConfiguration()
	config.HTTPClient = tlsClient

	clientset := swagger.NewAPIClient(config)

	if c.Path == "" {
		clientset.ChangeBasePath("https://localhost:6660/v1")
	} else {
		clientset.ChangeBasePath(c.Path)
	}

	client := &authClientV1{
		v1: *clientset,
	}

	var v1 AuthClientV1 = client
	return v1, nil
}

// Login logs is a user or an application by username and password.
// It will return a bearer token that can be used to authorize actions performed by he client
func (c *authClientV1) Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("invalid arguments: username or password is empty")
	}

	auth := context.WithValue(context.Background(), swagger.ContextBasicAuth, swagger.BasicAuth{
		UserName: username,
		Password: password,
	})

	tokenReview, resp, err := c.v1.AuthApi.Login(auth)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 500 {
		return "", errors.New("internal server error: authentication process failed on remote server")
	}

	if resp.StatusCode == 401 || tokenReview.Status.Authenticated == false {
		return "", errors.New("unauthorized: invalid authentication credentials")
	}

	return tokenReview.Spec.Token, nil
}

// Authenticate authenticates actions performed by the client
// It will return true, if the client in authenticated, false if not
func (c *authClientV1) Authenticate(bearerToken string) (bool, error) {
	if bearerToken == "" {
		return false, errors.New("invalid arguments: token is missing")
	}

	tokenReview, resp, err := c.v1.AuthApi.Authenticate(context.TODO(), swagger.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Spec: &swagger.TokenReviewSpec{
			Token: bearerToken,
		},
	})

	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 || tokenReview.Status.Authenticated == false {
		return false, nil
	}
	return true, nil
}