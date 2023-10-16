package db

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"gorm.io/gorm"
)

func TestMigrate(t *testing.T) {

	godotenv.Load("/etc/env/.blog.env")
	tdb, err := New(&gorm.Config{})
	if err != nil {
		t.Errorf("there was an error creating the database: %v", err)
	}

	var queryRepo models.Post

	var results []models.Post

	err = queryRepo.All(tdb, nil, &results)
	if err != nil {
		t.Errorf("there was an error querying the database: %v", err)
	}

	for _, post := range results {
		if strings.Contains(post.ContentUrl, post.Bucketname) {
			continue
		}

		post.ContentUrl = post.GenerateContentUrl()

		err = post.Update(tdb, int(post.ID), map[string]any{"contenturl": post.ContentUrl})
		if err != nil {
			t.Errorf("there was an error updating the post: %v", err)
		}
	}

	fmt.Printf("Migration complete")
}
