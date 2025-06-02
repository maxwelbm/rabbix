package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "rabbix",
	Short: "Rabbix é uma CLI para testar filas do RabbitMQ com JSON dinâmico",
	Long: `Rabbix é uma ferramenta de linha de comando para facilitar testes de filas RabbitMQ.
Você pode adicionar, listar e executar casos de teste baseados em JSON.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use um dos subcomandos. Ex: rabbix add --help")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
