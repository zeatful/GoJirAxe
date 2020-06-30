/*
	This script is used to convert Axe Browser Extension issues .json files into Jira Formatted .txt files
	Usage:
		go run axe.go /foldername/

	foldername is the folder containing a bunch of .json files exported from the Axe browser extensions
	
	Each json file should be named in the following format:
		Page_Name_With_Spaces-Component_With_Spaces.json
	_ represents spaces
	- separates page title from the sub component
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"strconv"
	"path/filepath"
)

type Issues struct {
	Collection []Issue
}

// Each issue and its content
type Issue struct {
	Id          string   `json:"id"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Help        string   `json:"help"`
	HelpUrl     string   `json:"help_url"`
	Selector    []string `json:"selector"`
}

func main() {
	readAllFilesInDirectory(os.Args[1])
}

func fileNameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// get all files in directory and parse json
func readAllFilesInDirectory(directory string) {
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := directory + file.Name()
			readJson(filePath)
		}
	}
}

// parse a Json and output values formatted for Jira
func readJson(filename string) {
	// open json file
	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	// read values as bytes
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// initialize an empty array of Issue
	issues := make([]Issue, 0)

	// convert json into issues array
	json.Unmarshal(byteValue, &issues)

	// split up issues and write them to a file
	splitAndWriteIssuesToFile(fileNameWithoutExtension(filename), issues)
}

func splitAndWriteIssuesToFile(filename string, issues []Issue) {
	numberOfIssues := len(issues)
	maxNumberPerTicket := 50

	// split issues into splices of maxNumberPerTicket
	for i, j:= 0, 1; i < numberOfIssues; i, j = i + maxNumberPerTicket, j + 1 {
		end := i + maxNumberPerTicket
		
		// ensure we don't go out of range
		if end > numberOfIssues {
			end = numberOfIssues
		}
		
		// write the chunk
		writeIssuesToFile(j, filename, issues[i:end])
	}	
}

// write a chunk of issues to a file
func writeIssuesToFile(filenumber int, filename string, issues []Issue){
	fileindex := ""
	if(filenumber > 1) {
		fileindex = "-" + strconv.Itoa(filenumber)
	}

	// create jira .txt file
	outputFile, err := os.Create(filename + fileindex + ".txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// iterate through the array and output the data formatted
	for i := 0; i < len(issues); i++ {
		var issue = issues[i]
		outputFile.WriteString("Issue " + strconv.Itoa(i+1) + ".\n")
		outputFile.WriteString("*Summary:*  " + issue.Summary)
		outputFile.WriteString("\n")
		outputFile.WriteString("*Description:*  _" + issue.Description + "_")
		outputFile.WriteString("\n")
		outputFile.WriteString("*Source:*  {code:html}" + issue.Source + "{code}")
		outputFile.WriteString("\n")
		outputFile.WriteString("*Selector:*  {code:css}" + issue.Selector[0] + "{code}")
		outputFile.WriteString("\n")
		outputFile.WriteString("*Help:*  _" + issue.Help + "_")
		outputFile.WriteString("\n")
		outputFile.WriteString("*Help URL:*  " + issue.HelpUrl)
		outputFile.WriteString("\n")
		outputFile.WriteString("*Fix:* _very generally state how it was addressed_")
		outputFile.WriteString("\n")
		outputFile.WriteString("----")
		outputFile.WriteString("\n")
		outputFile.Sync()
	}

	fmt.Println(outputFile.Name() + " successfully converted to jira")
}