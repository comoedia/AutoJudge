package config

import (
	"encoding/json"
	"os"

	Log "hoseo.dev/autojudge/src/log"
)

const FILENAME = "./autojudge.json"

type AutoJudgeConfig struct {
	Root FileRoot
}

func New() {
	// Define an instance of the Root struct
	root := FileRoot{
		Schema: "https://pastebin.com/raw/AtjMwGB1",
		Problem: Problem{
			Title:       "",
			Description: "",
			Limit: Limit{
				Time:   "",
				Memory: "",
			},
			Try: 0,
		},
		Submit: Submit{
			Lang: Lang{
				Index: 0,
				Str:   "",
			},
			Before: Before{
				Test: "",
				Run:  "",
			},
		},
		Endpoint: Endpoint{
			Host: "",
			Resources: Resources{
				Problem:     "",
				Submissions: "",
				Submit:      "",
			},
		},
		Credentials: Credentials{
			Username: "",
			Password: "",
		},
	}

	// Marshal the struct to JSON
	jsonData, err := json.MarshalIndent(root, "", "    ")
	if err != nil {
		Log.Error.Panic("New > Error marshaling JSON:", err)
		return
	}

	// Write the JSON to a file
	os.WriteFile(FILENAME, jsonData, 0o644)

	Log.Verbose.Printf("Generating new configuration... length: %d\n", len(jsonData))
}

func SaveAll(config FileRoot) {
	// Marshal the struct to JSON
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		Log.Error.Panic("SaveAll > Error marshaling JSON:", err)
		return
	}

	// Write the JSON to a file
	os.WriteFile(FILENAME, jsonData, 0o644)

	Log.Verbose.Printf("Saving configuration... length: %d\n", len(jsonData))
}
