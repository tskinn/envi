# envi
A CLI for managing application configuration backed by DynamoDB

# Installation

With a working [Golang installation](https://golang.org/doc/install) run the following command:

``` shell
go get github.com/tskinn/envi
```

Also, make sure that `$GOPATH/bin` is in your `$PATH` so that your shell can find the go binaries.

# Configuration

Since `envi` is backed by a DynamoDB table, the DynamoDB table must
exist before using `envi`. By default `envi` will try to use a table named
"envi" located in the "us-east-1" region.

The only other thing required before using `envi` is valid credentials
and permissions to the DynamodDB table. You can provide an AWS Access
Key ID and Access Key Secret via environment variables but the most
convenient way is to have AWS configigured locally through the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

Environment variables `ENVI_TABLE` and `ENVI_REGION` can be used to
change the default DynamoDB table and AWS region or they can be set
with a flag.


# Usage
``` text
NAME:
   envi

USAGE:
   envi set --application myapp --environment dev --variables=one=eno,two=owt,three=eerht
   envi s -a myapp -e dev -v one=eno,two=owt,three=eerht
   envi s -i myapp__dev -f path/to/file/with/exported/vars
   envi g -a myapp -e dev -o json

VERSION:
   0.0.0

DESCRIPTION:
   A simple application configuration store cli backed by dynamodb

COMMANDS:
     set, s     save application configuraton in dynamodb
     get, g     get the application configuration for a particular application
     update, u  update an applications configuration by inserting new vars and updating old vars if specified
     delete, d  delete the application configuration for a particular application
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Commands

### get

The `get` command is self expalatory. Besides the default output
format (simple key=value) there are two other options: `json` and
`sh`. The `json` option prints an array of objects containg and name
and value (this is the format used in AWS ECS Task Definition
templates). The `sh` option is the same as the regulare output but
`export ` is added to the beginning of each line so as to more
convenient in using the env vars in a cli.

``` text
NAME:
   envi get - get the application configuration for a particular application

USAGE:
   envi get [command options] [arguments...]

OPTIONS:
   --output value, -o value       format of the output of the variables (default: "text")
   --table value, -t value        name of the dynamodb to store values (default: "envi") [$ENVI_TABLE]
   --region value, -r value       name of the aws region in which dynamodb table resides (default: "us-east-1") [$ENVI_REGION]
   --id value, -i value           id of the application environment combo; if id is not provided then application__environment is used as the id
   --application value, -a value  name of the application
   --environment value, -e value  name of the environment
```

### set

Use with CAUTION.

Use `set` to create new configurations. Set overrides all variables so
if one attempts to set a config with only one variable, all current
variables will be deleted and replaced with the single new variable.

If not creating a new config, it is better to use the `update` command.

``` text
NAME:
   envi set - save application configuraton in dynamodb

USAGE:
   envi set [command options] [arguments...]

OPTIONS:
   --variables value, -v value    env variables to store in the form of key=value,key2=value2,key3=value3
   --file value, -f value         path to a shell file that exports env vars
   --table value, -t value        name of the dynamodb to store values (default: "envi") [$ENVI_TABLE]
   --region value, -r value       name of the aws region in which dynamodb table resides (default: "us-east-1") [$ENVI_REGION]
   --id value, -i value           id of the application environment combo; if id is not provided then application__environment is used as the id
   --application value, -a value  name of the application
   --environment value, -e value  name of the environment
```

### update

The `update` command will override provided variables while leaving
current variables untouched.

Here is an example of how to update a single variable:

``` text
envi u -i omega__staging -v OLD_VAR=new-value
```

The above command will replace the OLD_VAR value with "new-value" and
leave all other unmentioned variables untouched.

``` text
NAME:
   envi update - update an applications configuration by inserting new vars and updating old vars if specified

USAGE:
   envi update [command options] [arguments...]

OPTIONS:
   --variables value, -v value    env variables to store in the form of key=value,key2=value2,key3=value3
   --file value, -f value         path to a shell file that exports env vars
   --table value, -t value        name of the dynamodb to store values (default: "envi") [$ENVI_TABLE]
   --region value, -r value       name of the aws region in which dynamodb table resides (default: "us-east-1") [$ENVI_REGION]
   --id value, -i value           id of the application environment combo; if id is not provided then application__environment is used as the id
   --application value, -a value  name of the application
   --environment value, -e value  name of the environment
```

### delete

The `delete` command will delete variables from a config if variables
are provided. If no variables are provided the entire config will be
deleted.

When deleting variables provide only the names of the variables to be
deleted. The same goes for files containing the env vars. The file
should only have the names of the vars to be deleted.

``` text
NAME:
   envi delete - delete the application configuration for a particular application

USAGE:
   envi delete [command options] [arguments...]

OPTIONS:
   --variables value, -v value    env variables to delete in the form of key=value,key2=value2,key3=value3
   --file value, -f value         path to a shell file that contains env vars
   --table value, -t value        name of the dynamodb to store values (default: "envi") [$ENVI_TABLE]
   --region value, -r value       name of the aws region in which dynamodb table resides (default: "us-east-1") [$ENVI_REGION]
   --id value, -i value           id of the application environment combo; if id is not provided then application__environment is used as the id
   --application value, -a value  name of the application
   --environment value, -e value  name of the environment
```

## Testing

There is a script to run the go tests and to test the basic
functionality of envi called `run_tests.sh`. This script
assumes privileged AWS credentials are setup in the environment. From
the root of the project simply run:

``` text
bash run_tests.sh
```

To just run go tests:

``` text
cd store
go test
```

