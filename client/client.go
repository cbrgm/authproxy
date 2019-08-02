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

import "errors"

// ClientSet represents a new clientset for authproxy
type ClientSet interface {
	ClientV1() AuthClientV1
}

// AuthClient represents the concrete clientset implementation
type AuthClient struct {
	V1 AuthClientV1
}

// AuthClientConfig represents the clientset configuration
type AuthClientConfig struct {
	Path        string
	TLSCert     string
	TLSKey      string
	TLSClientCA string
}

// NewForConfig returns a new clientset for a given configuration
func NewForConfig(c *AuthClientConfig) (*AuthClient, error) {

	if c.TLSKey == "" {
		return nil, errors.New("invalid config: required tls key is missing")
	}
	if c.TLSCert == "" {
		return nil, errors.New("invalid config: required tls cert is missing")
	}
	if c.TLSClientCA == "" {
		return nil, errors.New("invalid config: required tls client ca is missing")
	}

	var cs AuthClient
	var err error

	cs.V1, err = newClientV1ForConfig(c)
	if err != nil {
		return nil, err
	}

	return &cs, nil
}

// ClientV1 returns the Clientsets V1 client
func (cs *AuthClient) ClientV1() AuthClientV1 {
	return cs.V1
}
