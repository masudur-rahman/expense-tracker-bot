package supabase

import (
	"context"
	"database/sql"

	isql "github.com/masudur-rahman/expense-tracker-bot/infra/database/sql"

	"github.com/nedpals/supabase-go"
)

// yRQX2p1ZGenUdlVb
type Supabase struct {
	ctx    context.Context
	table  string
	id     string
	client *supabase.Client
}

func NewSupabase(ctx context.Context, client *supabase.Client) Supabase {
	return Supabase{
		ctx:    ctx,
		client: client,
	}
}

func (s Supabase) Table(name string) isql.Database {
	s.table = name
	return s
}

func (s Supabase) ID(id string) isql.Database {
	s.id = id
	return s
}

func (s Supabase) FindOne(document interface{}, filter ...interface{}) (bool, error) {
	if err := checkIdOrFilterNonEmpty(s.id, filter); err != nil {
		return false, err
	}

	var kvs []keyValue
	if s.id != "" {
		kvs = []keyValue{{"id", s.id}}
	} else {
		kvs = generateFilters(filter[0])
	}

	cl := s.client.DB.From(s.table).Select("*").Single()
	for idx := range kvs {
		cl.Eq(kvs[idx].key, kvs[idx].value)
	}
	if err := cl.Execute(document); err != nil {
		return false, err
	}

	return true, nil
}

func (s Supabase) FindMany(documents interface{}, filter interface{}) error {
	kvs := generateFilters(filter)
	cl := s.client.DB.From(s.table).Select("*")

	for idx := range kvs {
		cl.Eq(kvs[idx].key, kvs[idx].value)
	}
	if err := cl.Execute(documents); err != nil {
		return err
	}

	return nil
}

func (s Supabase) InsertOne(document interface{}) (id string, err error) {
	docs := []Doc{}
	err = s.client.DB.From(s.table).Insert(document).Execute(&docs)
	if err != nil {
		return "", err
	}
	return docs[0].ID, nil
}

func (s Supabase) InsertMany(documents []interface{}) ([]string, error) {
	var ids = make([]string, 0, len(documents))
	for idx := range documents {
		id, err := s.InsertOne(documents[idx])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (s Supabase) UpdateOne(document interface{}) error {
	if err := checkIDNonEmpty(s.id); err != nil {
		return err
	}

	return s.client.DB.From(s.table).Update(document).Eq("id", s.id).Execute(&document)
}

func (s Supabase) DeleteOne(filter ...interface{}) error {
	if err := checkIdOrFilterNonEmpty(s.id, filter); err != nil {
		return err
	}

	var kvs []keyValue
	if s.id != "" {
		kvs = []keyValue{{"id", s.id}}
	} else {
		kvs = generateFilters(filter[0])
	}

	cl := s.client.DB.From(s.table).Delete()
	for idx := range kvs {
		cl.Eq(kvs[idx].key, kvs[idx].value)
	}

	rs := map[string]interface{}{}
	return cl.Execute(&rs)
}

func (s Supabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (s Supabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}
