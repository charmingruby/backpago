package cmd

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmingruby/backpago/pkg/requests"
	"github.com/spf13/cobra"
)

func upload() *cobra.Command {
	var (
		filename string
		folderID int32
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Faz upload de um novo arquivo",
		Run: func(cmd *cobra.Command, args []string) {
			if filename == "" {
				log.Println("Caminho do arquivo é obrigatório")
				os.Exit(1)
			}

			file, err := os.Open(filename)
			if err != nil {
				log.Println("%v", err)
				os.Exit(1)
			}
			defer file.Close()

			var body bytes.Buffer

			mw := multipart.NewWriter(&body)
			w, err := mw.CreateFormFile("file", filepath.Base(file.Name()))
			if err != nil {
				log.Println("%v", err)
				os.Exit(1)
			}
			io.Copy(w, file)

			if folderID > 0 {
				w, err := mw.CreateFormField("folder_id")
				if err != nil {
					log.Println("%v", err)
					os.Exit(1)
				}
				w.Write([]byte(strconv.Itoa(int(folderID))))
			}

			mw.Close()

			headers := map[string]string{
				"Content-Type": mw.FormDataContentType(),
			}

			_, err = requests.AuthenticatedPostWithHeaders("/files", &body, headers)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			log.Println("Arquivo enviado com sucesso!")
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "Caminho para o arquivo")
	cmd.Flags().Int32VarP(&folderID, "folder", "p", 0, "ID da pasta onde o arquivo será enviado")

	return cmd
}
