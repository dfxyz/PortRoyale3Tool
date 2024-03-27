package main

import (
	"PortRoyale3Tool/lib"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var dataFilePath *string = flag.String("data", "PortRoyale3Tool.json", "The path to the data file")

var data *lib.Data

func main() {
	flag.Parse()

	data = lib.NewData()
	data.Load(*dataFilePath)
	defer data.Save(*dataFilePath)

	reader := bufio.NewReader(os.Stdin)
	for {
		if !handleMainMenu(reader) {
			return
		}
	}
}

const helpInfo = `Commands:
* list: list all groups and cities
* list <cityName>: list the details of the given city
* list <groupIndex>: list the details of the city(s) in the given group
* set <cityName> <goodName> <buildingNum>: set the production building number of the given good in the given city
* set <groupIndex> <cityName> ...: associate the given city(s) to the given group
* unset <cityName> <goodName>: remove the production building of the given good in the given city
* unset <groupIndex> <cityName> ...: remove the association of the given city(s) from the given group
* remove <cityName>: remove the city by the given name
* remove <groupIndex>: remove the group by the given index
* save: save the data immediately
* help: print the help info about the commands
* exit: save and exit`

func handleMainMenu(reader *bufio.Reader) (continueLoop bool) {
	continueLoop = true
	fmt.Print(">> ")
	cmd, err := reader.ReadString('\n')
	if errors.Is(err, io.EOF) {
		return false
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to read the command: %v", err))
	}
	args := strings.Split(strings.TrimSpace(cmd), " ")
	switch args[0] {
	case "list":
		args = args[1:]
		if len(args) <= 0 {
			data.ListAll()
			return
		}
		if len(args) != 1 {
			fmt.Println("Invalid command; Available patterns: 'list', 'list <cityName>' or 'list <groupIndex>'")
			return
		}
		strOrInt := args[0]
		num, err := strconv.Atoi(strOrInt)
		if err != nil {
			cityName := strOrInt
			data.ListCity(cityName)
			return
		}
		groupIndex := num
		data.ListGroup(groupIndex)

	case "set":
		args = args[1:]
		if len(args) <= 0 {
			fmt.Println("Invalid command; Available patterns: 'set <cityName> <goodName> <buildingNum>' or 'set <groupIndex> <cityName> ...'")
			return
		}
		strOrInt := args[0]
		args = args[1:]
		num, err := strconv.Atoi(strOrInt)
		if err != nil {
			cityName := strOrInt
			if len(args) != 2 {
				fmt.Println("Invalid command; Available pattern: 'set <cityName> <goodName> <buildingNum>'")
				return
			}
			goodName := args[0]
			buildingNumStr := args[1]

			good, ok := lib.GoodFromStr(args[0])
			if !ok {
				fmt.Printf("Invalid good '%s'\n", goodName)
				return
			}
			buildingNum, err := strconv.Atoi(buildingNumStr)
			if err != nil {
				fmt.Printf("Invalid building number '%s'\n", buildingNumStr)
				return
			}
			if buildingNum <= 0 {
				fmt.Println("Building number must be positive")
				return
			}

			data.SetProduceBuilding(cityName, good, buildingNum)
			return
		}
		groupIndex := num
		if len(args) <= 0 {
			fmt.Println("Invalid command; Available pattern: 'set <groupIndex> <cityName> ...'")
			return
		}
		data.GroupAssociate(groupIndex, args)

	case "unset":
		args = args[1:]
		if len(args) <= 0 {
			fmt.Println("Invalid command; Available patterns: 'unset <cityName> <goodName>' or 'unset <groupIndex> <cityName> ...'")
			return
		}
		strOrInt := args[0]
		args = args[1:]
		num, err := strconv.Atoi(strOrInt)
		if err != nil {
			cityName := strOrInt
			if len(args) != 1 {
				fmt.Println("Invalid command; Available pattern: 'unset <cityName> <goodName>'")
				return
			}
			goodName := args[0]
			good, ok := lib.GoodFromStr(args[0])
			if !ok {
				fmt.Printf("Invalid good '%s'\n", goodName)
				return
			}
			data.UnsetProduceBuilding(cityName, good)
			return
		}
		groupIndex := num
		data.UnassociateGroup(groupIndex, args)

	case "remove":
		args = args[1:]
		if len(args) <= 0 {
			fmt.Println("Invalid command; Available patterns: 'remove <cityName>' or 'remove <groupIndex>'")
			return
		}
		strOrInt := args[0]
		num, err := strconv.Atoi(strOrInt)
		if err != nil {
			cityName := strOrInt
			data.RemoveCity(cityName)
			return
		}
		groupIndex := num
		data.RemoveGroup(groupIndex)

	case "save":
		data.Save(*dataFilePath)

	case "help":
		fmt.Println(helpInfo)

	case "exit":
		return false

	default:
		fmt.Println("Unknown command; Input 'help' to see the available commands")
	}
	return
}
