package backend

import (
	"context"
	"fmt"

	"appstore/constants"

	"github.com/olivere/elastic/v7"
)

var (
    ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
    client *elastic.Client
}

func InitElasticsearchBackend() {
//initializes an Elasticsearch client by connecting to the Elasticsearch server with the specified URL and 
//credentials. It then creates two indexes, constants.APP_INDEX and constants.USER_INDEX, 
    client, err := elastic.NewClient(
        elastic.SetURL(constants.ES_URL),
        elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
    if err != nil {
        panic(err)
    }

    exists, err := client.IndexExists(constants.APP_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "title":      { "type": "text"},
                    "description":  { "type": "text" },
                    "price":      { "type": "keyword", "index": false },
                    "url":     { "type": "keyword", "index": false }
                }
            }
        }`
        _, err := client.CreateIndex(constants.APP_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                     "mappings": {
                         "properties": {
                            "username": {"type": "keyword"},
                            "password": {"type": "keyword"},
                            "age": {"type": "long", "index": false},
                            "gender": {"type": "keyword", "index": false}
                         }
                    }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")
    //initialize ESbackend
    ESBackend = &ElasticsearchBackend{client: client}
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).//specifies the index to search on.
        Query(query).//specifies the query to execute on the specified index.
        Pretty(true).//formats the search results in a readable format for debugging purposes.
        Do(context.Background())//.Do(context.Background()) executes the search query and returns a *elastic.SearchResult and an error, if any.
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index(). 
        Index(index).// specify the name of the Elasticsearch index
        Id(id). //specify the unique identifier for the document being stored
        BodyJson(i).//specify the data to be stored
        Do(context.Background())
    return err
}