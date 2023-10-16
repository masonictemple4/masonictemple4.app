package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/masonictemple4/masonictemple4.app/db"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations for the masonictemple4 application. For now this will default to automigrate.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: In the future we'll specify for up and down migrations right now this
		// will default to automigrations.
		err := db.AutoMigrate(DB)
		if err != nil {
			log.Fatalf("There was a problem with the automigrations: %v", err)
		}

	},
}

func init() {
	migrateCmd.AddCommand(migrateBucketNameCmd)

	rootCmd.AddCommand(migrateCmd)
}

var migrateBucketNameCmd = &cobra.Command{
	Use:   "bucket <previous-bucketname> <new-bucket-name>",
	Short: "Change the bucketname anywhere it's persisted.",
	Long: `Will run through the dbs/ and any other location necessary replacing
	the previous storage bucket's name with the new one.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatalf("both previous and new bucket names are required.")
		}

		// params: prev, new
		err := migrateBucketNames(args[0], args[1])
		if err != nil {
			// because we previously printed the errors, we don't
			// need to do that again here.
			os.Exit(1)
		}

	},
}

// NOTE: Make sure that when adding fields to any models
// we remember to add them to the migrateBucketNames function
// Existing:
//
//	user: profilepicture
//	post: contenturl, bucketname, thumbnail
//	media: url, smallurl, mediumurl
func migrateBucketNames(prevBucket, newBucket string) error {

	userRepo := &models.User{}
	var users []models.User

	err := userRepo.All(DB, nil, &users)
	if err != nil {
		return err
	}

	var errors []error
	for _, user := range users {
		if strings.Contains(user.ProfilePicture, prevBucket) {
			user.ProfilePicture = strings.ReplaceAll(user.ProfilePicture, prevBucket, newBucket)
			err := user.UnsafeUpdate(DB)
			if err != nil {
				errors = append(errors, fmt.Errorf("userloop: failed to save update for user profile picture from %s bucket to %s new bucket. Result: %+v", prevBucket, newBucket, err))
				continue
			}
		}
	}

	postRepo := &models.Post{}
	var posts []models.Post
	err = postRepo.All(DB, nil, &posts)
	if err != nil {
		return err
	}

	// post contenturl, bucketname, thumbnail
	// media url, smallurl, mediumurl
	for _, post := range posts {
		if strings.Contains(post.ContentUrl, prevBucket) {
			post.ContentUrl = strings.ReplaceAll(post.ContentUrl, prevBucket, newBucket)
		}
		if strings.Contains(post.Thumbnail, prevBucket) {
			post.Thumbnail = strings.ReplaceAll(post.Thumbnail, prevBucket, newBucket)
		}
		if strings.Contains(post.Bucketname, prevBucket) {
			post.Bucketname = strings.ReplaceAll(post.Bucketname, prevBucket, newBucket)
		}

		err := post.UnsafeUpdate(DB)
		if err != nil {
			errors = append(errors, fmt.Errorf("postloop: failed to save update for post contenturl, thumbnail, bucketname from %s bucket to %s new bucket. Result: %+v", prevBucket, newBucket, err))
			continue
		}
	}

	mediaRepo := &models.Media{}
	var mediaList []models.Media

	err = mediaRepo.All(DB, nil, &mediaList)
	if err != nil {
		return err
	}

	for _, media := range mediaList {
		if strings.Contains(media.Url, prevBucket) {
			media.Url = strings.ReplaceAll(media.Url, prevBucket, newBucket)
		}
		if strings.Contains(media.SmallUrl, prevBucket) {
			media.SmallUrl = strings.ReplaceAll(media.SmallUrl, prevBucket, newBucket)
		}
		if strings.Contains(media.MediumUrl, prevBucket) {
			media.MediumUrl = strings.ReplaceAll(media.MediumUrl, prevBucket, newBucket)
		}

		err := media.UnsafeUpdate(DB)
		if err != nil {
			errors = append(errors, fmt.Errorf("medialoop: failed to save update for media url, smallurl, mediumurl from %s bucket to %s new bucket. Result: %+v", prevBucket, newBucket, err))
			continue
		}
	}

	if len(errors) > 0 {
		fmt.Println("There were a few errors during the migration:")
		for _, err := range errors {
			log.Println(err)
		}
		compositeErr := fmt.Errorf("there were %d errors during the migration\nHere they are %+v\n", len(errors), errors)

		return compositeErr
	}

	fmt.Println("Succesfully migrated all items.")
	return nil
}
