# database
Database Engine for different SQL and NoSQL databases

## Install
```shell
go get -u github.com/masudur-rahman/database
```

## Quickstart

```go
package main

import (
	"context"
	"time"

	"github.com/masudur-rahman/database/sql"
	"github.com/masudur-rahman/database/sql/sqlite"
	"github.com/masudur-rahman/database/sql/sqlite/lib"
)

type User struct {
	ID        int64     `db:"id,pk autoincr"`
	Name      string    `db:"name,uq"`
	FullName  string    `db:"full_name,uqs"`
	Email     string    `db:",uqs"`
	CreatedAt time.Time `db:"created_at"`
}

func main() {
	// Create sqlite connection
	conn, _ := lib.GetSQLiteConnection("test.db")

	// Start a database engine
	var db sql.Database
	db = sqlite.NewSqlite(context.Background(), conn)

	// Migrate database
	db.Sync(User{})

	db = db.Table("user")

	// Insert
	db.InsertOne(&User{Name: "masud", FullName: "Masudur Rahman", Email: "masud@example.com"})

	// Read
	var user User
	db.ID(1).FindOne(&user)
	db.Where("email=?", "masud@example.com").FindOne(&user)
	db.FindOne(&user, User{Name: "masud"})
	db.Columns("name", "email").FindOne(&user, User{Name: "masud"}) // fetch only name, email columns

	// Update
	db.ID(user.ID).UpdateOne(User{Email: "test@example.com"})
	db.Where("email=?", "test@example.com").UpdateOne(User{FullName: "Test User"})

	// Delete
	db.ID(1).DeleteOne()              // delete by id
	db.DeleteOne(User{Name: "masud"}) // delete using filter
}
```
