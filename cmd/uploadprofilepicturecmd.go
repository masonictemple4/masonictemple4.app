package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/masonictemple4/masonictemple4.app/internal/filestore"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/masonictemple4/masonictemple4.app/utils"
	"github.com/spf13/cobra"
)

var uploadprofilepicturecmd = &cobra.Command{
	Use:   `user-pic <username> <path-to-file>`,
	Short: "Uploads a profile picture for a user.",
	Long:  `Uploads a profile picture for a user. The user must exist in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("please provide a path to the file to upload")
		}

		if err := uploadPic(cmd.Context(), args[0], args[1]); err != nil {
			log.Fatalf("there was a problem uploading the profile picture: %v", err)
		}

		fmt.Printf("Profile picture uploaded successfully\n")
	},
}

func init() {
	rootCmd.AddCommand(uploadprofilepicturecmd)
}

func uploadPic(ctx context.Context, username, path string) error {
	var user models.User

	user.FindByUsername(DB, username, nil)

	if strings.Contains(path, "http") {
		updateParams := map[string]any{"profilepicture": path}
		if err := user.Update(DB, int(user.ID), updateParams); err != nil {
			return err
		}
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := utils.DetectFileType(path)

	fileHandler := filestore.NewGCPStore(false, 0)

	uploadPath := user.GenerateProfilePicturePath(ext)
	if uploadPath == "" {
		return errors.New("issue with generating path for picture upload")
	}

	written, err := fileHandler.Write(ctx, uploadPath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateParams := map[string]any{"profilepicture": user.GenerateProfilePictureURL(uploadPath)}

	if err := user.Update(DB, int(user.ID), updateParams); err != nil {
		return err
	}

	return nil
}
