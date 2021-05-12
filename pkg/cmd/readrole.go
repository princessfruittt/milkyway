package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func walk(rolePath string) {

	fileList := make([]string, 0)
	e := filepath.Walk(rolePath, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	if e != nil {
		log.Fatal(e)
	}

	for _, file := range fileList {
		fmt.Println(file)
	}

	fmt.Print(fileList)
}
