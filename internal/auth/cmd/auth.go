package cmd

import (
	"log"
	"os"

	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func authenticate() *cobra.Command {
	var (
		user string
		pass string
	)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Autentica usuário com API",
		Run: func(cmd *cobra.Command, args []string) {
			if user == "" || pass == "" {
				log.Println("usuário e senha são obrigatórios")
				os.Exit(1)
			}

			err := requests.Auth("/auth", user, pass)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&user, "user", "u", "", "nome do usuário")
	cmd.Flags().StringVarP(&pass, "pass", "p", "", "senha do usuário")

	return cmd
}
