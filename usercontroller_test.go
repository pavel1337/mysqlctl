package mysqlctl

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySQLController_CreateUser(t *testing.T) {
	c := createTestController()
	err := c.CreateUser(testUser, testPassword)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, "")
	assert.NoError(t, err)

	err = c.CreateUser(testUser, testPassword)
	assert.Error(t, err)
	assert.Equal(t, ErrUserExists, err)

	err = c.CreateUser("", "")
	assert.Error(t, err)

	for _, name := range baseUsers {
		err = c.CreateUser(name, testPassword)
		assert.Error(t, err)
	}

	err = c.DeleteUser(testUser)
	assert.NoError(t, err)
}

func TestMySQLController_UpdateUserPassword(t *testing.T) {
	c := createTestController()
	err := c.UpdateUserPassword(testUser, testPassword)
	assert.Error(t, err)
	assert.Equal(t, ErrUserDoesNotExist, err)

	err = c.CreateUser(testUser, testPassword)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, "")
	assert.NoError(t, err)

	err = c.UpdateUserPassword(testUser, testPassword)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, "")
	assert.NoError(t, err)

	err = c.UpdateUserPassword("", testPassword)
	assert.Error(t, err)

	err = c.UpdateUserPassword(testUser, "")
	assert.Error(t, err)

	for _, name := range baseUsers {
		err = c.UpdateUserPassword(name, testPassword)
		assert.Error(t, err)
	}

	err = c.DeleteUser(testUser)
	assert.NoError(t, err)
}

func TestMySQLController_DeleteUser(t *testing.T) {
	c := createTestController()
	err := c.DeleteUser(testUser)
	assert.Error(t, err)
	assert.Equal(t, ErrUserDoesNotExist, err)

	err = c.CreateUser(testUser, testPassword)
	assert.NoError(t, err)

	err = c.DeleteUser(testUser)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, "")
	assert.Error(t, err)

	err = c.DeleteUser("")
	assert.Error(t, err)

	for _, name := range baseUsers {
		err = c.DeleteUser(name)
		assert.Error(t, err)
	}
}

func TestMySQLController_ListUsers(t *testing.T) {
	c := createTestController()
	names, err := c.ListUsers()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(names))

	err = c.CreateUser(testUser, testPassword)
	assert.NoError(t, err)

	names, err = c.ListUsers()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(names))

	err = c.DeleteUser(testUser)
	assert.NoError(t, err)
}

func openMySQL(username, password, database string) error {
	connStr := fmt.Sprintf("%s:%s@tcp(127.0.0.1:6603)/%s", username, password, database)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return fmt.Errorf("error creating database connection: %s", err)
	}
	defer db.Close()

	return db.Ping()
}

func TestMySQLController_UserExists(t *testing.T) {
	c := createTestController()
	_, err := c.UserExists("")
	assert.Error(t, err)

	exists, err := c.UserExists(testUser)
	assert.NoError(t, err)
	assert.False(t, exists)

	err = c.CreateUser(testUser, testPassword)
	assert.NoError(t, err)

	exists, err = c.UserExists(testUser)
	assert.NoError(t, err)
	assert.True(t, exists)

	err = c.DeleteUser(testUser)
	assert.NoError(t, err)
}

func TestMySQLController_UsersMaxConnHandling(t *testing.T) {
	c := createTestController()
	err := c.CreateUserWithMaxConn(testUser, testPassword, 1)
	assert.NoError(t, err)

	err = openMySQL(testUser, testPassword, "")
	assert.NoError(t, err)

	err = c.CreateUserWithMaxConn(testUser, testPassword, 1)
	assert.Error(t, err)

	maxConn, err := c.GetUserMaxConn(testUser)
	assert.NoError(t, err)
	assert.Equal(t, 1, maxConn)

	err = c.UpdateUserMaxConn(testUser, 2)
	assert.NoError(t, err)

	maxConn, err = c.GetUserMaxConn(testUser)
	assert.NoError(t, err)
	assert.Equal(t, 2, maxConn)
}
