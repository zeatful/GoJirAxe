# GoJirAxe

Go scripts to support parsing Axe Violation JSON files exports from the Axe Browser Extensions and then turn them into an import.csv which can be bulk imported into Jira

## axe.go

Used to parse a directory of JSON files
Usage:
```go run axe.go /issues/```

## csv.go

Used to parse a directory of text files into a jira bulk import CSV
It directly points to an issues directory
Usage:
``` go run csv.go```