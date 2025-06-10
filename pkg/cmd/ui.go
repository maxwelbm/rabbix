package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// PublishMessageConfig contém as configurações para publicar uma mensagem
type PublishMessageConfig struct {
	Host       string
	Auth       string
	ClientKey  string
	Zone       string
	Exchange   string
}

// PublishMessage envia uma mensagem para o RabbitMQ usando a API HTTP
func PublishMessage(testCase TestCase) (*http.Response, error) {
	return PublishMessageWithConfig(testCase, nil)
}

// PublishMessageWithConfig envia uma mensagem para o RabbitMQ com configurações customizadas
func PublishMessageWithConfig(testCase TestCase, config *PublishMessageConfig) (*http.Response, error) {
	settings := loadSettings()
	
	// Usa configurações passadas ou carrega das configurações salvas
	var cfg PublishMessageConfig
	if config != nil {
		cfg = *config
	}
	
	// Aplica valores padrão se não estiverem definidos
	if cfg.Host == "" {
		cfg.Host = settings["host"]
		if cfg.Host == "" {
			cfg.Host = "http://localhost:15672"
		}
	}
	
	if cfg.Auth == "" {
		cfg.Auth = settings["auth"]
		if cfg.Auth == "" {
			cfg.Auth = "Basic Z3Vlc3Q6Z3Vlc3Q="
		}
	}
	
	if cfg.ClientKey == "" {
		cfg.ClientKey = settings["client"]
		if cfg.ClientKey == "" {
			cfg.ClientKey = "6b3c9fac-46e7-43ea-ad71-0641ee51e53d"
		}
	}
	
	if cfg.Zone == "" {
		cfg.Zone = settings["zone"]
		if cfg.Zone == "" {
			cfg.Zone = "issuer"
		}
	}
	
	if cfg.Exchange == "" {
		cfg.Exchange = "amq.default"
	}
	
	// Serializa o payload
	payloadBytes, err := json.Marshal(testCase.JSONPool)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	requestBody := map[string]interface{}{
		"properties":       map[string]interface{}{},
		"routing_key":      testCase.RouteKey,
		"payload":          string(payloadBytes),
		"payload_encoding": "string",
	}

	finalBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar request body: %w", err)
	}

	// Monta a URL final
	raptURL := strings.TrimRight(cfg.Host, "/") + "/api/exchanges/%2f/" + cfg.Exchange + "/publish"

	req, err := http.NewRequest("POST", raptURL, strings.NewReader(string(finalBody)))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}
	
	// Configura headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", cfg.Auth)
	req.Header.Set("x-ds-client-key", cfg.ClientKey)
	req.Header.Set("x-ds-zone", cfg.Zone)

	clientHttp := &http.Client{}
	return clientHttp.Do(req)
}

type TestCase struct {
	Name     string                 `json:"name"`
	RouteKey string                 `json:"route_key"`
	JSONPool map[string]interface{} `json:"json_pool"`
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Sobe uma interface web com os testes disponíveis",
	Run: func(cmd *cobra.Command, args []string) {

		// Carrega configuração para obter diretório de saída
		settings := loadSettings()
		outputDir := settings["output_dir"]
		if outputDir == "" {
			home, _ := os.UserHomeDir()
			outputDir = filepath.Join(home, ".rabbix", "tests")
		}

		files, err := os.ReadDir(outputDir)
		if err != nil {
			fmt.Println("Erro ao listar os testes:", err)
			return
		}

		var tests []TestCase
		for _, file := range files {
			if filepath.Ext(file.Name()) != ".json" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(outputDir, file.Name()))
			if err != nil {
				continue
			}
			var tc TestCase
			if err := json.Unmarshal(data, &tc); err == nil {
				tests = append(tests, tc)
			}
		}

		// Handlers
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="pt-br">
<head>
	<meta charset="UTF-8">
	<title>Rabbix UI</title>
	<style>
		body { font-family: sans-serif; background: #111; color: #eee; padding: 2rem; }
		.test-case { border: 1px solid #444; margin-bottom: 1rem; padding: 1rem; border-radius: 8px; background: #1a1a1a; }
		button { padding: 0.5rem 1rem; background: #1e90ff; color: white; border: none; border-radius: 5px; cursor: pointer; }
		button:hover { background: #0078d7; }
	</style>
</head>
<body>
	<h1>📡 Rabbix - Casos de Teste</h1>
	{{range .}}
	<div class="test-case">
		<h2>{{.Name}}</h2>
		<p><strong>Queue:</strong> {{.RouteKey}}</p>
		<form method="POST" action="/run/{{.Name}}">
			<button type="submit">▶ Executar</button>
		</form>
	</div>
	{{end}}
</body>
</html>`))
			tmpl.Execute(w, tests)
		})

		http.HandleFunc("/run/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
				return
			}

			testName := strings.TrimPrefix(r.URL.Path, "/run/")
			data, err := os.ReadFile(filepath.Join(outputDir, testName+".json"))
			if err != nil {
				http.Error(w, "Teste não encontrado", http.StatusNotFound)
				return
			}

			var tc TestCase
			if err := json.Unmarshal(data, &tc); err != nil {
				http.Error(w, "Erro ao carregar JSON", http.StatusInternalServerError)
				return
			}

			resp, err := PublishMessage(tc)
			if err != nil {
				http.Error(w, "Erro ao enviar requisição: "+err.Error(), http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			w.Header().Set("Content-Type", "application/json")
			w.Write(body)
		})

		// Inicia servidor
		server := &http.Server{Addr: ":7777"}
		go func() {
			fmt.Println("🌐 Rabbix UI rodando em http://localhost:7777 (Ctrl+C para sair)")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Println("Erro no servidor:", err)
			}
		}()

		// Aguarda sinal Ctrl+C
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		fmt.Println("\n⏹ Encerrando servidor Rabbix...")

		// Desliga com timeout
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println("Erro ao encerrar:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
