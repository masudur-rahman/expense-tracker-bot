package arangodb

import (
	"context"

	"github.com/masudur-rahman/expense-tracker-bot/infra/database/nosql"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"

	arango "github.com/arangodb/go-driver"
)

type ArangoDB struct {
	ctx            context.Context
	db             arango.Database
	id             string
	collectionName string
}

func NewArangoDB(ctx context.Context, db arango.Database) ArangoDB {
	return ArangoDB{
		db:  db,
		ctx: ctx,
	}
}

func (a ArangoDB) Collection(collection string) nosql.Database {
	a.collectionName = collection
	return a
}

func (a ArangoDB) ID(id string) nosql.Database {
	a.id = id
	return a
}

func (a ArangoDB) FindOne(document interface{}, filter ...interface{}) (bool, error) {
	if err := checkIdOrFilterNonEmpty(a.id, filter); err != nil {
		return false, err
	}

	collection, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return false, err
	}

	if filter == nil {
		meta, err := collection.ReadDocument(a.ctx, a.id, document)
		if arango.IsNotFoundGeneral(err) {
			return false, nil
		}
		return meta.ID != "", err
	}

	query := generateArangoQuery(a.collectionName, filter[0], false)
	results, err := executeArangoQuery(a.ctx, a.db, query, 1)
	if arango.IsNotFoundGeneral(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if len(results) != 1 {
		return false, nil
	}

	//reflect.ValueOf(documents).Elem().Set(reflect.ValueOf(results))
	if err = pkg.ParseInto(results[0], document); err != nil {
		return false, err
	}
	return true, nil
}

func (a ArangoDB) FindMany(documents interface{}, filter interface{}) error {
	_, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return err
	}

	query := generateArangoQuery(a.collectionName, filter, false)
	results, err := executeArangoQuery(a.ctx, a.db, query, -1)
	if err != nil {
		return err
	}

	return pkg.ParseInto(results, documents)
}

func (a ArangoDB) InsertOne(document interface{}) (id string, err error) {
	collection, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return "", err
	}

	meta, err := collection.CreateDocument(a.ctx, document)
	if err != nil {
		return "", err
	}

	return meta.Key, nil
}

func (a ArangoDB) InsertMany(documents []interface{}) ([]string, error) {
	collection, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return nil, err
	}

	metas, _, err := collection.CreateDocuments(a.ctx, documents)
	if err != nil {
		return nil, err
	}

	// Extract IDs of inserted documents
	ids := make([]string, len(metas))
	for i, result := range metas {
		ids[i] = string(result.ID)
	}

	return ids, nil
}

func (a ArangoDB) UpdateOne(document interface{}) error {
	if err := checkIDNonEmpty(a.id); err != nil {
		return err
	}

	collection, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return err
	}

	_, err = collection.UpdateDocument(a.ctx, a.id, document)
	return err
}

func (a ArangoDB) DeleteOne(filter ...interface{}) error {
	if err := checkIdOrFilterNonEmpty(a.id, filter); err != nil {
		return err
	}

	collection, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return err
	}

	if filter == nil {
		_, err = collection.RemoveDocument(a.ctx, a.id)
		return err
	}

	query := generateArangoQuery(a.collectionName, filter[0], true)
	_, err = executeArangoQuery(a.ctx, a.db, query, 1)
	if err != nil {
		return err
	}

	return nil
}

func (a ArangoDB) Query(query string, bindParams map[string]interface{}) (interface{}, error) {
	_, err := getDBCollection(a.ctx, a.db, a.collectionName)
	if err != nil {
		return nil, err
	}

	return executeArangoQuery(a.ctx, a.db, &Query{queryString: query, bindVars: bindParams}, -1)
}
