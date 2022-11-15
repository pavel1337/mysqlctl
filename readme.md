## Description
mysqlctl is a helper package that implements various CRUD operations for MySQL.
Currently, it implements these interfaces:
```go
type DBController interface {
	CreateDatabase(dbName string) error
	DeleteDatabase(dbName string) error
	ListDatabases() ([]string, error)
	DatabaseExists(dbName string) (bool, error)
}

type GrantController interface {
	Grant(grantName, dbName, username string) error
	GrantAll(dbName, username string) error
	RevokeAll(dbName, username string) error
}

type UserController interface {
	CreateUser(username, password string) error
	UpdateUser(username, password string) error
	DeleteUser(username string) error
	ListUsers() ([]string, error)
	UserExists(username string) (bool, error)
}
```
