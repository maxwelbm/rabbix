package batch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/maxwelbm/rabbix/pkg/cache"
	"github.com/maxwelbm/rabbix/pkg/request"
	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var (
	batchConcurrency int
	batchDelay       int
)

var BatchCmd = &cobra.Command{
	Use:   "batch [test-names...]",
	Short: "Executa múltiplos casos de teste em lote",
	Long: `Executa múltiplos casos de teste em lote com controle de concorrência.
Exemplos:
  rabbix batch teste1 teste2 teste3
  rabbix batch --concurrency 5 --delay 1000 teste1 teste2
  rabbix batch --all  # executa todos os testes disponíveis`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Sincroniza cache antes de fornecer sugestões
		cache.SyncCacheWithFileSystem()

		// Obtém lista de testes do cache
		cachedTests := cache.GetCachedTests()

		// Filtra testes que já foram especificados
		var suggestions []string
		for _, test := range cachedTests {
			alreadyUsed := false
			for _, arg := range args {
				if arg == test {
					alreadyUsed = true
					break
				}
			}
			if !alreadyUsed {
				suggestions = append(suggestions, test)
			}
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		settings := sett.LoadSettings()
		outputDir := settings["output_dir"]
		if outputDir == "" {
			home, _ := os.UserHomeDir()
			outputDir = filepath.Join(home, ".rabbix", "tests")
		}

		var testNames []string

		// Se --all foi especificado, carrega todos os testes
		if all, _ := cmd.Flags().GetBool("all"); all {
			files, err := os.ReadDir(outputDir)
			if err != nil {
				fmt.Printf("❌ Erro ao listar testes: %v\n", err)
				return
			}

			for _, file := range files {
				if filepath.Ext(file.Name()) == ".json" {
					name := file.Name()[:len(file.Name())-5] // remove .json
					testNames = append(testNames, name)
				}
			}
		} else {
			testNames = args
		}

		if len(testNames) == 0 {
			fmt.Println("❌ Nenhum teste especificado. Use 'rabbix batch --help' para ver as opções.")
			return
		}

		fmt.Printf("🚀 Executando %d teste(s) em lote\n", len(testNames))
		fmt.Printf("⚙️  Concorrência: %d | Delay: %dms\n", batchConcurrency, batchDelay)
		fmt.Println("─────────────────────────────────────")

		// Carrega todos os casos de teste
		var testCases []request.TestCase
		for _, testName := range testNames {
			testPath := filepath.Join(outputDir, testName+".json")
			data, err := os.ReadFile(testPath)
			if err != nil {
				fmt.Printf("⚠️  Pulando teste '%s': arquivo não encontrado\n", testName)
				continue
			}

			var tc request.TestCase
			if err := json.Unmarshal(data, &tc); err != nil {
				fmt.Printf("⚠️  Pulando teste '%s': erro no JSON: %v\n", testName, err)
				continue
			}
			testCases = append(testCases, tc)
		}

		if len(testCases) == 0 {
			fmt.Println("❌ Nenhum teste válido encontrado.")
			return
		}

		// Executa os testes com controle de concorrência
		results := executeBatch(testCases, batchConcurrency, time.Duration(batchDelay)*time.Millisecond)

		// Exibe resumo final
		fmt.Println("─────────────────────────────────────")
		fmt.Printf("📊 Resumo da execução:\n")

		success := 0
		failed := 0
		for _, result := range results {
			if result.Success {
				success++
			} else {
				failed++
			}
		}

		fmt.Printf("✅ Sucessos: %d\n", success)
		fmt.Printf("❌ Falhas: %d\n", failed)
		fmt.Printf("⏱️  Tempo total: %v\n", calculateTotalTime(results))

		if failed > 0 {
			fmt.Println("\n🔍 Detalhes das falhas:")
			for _, result := range results {
				if !result.Success {
					fmt.Printf("  • %s: %s\n", result.TestName, result.Error)
				}
			}
		}
	},
}

type BatchResult struct {
	TestName string
	Success  bool
	Error    string
	Duration time.Duration
	Status   int
	Response string
}

func executeBatch(testCases []request.TestCase, concurrency int, delay time.Duration) []BatchResult {
	var results []BatchResult
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Canal para controlar concorrência
	semaphore := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i, tc := range testCases {
		wg.Add(1)
		go func(index int, testCase request.TestCase) {
			defer wg.Done()

			// Adquire semáforo para controlar concorrência
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Aplica delay se não for o primeiro teste
			if index > 0 && delay > 0 {
				time.Sleep(delay)
			}

			// Executa o teste usando a função reutilizável
			testStart := time.Now()
			result := BatchResult{
				TestName: testCase.Name,
				Duration: 0,
			}

			fmt.Printf("🔄 [%d/%d] Executando: %s\n", index+1, len(testCases), testCase.Name)

			resp, err := request.PublishMessage(testCase)
			result.Duration = time.Since(testStart)

			if err != nil {
				result.Success = false
				result.Error = err.Error()
				fmt.Printf("❌ [%d/%d] %s: FALHOU (%v)\n", index+1, len(testCases), testCase.Name, err)
			} else {
				defer func() {
					err := resp.Body.Close()
					if err != nil {
						fmt.Printf("Erro ao fechar resposta: %v\n", err)
					}
				}()

				result.Status = resp.StatusCode

				body, _ := io.ReadAll(resp.Body)
				result.Response = string(body)

				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					result.Success = true
					fmt.Printf("✅ [%d/%d] %s: OK (Status: %d, %v)\n",
						index+1, len(testCases), testCase.Name, resp.StatusCode, result.Duration)
				} else {
					result.Success = false
					result.Error = fmt.Sprintf("Status HTTP %d", resp.StatusCode)
					fmt.Printf("⚠️  [%d/%d] %s: Status %d (%v)\n",
						index+1, len(testCases), testCase.Name, resp.StatusCode, result.Duration)
				}
			}

			// Thread-safe append
			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()

		}(i, tc)
	}

	wg.Wait()

	totalTime := time.Since(startTime)
	fmt.Printf("⏱️  Execução concluída em %v\n", totalTime)

	return results
}

func calculateTotalTime(results []BatchResult) time.Duration {
	var total time.Duration
	for _, result := range results {
		total += result.Duration
	}
	return total
}

func init() {

	// Flags para controlar a execução em lote
	BatchCmd.Flags().IntVarP(&batchConcurrency, "concurrency", "c", 3,
		"Número máximo de testes executados simultaneamente")
	BatchCmd.Flags().IntVarP(&batchDelay, "delay", "d", 500,
		"Delay em milissegundos entre execuções (0 = sem delay)")
	BatchCmd.Flags().BoolP("all", "a", false,
		"Executa todos os testes disponíveis")
}
