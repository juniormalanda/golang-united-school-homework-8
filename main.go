package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/juniormalanda/golang-united-school-homework-8/file"
)

type Arguments map[string]string

func parseArgs() (args Arguments) {
	argsList := []struct {
		name         string
		defaultValue string
		usage        string
	}{
		{
			name:         "operation",
			defaultValue: "list",
			usage:        "Needed operation type",
		},
		{
			name:         "item",
			defaultValue: "{}",
			usage:        "Body of item in json",
		},
		{
			name:         "id",
			defaultValue: "0",
			usage:        "Item identifier",
		},
		{
			name:         "fileName",
			defaultValue: ".",
			usage:        "Name of file",
		},
	}

	for _, arg := range argsList {
		args[arg.name] = *flag.String(arg.name, arg.defaultValue, arg.usage)
	}

	return
}

func Perform(args Arguments, writer io.Writer) error {
	fileName, ok := args["fileName"]

	if !ok || fileName == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	f := file.NewFile(fileName)

	operation, ok := args["operation"]

	if !ok || operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	switch operation {
	case "add":
		item, ok := args["item"]
		if !ok || item == "" {
			return fmt.Errorf("-item flag has to be specified")
		}

		user, err := f.AddUser(item)

		if errors.Is(err, file.ItemExistsError) {
			writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", user.Id)))

			return nil
		}

		if err != nil {
			return err
		}

		writer.Write([]byte("Item successfully added"))
	case "list":
		data, err := f.List()
		if err != nil {
			return err
		}
		writer.Write(data)
	case "findById":
		id, ok := args["id"]
		if !ok || id == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		data, err := f.FindById(id)

		if err != nil {
			return err
		}

		writer.Write(data)
	case "remove":
		id, ok := args["id"]
		if !ok || id == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		err := f.Remove(id)

		if errors.Is(err, file.NotFoundError) {
			writer.Write([]byte(fmt.Sprintf("Item with id %s not found", id)))

			return nil
		}

		if err != nil {
			return err
		}

		writer.Write([]byte("Item successfully removed"))
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
