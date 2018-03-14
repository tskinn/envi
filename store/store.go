package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	region    string
	tableName string
	db        *dynamodb.DynamoDB
)

// DynamodbItem is not what we want?
type DynamodbItem struct {
	Key   string `dynamodbav:"key"`
	Value string `dynamodbav:"value"`
}

// Init inits some mstuff
func Init(regionName, table string) {
	region = regionName
	tableName = table
	sesh := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	db = dynamodb.New(sesh)
}

// Get gets the thing ok
func Get(id string) (Item, error) {
	return get(id)
}

func get(id string) (Item, error) {
	var item Item
	var err error
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	resp, err := db.GetItem(params)
	if err != nil {
		return item, err
	}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &item)
	if err != nil {
		return item, err
	}
	item.decode()
	return item, err
}

// Save saves env vars given a string of vars in form of this=that,this2=that2
func Save(id, vars string) error {
	variables := parseVariables(vars, false)
	item := CreateItem(id, variables)
	return save(item)
}

// SaveFromFile gets env vars from a env file and saves to dynamo
func SaveFromFile(id, fileName string) error {
	variables, err := parseVariablesFromFile(fileName, false)
	if err != nil {
		return err
	}
	item := CreateItem(id, variables)
	return save(item)
}

func save(item Item) error {
	item.encode()
	atr, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      atr,
	}
	_, err = db.PutItem(params)
	return err
}

// Update updates configurate of given application with id
func Update(id, vars string) error {
	parsedVars := parseVariables(vars, false)
	return update(id, parsedVars)
}

// UpdateFromFile updates stuff from a file
func UpdateFromFile(id, fileName string) error {
	vars, err := parseVariablesFromFile(fileName, false)
	if err != nil {
		return err
	}
	return update(id, vars)
}

func update(id string, vars []Variable) error {
	item, err := get(id)
	if err != nil {
		if err.Error() != dynamodb.ErrCodeResourceNotFoundException {
			return err
		}
		// Save the item if it doesn't exist already
		return save(Item{
			ID:        id,
			Variables: vars,
		})
	}

	for i := 0; i < len(vars); i++ {
		found := false
		for j := 0; j < len(item.Variables); j++ {
			if vars[i].Name == item.Variables[j].Name {
				found = true
				item.Variables[j].Value = vars[i].Value
				break
			}
		}
		if !found { // add variable if not found already
			item.Variables = append(item.Variables, vars[i])
		}
	}
	return save(item)
}

// Delete deletes the thing
func Delete(id string) error {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := db.DeleteItem(params)
	return err
}

// DeleteVars deletes the given variables from the item with id of id
func DeleteVars(id, variables string) error {
	vars := parseVariables(variables, true)
	return deleteVars(id, vars)
}

// DeleteVarsFromFile deletes the variables found in the file filepath
func DeleteVarsFromFile(id, filePath string) error {
	vars, err := parseVariablesFromFile(filePath, true)
	if err != nil {
		return err
	}
	return deleteVars(id, vars)
}

func deleteVars(id string, vars []Variable) error {
	item, err := get(id)
	if err != nil {
		return err
	}
	for i := 0; i < len(vars); i++ {
		for j := 0; j < len(item.Variables); j++ {
			if vars[i].Name == item.Variables[j].Name {
				item.Variables = append(item.Variables[:j], item.Variables[j+1:]...)
			}
		}
	}
	return save(item)

}
