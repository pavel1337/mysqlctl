package mysqlctl

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// generateTestNames generates a list of test names.
func generateTestNames() []string {
	return []string{"test1", "test2", "test3"}
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
	err := c.CreateDatabase("test1")
	assert.NoError(t, err)

	err = c.CreateDatabase("test1")
	assert.Error(t, err)
	assert.Equal(t, ErrDBExists, err)

	err = c.CreateDatabase("")
	assert.Error(t, err)

	for _, name := range baseDBs {
		err = c.CreateDatabase(name)
		assert.Error(t, err)
	}

	err = c.DeleteDatabase("test1")
	assert.NoError(t, err)
}

func TestMySQLController_DeleteDatabase(t *testing.T) {
	c := createTestController()
	err := c.DeleteDatabase("test1")
	assert.Error(t, err)
	assert.Equal(t, ErrDBDoesNotExist, err)

	err = c.CreateDatabase("test1")
	assert.NoError(t, err)

	err = c.DeleteDatabase("test1")
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
	exists, err := c.DatabaseExists("test1")
	assert.NoError(t, err)
	assert.False(t, exists)

	err = c.CreateDatabase("test1")
	assert.NoError(t, err)

	exists, err = c.DatabaseExists("test1")
	assert.NoError(t, err)
	assert.True(t, exists)

	err = c.DeleteDatabase("test1")
	assert.NoError(t, err)
}
