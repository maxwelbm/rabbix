package conf

import (
	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

var (
	host      string
	outputDir string
	user      string
	password  string
)

type Conf struct {
	settings sett.SettItf
}

func New(settings sett.SettItf) *Conf {
	return &Conf{
		settings: settings,
	}
}

func (c *Conf) CmdConf() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conf",
		Short: "Define ou exibe configurações padrão da CLI",
	}

	cmd.AddCommand(c.CmdSet())
	cmd.AddCommand(c.CmdGet())

	return cmd
}
