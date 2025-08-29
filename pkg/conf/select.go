package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func (c *Conf) CmdSelect() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "select [nome]",
		Short: "Seleciona uma configuração existente ou cria uma nova",
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			baseDir := c.settings.GetBaseDir()
			files := listConfigFiles(baseDir)
			var res []string
			for _, f := range files {
				if strings.HasPrefix(strings.ToLower(f), strings.ToLower(toComplete)) {
					res = append(res, f)
				}
			}
			return res, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			baseDir := c.settings.GetBaseDir()
			_ = os.MkdirAll(baseDir, os.ModePerm)

			if len(args) == 0 {
				opts := listConfigFiles(baseDir)
				if len(opts) == 0 {
					fmt.Println("Nenhuma configuração encontrada. Informe um nome para criar uma nova, " +
						"por exemplo: rabbix conf select minha.json")

				} else {
					fmt.Println("Informe o nome da configuração. Disponíveis:")
					for _, o := range opts {
						fmt.Println("- " + o)
					}
				}
				return
			}

			name := args[0]
			if !strings.HasSuffix(strings.ToLower(name), ".json") {
				name += ".json"
			}
			target := filepath.Join(baseDir, name)

			// Create the configuration if it does not exist.
			if _, err := os.Stat(target); os.IsNotExist(err) {
				defaultCfg := map[string]string{
					"auth":       "Z3Vlc3Q6Z3Vlc3Q=",
					"host":       "http://localhost:15672",
					"output_dir": filepath.Join(baseDir, "tests"),
				}
				if data, err := json.MarshalIndent(defaultCfg, "", "  "); err == nil {
					_ = os.WriteFile(target, data, 0644)
				}
				fmt.Println("Criada nova configuração:", name)
			}

			// Updates settings.json with the selected file
			settPath := filepath.Join(baseDir, "settings.json")
			settings := map[string]string{"sett": name}
			if data, err := os.ReadFile(settPath); err == nil {
				_ = json.Unmarshal(data, &settings)
				settings["sett"] = name
			}
			if data, err := json.MarshalIndent(settings, "", "  "); err == nil {
				_ = os.WriteFile(settPath, data, 0644)
			}

			fmt.Println("Configuração ativa atualizada para:", name)
		},
	}

	return cmd
}

func listConfigFiles(baseDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return []string{}
	}

	var out []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		n := e.Name()
		if strings.EqualFold(n, "settings.json") || strings.EqualFold(n, "cache.json") {
			continue
		}

		if strings.HasSuffix(strings.ToLower(n), ".json") {
			out = append(out, strings.TrimSuffix(n, ".json"))
		}
	}

	return out
}
