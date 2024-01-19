# CIVAR

This is a simple commandline tool to show, create and update Gitlab CI/CD variables of your projects by using the the Gitlab API.

## Installation

### Download from Release Page
coming soon...

### Docker Image
coming soon...

## Configuration
You need to create a Gitlab token [here](https://gitlab.com/-/profile/personal_access_tokens) to work with civar.

### Basic command
You can set values as options in the command:
```shell
civar get --token $GITLAB_TOKEN --url https://gitlab.com --dotenv apps/project1
```

### YAML config
If you don't want to set options every time, create a config file in `$HOME/.civar.yaml`
```yaml
url: https://gitlab.com
token: [ gitlab token ]
format: [ table | dotenv | json ]
```
```shell
civar get apps/project1
```


## Usage
### Get all variables of a gitlab Project
#### Table format:
```shell
$ civar get -p apps/project1
|   SCOPE    |  KEY  |      VALUE       |  TYPE   | MASKED | PROTECTED |
|------------|-------|------------------|---------|--------|-----------|
| *          | VAR_1 | VALUE_1          | env_var | false  | false     |
| staging    | VAR_2 | VALUE_STAGING    | env_var | false  | false     |
| production | VAR_2 | VALUE_PRODUCTION | env_var | false  | false     |
```
#### Dotenv format
> :information: The K8S_SECRET_ prefix will be ignored by default
```shell
$ civar get -d apps/project1
# Scope: *
VAR_1="VALUE_1"

# Scope: staging
VAR_2="VALUE_STAGING"

# Scope: production
VAR_2="VALUE_PRODUCTION"
```
---


### Create variables from a .env file
```shell
$ cat .env
# Scope: *
VAR_1="VALUE_1"

# Scope: staging
VAR_2="VALUE_STAGING"

# Scope: production
VAR_2="VALUE_PRODUCTION"
```

```shell
# create the variable as defined in the file under the right scope
$ cat .env | civar create -d apps/project1

# if you want to add the K8S_SECRET_ prefix to be added to the variables use the -k option
$ cat .env | civar create -d -k apps/project1
```


### Copy vars from one project to another as a oneliner
```shell
$ civar get apps/project1 | civar create apps/project2
```

### Help Pages
#### General
```shell
$ civar help
Usage:
civar [command] [flags] [projectId | group/projectName]

Available Commands:
  get         Prints out all CI/CD variables for the given Gitlab project
  create      Creates CI/CD variables for a given Gitlab project


Global Flags:
      --config string   config file (default is $HOME/.civar.yaml)
  -h, --help            help page
  -u, --url             sets the url for gitlab
  -t, --token           sets a token for gitlab 

Use "civar [command]" for more information about a command.

Examples:

- Print all CI/CD variables in a table
	civar get -p 1

- Copy all CI/CD variables from one project to another:
	civar get 1 | civar create 2

- Save output to a file:
	civar get 1 > vars.txt

- Create CI/CD variables from a file:
	cat vars.txt | civar create 2
```

#### Get
```shell
$ civar get
Usage:
civar get [flags] [projectId | group/projectName]

Local Flags:
  -p, --pretty          output variables as a table 
  -d, --dotenv          output variables in dotenv format

Global Flags:
      --config string   config file (default is $HOME/.civar.yaml)
  -h, --help            help page
  -u, --url             sets the url for gitlab
  -t, --token           sets a token for gitlab
```

#### Create
```shell
$ civar create
Usage:
cat data.txt | civar create [flags] [projectId | group/projectName]

Local Flags:
  -d, --dotenv          parses variables in dotenv format
  -k, --k8s             creates variables with K8S_SECRET_ prefix
  -f --file [filename]  read input from a file

Global Flags:
      --config string   config file (default is $HOME/.civar.yaml)
  -h, --help            help for create
  -u, --url             sets the url for gitlab
  -t, --token           sets a token for gitlab
```

#### Update
```shell
$ civar help update
Reads data from stdin or file and updates already existing variables in a Gitlab project. Non existent variables will be skipped.

Usage:
  civar update group/project [flags]

Examples:
cat .env | civar update group/project -d

Flags:
  -d, --dotenv        output as dotenv
  -f, --file string   reads input from a file
  -h, --help          help for update
  -k, --k8s           update variables with K8S_SECRET_ prefix

Global Flags:
      --config string   config file (default is $HOME/.civar.7yml)
  -t, --token string    sets your token
  -u, --url string      sets your gitlab url

```

## Build from source
### with go on your system
```shell
$ bin/compile.sh
```

### via docker
```shell
$ bin/compile-via-docker.sh
```

## Temporary installation 
you can build and try it using:
```shell
$ go install
```

