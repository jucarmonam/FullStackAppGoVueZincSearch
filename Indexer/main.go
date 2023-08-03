package main

import (
	"Indexer/models"
	"bufio"
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
	"sync"
)

func readFolder(mainPath string) {
	maxWorkers := 4
	var wg sync.WaitGroup
	// Create a channel to receive file paths
	paths := make(chan string)

	//counter := 0

	// Start the worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				readFile(path)
			}
		}()
	}

	err := filepath.Walk(mainPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		paths <- path
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	close(paths)

	wg.Wait()
}

func readFile(name string) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var buf []byte
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

	jsonData, _ := json.Marshal(newEmail)
	zincSearchInsert(jsonData)
}

func zincSearchInsert(email []byte) {
	req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_doc", strings.NewReader(string(email)))
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

func main() {
	//path := "C:\\Users\\juan\\Desktop\\PruebaTruora\\enron_mail_20110402\\maildir"
	path := os.Args[1]

	//Inicio de profiling
	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()

	fmt.Println("Indexing")

	//read fyle system folder
	readFolder(path)

	//Se crean archivos finales de prifiling
	runtime.GC()
	mem, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Indexing finished")
	//go tool pprof -http=:8020 cpu.prof
}
