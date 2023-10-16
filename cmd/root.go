package cmd

import (
	"github.com/joho/godotenv"
	"github.com/masonictemple4/masonictemple4.app/db"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "/etc/env/.blog.env", "config file (default is /etc/env/.blog.env)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// the cli has failed it is okay to panic here.
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "masonictemple4",
	Short: "CLI Interface to interact with masonictemple4 backend",
	Long:  "CLI Interface to interact with masonictemple4 backend",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		config, _ := cmd.Flags().GetString("config")

		err = godotenv.Load(config)
		if err != nil {
			// Should panic here we failed to load the config file.
			panic(err)
		}
		// Add a config loader here.

		DB, err = db.New(&gorm.Config{})
		if err != nil {
			// Should panic here we failed to initiate a database connection.
			panic(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}
