package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

func writeJSONResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{Message: message}
	json.NewEncoder(w).Encode(response)
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		writeJSONResponse(w, "OK", http.StatusOK)
		return
	}
	http.NotFound(w, r)
}

func error403Handler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "Forbidden", http.StatusForbidden)
}

func error500Handler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "Internal Server Error", http.StatusInternalServerError)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "Redirecting to /redirect-two", http.StatusSeeOther)
}

func redirectTwoHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	message := fmt.Sprintf("Redirected from %s", referer)
	writeJSONResponse(w, message, http.StatusOK)
}

func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/timeout/")

	seconds, err := strconv.Atoi(path)
	if err != nil {
		writeJSONResponse(w, "Invalid timeout value. Must be a number.", http.StatusBadRequest)
		return
	}

	if seconds < 1 || seconds >= 300 {
		writeJSONResponse(w, "Timeout value must be between 1 and 299 seconds.", http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(seconds) * time.Second)
	message := fmt.Sprintf("Request completed after %d seconds", seconds)
	writeJSONResponse(w, message, http.StatusOK)
}

const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Best API",
    "description": "A simple HTTP API with various endpoints for testing",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:9999",
      "description": "Development server"
    }
  ],
  "paths": {
    "/": {
      "get": {
        "summary": "Health check endpoint",
        "description": "Returns OK status",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "message": "OK"
                }
              }
            }
          }
        }
      }
    },
    "/error403": {
      "get": {
        "summary": "Forbidden error endpoint",
        "description": "Returns 403 Forbidden status",
        "responses": {
          "403": {
            "description": "Forbidden error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "message": "Forbidden"
                }
              }
            }
          }
        }
      }
    },
    "/error500": {
      "get": {
        "summary": "Internal server error endpoint",
        "description": "Returns 500 Internal Server Error status",
        "responses": {
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "message": "Internal Server Error"
                }
              }
            }
          }
        }
      }
    },
    "/redirect": {
      "get": {
        "summary": "Redirect endpoint",
        "description": "Returns redirect to /redirect-two",
        "responses": {
          "303": {
            "description": "See Other redirect",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "message": "Redirecting to /redirect-two"
                }
              }
            }
          }
        }
      }
    },
    "/redirecttwo": {
      "get": {
        "summary": "Redirect target endpoint",
        "description": "Shows referrer information",
        "responses": {
          "200": {
            "description": "Successful response with referrer info",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "message": "Redirected from http://localhost:9999/redirect"
                }
              }
            }
          }
        }
      }
    },
    "/timeout/{seconds}": {
      "get": {
        "summary": "Timeout endpoint",
        "description": "Waits for specified number of seconds (1-299) before responding",
        "parameters": [
          {
            "name": "seconds",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer",
              "minimum": 1,
              "maximum": 299
            },
            "description": "Number of seconds to wait before responding"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response after timeout",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "message": "Request completed after 30 seconds"
                }
              }
            }
          },
          "400": {
            "description": "Bad request - invalid timeout value",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "examples": {
                  "invalid_number": {
                    "value": {
                      "message": "Invalid timeout value. Must be a number."
                    }
                  },
                  "out_of_range": {
                    "value": {
                      "message": "Timeout value must be between 1 and 299 seconds."
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Response": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string",
            "description": "Response message"
          }
        },
        "required": ["message"]
      }
    }
  }
}`

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/openapi.json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(openAPISpec))
		return
	}

	swaggerHTML := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/swagger/openapi.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(swaggerHTML))
}

func main() {
	http.HandleFunc("/", okHandler)
	http.HandleFunc("/error403", error403Handler)
	http.HandleFunc("/error500", error500Handler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/redirecttwo", redirectTwoHandler)
	http.HandleFunc("/timeout/", timeoutHandler)
	http.HandleFunc("/swagger/", swaggerHandler)
	http.HandleFunc("/swagger/openapi.json", swaggerHandler)

	log.Println("Server starting on :9999")
	log.Println("OpenAPI documentation available at http://localhost:9999/swagger/")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
