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

package api

import (
	"fmt"
	"github.com/cbrgm/authproxy/api/errors"
	"github.com/cbrgm/authproxy/api/v1/models"
	"github.com/cbrgm/authproxy/api/v1/restapi"
	"github.com/cbrgm/authproxy/api/v1/restapi/operations"
	"github.com/cbrgm/authproxy/api/v1/restapi/operations/auth"
	"github.com/cbrgm/authproxy/internal"
	"github.com/cbrgm/authproxy/provider"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-openapi/loads"
	restful "github.com/go-openapi/runtime/middleware"
	prom "github.com/prometheus/client_golang/prometheus"
	"net/http"
)

// NewV1 returns a new configured authproxy v1 multiplexer to be used by a router
func NewV1(prv *provider.Provider, logger log.Logger) (*chi.Mux, error) {
	router := chi.NewRouter()

	// load the metrics
	apiMetrics := apiMetrics()

	// load the swagger spec
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded swagger file: %s", err.Error())
	}

	api := operations.NewAuthproxyAPI(swaggerSpec)
	api.Middleware = func(b restful.Builder) http.Handler {
		return restful.Spec("", swaggerSpec.Raw(), api.Context().RoutesHandler(b))
	}

	api.BasicAuthAuth = func(username string, password string) (user *models.Principal, e error) {
		return &models.Principal{
			Username: username,
			Password: password,
		}, nil
	}

	// initialize services
	var sv internal.Service
	sv = internal.NewService(prv)
	sv = internal.NewLoggingService(log.WithPrefix(logger, "service", "provider"), sv)
	sv = internal.NewMetricsService(apiMetrics.LoginAttempts, sv)

	// initialize handlers

	api.AuthAuthenticateHandler = NewAuthenticationHandler(sv)
	api.AuthLoginHandler = NewLoginHandler(sv)

	router.Mount("/", api.Serve(nil))

	return router, nil

}

// APIMetrics represents all authproxy metrics
type APIMetrics struct {
	LoginAttempts metrics.Counter
}

// apiMetrics returns new metrics for metrics endpoint
func apiMetrics() *APIMetrics {
	namespace := "authproxy"

	return &APIMetrics{
		LoginAttempts: prometheus.NewCounterFrom(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "authentication",
			Name:      "login_attempts_total",
			Help:      "Number of login attempts that succeeded and failed",
		}, []string{"status"}),
	}
}

// NewAuthenticationHandler returns a new handler for /authenticate endpoint
func NewAuthenticationHandler(sv internal.Service) auth.AuthenticateHandlerFunc {
	return func(params auth.AuthenticateParams) restful.Responder {
		request := params.Body
		tokenReview, err := sv.Authenticate(request.Spec.Token)

		if errors.IsUnauthorized(err) {
			tokenReview = defaultResponse()
			return auth.NewAuthenticateUnauthorized().WithPayload(tokenReview)
		}

		if errors.IsInternalError(err) || err != nil {
			tokenReview = defaultResponse()
			return auth.NewLoginInternalServerError().WithPayload(tokenReview)
		}

		return auth.NewAuthenticateOK().WithPayload(tokenReview)
	}
}

// NewLoginHandler returns a new handler for /login endpoint
func NewLoginHandler(sv internal.Service) auth.LoginHandlerFunc {
	return func(params auth.LoginParams, user *models.Principal) restful.Responder {
		tokenReview, err := sv.Login(user.Username, user.Password)

		if errors.IsUnauthorized(err) {
			tokenReview = defaultResponse()
			return auth.NewLoginUnauthorized().WithPayload(tokenReview)
		}

		if errors.IsInternalError(err) || err != nil {
			tokenReview = defaultResponse()
			return auth.NewLoginInternalServerError().WithPayload(tokenReview)
		}

		return auth.NewAuthenticateOK().WithPayload(tokenReview)
	}
}

func defaultResponse() *models.TokenReviewRequest {
	return &models.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Status: &models.TokenReviewStatus{
			Authenticated: false,
		},
	}
}
