package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmingruby/backpago/internal/files"
	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func update() *cobra.Command {
	var id int32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Atualiza o nome de um arquivo",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || id <= 0 {
				log.Println("Nome do arquivo e ID são obrigatórios")
				os.Exit(1)
			}

			f := files.File{Name: name}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(f)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			path := fmt.Sprintf("/files/%d", id)
			_, err = requests.AuthenticatedPut(path, &body)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println("Arquivo atualizado com sucesso!")
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "ID do arquivo")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Nome do arquivo")

	return cmd
}
