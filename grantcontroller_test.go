package mysqlctl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySQLController_GrantAll(t *testing.T) {
	c := createTestController()
	err := c.GrantAll(testDB, testUser)
	assert.Error(t, err)
	assert.Equal(t, ErrUserDoesNotExist, err)

	c.CreateUser(testUser, testPassword)
	defer c.DeleteUser(testUser)

	err = c.GrantAll(testDB, testUser)
	assert.Error(t, err)
	assert.Equal(t, ErrDBDoesNotExist, err)

	c.CreateDatabase(testDB)
	defer c.DeleteDatabase(testDB)

	err = c.GrantAll(testDB, testUser)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, testDB)
	assert.NoError(t, err)

	err = c.RevokeAll(testDB, testUser)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, testDB)
	assert.Error(t, err)
}

func TestMySQLController_RevokeAll(t *testing.T) {
	c := createTestController()
	err := c.RevokeAll(testDB, testUser)
	assert.Error(t, err)

	c.CreateUser(testUser, testPassword)
	defer c.DeleteUser(testUser)

	err = c.RevokeAll(testDB, testUser)
	assert.Error(t, err)

	c.CreateDatabase(testDB)
	defer c.DeleteDatabase(testDB)

	err = c.RevokeAll(testDB, testUser)
	assert.Error(t, err)
}

func TestMySQLController_Grant(t *testing.T) {
	c := createTestController()
	err := c.Grant("select", "non-existing-db", "non-existing-user")
	assert.Error(t, err)

	c.CreateUser(testUser, testPassword)
	defer c.DeleteUser(testUser)

	c.CreateDatabase(testDB)
	defer c.DeleteDatabase(testDB)

	for g := range grants {
		err := c.Grant(g, testDB, testUser)
		assert.NoError(t, err)
	}
}

func TestMySQLController_GrantExists(t *testing.T) {
	c := createTestController()
	c.CreateUser(testUser, testPassword)
	defer c.DeleteUser(testUser)

	c.CreateDatabase(testDB)
	defer c.DeleteDatabase(testDB)

	err := c.Grant("select", testDB, testUser)
	assert.NoError(t, err)

	err = c.Grant("select", testDB, testUser)
	assert.NoError(t, err)

	b, err := c.GrantExists("select", testDB, testUser)
	assert.NoError(t, err)
	assert.True(t, b)

	for g := range grants {
		err := c.Grant(g, testDB, testUser)
		assert.NoError(t, err)

		b, err := c.GrantExists(g, testDB, testUser)
		assert.NoError(t, err)
		assert.True(t, b)
	}

	for g := range grants {
		err := c.Revoke(g, testDB, testUser)
		assert.NoError(t, err)

		b, err := c.GrantExists(g, testDB, testUser)
		assert.NoError(t, err)
		assert.False(t, b)
	}

	_, err = c.GrantExists("non-existing-grant", testDB, testUser)
	assert.ErrorIs(t, err, ErrInvalidGrant)
}

func TestMySQLController_Revoke(t *testing.T) {
	c := createTestController()
	err := c.Revoke("select", "non-existing-db", "non-existing-user")
	assert.Error(t, err)

	c.CreateUser(testUser, testPassword)
	defer c.DeleteUser(testUser)

	c.CreateDatabase(testDB)
	defer c.DeleteDatabase(testDB)

	for g := range grants {
		err := c.Revoke(g, testDB, testUser)
		assert.Error(t, err)
	}

	for g := range grants {
		err := c.Grant(g, testDB, testUser)
		assert.NoError(t, err)
	}

	for g := range grants {
		err := c.Revoke(g, testDB, testUser)
		assert.NoError(t, err)
	}

	for g := range grants {
		err := c.Revoke(g, testDB, testUser)
		assert.Error(t, err)
	}

}
