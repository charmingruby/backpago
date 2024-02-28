package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Gestão de arquivos",
	}

	cmd.AddCommand(upload())
	cmd.AddCommand(update())
	cmd.AddCommand(delete())

	c.AddCommand(cmd)
}
