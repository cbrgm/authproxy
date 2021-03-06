// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This is the api documentation for https://github.com/cbrgm/authproxy",
    "title": "authproxy OpenAPI",
    "termsOfService": "http://github.com/cbrgm/authproxy",
    "contact": {
      "email": "chris@cbrgm.net"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0"
  },
  "host": "cbrgm.net",
  "basePath": "/v1",
  "paths": {
    "/authenticate": {
      "post": {
        "description": "authenticates users",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "auth"
        ],
        "summary": "verifies user credentials",
        "operationId": "authenticate",
        "parameters": [
          {
            "description": "TokenReviewRequest object that needs to be verified",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK (successfully authenticated)",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "401": {
            "description": "unauthorized",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "500": {
            "description": "internal server error",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          }
        }
      }
    },
    "/login": {
      "post": {
        "security": [
          {
            "basicAuth": []
          }
        ],
        "description": "login users",
        "tags": [
          "auth"
        ],
        "summary": "issues tokens for cluster access",
        "operationId": "login",
        "responses": {
          "200": {
            "description": "OK (successfully authenticated)",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "401": {
            "description": "unauthorized",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "500": {
            "description": "internal server error",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Principal": {
      "description": "Principal contains information about the user",
      "type": "object",
      "properties": {
        "Password": {
          "description": "The pasword of the user",
          "type": "string",
          "example": "bar"
        },
        "Username": {
          "description": "The name that uniquely identifies this user among all active users",
          "type": "string",
          "example": "foo"
        }
      }
    },
    "TokenReviewRequest": {
      "description": "TokenReviewRequest is issued by K8s to this service",
      "type": "object",
      "properties": {
        "apiVersion": {
          "type": "string",
          "example": "authentication.k8s.io/v1beta1"
        },
        "kind": {
          "type": "string",
          "example": "TokenReview"
        },
        "spec": {
          "$ref": "#/definitions/TokenReviewSpec"
        },
        "status": {
          "$ref": "#/definitions/TokenReviewStatus"
        }
      }
    },
    "TokenReviewSpec": {
      "description": "TokenReviewSpec contains the token being reviewed",
      "type": "object",
      "properties": {
        "token": {
          "type": "string",
          "example": "12354234123141"
        }
      }
    },
    "TokenReviewStatus": {
      "description": "TokenReviewStatus is the result of the token authentication request",
      "type": "object",
      "properties": {
        "authenticated": {
          "description": "Authenticated is true if the token is valid",
          "type": "boolean",
          "example": "true"
        },
        "user": {
          "$ref": "#/definitions/UserInfo"
        }
      }
    },
    "UserInfo": {
      "description": "UserInfo contains information about the user",
      "type": "object",
      "properties": {
        "extra": {
          "description": "Any additional information provided by the authenticator",
          "type": "object",
          "additionalProperties": true
        },
        "groups": {
          "description": "The names of groups this user is a part of",
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "uid": {
          "description": "A unique value that identifies this user across time. If this user is deleted and another user by the same name is added, they will have different UIDs",
          "type": "string",
          "example": "43"
        },
        "username": {
          "description": "The name that uniquely identifies this user among all active users",
          "type": "string",
          "example": "foo"
        }
      }
    }
  },
  "responses": {
    "UnauthorizedError": {
      "description": "Authentication information is missing or invalid",
      "headers": {
        "WWW_Authenticate": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "basicAuth": {
      "type": "basic"
    }
  },
  "tags": [
    {
      "description": "provides endpoints for authentication against ldap",
      "name": "auth"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This is the api documentation for https://github.com/cbrgm/authproxy",
    "title": "authproxy OpenAPI",
    "termsOfService": "http://github.com/cbrgm/authproxy",
    "contact": {
      "email": "chris@cbrgm.net"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0"
  },
  "host": "cbrgm.net",
  "basePath": "/v1",
  "paths": {
    "/authenticate": {
      "post": {
        "description": "authenticates users",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "auth"
        ],
        "summary": "verifies user credentials",
        "operationId": "authenticate",
        "parameters": [
          {
            "description": "TokenReviewRequest object that needs to be verified",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK (successfully authenticated)",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "401": {
            "description": "unauthorized",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "500": {
            "description": "internal server error",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          }
        }
      }
    },
    "/login": {
      "post": {
        "security": [
          {
            "basicAuth": []
          }
        ],
        "description": "login users",
        "tags": [
          "auth"
        ],
        "summary": "issues tokens for cluster access",
        "operationId": "login",
        "responses": {
          "200": {
            "description": "OK (successfully authenticated)",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "401": {
            "description": "unauthorized",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          },
          "500": {
            "description": "internal server error",
            "schema": {
              "$ref": "#/definitions/TokenReviewRequest"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Principal": {
      "description": "Principal contains information about the user",
      "type": "object",
      "properties": {
        "Password": {
          "description": "The pasword of the user",
          "type": "string",
          "example": "bar"
        },
        "Username": {
          "description": "The name that uniquely identifies this user among all active users",
          "type": "string",
          "example": "foo"
        }
      }
    },
    "TokenReviewRequest": {
      "description": "TokenReviewRequest is issued by K8s to this service",
      "type": "object",
      "properties": {
        "apiVersion": {
          "type": "string",
          "example": "authentication.k8s.io/v1beta1"
        },
        "kind": {
          "type": "string",
          "example": "TokenReview"
        },
        "spec": {
          "$ref": "#/definitions/TokenReviewSpec"
        },
        "status": {
          "$ref": "#/definitions/TokenReviewStatus"
        }
      }
    },
    "TokenReviewSpec": {
      "description": "TokenReviewSpec contains the token being reviewed",
      "type": "object",
      "properties": {
        "token": {
          "type": "string",
          "example": "12354234123141"
        }
      }
    },
    "TokenReviewStatus": {
      "description": "TokenReviewStatus is the result of the token authentication request",
      "type": "object",
      "properties": {
        "authenticated": {
          "description": "Authenticated is true if the token is valid",
          "type": "boolean",
          "example": "true"
        },
        "user": {
          "$ref": "#/definitions/UserInfo"
        }
      }
    },
    "UserInfo": {
      "description": "UserInfo contains information about the user",
      "type": "object",
      "properties": {
        "extra": {
          "description": "Any additional information provided by the authenticator",
          "type": "object",
          "additionalProperties": true
        },
        "groups": {
          "description": "The names of groups this user is a part of",
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "uid": {
          "description": "A unique value that identifies this user across time. If this user is deleted and another user by the same name is added, they will have different UIDs",
          "type": "string",
          "example": "43"
        },
        "username": {
          "description": "The name that uniquely identifies this user among all active users",
          "type": "string",
          "example": "foo"
        }
      }
    }
  },
  "responses": {
    "UnauthorizedError": {
      "description": "Authentication information is missing or invalid",
      "headers": {
        "WWW_Authenticate": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "basicAuth": {
      "type": "basic"
    }
  },
  "tags": [
    {
      "description": "provides endpoints for authentication against ldap",
      "name": "auth"
    }
  ]
}`))
}
