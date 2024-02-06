package mysqlctl

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DBController interface {
	// CreateDatabase creates a database.
	CreateDatabase(dbName string) error
	// DeleteDatabase deletes a database.
	DeleteDatabase(dbName string) error
	// ListDatabases returns a list of databases.
	ListDatabases() ([]string, error)
	// DatabaseExists returns true if the database exists.
	DatabaseExists(dbName string) (bool, error)
	// Size returns the size of the database in Bytes.
	Size(dbName string) (int, error)
	// Tables returns a list of tables in the database.
	Tables(dbName string) ([]string, error)
}

var _ DBController = &MySQLController{}

var baseDBs = []string{"information_schema", "mysql", "performance_schema", "sys"}

var (
	ErrDBExists       = fmt.Errorf("database exists")
	ErrDBDoesNotExist = fmt.Errorf("database does not exist")
)

type MySQLController struct {
	db *sql.DB
}

// Option is a function that configures the MySQLController.
type Option func(*MySQLController)

// WithBadUsernames returns an Option that configures the MySQLController to
// ignore the given usernames when listing databases.
func WithBadUsernames(usernames []string) Option {
	return func(c *MySQLController) {
		baseUsers = append(baseUsers, usernames...)
	}
}

// NewMySQLController creates a new MySQLController.
func NewMySQLController(connStr string, opts ...Option) (*MySQLController, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection: %s", err)
	}

	c := &MySQLController{db: db}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *MySQLController) Close() error {
	return c.db.Close()
}

func (c *MySQLController) CreateDatabase(dbName string) error {
	err := validateDBName(dbName)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", dbName))
	if err != nil {
		if strings.Contains(err.Error(), "Error 1007") {
			return ErrDBExists
		}
	}
	return err
}

func (c *MySQLController) DeleteDatabase(dbName string) error {
	err := validateDBName(dbName)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(fmt.Sprintf("DROP DATABASE `%s`", dbName))
	if err != nil {
		if strings.Contains(err.Error(), "Error 1008") {
			return ErrDBDoesNotExist
		}
	}
	return err
}

func (c *MySQLController) ListDatabases() ([]string, error) {
	rows, err := c.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var databases []string
	for rows.Next() {
		var database string
		err = rows.Scan(&database)
		if err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}

	return filterBaseDatabases(databases), nil
}

func (c *MySQLController) DatabaseExists(dbName string) (bool, error) {
	err := validateDBName(dbName)
	if err != nil {
		return false, err
	}

	var count int
	err = c.db.QueryRow("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", dbName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Size returns the size of the database in Bytes.
func (c *MySQLController) Size(dbName string) (int, error) {
	err := validateDBName(dbName)
	if err != nil {
		return 0, err
	}

	var size *int
	err = c.db.QueryRow("SELECT SUM(data_length + index_length) FROM information_schema.tables WHERE table_schema = ?", dbName).Scan(&size)
	if err != nil {
		return 0, err
	}
	if size == nil {
		return 0, nil
	}

	return *size, nil
}

// Tables returns a list of tables in the database.
func (c *MySQLController) Tables(dbName string) ([]string, error) {
	err := validateDBName(dbName)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", dbName)
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func filterBaseDatabases(dbs []string) []string {
	var filtered []string
	for _, db := range dbs {
		if !contains(baseDBs, db) {
			filtered = append(filtered, db)
		}
	}
	return filtered
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// validateDBName validates a database name.
func validateDBName(dbName string) error {
	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	if contains(baseDBs, dbName) {
		return fmt.Errorf("%v is a disallowed database name", dbName)
	}
	return nil
}
