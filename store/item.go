package store

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Variable is a piece of configuration
type Variable struct {
	Name  string `dynamodbav:"name" json:"name"`
	Value string `dynamodbav:"value" json:"value"`
}

// Item is the format of the configuratoin stored in dynamodb
type Item struct {
	ID          string     `dynamodbav:"id" json:"id"`
	Application string     `dynamodbav:"application" json:"application"`
	Environment string     `dynamodbav:"environment" json:"environment"`
	Variables   []Variable `dynamodbav:"variables" json:"variables"`
}

// TODO this is pretty darn primitive so make it more robust
// support other formats and what not
func parseVariables(variables string) []Variable {
	split := strings.Split(variables, ",")
	vars := make([]Variable, len(split))
	for i := range split {
		j := strings.Index(split[i], "=")
		vars[i] = Variable{
			Name:  string(split[i][:j]),
			Value: base64.StdEncoding.EncodeToString([]byte(split[i][j+1:])),
		}
	}
	return vars
}

// CreateItem creates an item
func CreateItem(id, application, environment, variables string) Item {
	return Item{
		ID:          id,
		Application: application,
		Environment: environment,
		Variables:   parseVariables(variables),
	}
}

func (i *Item) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}

// attempt to decode the Variable values from base64
func (i *Item) decode() {
	for j := range i.Variables {
		tmp, err := base64.StdEncoding.DecodeString(i.Variables[j].Value)
		if err == nil {
			i.Variables[j].Value = string(tmp)
		}
	}
}
