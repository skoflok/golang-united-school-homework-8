package io_os_context

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var allowedOperations = []string{"add", "list", "findById", "remove"}

var idArg = flag.String("id", "", "Unique ID of the Item")
var operationArg = flag.String("operation", "", "One of operation:"+strings.Join(allowedOperations, ","))
var itemArg = flag.String("item", "", "Item to add. Json string with id(string), email(string), age(integer) properties")
var fileNameArg = flag.String("fileName", "", "Path to file(storage)")

type Element struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) (err error) {
	var payload string
	f, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("OpenFile error: %w", err)
	}
	defer f.Close()

	switch args["operation"] {
	case "list":
		payload, err = list(f)
	case "add":
		payload, err = add(f, args["item"])
	case "findById":
	case "remove":
	default:
		return fmt.Errorf("Operation " + args["operation"] + " not allowed!")
	}

	fmt.Println(payload)
	return err
}

func list(file *os.File) (payload string, err error) {
	elements, err := readElementsFromFile(file)
	if err != nil {
		return payload, err
	}

	b, err := json.Marshal(elements)
	if err != nil {
		return payload, err
	}
	return string(b), err
}

func add(file *os.File, item string) (payload string, err error) {

	if err = prepareitemArg(item); err != nil {
		return "", err
	}

	element := Element{}
	if err = json.Unmarshal([]byte(item), &element); err != nil {
		return "", err
	}

	elements, err := readElementsFromFile(file)
	if err != nil {
		return payload, err
	}

	elements = append(elements, element)

	b, err := json.Marshal(elements)
	if err != nil {
		return payload, err
	}

	if _, err = file.WriteAt(b, 0); err != nil {
		return payload, err
	}
	return string(b), err
}

func readElementsFromFile(file *os.File) (elements []Element, err error) {

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)

	if err != nil {
		return elements, fmt.Errorf("Read from file to buffer error: %w", err)
	}

	content := buf.Bytes()
	if err = json.Unmarshal(content, &elements); err != nil {
		return elements, err
	}
	return elements, nil
}

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

	if err = prepareOperationArg(*operationArg); err != nil {
		return args, err
	}
	args["operation"] = *operationArg

	if err = prepareFileNameArg(*fileNameArg); err != nil {
		return args, err
	}

	args["fileName"] = *fileNameArg
	args["item"] = *itemArg
	args["id"] = *idArg

	fmt.Println(args)
	return args, err
}

func prepareOperationArg(operation string) (err error) {
	if operation == "" {
		err = errors.New("-operation flag has to be specified")
	}
	return err
}

func prepareFileNameArg(fileName string) (err error) {
	if fileName == "" {
		err = errors.New("-fileName flag has to be specified")
	}
	return err
}

func prepareitemArg(item string) (err error) {
	if item == "" {
		err = errors.New("-item flag has to be specified")
	}
	return err
}

func main() {

	args, err := parseArgs()

	if err != nil {
		panic(err)
	}

	err = Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(args)
	}
}
