package reader

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ArtemGontar/betting/internal/app/models"
)

func ReadMatchResultsFromDataset(filePath string) []models.MatchResult {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	matchesResults := []models.MatchResult{}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 23
	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		matchesResults = append(matchesResults, arrToMatchResultsMapping(record))
	}

	return matchesResults
}
