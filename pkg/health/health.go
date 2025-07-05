package health

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var Health = &cobra.Command{
	Use:   "health",
	Short: "Verifica o status de saúde da API do RabbitMQ",
	Long:  `Faz uma requisição para o endpoint /api/overview para verificar se a API do RabbitMQ está funcionando corretamente.`,
	Run: func(cmd *cobra.Command, args []string) {
		settings := sett.LoadSettings()

		var auth = "Basic Z3Vlc3Q6Z3Vlc3Q="
		if settings["auth"] != "" {
			auth = "Basic " + settings["auth"]
		}

		var host = "http://localhost:15672"
		if settings["host"] != "" {
			host = settings["host"]
		}

		// Monta a URL do endpoint de overview
		url := strings.TrimRight(host, "/") + "/api/overview"

		fmt.Printf("🔍 Verificando saúde da API...\n")
		fmt.Printf("📡 URL: %s\n", url)

		// Cria a requisição
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("❌ Erro ao criar requisição: %v\n", err)
			return
		}

		// Adiciona o header Authorization
		req.Header.Add("Authorization", auth)

		// Faz a requisição
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Erro ao fazer requisição: %v\n", err)
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

		// Exibe resultado
		fmt.Printf("📊 Status: %s\n", resp.Status)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Printf("✅ API está funcionando corretamente!\n")
		} else {
			fmt.Printf("⚠️  API retornou status de erro\n")
		}

		fmt.Printf("📄 Resposta:\n%s\n", string(body))
	},
}
