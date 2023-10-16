package server

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

func TestGoogleAuth(t *testing.T) {
	err := godotenv.Load("/etc/env/.masonictemple4env")
	if err != nil {
		// Should panic here we failed to load the config file.
		t.Errorf("failed to load config.")
	}
	t.Run("Test Invalid Auth Token", func(t *testing.T) {
		// TODO: Put token here.
		token := ""

		_, err = idtoken.Validate(context.Background(), token, os.Getenv("GOOGLE_ID"))
		if err == nil {
			t.Errorf("Token info passed with invalid access token.")
		}
	})

	t.Run("Test Valid Auth Token", func(t *testing.T) {
		// TODO: Add token here
		token := ""

		_, err = idtoken.Validate(context.Background(), token, os.Getenv("GOOGLE_ID"))
		if err != nil {
			t.Errorf("Token info failed with valid access token.")
		}

	})

}

func TestGithubAuth(t *testing.T) {
	t.Run("Test Invalid Auth Token", func(t *testing.T) {
		// TODO: Add token here.
		token := ""

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		user, _, err := client.Users.Get(ctx, "")
		if err == nil {
			t.Errorf("Token info passed with invalid access token.")
		}

		if user != nil {
			t.Errorf("Token info passed with invalid access token.")
		}
	})
	t.Run("Test Valid Auth Token", func(t *testing.T) {
		// TODO: Add token here.
		token := ""

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			t.Errorf("Token info failed with valid access token.")
		}
		if user == nil {
			t.Errorf("Token info failed with valid access token.")
		}
	})
}
