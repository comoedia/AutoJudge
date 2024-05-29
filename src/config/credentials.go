package config

import (
	"encoding/json"
	"os"

	Log "hoseo.dev/autojudge/src/log"
)

func GetCredentials() Credentials {
	// Read the JSON file
	file, err := os.ReadFile(FILENAME)
	if os.IsNotExist(err) {
		Log.Verbose.Printf("File %s does not exist.\n", FILENAME)
		return Credentials{
			Username: "",
			Password: "",
		}
	}

	// Define an instance of the Root struct
	var root FileRoot

	// Unmarshal JSON into the Root struct
	err = json.Unmarshal(file, &root)
	if err != nil {
		Log.Error.Panic("GetProblem > Error unmarshaling JSON:", err)
	}

	return root.Credentials
}

func SetCredentials(config Credentials) {
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
	root.Credentials = config

	// Save changes
	SaveAll(root)
}
