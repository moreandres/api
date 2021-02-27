// Copyright (c) 2021 Andres More

// api_test

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func flow(client *httpexpect.Expect) {
	example := map[string]interface{}{
		"first_name": "George",
		"last_name":  "Washington",
		"birthday":   "22-02-1732",
		"address": map[string]interface{}{
			"street_address": "3200 Mount Vernon Memorial Highway",
			"city":           "Mount Vernon",
			"state":          "Virginia",
			"country":        "United States",
		},
	}

	object := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "objects",
			"attributes": example,
		},
	}

	health := client.GET("/v1/health").Expect().Status(http.StatusOK).JSON().Object().Value("data").Object()
	health.ValueEqual("type", "health").Value("attributes").Object()

	client.GET("/debug/pprof").Expect().Status(http.StatusOK)

	client.GET("/v1/objects").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Empty()
	client.GET("/v1/objects/1").Expect().Status(http.StatusNotFound).JSON().Object().Value("errors").Array().Element(0).Object().ValueEqual("status", http.StatusNotFound)
	client.POST("/v1/objects").WithJSON(object).Expect().Status(http.StatusCreated).JSON().Object().Value("data").Object().ValueEqual("type", "objects").Value("attributes").Object().NotEmpty()
	client.GET("/v1/objects/1").Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().ValueEqual("type", "objects").ValueEqual("ID", 1)
	client.POST("/v1/objects").WithJSON(object).Expect().Status(http.StatusCreated).JSON().Object().Value("data").Object().ValueEqual("type", "objects").Value("attributes").Object().NotEmpty()
	client.GET("/v1/objects/2").Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().ValueEqual("type", "objects").ValueEqual("ID", 2)
	client.POST("/v1/objects").WithJSON(object).Expect().Status(http.StatusCreated).JSON().Object().Value("data").Object().ValueEqual("type", "objects").Value("attributes").Object().NotEmpty()
	client.GET("/v1/objects/3").Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().ValueEqual("type", "objects").ValueEqual("ID", 3)
	client.GET("/v1/objects").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(3)
	client.GET("/v1/objects").WithQueryString("limit=1024").Expect().Status(http.StatusOK)
	client.GET("/v1/objects").WithQueryString("limit=invalid").Expect().Status(http.StatusBadRequest)
	client.GET("/v1/objects").WithQueryString("sort=-id").Expect().Status(http.StatusOK)
	client.GET("/v1/objects").WithQueryString("sort=id").Expect().Status(http.StatusOK)
	client.GET("/v1/objects").WithQueryString("sort=invalid").Expect().Status(http.StatusInternalServerError)
	client.PATCH("/v1/objects/4").WithJSON(object).Expect().Status(http.StatusNotFound)
	client.PATCH("/v1/objects/3").WithJSON(object).Expect().Status(http.StatusOK)
	client.DELETE("/v1/objects/3").Expect().Status(http.StatusNoContent)
	client.GET("/v1/objects").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(2)
	client.GET("/v1/objects/3").Expect().Status(http.StatusNotFound)
	client.DELETE("/v1/objects/3").Expect().Status(http.StatusNotFound)
	client.DELETE("/v1/objects/2").Expect().Status(http.StatusNoContent)
	client.GET("/v1/objects").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	client.DELETE("/v1/objects/1").Expect().Status(http.StatusNoContent)
	client.GET("/v1/objects").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Empty()

}

func TestFlow(t *testing.T) {

	router := setupRouter()

	server := httptest.NewServer(router)
	defer server.Close()

	client := httpexpect.New(t, server.URL)

	flow(client)
}

func BenchmarkFlow(b *testing.B) { // TODO: not passing

	for i := 0; i < b.N; i++ {
		router := setupRouter()

		server := httptest.NewServer(router)
		client := httpexpect.New(b, server.URL)

		flow(client)
		server.Close()
	}

}
