package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

type Tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

const mapping = `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store":true,
					"fielddata":true
				}
			}
		}
	}
}`

// Print error generated from function calls
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// Function to create client for ES server
func connectClient(ctx *context.Context, ip string) *elastic.Client {

	// Set up client and connect to ES defualt server
	client, err := elastic.NewClient()
	checkError(err)

	// Ping ES server
	info, code, err := client.Ping(ip).Do(*ctx)
	checkError(err)
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	return client
}

// Function to check if index exists in server
func checkIndexExists(ctx *context.Context, client *elastic.Client, name string) bool {
	exists, err := client.IndexExists(name).Do(*ctx)
	checkError(err)
	return exists
}

// Function to create an index in server
func createIndex(ctx *context.Context, client *elastic.Client, mapping string, name string) {
	index, err := client.CreateIndex(name).BodyString(mapping).Do(*ctx)
	checkError(err)
	if !index.Acknowledged {
		// Not Acknowledged
	} else {
		fmt.Printf("%s Index is created.\n", name)
	}
}

// Function to delete an index in server
func deleteIndex(ctx *context.Context, client *elastic.Client, name string) {
	delete, err := client.DeleteIndex(name).Do(*ctx)
	checkError(err)
	if !delete.Acknowledged {
		// Not acknowledged
	} else {
		fmt.Printf("%s Index is deleted.\n", name)
	}
}

func main() {
	// Set up ctx for each service
	ctx := context.Background()

	// Const
	const IP string = "http://127.0.0.1:9200"

	// Set up client and connect to default ES server
	client := connectClient(&ctx, IP)

	flag := 1

	for flag < 1000 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("What would you like to do?\n1. Check if index exists (1)\n2. Create an index (2)\n3. Delete an index (3)\nPress x to exit\n")
		option, _ := reader.ReadString('\n')
		optionSplice := strings.Split(option, "\n")
		if optionSplice[0] == "x" {
			break
			flag += 1000
		}
		fmt.Printf("Type in index name:\n")
		index, _ := reader.ReadString('\n')
		indexSplice := strings.Split(index, "\n")
		if optionSplice[0] == "1" {
			exists := checkIndexExists(&ctx, client, indexSplice[0])
			if exists {
				fmt.Printf("%s already exists.\n", indexSplice[0])
			} else {
				fmt.Printf("%s does not exists.\n", indexSplice[0])
			}
		} else if optionSplice[0] == "2" {
			exists := checkIndexExists(&ctx, client, indexSplice[0])
			if exists {
				fmt.Printf("Error: %s already exists.\n", indexSplice[0])
			} else {
				createIndex(&ctx, client, mapping, indexSplice[0])
			}
		} else {
			exists := checkIndexExists(&ctx, client, indexSplice[0])
			if exists {
				deleteIndex(&ctx, client, indexSplice[0])
			} else {
				fmt.Printf("Error: %s does not exists.\n", indexSplice[0])
			}
		}
	}
}
