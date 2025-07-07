package conf

import (
	"encoding/base64"
	"fmt"

	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var set = &cobra.Command{
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
		fmt.Println("Configuração atualizada com sucesso.")
	},
}

func init() {
	set.Flags().StringVar(&host, "host", "", "Host base do RabbitMQ (ex: http://localhost:15672)")
	set.Flags().StringVar(&outputDir, "output", "", "Diretório para salvar os testes")
	set.Flags().StringVar(&user, "user", "", "Usuário do RabbitMQ (texto puro ou base64)")
	set.Flags().StringVar(&password, "password", "", "Senha do RabbitMQ (texto puro ou base64)")
}
