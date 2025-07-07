package list

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

func CmdList(settings sett.SettItf) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lista todos os casos de teste salvos",
		Run: func(_ *cobra.Command, args []string) {
			settings := settings.LoadSettings()
			outputDir := settings["output_dir"]
			if outputDir == "" {
				home, _ := os.UserHomeDir()
				outputDir = filepath.Join(home, ".rabbix", "tests")
			}

			files, err := os.ReadDir(outputDir)
			if err != nil {
				fmt.Printf("Erro ao acessar diretório: %v\n", err)
				return
			}

			fmt.Println("Casos de teste:")

			for _, file := range files {
				if filepath.Ext(file.Name()) == ".json" {
					path := filepath.Join(outputDir, file.Name())
					data, err := os.ReadFile(path)
					if err != nil {
						continue
					}

					var test map[string]any
					if err := json.Unmarshal(data, &test); err != nil {
						continue
					}

					name := test["name"]
					rk := test["route_key"]
					fmt.Printf(" %s  (routeKey: %s)\n", name, rk)
				}
			}
		},
	}
}
