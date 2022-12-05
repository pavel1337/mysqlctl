package mysqlctl

import (
	"fmt"
	"strings"
)

type GrantController interface {
	Grant(grantName, dbName, username string) error
	GrantAll(dbName, username string) error
	RevokeAll(dbName, username string) error
}

var validGrants = []string{
	"INSERT", "SELECT", "UPDATE", "DELETE",
}

// just checking if the database name is valid
var _ GrantController = &MySQLController{}

var (
	ErrInvalidGrant = fmt.Errorf("invalid grant")
)

// GrantAll grants all privileges for the given database and user
func (c *MySQLController) GrantAll(dbName, username string) error {
	ok, err := c.UserExists(username)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %s", err)
	}
	if !ok {
		return ErrUserDoesNotExist
	}

	ok, err = c.DatabaseExists(dbName)
	if err != nil {
		return fmt.Errorf("error checking if database exists: %s", err)
	}
	if !ok {
		return ErrDBDoesNotExist
	}

	_, err = c.db.Exec("GRANT ALL PRIVILEGES ON `" + dbName + "`.* TO '" + username + "'@'%'")
	if err != nil {
		return fmt.Errorf("error granting privileges: %s", err)
	}
	return nil
}

// RevokeAll revokes all privileges for the given database and user
func (c *MySQLController) RevokeAll(dbName, username string) error {
	err := validateDBName(dbName)
	if err != nil {
		return fmt.Errorf("error validating database name: %s", err)
	}

	err = validateUsername(username)
	if err != nil {
		return fmt.Errorf("error validating username: %s", err)
	}

	_, err = c.db.Exec("REVOKE ALL PRIVILEGES ON `" + dbName + "`.* FROM '" + username + "'@'%'")
	if err != nil {
		return fmt.Errorf("error revoking privileges: %s", err)
	}
	return nil
}

// Grant grants the given grant to the given database and user
func (c *MySQLController) Grant(grantName, dbName, username string) error {
	err := validateDBName(dbName)
	if err != nil {
		return fmt.Errorf("error validating database name: %s", err)
	}

	err = validateUsername(username)
	if err != nil {
		return fmt.Errorf("error validating username: %s", err)
	}

	err = validateGrant(grantName)
	if err != nil {
		return fmt.Errorf("error validating grant: %s", err)
	}

	_, err = c.db.Exec("GRANT " + grantName + " ON `" + dbName + "`.* TO '" + username + "'@'%'")
	if err != nil {
		return fmt.Errorf("error granting privileges: %s", err)
	}
	return nil
}

// validateGrant checks if the given grant is valid
func validateGrant(grantName string) error {
	if grantName == "" {
		return ErrInvalidGrant
	}

	if !contains(validGrants, strings.ToUpper(grantName)) {
		return ErrInvalidGrant
	}

	return nil
}
