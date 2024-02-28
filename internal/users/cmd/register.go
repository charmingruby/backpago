package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Gestão de usuários",
	}

	cmd.AddCommand(create())
	cmd.AddCommand(update())
	cmd.AddCommand(delete())
	cmd.AddCommand(list())
	cmd.AddCommand(get())

	c.AddCommand(cmd)
}
