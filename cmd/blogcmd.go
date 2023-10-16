package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
	"github.com/masonictemple4/masonictemple4.app/internal/filestore"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/masonictemple4/masonictemple4.app/internal/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var blogCmd = &cobra.Command{
	Use:   `blog <path-to-file>`,
	Short: "Creates a new blog post from markdown file.",
	Long: `Creates a new blog record and uploads the markdown file (content)
	to the storage bucket.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		if len(args) < 1 {
			panic("please provide a file path with the markdown for the post")
		}
		if err := createBlog(cmd.Context(), args[0], flags); err != nil {
			log.Fatalf("there was a problem creating the blog: %v", err)
		}
	},
}

var updateBlogCmd = &cobra.Command{
	Use:   `update <path-to-file>`,
	Short: "Updates a blog post from markdown file.",
	Long: `Updates a blog record and uploads the markdown file (content)
	to the storage bucket.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			panic("please provide a file path with the markdown for the post")
		}

		bid, _ := cmd.Flags().GetInt("id")

		if err := updateBlog(cmd.Context(), bid, args[0]); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully updated blog!")
	},
}

func init() {
	blogCmd.AddCommand(updateBlogCmd)
	blogCmd.PersistentFlags().IntP("id", "i", 0, "The id of the blog.")
	blogCmd.MarkPersistentFlagRequired("id")

	rootCmd.AddCommand(blogCmd)
}

func createBlog(ctx context.Context, path string, flags *pflag.FlagSet) error {
	// What if instead of passing a path to parser here and eventually would have to be
	// the filestore too if i read file here and pass bytes to the parsr and writer.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result dtos.PostInput
	err = parser.ParseFile(path, &result)
	if err != nil {
		return err
	}

	fmt.Printf("\nThe parsed result is:  %+v\n", result)

	post := &models.Post{
		Bucketname: os.Getenv("STORAGE_BUCKET"),
		State:      models.BlogStateDraft,
	}

	err = post.FromPostInput(DB, &result)
	if err != nil {
		return err
	}

	// generate a slug because now we have a title.
	post.Slug = post.GenerateSlug("")

	fileHandler := filestore.NewGCPStore(false, 0)

	post.Docpath, err = post.GenerateDocPath()
	if err != nil {
		return err
	}

	written, err := fileHandler.Write(ctx, post.Docpath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateBody := map[string]any{"contenturl": post.GenerateContentUrl(), "docpath": post.Docpath, "state": models.BlogStatePublished, "slug": post.Slug}
	err = post.Update(DB, int(post.ID), updateBody)
	if err != nil {
		return err
	}

	return nil
}

func updateBlog(ctx context.Context, bid int, path string) error {
	var blog models.Post

	err := blog.FindByID(DB, bid, nil)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result dtos.PostInput
	err = parser.ParseFile(path, &result)
	if err != nil {
		return err
	}

	fmt.Printf("\nThe parsed result is:  %+v\n", result)

	err = blog.FromPostInput(DB, &result)
	if err != nil {
		return err
	}

	// generate a slug because now we have a title.
	blog.Slug = blog.GenerateSlug("")

	fileHandler := filestore.NewGCPStore(false, 0)

	// TODO: Should probably delete here to be safe we're not
	// gathering too unused files.
	blog.Docpath, err = blog.GenerateDocPath()
	if err != nil {
		return err
	}

	written, err := fileHandler.Write(ctx, blog.Docpath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateBody := map[string]any{
		"contenturl": blog.GenerateContentUrl(),
		"docpath":    blog.Docpath,
		"state":      models.BlogStatePublished,
		"slug":       blog.Slug,
	}
	err = blog.Update(DB, int(blog.ID), updateBody)
	if err != nil {
		return err
	}

	return nil

}
