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
	v1 "github.com/cbrgm/authproxy/client/v1"
	httpclient "github.com/go-openapi/runtime/client"
)

// ClientSet represents a v1 authproxy client
type ClientSet interface {
	Login(username, password string) (string, error)
	Authenticate(bearerToken string) (*v1.TokenReviewRequest, error)
}

// AuthClientConfig represents the clientSet configuration
type AuthClientConfig struct {
	Path        string
	TLSCert     string
	TLSKey      string
	TLSClientCA string
}

// clientSet represents the v1 authproxy client implementation
type clientSet struct {
	client *v1.APIClient
}

// newClientV1ForConfig returns a new v1 client for a given config
func NewForConfig(c *AuthClientConfig) (ClientSet, error) {

	if c.TLSKey == "" {
		return nil, errors.New("invalid config: required tls key is missing")
	}
	if c.TLSCert == "" {
		return nil, errors.New("invalid config: required tls cert is missing")
	}
	if c.TLSClientCA == "" {
		return nil, errors.New("invalid config: required tls client ca is missing")
	}

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

	// create a new clientSet using tls and basepath
	config := v1.NewConfiguration()
	config.HTTPClient = tlsClient

	swg := v1.NewAPIClient(config)

	if c.Path == "" {
		swg.ChangeBasePath("https://localhost:6660/v1")
	} else {
		swg.ChangeBasePath(c.Path)
	}

	cl := clientSet{client: swg}

	var res ClientSet = &cl
	return res, nil
}

// Login logs is a user or an application by username and password.
// It will return a bearer token that can be used to authorize actions performed by he client
func (c *clientSet) Login(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("invalid arguments: username or password is empty")
	}

	auth := context.WithValue(context.Background(), v1.ContextBasicAuth, v1.BasicAuth{
		UserName: username,
		Password: password,
	})

	tokenReview, resp, err := c.client.AuthApi.Login(auth)
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
func (c *clientSet) Authenticate(bearerToken string) (*v1.TokenReviewRequest, error) {
	if bearerToken == "" {
		return nil, errors.New("invalid arguments: token is missing")
	}

	tokenReview, resp, err := c.client.AuthApi.Authenticate(context.TODO(), v1.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Spec: &v1.TokenReviewSpec{
			Token: bearerToken,
		},
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 || tokenReview.Status.Authenticated == false {
		return &tokenReview, nil
	}
	return &tokenReview, nil
}
