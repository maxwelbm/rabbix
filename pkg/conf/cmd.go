package conf

import (
	"github.com/spf13/cobra"
)

var (
	host      string
	outputDir string
	user      string
	password  string
)

var ConfCmd = &cobra.Command{
	Use:   "conf",
	Short: "Define ou exibe configurações padrão da CLI",
}

func init() {
	ConfCmd.AddCommand(set)
	ConfCmd.AddCommand(get)
}
