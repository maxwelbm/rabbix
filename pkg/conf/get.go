package conf

import (
	"fmt"

	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var get = &cobra.Command{
	Use:   "get",
	Short: "Exibe a configuração atual",
	Run: func(cmd *cobra.Command, args []string) {
		settings := sett.LoadSettings()
		fmt.Println("Configuração atual:")
		for k, v := range settings {
			fmt.Printf("%s: %s\n", k, v)
		}
	},
}
