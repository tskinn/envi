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
   envi get --application myapp --environment dev
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

