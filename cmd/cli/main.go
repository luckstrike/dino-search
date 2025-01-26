package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/luckstrike/dino-search/internal/crawler"
	"github.com/luckstrike/dino-search/internal/scraper"
	"github.com/luckstrike/dino-search/internal/storage"
)

const (
	prompt = "search> "
)

func main() {
	db, err := storage.InitDB()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}
	defer db.Close()

	fmt.Println("Welcome to the Search Engine")
	fmt.Println("Type 'quit' or 'exit' to close the program")
	fmt.Println("Type 'help' for available commands")

	runCLI()
}

func runCLI() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if err := handleCommand(input); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func handleCommand(input string) error {
	// Convert to lowercase for easier comparison
	command := strings.ToLower(input)

	switch command {
	case "":
		return nil
	case "quit", "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "help":
		printHelp()
		return nil
	default:
		return performSearch(input)
	}

	return nil
}

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("\thelp - Show this help message")
	fmt.Println("\tquit - Exit the program")
	fmt.Println("\text - Exit the program")
	fmt.Println("\nOr simply type your search query")
	fmt.Println()
}

func performSearch(query string) error {
	fmt.Printf("Searching for %s\n", query)

	crawler.Crawl(query)
	return nil
}

// This could probably be combined with the performSearch function in the future
func processURL(url string) error {
	scraper := scraper.NewScraper()
	content, err := scraper.Scrape(url)

	// TODO: Change this later when you add in a way to save website content
	if false {
		fmt.Println(content)
	}

	if err != nil {
		return err
	}

	return nil
}
