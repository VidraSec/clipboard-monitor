package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// helper variable to prevent writing to file when the program has just written to it
	programWroteToFile := false
	// helper variable to prevent writing to clipboard when the program has just written to it
	programWroteToClipboard := false

	filePath := "clipboard.txt"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %s", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && !programWroteToFile {
					fmt.Println("File change: file -> clipboard")
					content, err := os.ReadFile(filePath)
					if err != nil {
						log.Printf("Error reading file: %s", err)
						continue
					}
					if err := clipboard.WriteAll(string(content)); err != nil {
						log.Printf("Error writing to clipboard: %s", err)
					}
					programWroteToClipboard = true
				} else if programWroteToFile {
					// if the program wrote to the file, reset the flag
					programWroteToFile = false
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %s", err)
			}
		}
	}()

	// Ensure the file exists before starting to watch it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist, creating file.")
		// Create the file with empty content
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			log.Fatalf("Failed to create %s: %s", filePath, err)
		}
	}

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatalf("Error adding watcher to file: %s", err)
	}

	// Initial read and clipboard set
	initialContent, err := os.ReadFile(filePath)
	if err == nil {
		clipboard.WriteAll(string(initialContent))
	}

	// Start a ticker for periodic file check as a fallback
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		var lastModTime time.Time
		for range ticker.C {
			fi, err := os.Stat(filePath)
			if err != nil {
				log.Printf("Error stating file: %s", err)
				continue
			}
			if fi.ModTime().After(lastModTime) && !programWroteToFile {
				fmt.Println("Fallback file check: file -> clipboard")
				content, err := os.ReadFile(filePath)
				if err != nil {
					log.Printf("Error reading file: %s", err)
					continue
				}
				if err := clipboard.WriteAll(string(content)); err != nil {
					log.Printf("Error writing to clipboard: %s", err)
				}
				lastModTime = fi.ModTime()
				programWroteToClipboard = true
			}
		}
	}()

	// Monitor clipboard changes and update the file
	previousClipboardContent, _ := clipboard.ReadAll()
	for {
		time.Sleep(1 * time.Second)
		currentClipboardContent, _ := clipboard.ReadAll()
		if currentClipboardContent != previousClipboardContent && !programWroteToClipboard {
			if err := os.WriteFile(filePath, []byte(currentClipboardContent), 0644); err != nil {
				log.Printf("Error writing to file: %s", err)
			}
			fmt.Println("Clipboard change: clipboard -> file")
			previousClipboardContent = currentClipboardContent
			programWroteToFile = true
		} else if programWroteToClipboard {
			// if the program wrote to the clipboard, reset the flag
			programWroteToClipboard = false
		}
	}
}
