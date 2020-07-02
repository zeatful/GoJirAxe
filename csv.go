package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"encoding/csv"
	"path/filepath"
)

// store information for each file / ticket
type Ticket struct {
	Summary		string
	IssueType	string
	EpicLink	string
	Description string
}

// create array of Tickets
var tickets []Ticket
var accessibilityEpic string
var issueType string
var descriptionTemplate string

// custom error check
func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}

func main() {
	// default placeholder values
	accessibilityEpic = "ADS - Web Content Accessibility Guidelines"
	issueType = "Improvement"
	descriptionTemplate = "+Need:+\n"+
	"-update the Home page of ADS to be WCAG 2.1 AAA compliant and address Axe violations\n"+
	"+Steps:+\n"+
	"- run the Axe tools against the page\n"+
	"- fix any issues reported by axe or denote / justify any issue that is not valid or applicable\n"
	
	getTextFiles()
	writeCSV()
}

func getTextFiles() {
    root := "./issues"
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// don't capture folders
		if !info.IsDir() {
			// only read text files
			if filepath.Ext(path) == ".txt" {
				// reader in each file and parse it
				readIssue(path, info)
			}
		}
        return nil
	})
	
    checkError("Error getting files from " + root, err)
}

func readIssue(path string, info os.FileInfo) {
	content, err := ioutil.ReadFile(path)
	checkError("Error reading: " + path, err)

	// format summary
	summary := strings.ReplaceAll(info.Name(), "_", " ")
	summary = strings.ReplaceAll(summary, "-", " - ")
	summary = strings.ReplaceAll(summary, ".txt", "")

	fmt.Println(summary)

    // Convert []byte to string and print to screen
	text := string(content)	
	tickets = append(tickets, Ticket{summary, issueType, accessibilityEpic, text})
}

func writeCSV() {
	// // create export file and write out
	file, err := os.Create("import.csv")
    checkError("Cannot create file", err)
    defer file.Close()

    writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Summary", "Issue Type", "Epic Link", "Description"})
	checkError("Cannot write to file", err)

    for _, ticket := range tickets {
		row := []string{ticket.Summary, ticket.IssueType, ticket.EpicLink, descriptionTemplate + ticket.Description}
        err = writer.Write(row)
        checkError("Cannot write to file", err)
    }
}