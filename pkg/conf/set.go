package conf

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Conf) CmdSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Atualiza o host e/ou diretório onde os testes são salvos",
		Run: func(cmd *cobra.Command, args []string) {
			settings := c.settings.LoadSettings()

			if host != "" {
				settings["host"] = host
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

			c.settings.SaveSettings(settings)
			fmt.Println("✅ Configuração atualizada com sucesso.")
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Host base do RabbitMQ (ex: http://localhost:15672)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Diretório para salvar os testes")
	cmd.Flags().StringVar(&user, "user", "", "Usuário do RabbitMQ (texto puro ou base64)")
	cmd.Flags().StringVar(&password, "password", "", "Senha do RabbitMQ (texto puro ou base64)")

	return cmd
}
