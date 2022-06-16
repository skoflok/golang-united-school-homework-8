package main

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

var allowedOperations = []string{"add", "list", "findById", "remove", "info"}

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

	if err = prepareOperationArg(args["operation"]); err != nil {
		return err
	}

	if err = prepareFileNameArg(args["fileName"]); err != nil {
		return err
	}

	f, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("OpenFile error: %w", err)
	}
	defer f.Close()

	filesize, err := fsize(f)

	if err != nil {
		return fmt.Errorf("FileStat Error: %w", err)
	}

	if filesize == 0 {
		f.WriteAt([]byte("[]"), 0)
	}

	switch args["operation"] {
	case "list":
		payload, err = list(f)
	case "add":
		payload, err = add(f, args["item"])
	case "findById":
		payload, err = findById(f, args["id"])
	case "remove":
	default:
		return fmt.Errorf("Operation " + args["operation"] + " not allowed!")
	}

	writer.Write([]byte(payload))
	return err
}

func fsize(file *os.File) (size int64, err error) {
	stat, err := file.Stat()
	if err != nil {
		return size, err
	}

	size = stat.Size()
	return size, err
}

func findById(file *os.File, id string) (payload string, err error) {
	if err = prepareIdArg(id); err != nil {
		return payload, err
	}

	elements, err := readElementsFromFile(file)
	if err != nil {
		return payload, err
	}

	element, ok := getElementsById(elements, id)
	if !ok {
		return "Item with id " + id + " not found", err
	}

	b, err := json.Marshal(element)
	if err != nil {
		return payload, err
	}
	return string(b), err
}

func getElementsById(elements []Element, id string) (element Element, ok bool) {
	for _, e := range elements {
		if e.Id == id {
			return e, true
		}
	}
	return element, false
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

	for _, e := range elements {
		if e.Id == element.Id {
			return "Item with id " + element.Id + " already exists", nil
		}
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

func parseArgs() (args Arguments) {
	flag.Parse()

	args = getDefaultArguments()

	args["operation"] = *operationArg
	args["fileName"] = *fileNameArg
	args["item"] = *itemArg
	args["id"] = *idArg
	return args
}

func prepareOperationArg(operation string) (err error) {
	if operation == "" {
		err = errors.New("-operation flag has to be specified")
	}
	return err
}

func prepareIdArg(id string) (err error) {
	if id == "" {
		err = errors.New("-id flag has to be specified")
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

	args := parseArgs()

	err := Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	}
}
