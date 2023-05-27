package mock

type Database interface {
	Entity(name string) Database

	ID(id string) Database

	FindOne(document interface{}, filter ...interface{}) (bool, error)
	FindMany(documents interface{}, filter interface{}) error

	InsertOne(document interface{}) (id string, err error)
	InsertMany(documents []interface{}) ([]string, error)

	UpdateOne(document interface{}) error

	DeleteOne(filter ...interface{}) error
}
