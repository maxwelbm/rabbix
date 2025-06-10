package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/maxwelbm/rabbix/web"
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

type BatchExecution struct {
	ID            string    `json:"id"`
	Tests         []string  `json:"tests"`
	Concurrency   int       `json:"concurrency"`
	Delay         int       `json:"delay"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalTests    int       `json:"total_tests"`
	SuccessCount  int       `json:"success_count"`
	FailureCount  int       `json:"failure_count"`
	Results       []TestResult `json:"results"`
}

type TestResult struct {
	TestName   string    `json:"test_name"`
	Status     string    `json:"status"`
	Duration   int64     `json:"duration_ms"`
	HTTPStatus int       `json:"http_status"`
	Response   string    `json:"response"`
	Error      string    `json:"error"`
	Timestamp  time.Time `json:"timestamp"`
}

type LogMessage struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	activeExecutions = make(map[string]*BatchExecution)
	executionsMutex  = sync.RWMutex{}
	logClients       = make(map[string]chan LogMessage)
	logClientsMutex  = sync.RWMutex{}
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Sobe uma interface web avançada com os testes disponíveis",
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

		// Handler para página principal
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := web.GetTemplate("index.html")
			if err != nil {
				http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
				return
			}
			
			data := struct {
				Tests []TestCase
			}{
				Tests: tests,
			}
			
			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Erro ao executar template", http.StatusInternalServerError)
				return
			}
		})

		// Handler para arquivos estáticos
		http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/static/")
			fullPath := "static/" + path
			
			data, err := web.GetStaticFile(fullPath)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			
			// Set content type based on extension
			if strings.HasSuffix(path, ".css") {
				w.Header().Set("Content-Type", "text/css; charset=utf-8")
			} else if strings.HasSuffix(path, ".js") {
				w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			} else {
				w.Header().Set("Content-Type", "application/octet-stream")
			}
			
			w.Write(data)
		})

		// API para listar testes
		http.HandleFunc("/api/tests", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tests)
		})

		// API para executar teste individual
		http.HandleFunc("/api/run/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
				return
			}

			testName := strings.TrimPrefix(r.URL.Path, "/api/run/")
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

			start := time.Now()
			resp, err := PublishMessage(tc)
			duration := time.Since(start)

			result := TestResult{
				TestName:  tc.Name,
				Duration:  duration.Milliseconds(),
				Timestamp: time.Now(),
			}

			if err != nil {
				result.Status = "error"
				result.Error = err.Error()
			} else {
				defer resp.Body.Close()
				result.HTTPStatus = resp.StatusCode
				body, _ := io.ReadAll(resp.Body)
				result.Response = string(body)
				
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					result.Status = "success"
				} else {
					result.Status = "failure"
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		// API para executar batch
		http.HandleFunc("/api/batch", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
				return
			}

			var request struct {
				Tests       []string `json:"tests"`
				Concurrency int      `json:"concurrency"`
				Delay       int      `json:"delay"`
			}

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "JSON inválido", http.StatusBadRequest)
				return
			}

			// Cria execução
			execution := &BatchExecution{
				ID:          generateID(),
				Tests:       request.Tests,
				Concurrency: request.Concurrency,
				Delay:       request.Delay,
				Status:      "running",
				StartTime:   time.Now(),
				TotalTests:  len(request.Tests),
				Results:     make([]TestResult, 0),
			}

			executionsMutex.Lock()
			activeExecutions[execution.ID] = execution
			executionsMutex.Unlock()

			// Executa em background
			go executeBatchAsync(execution, outputDir, tests)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"execution_id": execution.ID})
		})

		// API para status da execução
		http.HandleFunc("/api/execution/", func(w http.ResponseWriter, r *http.Request) {
			executionID := strings.TrimPrefix(r.URL.Path, "/api/execution/")
			
			executionsMutex.RLock()
			execution, exists := activeExecutions[executionID]
			executionsMutex.RUnlock()

			if !exists {
				http.Error(w, "Execução não encontrada", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(execution)
		})

		// Server-Sent Events para logs em tempo real
		http.HandleFunc("/api/logs/", func(w http.ResponseWriter, r *http.Request) {
			executionID := strings.TrimPrefix(r.URL.Path, "/api/logs/")
			
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			clientChan := make(chan LogMessage, 100)
			logClientsMutex.Lock()
			logClients[executionID] = clientChan
			logClientsMutex.Unlock()

			defer func() {
				logClientsMutex.Lock()
				delete(logClients, executionID)
				logClientsMutex.Unlock()
				close(clientChan)
			}()

			for {
				select {
				case msg := <-clientChan:
					data, _ := json.Marshal(msg)
					fmt.Fprintf(w, "data: %s\n\n", data)
					w.(http.Flusher).Flush()
				case <-r.Context().Done():
					return
				}
			}
		})

		// Inicia servidor
		server := &http.Server{Addr: ":7777"}
		go func() {
			fmt.Println("🌐 Rabbix UI Avançada rodando em http://localhost:7777 (Ctrl+C para sair)")
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

func executeBatchAsync(execution *BatchExecution, outputDir string, allTests []TestCase) {
	defer func() {
		execution.EndTime = time.Now()
		execution.Status = "completed"
	}()

	// Canal para controlar concorrência
	semaphore := make(chan struct{}, execution.Concurrency)
	var wg sync.WaitGroup
	var resultsMutex sync.Mutex

	sendLog := func(level, message string) {
		logMsg := LogMessage{
			Level:     level,
			Message:   message,
			Timestamp: time.Now(),
		}
		
		logClientsMutex.RLock()
		if client, exists := logClients[execution.ID]; exists {
			select {
			case client <- logMsg:
			default:
			}
		}
		logClientsMutex.RUnlock()
	}

	sendLog("info", fmt.Sprintf("Iniciando execução em lote de %d testes", len(execution.Tests)))
	sendLog("info", fmt.Sprintf("Configurações: Concorrência=%d, Delay=%dms", execution.Concurrency, execution.Delay))

	for i, testName := range execution.Tests {
		wg.Add(1)
		go func(index int, name string) {
			defer wg.Done()
			
			// Controla concorrência
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Aplica delay
			if index > 0 && execution.Delay > 0 {
				time.Sleep(time.Duration(execution.Delay) * time.Millisecond)
			}

			sendLog("info", fmt.Sprintf("[%d/%d] Executando: %s", index+1, len(execution.Tests), name))

			// Encontra o caso de teste
			var testCase *TestCase
			for _, tc := range allTests {
				if tc.Name == name {
					testCase = &tc
					break
				}
			}

			result := TestResult{
				TestName:  name,
				Timestamp: time.Now(),
			}

			if testCase == nil {
				result.Status = "error"
				result.Error = "Teste não encontrado"
				sendLog("error", fmt.Sprintf("[%d/%d] %s: Teste não encontrado", index+1, len(execution.Tests), name))
			} else {
				start := time.Now()
				resp, err := PublishMessage(*testCase)
				result.Duration = time.Since(start).Milliseconds()

				if err != nil {
					result.Status = "error"
					result.Error = err.Error()
					sendLog("error", fmt.Sprintf("[%d/%d] %s: ERRO - %s", index+1, len(execution.Tests), name, err.Error()))
				} else {
					defer resp.Body.Close()
					result.HTTPStatus = resp.StatusCode
					body, _ := io.ReadAll(resp.Body)
					result.Response = string(body)

					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						result.Status = "success"
						sendLog("success", fmt.Sprintf("[%d/%d] %s: OK (Status: %d, %dms)", index+1, len(execution.Tests), name, resp.StatusCode, result.Duration))
					} else {
						result.Status = "failure"
						sendLog("warning", fmt.Sprintf("[%d/%d] %s: Status %d (%dms)", index+1, len(execution.Tests), name, resp.StatusCode, result.Duration))
					}
				}
			}

			// Adiciona resultado thread-safe
			resultsMutex.Lock()
			execution.Results = append(execution.Results, result)
			if result.Status == "success" {
				execution.SuccessCount++
			} else {
				execution.FailureCount++
			}
			resultsMutex.Unlock()

		}(i, testName)
	}

	wg.Wait()
	
	sendLog("info", fmt.Sprintf("Execução concluída! Sucessos: %d, Falhas: %d, Total: %d", 
		execution.SuccessCount, execution.FailureCount, execution.TotalTests))
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func init() {
	rootCmd.AddCommand(uiCmd)
}