package reader

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ArtemGontar/betting/internal/app/model"
)

func ReadMatchResultsFromDataset(filePath string) ([]model.MatchResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	matchesResults := []model.MatchResult{}
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

	return matchesResults, nil
}
