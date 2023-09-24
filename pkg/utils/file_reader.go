package utils

import (
	"log"
	"os"
)

func GetBytes(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}
