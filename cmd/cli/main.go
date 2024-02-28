package main

import (
	"log"

	authCmd "github.com/charmingruby/backpago/internal/auth/cmd"
	filesCmd "github.com/charmingruby/backpago/internal/files/cmd"
	foldersCmd "github.com/charmingruby/backpago/internal/folders/cmd"
	usersCmd "github.com/charmingruby/backpago/internal/users/cmd"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{}

func main() {
	authCmd.Register(RootCmd)
	filesCmd.Register(RootCmd)
	foldersCmd.Register(RootCmd)
	usersCmd.Register(RootCmd)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
