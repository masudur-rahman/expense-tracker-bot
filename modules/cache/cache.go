package cache

import (
	"fmt"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/xid"
)

var c *cache.Cache

func init() {
	c = cache.New(6*time.Hour, 24*time.Hour)
}

func StoreData(obj any) (key string) {
	key = xid.New().String()
	c.Set(key, obj, 0)
	return key
}

func FetchData(key string, obj any) error {
	data, ok := c.Get(key)
	if !ok {
		return fmt.Errorf("no data found")
	}

	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(data))
	c.Delete(key)
	return nil
}

func FetchDataWithCustomFunc(key string, obj any, getFunc func() (any, error)) error {
	if err := FetchData(key, obj); err == nil {
		return nil
	}

	obj, err := getFunc()
	if err != nil {
		return err
	}

	c.Set(key, obj, 5*time.Minute)
	return nil
}
