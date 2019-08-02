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

package main

import (
	"fmt"
	"github.com/cbrgm/authproxy/authproxy"
	"github.com/cbrgm/authproxy/provider/mock"
	"github.com/urfave/cli"
	"os"
)

const (
	FlagHTTPAddr        = "http-addr"
	FlagHTTPPrivateAddr = "http-internal-addr"
	FlagTLSCert         = "tls-cert"
	FlagTLSKey          = "tls-key"
	FlagTLSClientCA     = "tls-ca-cert"
	FlagLogJSON         = "log-json"
	FlagLogLevel        = "log-level"

	EnvHTTPAddr = "API_HTTP_ADDR"
	EnvLogJSON  = "API_LOG_JSON"
	EnvLogLevel = "API_LOG_LEVEL"
)

type apiConf struct {
	HTTPAddr        string
	HTTPPrivateAddr string
	TLSCert         string
	TLSKey          string
	TLSClientCA     string
	LogJSON         bool
	LogLevel        string
}

var (
	apiConfig = apiConf{}

	apiFlags = []cli.Flag{
		cli.StringFlag{
			Name:        FlagHTTPAddr,
			EnvVar:      EnvHTTPAddr,
			Usage:       "The address the proxy runs on",
			Value:       ":6660",
			Destination: &apiConfig.HTTPAddr,
		},
		cli.StringFlag{
			Name:        FlagHTTPPrivateAddr,
			Usage:       "The address authproxy runs a http server only for internal access",
			Value:       ":6661",
			Destination: &apiConfig.HTTPPrivateAddr,
		},
		cli.StringFlag{
			Name:        FlagTLSKey,
			Usage:       "The tls key file to be used",
			Destination: &apiConfig.TLSKey,
		},
		cli.StringFlag{
			Name:        FlagTLSCert,
			Usage:       "The tls cert file to be used",
			Destination: &apiConfig.TLSCert,
		},
		cli.StringFlag{
			Name:        FlagTLSClientCA,
			Usage:       "The tls client ca file to be used",
			Destination: &apiConfig.TLSClientCA,
		},
		cli.BoolFlag{
			Name:        FlagLogJSON,
			EnvVar:      EnvLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &apiConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        FlagLogLevel,
			EnvVar:      EnvLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &apiConfig.LogLevel,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "authproxy"
	app.Usage = "kubernetes compatible webhook authentication proxy"
	app.Action = apiAction
	app.Flags = apiFlags

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("failed to run api: %s", err)
		os.Exit(1)
	}
}

func apiAction(c *cli.Context) error {

	// create the config from command line flags
	config := authproxy.ProxyConfig{
		HTTPAddr:        apiConfig.HTTPAddr,
		HTTPPrivateAddr: apiConfig.HTTPPrivateAddr,
		TLSKey:          apiConfig.TLSKey,
		TLSCert:         apiConfig.TLSCert,
		TLSClientCA:     apiConfig.TLSClientCA,
		LogJSON:         apiConfig.LogJSON,
		LogLevel:        apiConfig.LogLevel,
	}

	// initialize the authentication provider
	fake := mock.NewMockProvider()

	// add the provider and config to the proxy
	prx := authproxy.NewWithProvider(fake, config)

	if err := prx.ListenAndServe(); err != nil {
		fmt.Printf("something went wrong: %s", err)
		os.Exit(1)
	}
	return nil
}
