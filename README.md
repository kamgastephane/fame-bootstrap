# Fame bootstrap

Command line utility to bootstrap fame entities built with in GO with [viper](https://github.com/spf13/viper) and [cobra](https://github.com/spf13/cobra)

## Configuration file

#### Location
A configuration file is required in order to run this application.
The default path for the configuration file is `$HOME/.famebootstrap.yaml`
You can also pass an absolute path to your custom configuration file using the flag `config`

#### Content
A default configuration file can be found in this repo at `.famebootstrap.yaml`.
The baseUrl has to be setup for all three environments.

## How to use
Assuming you have access to the executable binary named `famebootstrap`, the list of available commands can be found:

### Help
You can get a list of all available commands
```
$ famebootstrap help 
$ famebootstrap bootstrap --help 

```  
### Bootstrap
You can bootstrap a entire fame entity or a specific resource belonging to a fame entity
```
$ famebootstrap bootstrap [entityName] [resourceId]
```
- The **entityName** argument has to be the name of a valid entity from FAME.
- If the **entityName** argument is not passed as input, you will be required to enter an entity name during the program execution.
- The **entityName** argument has to be specified in order to use the resourceId argument
- If both the **entityName** and the **resourceId** arguments are specified, only that specific **resourceId** will be bootstrapped


The following flags can be passed as input:
- **-q** : The destination queue for the messages you are bootstrapping. If not specified, you will be prompted to enter one during the execution
- **-e** : The environment to be used for the bootstrap. If not specified, you will be prompted to enter one during the execution
- **-u** : A custom base url can be supplied here and it will override the base url provided in the configuration file
```
 $ famebootstrap bootstrap
 $ famebootstrap bootstrap match
 $ famebootstrap bootstrap stadium
 $ famebootstrap bootstrap match -e=GOLD -q=fsp
```
 
### Status
You can get the status of a running bootstrap
The argument **commandId** can be retrieved from the output of the bootstrap command
```
$ famebootstrap status {commandId}
```

## How to build
```
$ go get
$ go fmt path_to_source_code\fame-bootstrap
```
- Windows: 
```
$ env GOOS=windows GOARCH=amd64 go build -o build/famebootstrap.exe
```
- Linux: 
```
$ env GOOS=linux GOARCH=amd64 go build -o build/famebootstrap
```
Once the executable binary has been created, you can run it directly from the command line as explained in [how to use](#How-to-use) section
## how to run

The program can be run directly from the code using:
```
$ go run main.go {name_of_the_command} [args] [flags]
```
e.g.
`
$ go run main.go bootstrap match -q=fsp -e=PROD
`