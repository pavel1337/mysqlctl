package mysqlctl

import (
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
