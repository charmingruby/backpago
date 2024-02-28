package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func delete() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deleta um usuário",
		Run: func(cmd *cobra.Command, args []string) {
			if id <= 0 {
				log.Println("ID é obrigatório")
				os.Exit(1)
			}

			path := fmt.Sprintf("/users/%d", id)

			err := requests.AuthenticatedDelete(path)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println("Usuário deletada com sucesso!")
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "ID do usuário")

	return cmd
}
