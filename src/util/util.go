package util

import (
	"log"
	"os"
	"strings"

	Log "hoseo.dev/autojudge/src/log"
)

func TrimString(str string) string {
	removedLn := strings.Replace(str, "\n", "", -1)
	return strings.TrimSpace(removedLn)
}

func GetTextFromFile(filename string) string {
	file, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		Log.Verbose.Printf("File %s does not exist.\n", filename)
		return ""
	}
	if err != nil {
		log.Fatal(err)
	}
	Log.Verbose.Printf("Reading %s... (size: %d)\n", filename, len(file))

	return string(file)
}

func IsExistFile(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	return false
}
