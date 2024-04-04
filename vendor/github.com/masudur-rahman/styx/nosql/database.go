package nosql

type Engine interface {
	Collection(name string) Engine

	ID(id string) Engine

	FindOne(document interface{}, filter ...interface{}) (bool, error)
	FindMany(documents interface{}, filter interface{}) error

	InsertOne(document interface{}) (id string, err error)
	InsertMany(documents []interface{}) ([]string, error)

	UpdateOne(document interface{}) error

	DeleteOne(filter ...interface{}) error

	Query(query string, bindParams map[string]interface{}) (interface{}, error)
}
