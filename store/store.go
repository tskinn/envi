package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strings"
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
func Get(id string) (item Item, err error) {
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
		return
	}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &item)
	if err != nil {
		return
	}
	item.decode()
	return
}

func Put(db *dynamodb.DynamoDB, table, key, value string) error {
	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"key": &dynamodb.AttributeValue{
				S: aws.String(key),
			},
			"value": &dynamodb.AttributeValue{
				S: aws.String(value),
			},
		},
		TableName: aws.String(table),
	}
	_, err := db.PutItem(params)
	return err
}

func get(db *dynamodb.DynamoDB, cluster, key string) string {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(cluster),
	}
	resp, err := db.GetItem(params)
	if err != nil {
		return ""
	}
	if resp.Item["value"].S != nil {
		return *resp.Item["value"].S
	}
	return ""
}

func GetAll(db *dynamodb.DynamoDB, table, searchKey string) ([]string, []string) {
	items, err := scan(db, table)
	if err != nil {
		return make([]string, 0), make([]string, 0)
	}
	keys := make([]string, 0)
	values := make([]string, 0)
	for _, value := range items {
		k := *value["key"].S
		v := *value["value"].S
		if strings.HasPrefix(k, searchKey) {
			tKey := strings.TrimPrefix(k, searchKey)
			keys = append(keys, tKey)
			values = append(values, v)
		}
	}
	return keys, values
}

func scan(db *dynamodb.DynamoDB, table string) ([]map[string]*dynamodb.AttributeValue, error) {
	// TODO use query instead?
	params := &dynamodb.ScanInput{
		TableName: aws.String(table),
	}
	items := make([]map[string]*dynamodb.AttributeValue, 0)
	err := db.ScanPages(params,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			items = append(items, page.Items...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func Save(id, app, env, vars string) error {
	variables := parseVariables(vars)
	item := CreateItem(id, app, env, variables)
	return save(item)
}

func SaveFromFile(id, app, env, fileName string) error {
	variables, err := parseVariablesFromFile(fileName)
	if err != nil {
		return err
	}
	item := CreateItem(id, app, env, variables)
	return save(item)
}

func save(item Item) error {
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
