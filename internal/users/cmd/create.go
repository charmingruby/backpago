package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/charmingruby/backpago/internal/users"
	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func create() *cobra.Command {
	var (
		name  string
		login string
		pass  string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Cria um novo usuário",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || login == "" || pass == "" {
				log.Println("Nome, login e senha são obrigatórios")
				os.Exit(1)
			}

			u := users.User{
				Name:     name,
				Login:    login,
				Password: pass,
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(u)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}
			_, err = requests.Post("/users", &body)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println("Usuário criado com sucesso!")
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Nome do usuário")
	cmd.Flags().StringVarP(&login, "login", "l", "", "Login do usuário")
	cmd.Flags().StringVarP(&pass, "pass", "p", "", "Senha do usuário")

	return cmd
}
