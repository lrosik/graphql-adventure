package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Post struct {
	ID      int
	Title   string
	Author  *Author
	Content string
}

type Author struct {
	ID        int
	FirstName string
	LastName  string
}

var authorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"firstName": &graphql.Field{Type: graphql.String},
		"lastName":  &graphql.Field{Type: graphql.String},
	},
})

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id":      &graphql.Field{Type: graphql.Int},
		"title":   &graphql.Field{Type: graphql.String},
		"author":  &graphql.Field{Type: authorType},
		"content": &graphql.Field{Type: graphql.String},
	},
})

var author = Author{ID: 1, FirstName: "Łukasz", LastName: "Rosik"}

var posts = []Post{Post{ID: 1, Title: "Hello World", Author: &author, Content: "My first blog post."}}

func main() {
	// Schema
	fields := graphql.Fields{
		"post": &graphql.Field{
			Type:        postType,
			Description: "Get post by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int)
				if ok {
					for _, post := range posts {
						if post.ID == id {
							return post, nil
						}
					}
				}
				return nil, nil
			},
		},
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Get all posts",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return posts, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	query := `
		{
			posts {
				id
			}
		}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	request := graphql.Do(params)

	if len(request.Errors) > 0 {
		log.Fatalf("Failed to execute graphQL query, errors: %+v", request.Errors)
	}

	requestJSON, _ := json.Marshal(request)
	fmt.Printf("%s \n", requestJSON)
}
