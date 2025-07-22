package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/maxwelbm/rabbix/pkg/rabbix"
	"github.com/maxwelbm/rabbix/pkg/sett"
)

type Request struct {
	settings sett.SettItf
}

type RequestItf interface {
	Request(testCase rabbix.TestCase) (*http.Response, error)
}

func New(settings sett.SettItf) RequestItf {
	return &Request{
		settings: settings,
	}
}

// PublishMessage envia uma mensagem para o RabbitMQ usando a API HTTP
func (r *Request) Request(testCase rabbix.TestCase) (*http.Response, error) {
	settings := r.settings.LoadSettings()

	var auth = settings["auth"]
	if auth == "" {
		fmt.Printf("necessario configurar user e password com o comando 'rabbix conf set --user <user> --password <password>'\n")
		return nil, fmt.Errorf("autenticação não configurada")
	}

	auth = "Basic " + auth

	var host = "http://localhost:15672" // host default
	if settings["host"] != "" {
		host = settings["host"]
	}

	payloadBytes, err := json.Marshal(testCase.JSONPool)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	properties := map[string]any{}
	if len(testCase.Headers) > 0 {
		properties["headers"] = testCase.Headers
	}

	requestBody := map[string]any{
		"properties":       properties,
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

	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	req.Header.Set("Content-Type", "application/json")

	clientHttp := &http.Client{}
	return clientHttp.Do(req)
}
