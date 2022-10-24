package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type generated struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	f, err := os.Open("./stories")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fmt.Println("File successfully created.")
	}()

	files, err := f.ReadDir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		f, err := os.Open(fmt.Sprintf("./stories/%s", v.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer func(f *os.File) {
			f.Close()
		}(f)

		results := make([]byte, 0)

		chunkSize := 512

		buf := make([]byte, chunkSize)
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			results = append(results, buf[:n]...)
		}

		fileName := strings.TrimSuffix(v.Name(), filepath.Ext(fmt.Sprintf("./stories/%s", v.Name())))

		results, err = json.Marshal(generated{
			Title:   fileName,
			Content: string(results),
		})
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(fmt.Sprintf("./results/%s.json", fileName), results, 0644)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully generated %s.json\n", fileName)
	}
}
