package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	inputFile string
	routeKey  string
	testName  string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adiciona um novo caso de teste",
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" || routeKey == "" || testName == "" {
			fmt.Println("Você precisa informar --file, --routeKey e --name")
			return
		}

		// Lê o arquivo JSON de entrada
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Erro ao ler o arquivo: %v\n", err)
			return
		}

		// Verifica se é JSON válido
		var temp interface{}
		if err := json.Unmarshal(data, &temp); err != nil {
			fmt.Printf("JSON inválido: %v\n", err)
			return
		}

		// Carrega configuração para obter diretório de saída
		settings := loadSettings()
		outputDir := settings["output_dir"]
		if outputDir == "" {
			home, _ := os.UserHomeDir()
			outputDir = filepath.Join(home, ".rabbix", "tests")
		}

		// Cria estrutura e salva o arquivo
		os.MkdirAll(outputDir, os.ModePerm)

		// Cria estrutura de caso de teste
		testCase := map[string]interface{}{
			"name":      testName,
			"route_key": routeKey,
			"json_pool": temp,
		}

		// Salva no arquivo
		outPath := filepath.Join(outputDir, testName+".json")
		outData, _ := json.MarshalIndent(testCase, "", "  ")
		if err := os.WriteFile(outPath, outData, 0644); err != nil {
			fmt.Printf("Erro ao salvar caso de teste: %v\n", err)
			return
		}

		fmt.Printf("✅ Caso de teste \"%s\" salvo com sucesso em %s\n", testName, outPath)

		// Adiciona ao cache para autocomplete
		addToCache(testName, routeKey)
	},
}

func init() {
	addCmd.Flags().StringVar(&inputFile, "file", "", "Caminho para o arquivo JSON de entrada")
	addCmd.Flags().StringVar(&routeKey, "routeKey", "", "Routing key do RabbitMQ")
	addCmd.Flags().StringVar(&testName, "name", "", "Nome do caso de teste")
	rootCmd.AddCommand(addCmd)
}
