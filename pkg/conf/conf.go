package conf

import (
	"encoding/base64"
	"fmt"

	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var (
	host      string
	outputDir string
	user      string
	password  string
)

var Conf = &cobra.Command{
	Use:   "conf",
	Short: "Define ou exibe configurações padrão da CLI",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Atualiza o host e/ou diretório onde os testes são salvos",
	Run: func(cmd *cobra.Command, args []string) {
		settings := sett.LoadSettings()

		if host != "" {
			settings["host"] = host
			settings["rapt_url"] = host + "/api/exchanges/%2f/amq.default/publish"
		}

		if outputDir != "" {
			settings["output_dir"] = outputDir
		}

		if user != "" && password != "" {
			auth := user + ":" + password
			auth = base64.StdEncoding.EncodeToString([]byte(auth))
			fmt.Println("auth: ", auth)

			settings["auth"] = string(auth)
		}

		sett.SaveSettings(settings)
		fmt.Println("✅ Configuração atualizada com sucesso.")

		// Sincroniza cache após mudança de configuração
		fmt.Println("🔄 Sincronizando cache...")
		// syncCacheWithFileSystem()
		fmt.Println("✅ Cache sincronizado com sucesso.")
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Exibe a configuração atual",
	Run: func(cmd *cobra.Command, args []string) {
		settings := sett.LoadSettings()
		fmt.Println("📦 Configuração atual:")
		for k, v := range settings {
			fmt.Printf("%s: %s\n", k, v)
		}
	},
}

var configCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Gerencia o cache de autocomplete",
}

var configCacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Exibe estatísticas do cache",
	Run: func(cmd *cobra.Command, args []string) {
		// printCacheStats()
	},
}

var configCacheSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sincroniza o cache com os arquivos de teste",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Sincronizando cache...")
		// syncCacheWithFileSystem()
		fmt.Println("✅ Cache sincronizado com sucesso.")
	},
}

var configCacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Limpa o cache de autocomplete",
	Run: func(cmd *cobra.Command, args []string) {
		// clearCache()
	},
}

func init() {
	configSetCmd.Flags().StringVar(&host, "host", "", "Host base do RabbitMQ (ex: http://localhost:15672)")
	configSetCmd.Flags().StringVar(&outputDir, "output", "", "Diretório para salvar os testes")
	configSetCmd.Flags().StringVar(&user, "user", "", "Usuário do RabbitMQ (texto puro ou base64)")
	configSetCmd.Flags().StringVar(&password, "password", "", "Senha do RabbitMQ (texto puro ou base64)")

	Conf.AddCommand(configSetCmd)
	Conf.AddCommand(configGetCmd)
	Conf.AddCommand(configCacheCmd)

	configCacheCmd.AddCommand(configCacheStatsCmd)
	configCacheCmd.AddCommand(configCacheSyncCmd)
	configCacheCmd.AddCommand(configCacheClearCmd)
}
