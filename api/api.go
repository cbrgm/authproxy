package api

import (
	"fmt"
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

type APIMetrics struct {
	LoginAttempts metrics.Counter
}

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

func NewAuthenticationHandler(sv internal.Service) auth.AuthenticateHandlerFunc {
	return func(params auth.AuthenticateParams) restful.Responder {

		trr := params.Body

		tokenReview, err := sv.Authenticate(trr.Spec.Token)
		if err != nil {
			return auth.NewLoginInternalServerError()
		}

		return auth.NewAuthenticateOK().WithPayload(tokenReview)
	}
}
func NewLoginHandler(sv internal.Service) auth.LoginHandlerFunc {
	return func(params auth.LoginParams, user *models.Principal) restful.Responder {

		tokenReview, err := sv.Login(user.Username, user.Password)
		if err != nil {
			return auth.NewLoginInternalServerError()
		}

		// TODO: not yet implemented
		return auth.NewAuthenticateOK().WithPayload(tokenReview)
	}
}
