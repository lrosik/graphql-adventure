package main

import (
	"fmt"
	"graphql-adventure/schema"
)

func main() {
	gqlSchema := schema.PrepareSchema()

	query := `
		{
			posts {
				id
				title
				content
			}
		}
	`

	fmt.Printf("%s \n", schema.ParseQueryToJson(gqlSchema, query))
}
