package schema

import (
	"encoding/json"
	"graphql-adventure/models"
	"log"

	"github.com/graphql-go/graphql"
)

var author = models.Author{ID: 1, FirstName: "Łukasz", LastName: "Rosik"}

var posts = map[int]models.Post{1: {ID: 1, Title: "Hello World", Author: &author, Content: "My first blog post."},
	2: {ID: 2, Title: "My second post", Author: &author, Content: "Welcome in my second post on this blog."}}

func PrepareSchema() graphql.Schema {
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
					post, postOk := posts[id]
					if postOk {
						return post, nil
					}
				}
				return nil, nil
			},
		},
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Get all posts",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var postsSlice []models.Post

				for _, value := range posts {
					postsSlice = append(postsSlice, value)
				}

				return postsSlice, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	return schema
}

func ParseQueryToJson(schema graphql.Schema, query string) []byte {
	params := graphql.Params{Schema: schema, RequestString: query}
	request := graphql.Do(params)

	if len(request.Errors) > 0 {
		log.Fatalf("Failed to execute graphQL query, errors: %+v", request.Errors)
	}

	requestJSON, _ := json.Marshal(request)

	return requestJSON
}
