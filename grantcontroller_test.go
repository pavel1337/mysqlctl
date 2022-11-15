package mysqlctl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySQLController_GrantAll(t *testing.T) {
	c := createTestController()
	err := c.GrantAll("test1", "test1")
	assert.Error(t, err)
	assert.Equal(t, ErrUserDoesNotExist, err)

	c.CreateUser("test1", "password")
	defer c.DeleteUser("test1")

	err = c.GrantAll("test1", "test1")
	assert.Error(t, err)
	assert.Equal(t, ErrDBDoesNotExist, err)

	c.CreateDatabase("test1")
	defer c.DeleteDatabase("test1")

	err = c.GrantAll("test1", "test1")
	assert.NoError(t, err)

	err = openMySQL("test1", "password", "test1")
	assert.NoError(t, err)

	err = c.RevokeAll("test1", "test1")
	assert.NoError(t, err)

	err = openMySQL("test1", "password", "test1")
	assert.Error(t, err)
}

func TestMySQLController_RevokeAll(t *testing.T) {
	c := createTestController()
	err := c.RevokeAll("test1", "test1")
	assert.Error(t, err)

	c.CreateUser("test1", "password")
	defer c.DeleteUser("test1")

	err = c.RevokeAll("test1", "test1")
	assert.Error(t, err)

	c.CreateDatabase("test1")
	defer c.DeleteDatabase("test1")

	err = c.RevokeAll("test1", "test1")
	assert.Error(t, err)
}

func TestMySQLController_Grant(t *testing.T) {
	c := createTestController()
	err := c.Grant("test1", "test1", "test1")
	assert.Error(t, err)

	c.CreateUser("test", "password")
	defer c.DeleteUser("test")

	c.CreateDatabase("test")
	defer c.DeleteDatabase("test")

	err = c.Grant("select", "test", "test")
	assert.NoError(t, err)
}
