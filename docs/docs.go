// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/papaya/predict": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "papaya"
                ],
                "summary": "get papaya ripeness prediction",
				"parameters": [
                    {
                        "in": "formData",
                        "name": "image",
                        "type": "file",
                        "description": "Image file of papaya to predict ripeness"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "the prediction confidence and classification",
                        "schema": {
                            "$ref": "#/definitions/handler.PredictionResponseDto"
                        }
                    },
                    "404": {
                        "description": "Multipart-form data error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorMessage"
                        }
                    },
                    "413": {
                        "description": "Payload too large (10MB limit)",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.PredictionResponseDto": {
            "type": "object",
            "properties": {
                "classification": {
                    "type": "string"
                },
                "confidence": {
                    "type": "string"
                }
            }
        },
        "main.ErrorMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "papaya.cscms.me",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Papaya Ripeness Prediction API",
	Description: "This is a sample of papaya ripeness prediction api for CSC340",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
