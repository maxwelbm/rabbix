package run

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/maxwelbm/rabbix/pkg/cache"
	"github.com/maxwelbm/rabbix/pkg/rabbix"
	"github.com/maxwelbm/rabbix/pkg/request"
	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

type Run struct {
	settings sett.SettItf
	Cache    cache.CacheItf
	request  request.RequestItf
}

func New(
	settings sett.SettItf,
	cache cache.CacheItf,
	request request.RequestItf,
) *Run {
	return &Run{
		settings: settings,
		Cache:    cache,
		request:  request,
	}
}

func (r *Run) CmdRun() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run [test-name]",
		Short: "Executa um caso de teste específico",
		Long: `Executa um caso de teste específico salvamento previamente.
Exemplo: rabbix run meu-teste`,
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Sincroniza cache antes de fornecer sugestões
			r.Cache.SyncCacheWithFileSystem()

			// Obtém lista de testes do cache
			cachedTests := r.Cache.GetCachedTests()

			return cachedTests, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			testName := args[0]

			// Carrega configuração para obter diretório de saída
			settings := r.settings.LoadSettings()
			outputDir := settings["output_dir"]
			if outputDir == "" {
				home, _ := os.UserHomeDir()
				outputDir = filepath.Join(home, ".rabbix", "tests")
			}

			// Lê o arquivo do teste
			testPath := filepath.Join(outputDir, testName+".json")
			data, err := os.ReadFile(testPath)
			if err != nil {
				fmt.Printf("❌ Erro: Teste '%s' não encontrado em %s\n", testName, testPath)
				fmt.Println("💡 Use 'rabbix list' para ver os testes disponíveis")
				return
			}

			var tc rabbix.TestCase
			if err := json.Unmarshal(data, &tc); err != nil {
				fmt.Printf("❌ Erro ao carregar JSON do teste '%s': %v\n", testName, err)
				return
			}

			fmt.Printf("🚀 Executando teste: %s\n", tc.Name)
			fmt.Printf("📤 Route Key: %s\n", tc.RouteKey)

			// Usa a função reutilizável PublishMessage
			resp, err := r.request.Request(tc)
			if err != nil {
				fmt.Printf("❌ Erro ao enviar mensagem: %v\n", err)
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Printf("❌ Erro ao fechar corpo da resposta: %v\n", err)
				}
			}()

			// Lê a resposta
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("❌ Erro ao ler resposta: %v\n", err)
				return
			}

			// Exibe o resultado
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Printf("✅ Mensagem enviada com sucesso! (Status: %d)\n", resp.StatusCode)
			} else {
				fmt.Printf("⚠️  Resposta com status %d\n", resp.StatusCode)
			}

			fmt.Printf("📥 Resposta do RabbitMQ:\n%s\n", string(body))
		},
	}

	return cmd
}
