package cmd

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/spf13/cobra"
)

var tagStressTestCmd = &cobra.Command{
	Use:   "bulk-tags",
	Short: "Bulk create tags",
	Long:  `Bulk create tags`,
	Run: func(cmd *cobra.Command, args []string) {
		const total = 8_000_000

		numWorkers := runtime.NumCPU()
		var wg sync.WaitGroup
		wg.Add(numWorkers)

		jobs := make(chan int, total)
		for i := 0; i < numWorkers; i++ {
			go worker(jobs, &wg)
		}

		close(jobs)
		wg.Wait()

	},
}

func worker(jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		var t models.Tag
		if err := DB.FirstOrCreate(&t, models.Tag{Name: fmt.Sprintf("tag-%d", job)}).Error; err != nil {
			log.Fatalf("There was a problem creating the tag: %v", err)
		}

		fmt.Printf("Created tag %d\n", job)
	}
}

func init() {
	rootCmd.AddCommand(tagStressTestCmd)
}
