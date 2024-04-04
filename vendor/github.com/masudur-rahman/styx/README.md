# Styx
Database Engine for different SQL and NoSQL databases

## Install
```shell
go get -u github.com/masudur-rahman/styx
```

## Quickstart

```go
package main

import (
	"context"
	"time"

	"github.com/masudur-rahman/styx/sql"
	"github.com/masudur-rahman/styx/sql/sqlite"
	"github.com/masudur-rahman/styx/sql/sqlite/lib"
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
	var db sql.Engine
	db = sqlite.NewSQLite(context.Background(), conn)

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

<br>
<hr>

### Why `styx` name is chosen as this go orm

* **Mythological Connection:** In Greek mythology, the River Styx separates the world of the living from the world of the dead.
Similarly, this ORM acts as a bridge  between your application code (the living) and the database (the "dead" storage).
It facilitates the flow of data between these two realms.

* **Focus on Data Access:**  The Styx was also considered a barrier or boundary.  Similarly, this ORM acts as a controlled point of access for your application to interact with the database.
It ensures data integrity and prevents unauthorized modification.

* **Symbolism:** The Styx is often depicted as a dark and mysterious river. This can be seen as a metaphor for the complexity of database interactions that this ORM simplifies for developers.

Overall, Styx evokes a sense of connection, control, and hidden power, all of which are relevant functionalities of an ORM.
