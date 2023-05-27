package _map

import (
	"reflect"

	"github.com/masudur-rahman/expense-tracker-bot/infra/database/mock"

	"github.com/rs/xid"
)

type MockDB struct {
	db     map[string]map[string]interface{}
	entity string
	id     string
}

func NewMockDB() *MockDB {
	return &MockDB{
		db: map[string]map[string]interface{}{},
	}
}

func (m *MockDB) Entity(name string) mock.Database {
	m.entity = name
	return m
}

func (m *MockDB) ID(id string) mock.Database {
	m.id = id
	return m
}

func (m *MockDB) FindOne(document interface{}, filter ...interface{}) (bool, error) {
	if err := checkEntityNameNonEmpty(m.entity); err != nil {
		return false, err
	}
	if err := checkIdOrFilterNonEmpty(m.id, filter); err != nil {
		return false, err
	}
	if filter == nil {
		doc, ok := m.db[m.entity][m.id]
		if !ok {
			return ok, nil
		}
		reflect.ValueOf(document).Elem().Set(reflect.ValueOf(doc))
	}

	return false, nil
}

func (m *MockDB) FindMany(documents interface{}, filter interface{}) error {
	if err := checkEntityNameNonEmpty(m.entity); err != nil {
		return err
	}

	reflect.ValueOf(documents).Elem().Set(reflect.ValueOf(documents))
	return nil
}

func (m *MockDB) InsertOne(document interface{}) (id string, err error) {
	id = xid.New().String()
	reflect.ValueOf(document).Elem().FieldByName("ID").SetString(id)
	m.db[m.entity][id] = document
	return
}

func (m *MockDB) InsertMany(documents []interface{}) ([]string, error) {
	var ids []string
	for document := range documents {
		id := xid.New().String()
		reflect.ValueOf(document).Elem().FieldByName("ID").SetString(id)
		m.db[m.entity][id] = document
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *MockDB) UpdateOne(document interface{}) error {
	if err := checkEntityNameNonEmpty(m.entity); err != nil {
		return err
	}
	if err := checkIDNonEmpty(m.id); err != nil {
		return err
	}
	m.db[m.entity][m.id] = document
	return nil
}

func (m *MockDB) DeleteOne(filter ...interface{}) error {
	delete(m.db[m.entity], m.id)
	return nil
}
