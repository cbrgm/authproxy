# authproxy

**A Kubernetes compatible webhook authentication proxy framework and clientset.**

[![](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://github.com/cbrgm/authproxy/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/cbrgm/authproxy)](https://goreportcard.com/report/github.com/cbrgm/authproxy)

## Features

* Full-featured, Kubernetes compatible authentication middleware framework
    * Highly configurable due to plugable identity provider implementations
    * Integrates with Kubernetes apiserver webhook mode
    * Supports mutual TLS, Prometheus Metrics, Log Levels, ...
* Client for easy implementation of authentication mechanisms for applications
* Lightweight, extensible, and built with developer UX in mind

## Overview

authproxy is a framework that allows to build middleware in order to authenticate users and applications against different identity providers.
It is built for use with Kubernetes in mind and supports its token and auhorizations webhook modes but it can also be used standalone.


![overview](./docs/images/authproxy.png "Overview")  



authproxy offers two endpoints for issuing bearer tokens and validating them. 

* POST `/login` - issues bearer token
* POST `/authenticate`- validates bearer tokens

The concrete behavior of the endpoints is determined by an implementation for a specific backend, a so-called identity provider. 

An identity provider only has to implement the [provider interface](https://github.com/cbrgm/authproxy/blob/master/provider/provider.go) specification and authproxy takes care of handling requests, marshalling data types, managing encryption and allowing you to focus on your provider implementation.

## Getting started

The following explains how to implement your own identity provider using authproxy.

### Building a custom identity provider implementation

Your provider should implement the following interface:

***Provider Interface***:
```go
type Provider interface {
	Login(username, password string) (*models.TokenReviewRequest, error)
	Authenticate(bearerToken string) (*models.TokenReviewRequest, error)
}
```

Here is an example of a fake provider, that creates a BearerToken for the user `foo` with password `bar` and can validate it.
This serves only as inspiration. Of course you can easily implement other providers like database queries or third party services.

***Example Provider***:
```go
type Mock struct {
	Name string
}

func NewMockProvider() *Mock {
	return &Mock{
		Name: "mock-authenticator",
	}
}

func (provider *Mock) Login(username, password string) (*models.TokenReviewRequest, error) {
	var isAuthenticated = false
	var bearerToken = ""

	if username == "foo" && password == "bar" {
		isAuthenticated = true

		// generate bearerToken to be used for authentication
		bearerToken = "AbCdEf123456"
	}

	return &models.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Status: &models.TokenReviewStatus{
			// Required: let the client know if the user has successfully authenticated
			Authenticated: isAuthenticated,

			// optional: add user information for the client
			User: &models.UserInfo{
				Username: username,
				UID:      "1",
				Groups:   []string{"developers"},
			},
		},
		// Required: return the token for the client
		Spec: &models.TokenReviewSpec{
			Token: bearerToken,
		},
	}, nil
}

func (provider *Mock) Authenticate(bearerToken string) (*models.TokenReviewRequest, error) {
	var isTokenValid = false

	if bearerToken == "AbCdEf123456" {
		isTokenValid = true
	}

	return &models.TokenReviewRequest{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		// Required: let the client know that the token is valid or not
		Status: &models.TokenReviewStatus{
			Authenticated: isTokenValid,
		},
	}, nil
}
```

Start the authproxy with the fake provider:
```go
// create the config from command line flags
config := authproxy.ProxyConfig{
	HTTPAddr:         ":6660",
	HTTPPrivateAddr:  ":6661",
	TLSKey:           "./server.key",
	TLSCert:          "./server.crt",
	TLSClientCA:      "./ca.crt",
	LogJSON:          false,
	LogLevel:         "info",
}

// initialize the identity provider
fake := mock.NewMockProvider()

// add the provider and config to the proxy
prx := authproxy.NewWithProvider(fake, config)

if err := prx.ListenAndServe(); err != nil {
    fmt.Printf("something went wrong: %s", err)
    os.Exit(1)
}
return nil
```

For the complete example please see here.

### Client usage

authproxy provides a client to communicate with the API. It can be used to build authentication mechanisms into apps.
The use is kept very simple:

```go
username, password := "foo", "bar"

cfg := client.AuthClientConfig{
	TLSKey:      "./client.key",
	TLSCert:     "./client.crt",
	TLSClientCA: "./ca.crt",
}

cl, err := client.NewForConfig(&cfg)
if err != nil {
	return err
}

// receive a bearer token
token, err := cl.V1.Login(username, password)
if err != nil {
	return err
}

// authenticate the bearer token
ok, err := cl.V1.Authenticate(token)
if err != nil {
	return err
}

if !ok {
	fmt.Println("client unauthenticated, token is invalid")
	return nil
}

fmt.Println("client successfully authenticated, token is valid")
return nil
```

## Kubernetes and authproxy

When a client attempts to authenticate with the API server using a bearer token, the apiservers authentication webhook POSTs a JSON-serialized authentication.k8s.io/v1beta1 TokenReview object containing the token to authproxy.

In order to use authproxy as a service for Kubernetes Webhook Token Authentication, you must configure the api server according to the documentation [here](https://kubernetes.io/docs/reference/access-authn-authz/authentication/).

## Projects using authproxy

updated soon...


## Credit & License

authproxy is open-source and is developed under the terms of the [Apache 2.0 License](https://github.com/cbrgm/authproxy/blob/master/LICENSE).

Maintainer of this repository is:

-   [@cbrgm](https://github.com/cbrgm) | Christian Bargmann <mailto:chris@cbrgm.net>

Please refer to the git commit log for a complete list of contributors.

## Contributing

See the [Contributing Guide](https://github.com/cbrgm/authproxy/blob/master/CONTRIBUTING.md).