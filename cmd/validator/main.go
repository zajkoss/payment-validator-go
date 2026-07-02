package main

import (
	"fmt"
	"os"
	"payment-validator-go/internal/validator"
	"strings"
)

func main() {
	fmt.Println("=== Validator start ===")
	if len(os.Args) < 2 {
		fmt.Println("Required argument missing - directory")
		os.Exit(1)
	}

	directory := os.Args[1]
	fmt.Printf("\tDirectory: %s\n", directory)

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		os.Exit(1)
	}

	//Standard solution, iterations
	/*
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".xml") {
				fmt.Printf("\t\tValidation of: %s\n", file.Name())
				result, err := validator.ValidateFile(directory + "/" + file.Name()) //Required to change validateFile to public
				if err != nil {
					fmt.Printf("\t\t\tError: %s\n", err)
					continue
				}
				fmt.Printf("\t\t\tFile: %s, Valid: %v, Errors: %v\n", result.Filename, result.IsValid, result.Error)

			}
		}
	*/
	var paths []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".xml") {
			paths = append(paths, directory+"/"+file.Name())
		}
	}
	results := validator.RunWorkerPool(paths, 3)

	totalFiles := len(paths)
	totalValid := 0
	totalInvalid := 0

	fmt.Println("=== Validation Report ===")
	for _, result := range results {
		status := "OK"
		if result.IsValid {
			totalValid++
			fmt.Printf("\t%s - %s\n", result.Filename, status)
		} else {
			totalInvalid++
			status = "FAIL"
			fmt.Printf("\t%s - %s\n", result.Filename, status)
			for _, e := range result.Error {
				fmt.Printf("\t\t- %s\n", e)
			}
		}
	}
	fmt.Printf("Total files: %d Total valid: %d, Total invalid: %d\n", totalFiles, totalValid, totalInvalid)

}
