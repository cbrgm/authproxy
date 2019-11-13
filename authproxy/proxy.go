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

package authproxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/cbrgm/authproxy/api"
	"github.com/cbrgm/authproxy/provider"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	prom "github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// ProxyConfig represents a the proxy configuration parameters
type ProxyConfig struct {
	HTTPAddr        string
	HTTPPrivateAddr string
	TLSCert         string
	TLSKey          string
	TLSClientCA     string
	LogJSON         bool
	LogLevel        string
}

// Proxy represents the authproxy instance
type Proxy struct {
	Provider provider.Provider
	Config   ProxyConfig
}

// NewConfiguration returns a new default configuration
func NewConfiguration() ProxyConfig {
	return ProxyConfig{
		HTTPAddr:        ":6660",
		HTTPPrivateAddr: ":6661",
		TLSCert:         "server.crt",
		TLSKey:          "server.key",
		TLSClientCA:     "ca.crt",
		LogJSON:         false,
		LogLevel:        "info",
	}
}

// NewWithProvider returns a new proxy instance using a provider implementation as backend
func NewWithProvider(provider provider.Provider, cfg ProxyConfig) *Proxy {
	return &Proxy{
		Provider: provider,
		Config:   cfg,
	}
}

// ListenAndServe starts the proxy
func (p *Proxy) ListenAndServe() error {

	// validate config
	if p.Config.TLSKey == "" {
		return errors.New("invalid config: no private key specified for HTTPS")
	}
	if p.Config.TLSCert == "" {
		return errors.New("invalid config: no cert specified for HTTPS")
	}
	if p.Config.TLSClientCA == "" {
		return errors.New("invalid config: no client ca cert for HTTPS")
	}
	if p.Provider == nil {
		return errors.New("invalid config: no provider registered")
	}

	// initialize logger
	logger := newLogger(p.Config.LogJSON, p.Config.LogLevel)
	logger = log.WithPrefix(logger, "app", "authproxy")

	var gr run.Group
	{
		apiV1, err := api.NewV1(&p.Provider, log.WithPrefix(logger, "component", "api"))
		if err != nil {
			return err
		}

		router := chi.NewRouter()
		router.Use(requestLogger(logger))
		router.Mount("/", apiV1)

		//parse certificates from cert and key file for the authproxy server
		cert, err := tls.LoadX509KeyPair(p.Config.TLSCert, p.Config.TLSKey)
		if err != nil {
			return fmt.Errorf("invalid config: error parsing tls certificate file: %v", err)
		}

		tlsConfig := tls.Config{
			Certificates:             []tls.Certificate{cert},
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		}
		// parse certificates from certificate authority file to a new CertPool.
		cPool := x509.NewCertPool()
		clientCert, err := ioutil.ReadFile(p.Config.TLSClientCA)
		if err != nil {
			return fmt.Errorf("invalid config: error reading CA file: %v", err)
		}
		if cPool.AppendCertsFromPEM(clientCert) != true {
			return errors.New("invalid config: failed to parse client CA")
		}

		tlsConfig.ClientCAs = cPool

		server := http.Server{
			Addr:      p.Config.HTTPAddr,
			Handler:   router,
			TLSConfig: &tlsConfig,
		}

		// private router initialization

		privateRouter := chi.NewRouter()
		privateRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, http.StatusText(http.StatusOK))
		})

		privateRouter.Mount("/metrics", prom.UninstrumentedHandler())

		privateServer := &http.Server{
			Addr:    p.Config.HTTPPrivateAddr,
			Handler: privateRouter,
		}

		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "running api",
				"addr", server.Addr,
			)
			return server.ListenAndServeTLS(p.Config.TLSCert, p.Config.TLSKey)
		}, func(err error) {
			_ = server.Shutdown(context.TODO())
		})

		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "running internal api",
				"addr", privateServer.Addr,
			)
			return privateServer.ListenAndServe()
		}, func(err error) {
			_ = privateServer.Shutdown(context.TODO())
		})
	}

	if err := gr.Run(); err != nil {
		return fmt.Errorf("error running: %s", err)
	}

	return nil
}

// requestLogger proxies incoming requests and logs them
func requestLogger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			level.Debug(logger).Log(
				"proto", r.Proto,
				"method", r.Method,
				"status", ww.Status(),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}

// newLogger returns a new logger
func newLogger(json bool, loglevel string) log.Logger {
	var logger log.Logger

	if json {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}

	switch strings.ToLower(loglevel) {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
}
