package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmingruby/backpago/internal/folders"
	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lista conteúdo de uma pasta",
		Run: func(cmd *cobra.Command, args []string) {
			path := "/folders"
			if id > 0 {
				path = fmt.Sprintf("/folders/%d", id)
			}

			data, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			var fc folders.FolderContent
			err = json.Unmarshal(data, &fc)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println(fc.Folder.Name)
			log.Println("================")
			for _, c := range fc.Content {
				log.Println(c.ID, " - ", c.Type, " - ", c.Name)
			}
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "ID da pasta")

	return cmd
}
