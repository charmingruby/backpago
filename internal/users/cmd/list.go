package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/charmingruby/backpago/internal/users"
	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lista usuários",
		Run: func(cmd *cobra.Command, args []string) {
			data, err := requests.AuthenticatedGet("/users")
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			var us []users.User
			err = json.Unmarshal(data, &us)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			for _, u := range us {
				log.Println(u.Name, u.Login, u.LastLogin)
			}
		},
	}

	return cmd
}
