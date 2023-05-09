## Description

mysqlctl is a helper package that implements various CRUD operations for MySQL.
Currently, it implements these interfaces:

```go
type DBController interface {
	CreateDatabase(dbName string) error
	DeleteDatabase(dbName string) error
	ListDatabases() ([]string, error)
	DatabaseExists(dbName string) (bool, error)
	Size(dbName string) (int, error)
}

type GrantController interface {
	Grant(grantName, dbName, username string) error
	GrantExists(grantName, dbName, username string) (bool, error)
	GrantAll(dbName, username string) error
	RevokeAll(dbName, username string) error
	Revoke(grantName, dbName, username string) error
}

type UserController interface {
	CreateUser(username, password string) error
	UpdateUserPassword(username, password string) error
	DeleteUser(username string) error
	ListUsers() ([]string, error)
	UserExists(username string) (bool, error)
	CreateUserWithMaxConn(username, password string, maxConn int) error
	UpdateUserMaxConn(username string, maxConn int) error
	GetUserMaxConn(username string) (int, error)
}
```

List of supported GRANTS:

```go
ALTER
ALTER ROUTINE
CREATE
CREATE ROUTINE
CREATE TEMPORARY TABLES
CREATE VIEW
DELETE
DROP
EVENT
EXECUTE
INDEX
INSERT
LOCK TABLES
REFERENCES
SELECT
SHOW VIEW
TRIGGER
UPDATE
```
