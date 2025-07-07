package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/maxwelbm/rabbix/pkg/sett"
)

// PublishMessage envia uma mensagem para o RabbitMQ usando a API HTTP
func PublishMessage(testCase TestCase) (*http.Response, error) {
	settings := sett.LoadSettings()

	var auth = "Basic Z3Vlc3Q6Z3Vlc3Q="
	if settings["auth"] != "" {
		auth = "Basic " + settings["auth"]
	}

	var host = "http://localhost:15672"
	if settings["host"] != "" {
		host = settings["host"]
	}

	clientKey := "6b3c9fac-46e7-43ea-ad71-0641ee51e53d"
	if settings["client"] != "" {
		clientKey = settings["client"]
	}

	zone := "issuer"
	if settings["zone"] != "" {
		zone = settings["zone"]
	}

	payloadBytes, err := json.Marshal(testCase.JSONPool)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	requestBody := map[string]any{
		"properties":       map[string]any{},
		"routing_key":      testCase.RouteKey,
		"payload":          string(payloadBytes),
		"payload_encoding": "string",
	}

	finalBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar request body: %w", err)
	}

	raptURL := strings.TrimRight(host, "/") + "/api/exchanges/%2f/amq.default/publish"

	req, err := http.NewRequest("POST", raptURL, bytes.NewBuffer(finalBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Configura headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	req.Header.Set("x-ds-client-key", clientKey)
	req.Header.Set("x-ds-zone", zone)

	clientHttp := &http.Client{}
	return clientHttp.Do(req)
}

type TestCase struct {
	Name     string         `json:"name"`
	RouteKey string         `json:"route_key"`
	JSONPool map[string]any `json:"json_pool"`
}
