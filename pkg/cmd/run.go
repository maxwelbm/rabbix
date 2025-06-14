package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [test-name]",
	Short: "Executa um caso de teste específico",
	Long: `Executa um caso de teste específico salvamento previamente.
Exemplo: rabbix run meu-teste`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Sincroniza cache antes de fornecer sugestões
		syncCacheWithFileSystem()

		// Obtém lista de testes do cache
		cachedTests := getCachedTests()

		return cachedTests, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		testName := args[0]

		// Carrega configuração para obter diretório de saída
		settings := loadSettings()
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

		var tc TestCase
		if err := json.Unmarshal(data, &tc); err != nil {
			fmt.Printf("❌ Erro ao carregar JSON do teste '%s': %v\n", testName, err)
			return
		}

		fmt.Printf("🚀 Executando teste: %s\n", tc.Name)
		fmt.Printf("📤 Route Key: %s\n", tc.RouteKey)

		// Usa a função reutilizável PublishMessage
		resp, err := PublishMessage(tc)
		if err != nil {
			fmt.Printf("❌ Erro ao enviar mensagem: %v\n", err)
			return
		}
		defer resp.Body.Close()

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

func init() {
	rootCmd.AddCommand(runCmd)
}
