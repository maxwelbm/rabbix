package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

type TestCase struct {
	Name     string                 `json:"name"`
	RouteKey string                 `json:"route_key"`
	JSONPool map[string]interface{} `json:"json_pool"`
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Sobe uma interface web com os testes disponíveis",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		testDir := filepath.Join(homeDir, ".rabbix", "tests")

		files, err := os.ReadDir(testDir)
		if err != nil {
			fmt.Println("Erro ao listar os testes:", err)
			return
		}

		var tests []TestCase
		for _, file := range files {
			if filepath.Ext(file.Name()) != ".json" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(testDir, file.Name()))
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
			data, err := ioutil.ReadFile(filepath.Join(testDir, testName+".json"))
			if err != nil {
				http.Error(w, "Teste não encontrado", http.StatusNotFound)
				return
			}

			var tc TestCase
			if err := json.Unmarshal(data, &tc); err != nil {
				http.Error(w, "Erro ao carregar JSON", http.StatusInternalServerError)
				return
			}

			payloadBytes, _ := json.Marshal(tc.JSONPool)

			requestBody := map[string]interface{}{
				"properties":       map[string]interface{}{},
				"routing_key":      tc.RouteKey,
				"payload":          string(payloadBytes),
				"payload_encoding": "string",
			}

			finalBody, _ := json.Marshal(requestBody)

			req, err := http.NewRequest("POST", "http://localhost:15672/api/exchanges/%2f/amq.default/publish", strings.NewReader(string(finalBody)))
			if err != nil {
				http.Error(w, "Erro ao montar requisição", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Basic Z3Vlc3Q6Z3Vlc3Q=")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, "Erro ao enviar requisição", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
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
