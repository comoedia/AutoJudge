package config

import (
	"encoding/json"
	"os"

	Log "hoseo.dev/autojudge/src/log"
)

func GetProblem() Problem {
	// Read the JSON file
	file, err := os.ReadFile(FILENAME)
	if os.IsNotExist(err) {
		Log.Verbose.Printf("File %s does not exist.\n", FILENAME)
		return Problem{
			Title:       "",
			Description: "",
			Limit: Limit{
				Time:   "",
				Memory: "",
			},
			Try: 0,
		}
	}

	// Define an instance of the Root struct
	var root FileRoot

	// Unmarshal JSON into the Root struct
	err = json.Unmarshal(file, &root)
	if err != nil {
		Log.Error.Panic("GetProblem > Error unmarshaling JSON:", err)
	}

	return root.Problem
}

func SetProblem(config Problem) {
	// Read the JSON file
	file, err := os.ReadFile(FILENAME)
	if os.IsNotExist(err) {
		Log.Verbose.Printf("File %s does not exist.\n", FILENAME)
		return
	}

	// Define an instance of the Root struct
	var root FileRoot

	// Unmarshal JSON into the Root struct
	err = json.Unmarshal(file, &root)
	if err != nil {
		Log.Error.Panic("SetProblem > Error unmarshaling JSON:", err)
	}

	// Set the Problem field of the Root struct (overwrite)
	root.Problem = config

	// Save changes
	SaveAll(root)
}
