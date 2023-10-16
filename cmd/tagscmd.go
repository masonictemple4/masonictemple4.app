package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/masonictemple4/masonictemple4.app/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm/clause"
)

var tagsCmd = &cobra.Command{
	Use:   "tags [tags]",
	Short: "Manage pre-defined tags",
	Long:  `Manage pre-defined tags on your platform`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			processTagsFile(file)
		} else if len(args) == 0 {
			log.Fatal("Please provide at least 1 tag to add.")
		} else {
			for _, tag := range args {
				var t models.Tag
				if err := DB.FirstOrCreate(&t, models.Tag{Name: tag}).Error; err != nil {
					log.Fatalf("There was a problem creating the tag: %v", err)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
	rootCmd.PersistentFlags().StringP("file", "f", "", "The file you would like to load the tags from.")
}

func processTagsFile(file string) {
	log.Printf("Processing tags from file %s", file)
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("there was a problem reading the file: %v", err)
	}

	if len(data) == 0 {
		log.Printf("The file was empty.")
	}

	var results []dtos.TagInput
	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Fatalf("Invalid format. Make sure the json file conatins a list of TagInput objects.\nErr: %v", err)
	}

	var objects []models.Tag
	err = utils.Convert(results, &objects)
	if err != nil {
		log.Fatalf("There was a problem converting the results to Tag objects.\nErr: %v", err)
	}

	// The Do nothing clause will allow the insert to continue ignoring existing unique records.
	err = DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&objects).Error
	if err != nil {
		log.Fatalf("There was a problem saving the tags to the database.\nErr: %v", err)
	}

}
