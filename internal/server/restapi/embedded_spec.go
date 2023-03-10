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
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Users API for Racoon Media Server Project",
    "title": "RMS Users API",
    "version": "1.0.0"
  },
  "host": "136.244.108.126",
  "paths": {
    "/users": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "users"
        ],
        "summary": "Получить список пользователей и информацию по ним",
        "operationId": "getUsers",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "id": {
                        "type": "string"
                      },
                      "info": {
                        "type": "string"
                      },
                      "isAdmin": {
                        "type": "boolean"
                      },
                      "lastRequestTime": {
                        "type": "integer"
                      },
                      "reqPerDay": {
                        "type": "number"
                      }
                    }
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      },
      "post": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "users"
        ],
        "summary": "Создать новый ключ пользователя",
        "operationId": "createUser",
        "parameters": [
          {
            "name": "user",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "info"
              ],
              "properties": {
                "info": {
                  "type": "string"
                },
                "isAdmin": {
                  "type": "boolean",
                  "default": false
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "id"
              ],
              "properties": {
                "id": {
                  "type": "string"
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/users/{id}": {
      "delete": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "users"
        ],
        "summary": "Удалить ключ пользователя",
        "operationId": "deleteUser",
        "parameters": [
          {
            "type": "string",
            "description": "Ключ пользователя",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "404": {
            "description": "Ключ не найден"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    }
  },
  "definitions": {
    "principal": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "x-token",
      "in": "header"
    }
  },
  "tags": [
    {
      "description": "Администрирование пользователей",
      "name": "users"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Users API for Racoon Media Server Project",
    "title": "RMS Users API",
    "version": "1.0.0"
  },
  "host": "136.244.108.126",
  "paths": {
    "/users": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "users"
        ],
        "summary": "Получить список пользователей и информацию по ним",
        "operationId": "getUsers",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/ResultsItems0"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      },
      "post": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "users"
        ],
        "summary": "Создать новый ключ пользователя",
        "operationId": "createUser",
        "parameters": [
          {
            "name": "user",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "info"
              ],
              "properties": {
                "info": {
                  "type": "string"
                },
                "isAdmin": {
                  "type": "boolean",
                  "default": false
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "id"
              ],
              "properties": {
                "id": {
                  "type": "string"
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/users/{id}": {
      "delete": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "users"
        ],
        "summary": "Удалить ключ пользователя",
        "operationId": "deleteUser",
        "parameters": [
          {
            "type": "string",
            "description": "Ключ пользователя",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "404": {
            "description": "Ключ не найден"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    }
  },
  "definitions": {
    "ResultsItems0": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "info": {
          "type": "string"
        },
        "isAdmin": {
          "type": "boolean"
        },
        "lastRequestTime": {
          "type": "integer"
        },
        "reqPerDay": {
          "type": "number"
        }
      }
    },
    "principal": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "x-token",
      "in": "header"
    }
  },
  "tags": [
    {
      "description": "Администрирование пользователей",
      "name": "users"
    }
  ]
}`))
}
