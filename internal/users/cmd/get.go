package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmingruby/backpago/internal/users"
	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func get() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Obtém informações do usuário",
		Run: func(cmd *cobra.Command, args []string) {
			if id <= 0 {
				log.Println("ID do usuário é obrigatório")
				os.Exit(1)
			}

			path := fmt.Sprintf("/users/%d", id)
			data, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			var u users.User
			err = json.Unmarshal(data, &u)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println(u.Name)
			log.Println(u.Login)
			log.Println(u.LastLogin)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "ID do usuário")

	return cmd
}
