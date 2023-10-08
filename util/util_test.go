package util

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestSaveImageFromBase64(t *testing.T) {
	data, err := os.ReadFile("b64image.txt")
	if err != nil {
		log.Fatalln(err)
	}

	if err := SaveImageFromBase64(string(data), "./test.png"); err != nil {
		fmt.Println(err)
	}
}
