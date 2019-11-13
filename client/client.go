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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	v1 "github.com/cbrgm/authproxy/client/v1"
	"io/ioutil"
	"net/http"
)

// ClientSet represents a v1 authproxy client
type ClientSet interface {
	Login(username, password string) (string, error)
	Authenticate(bearerToken string) (*v1.TokenReviewRequest, error)
}

// AuthClientConfig represents the clientSet configuration
type AuthClientConfig struct {
	Path string
	CA   string
}

// clientSet represents the v1 authproxy client implementation
type clientSet struct {
	client *v1.APIClient
}

// newClientV1ForConfig returns a new v1 client for a given config
func NewForConfig(c *AuthClientConfig) (ClientSet, error) {

	if c.CA == "" {
		return nil, errors.New("invalid config: required authproxy ca is missing")
	}

	// try to load the ca file
	caPool, err := LoadCAFile(c.CA)
	if err != nil {
		return nil, err
	}

	// Trust the augmented cert pool in our client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            caPool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	// create a new clientSet using tls and basepath
	config := v1.NewConfiguration()
	config.HTTPClient = client

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

// LoadCAFile loads a single PEM-encoded file from the path specified.
func LoadCAFile(caFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading CA File: %s", err)
	}

	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, fmt.Errorf("Error loading CA File: Couldn't parse PEM in: %s", caFile)
	}

	return pool, nil
}
