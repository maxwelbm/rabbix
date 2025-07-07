package conf

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Conf) CmdGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Exibe a configuração atual",
		Run: func(cmd *cobra.Command, args []string) {
			settings := c.settings.LoadSettings()
			fmt.Println("Configuração atual:")
			for k, v := range settings {
				fmt.Printf("%s: %s\n", k, v)
			}
		},
	}
}
