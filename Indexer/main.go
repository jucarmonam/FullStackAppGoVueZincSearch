package main

import (
	"Indexer/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//read fyle system folder
	//var finalEmails []string
	path := "C:\\Users\\juan\\Desktop\\Prueba truora\\enron_mail_20110402\\maildir"

	fmt.Println("Indexing")

	allEmails := readFolder(path)
	jsonData, err := json.Marshal(allEmails)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func readFolder(mainPath string) []models.Email {
	var folderList []string
	emails := make([]models.Email, 0)
	counter := 0
	err := filepath.Walk(mainPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		// Calculate the depth of the current directory
		depth := strings.Count(path, string(os.PathSeparator)) - strings.Count(mainPath, string(os.PathSeparator))
		if depth == 1 {
			folderList = append(folderList, info.Name())
		}

		if counter == 1 {
			return filepath.SkipDir // Exit the loop
		}

		if info.IsDir() {
			// Skip directories
			return nil
		}

		emails = append(emails, readFile(path))

		counter++
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return emails
}

func readFile(name string) models.Email {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	newEmail := models.Email{}

	for scanner.Scan() {
		line := scanner.Text()

		index := strings.Index(line, ":")
		if index != -1 && index < len(line)-1 {
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index+1:])

			// Fill the email object based on the key
			switch key {
			case "Message-ID":
				newEmail.Message_ID = value
			case "Date":
				newEmail.Date = value
			case "From":
				newEmail.From = value
			case "To":
				newEmail.To = value
			case "Subject":
				newEmail.Subject = value
			}
		} else {
			newEmail.Body = newEmail.Body + line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return newEmail
}
