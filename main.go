package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

var allowedOperations = []string{"add", "list", "findById", "remove"}

var id = flag.String("id", "", "Unique ID of the Item")
var operation = flag.String("operation", "", "One of operation:"+strings.Join(allowedOperations, ","))
var item = flag.String("item", "", "Item to add. Json string with id(string), email(string), age(integer) properties")
var fileName = flag.String("fileName", "", "Path to file(storage)")

type Arguments map[string]string

// func Perform(args Arguments, writer io.Writer) error {

// }

func getDefaultArguments() Arguments {
	args := Arguments{
		"id":        "",
		"operation": "",
		"item":      "",
		"fileName":  "",
	}
	return args
}

func parseArgs() (args Arguments, err error) {
	flag.Parse()

	args = getDefaultArguments()

	if err = prepareOperationArg(*operation); err != nil {
		return args, err
	}
	args["operation"] = *operation

	if err = prepareFileNameArg(*fileName); err != nil {
		return args, err
	}
	args["fileName"] = *fileName
	return args, err
}

func prepareOperationArg(operation string) (err error) {
	allowed := false

	if operation == "" {
		err = errors.New("operation flag has to be specified")
	}
	for _, v := range allowedOperations {
		if v == operation {
			allowed = true
		}
	}

	if !allowed {
		err = errors.New("Operation " + operation + " not allowed!")
	}

	return err
}

func prepareFileNameArg(fileName string) (err error) {
	if fileName == "" {
		err = errors.New("-fileName flag has to be specified")
	}
	return err
}

func main() {

	_, err := parseArgs()

	if err != nil {
		panic(err)
	}

	// err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Some")
	}
}
