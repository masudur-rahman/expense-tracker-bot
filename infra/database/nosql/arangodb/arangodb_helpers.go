package arangodb

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/masudur-rahman/expense-tracker-bot/configs"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	arango "github.com/arangodb/go-driver"
	ahttp "github.com/arangodb/go-driver/http"
	"github.com/iancoleman/strcase"
)

func InitializeArangoDB(ctx context.Context) (arango.Database, error) {
	cfg := configs.PurrfectConfig.Database.ArangoDB
	conn, err := ahttp.NewConnection(ahttp.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)},
		TLSConfig: &tls.Config{ /*...*/ },
	})
	if err != nil {
		return nil, err
	}

	c, err := arango.NewClient(arango.ClientConfig{
		Connection:     conn,
		Authentication: arango.BasicAuthentication(cfg.User, cfg.Password),
	})
	if err != nil {
		return nil, err
	}

	db, err := c.Database(ctx, cfg.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Query struct {
	queryString string
	bindVars    map[string]interface{}
}

func getDBCollection(ctx context.Context, db arango.Database, col string) (arango.Collection, error) {
	collection, err := db.Collection(ctx, col)
	if err != nil {
		if arango.IsNotFoundGeneral(err) {
			return db.CreateCollection(ctx, col, &arango.CreateCollectionOptions{})
		}
		return nil, err
	}

	return collection, nil
}

func generateArangoQuery(collection string, filter interface{}, removeQuery bool) *Query {
	queryString := "FOR doc IN " + collection // + " FILTER "
	bindVars := map[string]interface{}{}

	var filters []string

	val := reflect.ValueOf(filter)
	for idx := 0; idx < val.NumField(); idx++ {
		field := val.Type().Field(idx)
		if val.Field(idx).IsZero() {
			continue
		}

		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = strcase.ToLowerCamel(field.Name)
		}
		filters = append(filters, fmt.Sprintf("doc.%s == @%s", fieldName, fieldName))
		bindVars[fieldName] = val.Field(idx).Interface()
	}

	if len(filters) > 0 {
		queryString += " FILTER "
		queryString += strings.Join(filters, " AND ")
	}

	if removeQuery {
		queryString += " REMOVE doc IN " + collection
	} else {
		queryString += " RETURN doc"
	}

	return &Query{
		queryString: queryString,
		bindVars:    bindVars,
	}
}

func executeArangoQuery(ctx context.Context, db arango.Database, query *Query, lim int64) ([]interface{}, error) {
	cursor, err := db.Query(ctx, query.queryString, query.bindVars)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	var results []interface{}
	for {
		var doc interface{}
		_, err = cursor.ReadDocument(ctx, &doc)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		results = append(results, doc)
		if lim > 0 && int64(len(results)) >= lim {
			break
		}
	}

	return results, nil
}

func checkEntityNameNonEmpty(entity string) error {
	if entity == "" {
		return models.StatusError{
			Message: "entity name must be set",
		}
	}
	return nil
}

func checkIDNonEmpty(id string) error {
	if id == "" {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "must provide document id",
		}
	}
	return nil
}

func checkIdOrFilterNonEmpty(id string, filter interface{}) error {
	if id == "" && filter == nil {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "must provide id and/or filter",
		}
	}
	return nil
}
