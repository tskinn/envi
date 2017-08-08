package main

import (
	"fmt"
	"os"

	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tskinn/envi-cli/store"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Description = "A simple key-value store cli for dynamodb"
	app.Name = "envi"
	app.Usage = ""
	app.UsageText = "envi --put --key <key> --value <value>\n   envi --get --key <key>\n   envi --getall --key <key>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output",
			Value: "text",
			Usage: "json | text. only for getall option",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "us-east-1",
			Usage: "aws region in which the dynamodb table resides",
		},
		cli.StringFlag{
			Name:  "table",
			Value: "envi",
			Usage: "name of the dynamodb table used for envi",
		},
		cli.StringFlag{
			Name:  "key",
			Value: "",
			Usage: "key to store in dynamodb",
		},
		cli.StringFlag{
			Name:  "value",
			Value: "",
			Usage: "value of the value to be stored",
		},
		cli.BoolFlag{
			Name:  "get",
			Usage: "get the value of the given key",
		},
		cli.BoolFlag{
			Name:  "getall",
			Usage: "get all the keys with the prefix of key",
		},
		cli.BoolFlag{
			Name:  "put",
			Usage: "put the key value pair in the cluster",
		},
	}

	app.Action = func(c *cli.Context) error {
		var err error
		sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(c.String("region"))}))
		db := dynamodb.New(sess)
		if c.Bool("get") && c.String("key") != "" {
			// get the key!
			value := store.Get(db, c.String("table"), c.String("key"))
			if value != "" {
				print(c.String("output"), []string{c.String("key")}, []string{value})
			}
		}
		if c.Bool("getall") && c.String("key") != "" {
			// get all keys-values under key
			keys, values := store.GetAll(db, c.String("table"), c.String("key"))
			print(c.String("output"), keys, values)
		}
		if c.Bool("put") && c.String("value") != "" && c.String("key") != "" {
			// put the key-value!
			err = store.Put(db, c.String("table"), c.String("key"), c.String("value"))
			if err != nil {
				fmt.Println("failed to put key-value")
			} else {
				fmt.Printf("Successfully put key-value:\n\t%s:\t%s\n", c.String("key"), c.String("value"))
			}
		}
		return err
	}

	app.Run(os.Args)
}

func print(output string, keys, values []string) {
	switch output {
	case "json":
		printJSON(keys, values)
	case "text":
		printText(keys, values)
	default:
		printText(keys, values)
	}
}

func printJSON(keys, values []string) {
	list := make([]map[string]string, 0)
	for i := 0; i < len(keys); i++ {
		tmp := map[string]string{"key": keys[i], "value": values[i]}
		list = append(list, tmp)
	}
	raw, err := json.Marshal(list)
	if err != nil {
		// print something
		return
	}
	fmt.Printf("%s", raw)
}

func printText(keys, values []string) {
	for i := 0; i < len(keys); i++ {
		fmt.Printf("%s=%s\n", keys[i], values[i])
	}
}
