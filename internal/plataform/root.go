package plataform

import (
	"fmt"
	"os"

	"github.com/maxwelbm/rabbix/pkg/batch"
	"github.com/maxwelbm/rabbix/pkg/cache"
	"github.com/maxwelbm/rabbix/pkg/conf"
	"github.com/maxwelbm/rabbix/pkg/health"
	"github.com/maxwelbm/rabbix/pkg/list"
	"github.com/maxwelbm/rabbix/pkg/run"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "rabbix",
	Short: "Rabbix é uma CLI para testar filas do RabbitMQ com JSON dinâmico",
	Long: `Rabbix é uma ferramenta de linha de comando para facilitar testes de filas RabbitMQ.
Você pode adicionar, listar e executar casos de teste baseados em JSON.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use um dos subcomandos. Ex: rabbix add --help")
	},
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(conf.ConfCmd)
	root.AddCommand(health.HealthCmd)
	root.AddCommand(cache.CacheCmd)
	root.AddCommand(batch.BatchCmd)
	root.AddCommand(list.ListCmd)
	root.AddCommand(run.RunCmd)
}
