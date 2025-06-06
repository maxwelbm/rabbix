package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	host      string
	outputDir string
	user      string
	password  string
	zone      string
	client    string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Define ou exibe configurações padrão da CLI",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Atualiza o host e/ou diretório onde os testes são salvos",
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadSettings()

		if host != "" {
			settings["host"] = host
			settings["rapt_url"] = host + "/api/exchanges/%2f/amq.default/publish"
		}

		if outputDir != "" {
			settings["output_dir"] = outputDir
		}

		if zone != "" {
			settings["zone"] = zone
		}

		if client != "" {
			settings["client"] = client
		}

		var decodedUser, decodedPassword string
		
		if user != "" {
			decoded, err := base64.StdEncoding.DecodeString(user)
			if err != nil {
				fmt.Printf("Erro ao decodificar usuário base64: %v\n", err)
				return
			}
			decodedUser = string(decoded)
		}

		if password != "" {
			decoded, err := base64.StdEncoding.DecodeString(password)
			if err != nil {
				fmt.Printf("Erro ao decodificar senha base64: %v\n", err)
				return
			}
			decodedPassword = string(decoded)
		}

		// Se ambos user e password foram fornecidos, cria o auth
		if decodedUser != "" && decodedPassword != "" {
			auth := decodedUser + ":" + decodedPassword
			encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			settings["auth"] = encodedAuth
		}

		saveSettings(settings)
		fmt.Println("✅ Configuração atualizada com sucesso.")
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Exibe a configuração atual",
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadSettings()
		fmt.Println("📦 Configuração atual:")
		for k, v := range settings {
			fmt.Printf("%s: %s\n", k, v)
		}
	},
}

func init() {
	configSetCmd.Flags().StringVar(&host, "host", "", "Host base do RabbitMQ (ex: http://localhost:15672)")
	configSetCmd.Flags().StringVar(&outputDir, "output", "", "Diretório para salvar os testes")
	configSetCmd.Flags().StringVar(&user, "user", "", "Usuário do RabbitMQ (codificado em base64)")
	configSetCmd.Flags().StringVar(&password, "password", "", "Senha do RabbitMQ (codificada em base64)")
	configSetCmd.Flags().StringVar(&zone, "zone", "", "Zona para requisições")
	configSetCmd.Flags().StringVar(&client, "client", "", "Cliente para requisições")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
}

func getSettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rabbix", "settings.json")
}

func loadSettings() map[string]string {
	path := getSettingsPath()
	settings := map[string]string{}

	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	return settings
}

func saveSettings(settings map[string]string) {
	path := getSettingsPath()
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)

	data, _ := json.MarshalIndent(settings, "", "  ")
	_ = os.WriteFile(path, data, 0644)
}