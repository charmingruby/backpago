package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/charmingruby/backpago/internal/folders"
	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func create() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Cria uma nova pasta",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				log.Println("Nome da pasta é obrigatório")
				os.Exit(1)
			}

			folder := folders.Folder{Name: name}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(folder)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}
			_, err = requests.AuthenticatedPost("/folders", &body)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println("Pasta criada com sucesso!")
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Nome da pasta")

	return cmd
}
