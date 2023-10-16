package cmd

import (
	"log"

	"github.com/masonictemple4/masonictemple4.app/internal/auth"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/spf13/cobra"
)

var createSuperUserCmd = &cobra.Command{
	Use:   "create-super-user [email] [username] [password]",
	Short: "Creates a new user that is allowed to modify content on the blog.",
	Long:  `Creates a new user that is allowed to modify/manage content on the blog. This should be reserved for admins only.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Fatal("Please provide both a username and password.")
		}

		err := createNewUser(args[0], args[1], args[2])
		if err != nil {
			log.Fatalf("There was a problem creating the user: %v", err)
		}

		log.Printf("\nSuccessfully created user: %s", args[1])
	},
}

func init() {
	rootCmd.AddCommand(createSuperUserCmd)
}

func createNewUser(email, username, password string) error {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	return DB.Create(&models.User{
		Username: username,
		Password: hash,
	}).Error
}
