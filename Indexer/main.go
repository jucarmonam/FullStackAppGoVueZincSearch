package main

import (
	"Indexer/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
)

func main() {
	path := "C:\\Users\\juan\\Desktop\\PruebaTruora\\enron_mail_20110402\\maildir"
	//path := os.Args[1]

	////Proceso de rendimiento de la aplicación/////////////
	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()
	////////Fin prceso de rendimiento de la aplicación/////////

	fmt.Println("Indexing")
	//read fyle system folder
	allEmails := readFolder(path)

	var ndjson bytes.Buffer

	// Generate the NDJSON variable
	for _, email := range allEmails {
		jsonData, err := json.Marshal(email)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		ndjson.Write(jsonData)
		ndjson.WriteString("\n")
	}

	//Inserting in zinsearch
	zincSearchInsert(ndjson.String())

	////Proceso de rendimiento de la aplicación/////////////
	runtime.GC()
	mem, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal(err)
	}
	////Fin proceso de rendimiento de la aplicación/////////////
}

func readFolder(mainPath string) []models.Email {
	emails := make([]models.Email, 0)
	counter := 0
	err := filepath.Walk(mainPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if counter == 1000 {
			return filepath.SkipDir
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

	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 2048*1024)
	newEmail := models.Email{}

	lineNumber := 0
	counter := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		index := strings.Index(line, ":")
		if index != -1 && index < len(line)-1 && counter < 5 {
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
			counter++
		} else {
			newEmail.Body = newEmail.Body + line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("something bad happened in the line %v: %v", lineNumber, err)
	}

	return newEmail
}

func zincSearchInsert(emails string) {
	req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_multi", strings.NewReader(emails))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
