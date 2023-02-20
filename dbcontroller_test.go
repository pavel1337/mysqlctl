package mysqlctl

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	testDB       = "test-db"
	testUser     = "test-user"
	testPassword = "test-password"
)

// generateTestNames generates a list of test names.
func generateTestNames() []string {
	return []string{"test1", "test2", "test3", "test4-test", "test5-test"}
}

// createTestController creates a test controller.
func createTestController() *MySQLController {
	c, err := NewMySQLController("root:password@tcp(127.0.0.1:6603)/")
	if err != nil {
		panic(err)
	}
	return c
}

func TestMySQLController_CreateDatabase(t *testing.T) {
	c := createTestController()
	err := c.CreateDatabase(testDB)
	assert.NoError(t, err)

	err = c.CreateDatabase(testDB)
	assert.Error(t, err)
	assert.Equal(t, ErrDBExists, err)

	err = c.CreateDatabase("")
	assert.Error(t, err)

	for _, name := range baseDBs {
		err = c.CreateDatabase(name)
		assert.Error(t, err)
	}

	err = c.DeleteDatabase(testDB)
	assert.NoError(t, err)
}

func TestMySQLController_DeleteDatabase(t *testing.T) {
	c := createTestController()
	err := c.DeleteDatabase(testDB)
	assert.Error(t, err)
	assert.Equal(t, ErrDBDoesNotExist, err)

	err = c.CreateDatabase(testDB)
	assert.NoError(t, err)

	err = c.DeleteDatabase(testDB)
	assert.NoError(t, err)

	err = c.DeleteDatabase("")
	assert.Error(t, err)

	for _, name := range baseDBs {
		err = c.DeleteDatabase(name)
		assert.Error(t, err)
	}
}

func TestMySQLController_ListDatabases(t *testing.T) {
	c := createTestController()
	dbs, err := c.ListDatabases()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(dbs))

	for _, name := range generateTestNames() {
		err = c.CreateDatabase(name)
		assert.NoError(t, err)
	}

	dbs, err = c.ListDatabases()
	assert.NoError(t, err)
	assert.Equal(t, generateTestNames(), dbs)

	for _, name := range generateTestNames() {
		err = c.DeleteDatabase(name)
		assert.NoError(t, err)
	}
}

func TestMySQLController_DatabaseExists(t *testing.T) {
	c := createTestController()
	exists, err := c.DatabaseExists(testDB)
	assert.NoError(t, err)
	assert.False(t, exists)

	err = c.CreateDatabase(testDB)
	assert.NoError(t, err)

	exists, err = c.DatabaseExists(testDB)
	assert.NoError(t, err)
	assert.True(t, exists)

	err = c.DeleteDatabase(testDB)
	assert.NoError(t, err)
}

func TestMySQLController_Size(t *testing.T) {
	c := createTestController()
	size, err := c.Size(testDB)
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = c.CreateDatabase(testDB)
	assert.NoError(t, err)

	size, err = c.Size(testDB)
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = c.CreateUser(testUser, testPassword)
	assert.NoError(t, err)
	defer c.DeleteUser(testUser)

	err = c.GrantAll(testDB, testUser)
	assert.NoError(t, err)

	db, err := openMySQLWithDB(testUser, testPassword, testDB)
	assert.NoError(t, err)
	defer db.Close()

	err = createTestTable(db)
	assert.NoError(t, err)

	size, err = c.Size(testDB)
	assert.NoError(t, err)
	assert.Equal(t, 16384, size)

	err = c.DeleteDatabase(testDB)
	assert.NoError(t, err)
}

func openMySQLWithDB(username, password, database string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(127.0.0.1:6603)/%s", username, password, database)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %s", err)
	}

	return db, nil
}

// createTestTable creates a test table.
func createTestTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE test (id INT, name VARCHAR(255))")
	return err
}

// randomString generates a random string of a length n.
func randomString(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = byte(65 + i)
	}
	return string(b)
}

func Test_randomString(t *testing.T) {
	s := randomString(10)
	assert.Equal(t, 10, len(s))

	s = randomString(0)
	assert.Equal(t, 0, len(s))

	s = randomString(-1)
	assert.Equal(t, 0, len(s))

	s = randomString(100)
	assert.Equal(t, 100, len(s))
}
