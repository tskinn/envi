package main

import (
	"fmt"
	"os"

	"github.com/tskinn/envi/store"
	"github.com/urfave/cli"
)

func main() {
	var tableName, awsRegion, application, environment, id, variables, filePath, output string
	app := cli.NewApp()

	app.Description = "A simple application configuration store cli backed by dynamodb"
	app.Name = "envi"
	app.Usage = ""
	app.UsageText = `envi set --application myapp --environment dev --variables=one=eno,two=owt,three=eerht
   envi s -a myapp -e dev -v one=eno,two=owt,three=eerht
   envi s --id myapp__dev -f path/to/file/with/exported/vars
   envi get --application myapp --environment dev
   envi g -a myapp -e dev -o json`

	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "table, t",
			Value:       "envi",
			Usage:       "name of the dynamodb to store values",
			EnvVar:      "ENVI_TABLE",
			Destination: &tableName,
		},
		cli.StringFlag{
			Name:        "region, r",
			Value:       "us-east-1",
			Usage:       "name of the aws region in which dynamodb table resides",
			EnvVar:      "ENVI_REGION",
			Destination: &awsRegion,
		},
		cli.StringFlag{
			Name:        "id, i",
			Value:       "",
			Usage:       "id of the application environment combo; if id is not provided then application__environment is used as the id",
			Destination: &id,
		},
		cli.StringFlag{
			Name:        "application, a",
			Value:       "",
			Usage:       "name of the application",
			Destination: &application,
		},
		cli.StringFlag{
			Name:        "environment, e",
			Value:       "",
			Usage:       "name of the environment",
			Destination: &environment,
		},
	}

	setCommand := cli.Command{
		Name:    "set",
		Aliases: []string{"s"},
		Usage:   "save application configuraton in dynamodb",
		Action: func(c *cli.Context) error {
			var err error
			id, err = getID(id, application, environment)
			if err != nil {
				return err
			}
			store.Init(awsRegion, tableName)
			if filePath != "" {
				return store.SaveFromFile(id, filePath)
			} else if variables != "" {
				return store.Save(id, variables)
			}
			return fmt.Errorf("must provide variables or a path to a file containing variables")
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "variables, v",
				Value:       "",
				Usage:       "env variables to store in the form of key=value,key2=value2,key3=value3",
				Destination: &variables,
			},
			cli.StringFlag{
				Name:        "file, f",
				Value:       "",
				Usage:       "path to a shell file that exports env vars",
				Destination: &filePath,
			},
		},
	}
	setCommand.Flags = append(setCommand.Flags, globalFlags...)

	updateCommand := cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "update an applications configuration by inserting new vars and updating old vars if specified",
		Action: func(c *cli.Context) error {
			var err error
			id, err = getID(id, application, environment)
			if err != nil {
				return err
			}
			store.Init(awsRegion, tableName)
			if filePath != "" {
				return store.UpdateFromFile(id, filePath)
			} else if variables != "" {
				return store.Update(id, variables)
			}
			return fmt.Errorf("must provide variables or a path to a file containing variables")
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "variables, v",
				Value:       "",
				Usage:       "env variables to store in the form of key=value,key2=value2,key3=value3",
				Destination: &variables,
			},
			cli.StringFlag{
				Name:        "file, f",
				Value:       "",
				Usage:       "path to a shell file that exports env vars",
				Destination: &filePath,
			},
		},
	}
	updateCommand.Flags = append(updateCommand.Flags, globalFlags...)

	getCommand := cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get the application configuration for a particular application",
		Action: func(c *cli.Context) error {
			var item store.Item
			var err error
			id, err = getID(id, application, environment)
			if err != nil {
				return err
			}
			store.Init(awsRegion, tableName)
			item, err = store.Get(id)
			if err != nil {
				return err
			}
			item.PrintVars(output)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "output, o",
				Value:       "text",
				Usage:       "format of the output of the variables",
				Destination: &output,
			},
		},
	}
	getCommand.Flags = append(getCommand.Flags, globalFlags...)

	deleteCommand := cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete the application configuration for a particular application",
		Action: func(c *cli.Context) error {
			var err error
			id, err = getID(id, application, environment)
			if err != nil {
				return err
			}
			store.Init(awsRegion, tableName)
			if filePath != "" {
				return store.DeleteVarsFromFile(id, filePath)
			} else if variables != "" {
				return store.DeleteVars(id, variables)
			} else {
				return store.Delete(id)
			}
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "variables, v",
				Value:       "",
				Usage:       "env variables to delete in the form of key=value,key2=value2,key3=value3",
				Destination: &variables,
			},
			cli.StringFlag{
				Name:        "file, f",
				Value:       "",
				Usage:       "path to a shell file that contains env vars",
				Destination: &filePath,
			},
		},
	}
	deleteCommand.Flags = append(deleteCommand.Flags, globalFlags...)

	app.Commands = []cli.Command{
		setCommand,
		getCommand,
		updateCommand,
		deleteCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func getID(id, application, environment string) (string, error) {
	if id == "" && (application == "" || environment == "") {
		return "", fmt.Errorf("must provide an id or the application name and environment")
	}

	if id != "" {
		return id, nil
	}
	return application + "__" + environment, nil
}
