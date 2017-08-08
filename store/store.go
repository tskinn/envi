package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

// DynamodbItem is not what we want?
type DynamodbItem struct {
	Key   string `dynamodbav:"key"`
	Value string `dynamodbav:"value"`
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

func Get(db *dynamodb.DynamoDB, cluster, key string) string {
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
