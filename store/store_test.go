package store

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	testItemOne = Item{
		ID: "app__one",
		Variables: []Variable{
			{
				Name:  "one",
				Value: "two",
			},
			{
				Name:  "three",
				Value: "four",
			},
			{
				Name:  "five",
				Value: "six",
			},
		},
	}
	testItemOneEncoded = Item{
		ID: "app__one",
		Variables: []Variable{
			{
				Name:  "one",
				Value: "dHdv",
			},
			{
				Name:  "three",
				Value: "Zm91cg==",
			},
			{
				Name:  "five",
				Value: "c2l4",
			},
		},
	}

	// doesn't handle comments in middle or end of lines
	testFileContent = `  one=two

	three=four 
#^ tab at front

five=six   
#       ^ spaces at end
`

	// should error out if no '=' sign found on line of text
	testFileContentBadFormat = `two
	three=four 
#^ tab at front

five=six   
#       ^ spaces at end
`

	testFileContentJustNames = `one
three
#  
five`

	testRawVariables = "one=two,three=four,five=six"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	items map[string]map[string]*dynamodb.AttributeValue
}

// TODO maybe use dynamodbattribute to make this easier to read at a
// glance for those who don't speak dynamodb
func (m mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	output := &dynamodb.GetItemOutput{}
	item, exists := m.items[*input.Key["id"].S]
	if !exists {
		return nil, fmt.Errorf(dynamodb.ErrCodeResourceNotFoundException)
	}
	output.Item = item
	return output, nil
}

func (m mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.items[*input.Item["id"].S] = input.Item
	return &dynamodb.PutItemOutput{}, nil
}

func (m mockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	delete(m.items, *input.Key["id"].S)
	return &dynamodb.DeleteItemOutput{}, nil
}

func TestParseVariables(t *testing.T) {
	variables := parseVariables(testRawVariables, false)
	if len(variables) != len(testItemOne.Variables) {
		t.Fatalf("length of variables not expected value")
	}
	if !variablesEqual(variables, testItemOne.Variables) {
		t.Fatalf("variables don't match expected")
	}
}

func TestParseVariablesFromFileSuccess(t *testing.T) {
	buffer := bytes.NewBufferString(testFileContent)
	scanner := bufio.NewScanner(buffer)
	variables, err := parseVariablesFromScanner(scanner, false)
	if err != nil {
		t.Fatalf("error parsing from scanner: %s", err)
	}
	if len(variables) != 3 {
		t.Fatalf("length of variables should be three\n%v", variables)
	}
	if !variablesEqual(variables, testItemOne.Variables) {
		t.Fatalf("variables do not match expected")
	}
}

func TestParseVariablesFromFileFailure(t *testing.T) {
	buffer := bytes.NewBufferString(testFileContentBadFormat)
	scanner := bufio.NewScanner(buffer)
	_, err := parseVariablesFromScanner(scanner, false)
	if err == nil {
		t.Fatalf("expected error on parsing bad format")
	}
}

func TestParseVariablesFromFileJustNames(t *testing.T) {
	buffer := bytes.NewBufferString(testFileContentJustNames)
	scanner := bufio.NewScanner(buffer)
	variables, err := parseVariablesFromScanner(scanner, true)
	if err != nil {
		t.Fatalf("error parsing from scanner: %s", err)
	}
	if len(variables) != 3 {
		t.Fatalf("length of variables expected to be three\n%v", variables)
	}
	for _, variable := range variables {
		if variable.Value != "" {
			t.Fatalf("expected values of variables to be empty strings")
		}
	}
}

func TestEncoding(t *testing.T) {
	testItem := testItemOne
	testItem.encode()
	if !variablesEqual(testItem.Variables, testItemOneEncoded.Variables) {
		t.Fatalf("encoding failed to produce the same variable encodings\n%v", testItem.Variables)
	}
	testItem.decode()
	if !variablesEqual(testItem.Variables, testItemOne.Variables) {
		t.Fatalf("decoding failed to produce the same variable values")
	}
}

func TestCreateItem(t *testing.T) {
	newItem := CreateItem("app_one", parseVariables(testRawVariables, false))
	if newItem.ID != "app_one" {
		t.Fatalf("Wow thats embarrassing")
	}
	if !variablesEqual(newItem.Variables, testItemOne.Variables) {
		t.Fatalf("created items variables don't match expected")
	}
}

func variablesEqual(one, two []Variable) bool {
	if len(one) != len(two) {
		return false
	}
	for i := range one {
		if one[i].Name != two[i].Name || one[i].Value != two[i].Value {
			return false
		}
	}
	return true
}

func variableExists(vars []Variable, name string) bool {
	for _, variable := range vars {
		if variable.Name == name {
			return true
		}
	}
	return false
}

func TestBasicSaveAndGetAndDelete(t *testing.T) {
	mock := mockDynamoDBClient{items: map[string]map[string]*dynamodb.AttributeValue{}}
	SetDB(mock)
	err := Save("app__test", "one=two,three=four,five=six")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	item, err := Get("app__test")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	if !variablesEqual(item.Variables, testItemOne.Variables) {
		t.Fatalf("variables don't match expected %v", item)
	}
	err = Delete("app__test")
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = Get("app__test")
	if err == nil {
		t.Fatalf("error expected to be non-nil")
	}
}

func TestDeleteVars(t *testing.T) {
	mock := mockDynamoDBClient{items: map[string]map[string]*dynamodb.AttributeValue{}}
	SetDB(mock)
	err := Save("app__test", "one=two,three=four,five=six")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	variableNameToDelete := "one"
	err = DeleteVars("app__test", variableNameToDelete)
	if err != nil {
		t.Fatalf("%s", err)
	}
	item, err := Get("app__test")
	if err != nil {
		t.Fatalf("error getting item %s", err)
	}
	if variableExists(item.Variables, variableNameToDelete) {
		t.Fatalf("expected variable to be deleted %s", variableNameToDelete)
	}
}

func TestDeleteMultipleVars(t *testing.T) {
	mock := mockDynamoDBClient{items: map[string]map[string]*dynamodb.AttributeValue{}}
	SetDB(mock)
	err := Save("app__test", "one=two,three=four,five=six")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	variableNamesToDelete := []string{"one", "three"}
	err = DeleteVars("app__test", strings.Join(variableNamesToDelete, ","))
	if err != nil {
		t.Fatalf("%s", err)
	}
	item, err := Get("app__test")
	if err != nil {
		t.Fatalf("error getting item %s", err)
	}
	for _, variableName := range variableNamesToDelete {
		if variableExists(item.Variables, variableName) {
			t.Fatalf("expected variable to be deleted %s", variableName)
		}
	}
}

func TestUpdateVariable(t *testing.T) {
	mock := mockDynamoDBClient{items: map[string]map[string]*dynamodb.AttributeValue{}}
	SetDB(mock)
	err := Save("app__test", "one=two,three=four,five=six")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	err = Update("app__test", "one=ten")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	item, err := Get("app__test")
	if err != nil {
		t.Fatalf("error getting item %s", err)
	}
	for _, variable := range item.Variables {
		if variable.Name == "one" && variable.Value != "ten" {
			t.Fatalf("expected var 'one' to equal 'ten'")
		}
	}
}

func TestUpdateAddVariable(t *testing.T) {
	mock := mockDynamoDBClient{items: map[string]map[string]*dynamodb.AttributeValue{}}
	SetDB(mock)
	err := Save("app__test", "one=two,three=four,five=six")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	err = Update("app__test", "seven=eight")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	item, err := Get("app__test")
	if err != nil {
		t.Fatalf("error getting item %s", err)
	}
	found := false
	for _, variable := range item.Variables {
		if variable.Name == "seven" {
			found = true
			if variable.Value != "eight" {
				t.Fatalf("expected var 'seven' to equal 'eight'")
			}
		}
	}
	if !found {
		t.Fatalf("expected to find var 'seven")
	}
}

func TestUpdateAddItem(t *testing.T) {
	mock := mockDynamoDBClient{items: map[string]map[string]*dynamodb.AttributeValue{}}
	SetDB(mock)
	err := Update("app__testNew", "seven=eight")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	_, err = Get("app__testNew")
	if err != nil {
		t.Fatalf("expected item to exist %s", err)
	}
}

func TestInit(t *testing.T) {
	db = nil
	Init("us-east-1", "envi")
	if db == nil {
		t.Fatalf("expected db to initialized")
	}
}
