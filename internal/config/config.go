package config

import (
	"encoding/json"
	"os"
	"path"
	"sort"

	"github.com/ImDevinC/go-pd3/internal/models"
	"github.com/adrg/xdg"
)

var configDir = path.Join(xdg.DataHome, "pd3-challenges")
var challengesFile = "challenges.json"

func LoadSaved() ([]models.PD3DataResponse, error) {
	challenges, err := loadExisting(path.Join(configDir, challengesFile))
	if err != nil {
		challenges, err = loadExisting("default.json")
	}
	if err != nil {
		return challenges, err
	}
	sort.Slice(challenges, func(i, j int) bool {
		return challenges[i].Challenge.Name < challenges[j].Challenge.Name
	})
	return challenges, nil
}

func SaveChallenges(challenges []models.PD3DataResponse) error {
	err := os.MkdirAll(configDir, 0644)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(challenges)
	if err != nil {
		return err
	}
	file := path.Join(configDir, challengesFile)
	err = os.WriteFile(file, payload, 0644)
	return err
}

func loadExisting(file string) ([]models.PD3DataResponse, error) {
	response := []models.PD3DataResponse{}
	payload, err := os.ReadFile(file)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(payload, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
