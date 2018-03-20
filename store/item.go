package store

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Variable is a piece of configuration
type Variable struct {
	Name  string `dynamodbav:"name" json:"name"`
	Value string `dynamodbav:"value" json:"value"`
}

// Item is the format of the configuratoin stored in dynamodb
type Item struct {
	ID        string     `dynamodbav:"id" json:"id"`
	Variables []Variable `dynamodbav:"variables" json:"variables"`
}

// PrintVars prints the variables in the item
func (item *Item) PrintVars(format string) {
	format = strings.ToLower(format)
	if format == "json" {
		item.printJSON()
	} else if format == "sh" {
		item.printShell()
	} else {
		item.printPlain()
	}
}

func (item *Item) printPlain() {
	for i := range item.Variables {
		fmt.Printf("%s=%s\n", item.Variables[i].Name, item.Variables[i].Value)
	}
}

func (item *Item) printJSON() {
	// The default json.unmarshal HTML escapes the string
	// We create a custom encoder so we don't have to HTML escape
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	err := encoder.Encode(item.Variables)
	if err != nil {
		// TODO debug print error or something
		fmt.Println("ERROR: ", err)
	}
}

func (item *Item) printShell() {
	for i := range item.Variables {
		fmt.Printf("export %s=%s\n", item.Variables[i].Name, item.Variables[i].Value)
	}
}

// TODO this is pretty darn primitive so make it more robust
// support other formats and what not
func parseVariables(variablesRaw string, nameOnly bool) []Variable {
	keyValuePairs := strings.Split(variablesRaw, ",")
	variables := make([]Variable, len(keyValuePairs))
	for i, keyValuePair := range keyValuePairs {
		if nameOnly {
			variables[i] = Variable{
				Name:  string(keyValuePair),
				Value: "",
			}
		} else {
			indexOfSeparator := strings.Index(keyValuePair, "=")
			variables[i] = Variable{
				Name:  string(keyValuePair[:indexOfSeparator]),
				Value: string(keyValuePair[indexOfSeparator+1:]),
			}
		}
	}
	return variables
}

func parseVariablesFromFile(fileName string, nameOnly bool) ([]Variable, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	return parseVariablesFromScanner(fileScanner, nameOnly)
}

func parseVariablesFromScanner(scanner *bufio.Scanner, nameOnly bool) ([]Variable, error) {
	variables := make([]Variable, 0)

	for scanner.Scan() {
		// Parse line
		line := scanner.Text()
		line = strings.TrimPrefix(line, "export") // remove export if exists
		line = strings.TrimLeft(line, " \t")      // remove all spaces on left
		line = strings.TrimRight(line, " \t\n")   // trim right
		if line == "" || line[0] == byte('#') {   // try to skip comments and empty lines
			continue
		}
		var name, value string
		if nameOnly {
			name = line
		} else {
			indexOfSeparator := strings.Index(line, "=") // split into no more than 2 strings
			if indexOfSeparator < 0 {
				return variables, fmt.Errorf("error parsing file")
			}
			name = line[:indexOfSeparator]
			value = line[indexOfSeparator+1:]
		}
		variables = append(variables, Variable{Name: name, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return variables, err
	}
	return variables, nil
}

// CreateItem creates an item
func CreateItem(id string, variables []Variable) Item {
	return Item{
		ID:        id,
		Variables: variables,
	}
}

func (item *Item) String() string {
	b, _ := json.MarshalIndent(item, "", "\t")
	return string(b)
}

// attempt to decode the Variable values from base64
func (item *Item) decode() {
	// TODO add debug maybe?
	for i := range item.Variables {
		decodedValue, err := base64.StdEncoding.DecodeString(item.Variables[i].Value)
		if err == nil {
			item.Variables[i].Value = string(decodedValue)
		}
	}
}

func (item *Item) encode() {
	// TODO add debug maybe?
	for i := range item.Variables {
		encodedValue := base64.StdEncoding.EncodeToString([]byte(item.Variables[i].Value))
		item.Variables[i].Value = encodedValue
	}
}
