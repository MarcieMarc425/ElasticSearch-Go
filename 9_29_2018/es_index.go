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

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("What would you like to do?\n1. Check if index exists (1)\n2. Create an index (2)\n3. Delete an index (3)\n")
	text, _ := reader.ReadString('\n')
	if text == "1\n" {
		fmt.Printf("Type in index name:\n")
		text, _ := reader.ReadString('\n')
		exists := checkIndexExists(&ctx, client, text)
		if exists {
			fmt.Printf("%s already exists.\n", text)
		} else {
			fmt.Printf("%s does not exists.\n", text)
		}
	} else if text == "2" {

	} else {

	}

	// // Check index in ES server, if not exists, create index
	// exists := checkIndexExists(&ctx, client, INDEX)
	// if exists {
	// 	fmt.Printf("%s already exists.\n", INDEX)
	// } else {
	// 	fmt.Printf("%s does not exists.\n", INDEX)
	// 	createIndex(&ctx, client, mapping, INDEX)
	// }

	// // Check index in ES server, if exists, delete index
	// exists := checkIndexExists(&ctx, client, INDEX)
	// if exists {
	// 	fmt.Printf("%s already exists.\n", INDEX)
	// 	deleteIndex(&ctx, client, INDEX)
	// } else {
	// 	fmt.Printf("%s does not exists.\n", INDEX)
	// }
}
