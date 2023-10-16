package cmd

import (
	"log"

	"github.com/masonictemple4/masonictemple4.app/db"
	"github.com/masonictemple4/masonictemple4.app/internal/server"
	"github.com/spf13/cobra"
)

var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run the masonictemple4 server",
	Long:  `Run the masonictemple4 server`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		runServer(host, port)

	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)
	runserverCmd.PersistentFlags().StringP("port", "p", "8080", "The port to run on. Default: 8080")
	runserverCmd.PersistentFlags().StringP("host", "H", "0.0.0.0", "The host to run on. Default: 0.0.0.0")

}

func runServer(host, port string) {
	srv := server.NewServer(DB)
	// TODO: Might make more sense to put this in the newserver method.
	err := db.AutoMigrate(DB)
	if err != nil {
		log.Fatalf("There was a problem with the automigrations: %v", err)
	}
	srv.Run(host + ":" + port)
}
