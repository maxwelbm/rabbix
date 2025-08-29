package main

import (
	"fmt"
	"os"

	"github.com/maxwelbm/rabbix/pkg/batch"
	"github.com/maxwelbm/rabbix/pkg/cache"
	"github.com/maxwelbm/rabbix/pkg/conf"
	"github.com/maxwelbm/rabbix/pkg/health"
	"github.com/maxwelbm/rabbix/pkg/list"
	"github.com/maxwelbm/rabbix/pkg/request"
	"github.com/maxwelbm/rabbix/pkg/run"
	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "rabbix",
	Short: "Rabbix é uma CLI para testar filas do RabbitMQ com JSON dinâmico",
	Long: `Rabbix é uma ferramenta de linha de comando para facilitar testes de filas RabbitMQ.
Você pode adicionar, listar e executar casos de teste baseados em JSON.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use um dos subcommands. Ex: rabbix add --help")
	},
}

func init() {
	settings := sett.New()

	cached := cache.New(settings)
	requested := request.New(settings)
	batched := batch.New(settings, cached, requested)
	r := run.New(settings, cached, requested)
	c := conf.New(settings)

	root.AddCommand(c.CmdConf())
	root.AddCommand(health.CmdHealth(settings))
	root.AddCommand(cached.CmdCache())
	root.AddCommand(batched.CmdBatch())
	root.AddCommand(list.CmdList(settings))
	root.AddCommand(r.CmdRun())
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
