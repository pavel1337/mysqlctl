package mysqlctl

import (
	"fmt"
	"strings"
)

type GrantController interface {
	Grant(grantName, dbName, username string) error
	GrantExists(grantName, dbName, username string) (bool, error)
	GrantAll(dbName, username string) error
	RevokeAll(dbName, username string) error
	Revoke(grantName, dbName, username string) error
}

// just checking if the database name is valid
var _ GrantController = &MySQLController{}

var (
	ErrInvalidGrant = fmt.Errorf("invalidwgrant")
)

// GrantAll grants all privileges for the given database and user
func (c *MySQLController) GrantAll(dbName, username string) error {
	ok, err := c.UserExists(username)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if !ok {
		return ErrUserDoesNotExist
	}

	ok, err = c.DatabaseExists(dbName)
	if err != nil {
		return fmt.Errorf("error checking if database exists: %w", err)
	}
	if !ok {
		return ErrDBDoesNotExist
	}

	_, err = c.db.Exec("GRANT ALL PRIVILEGES ON `" + dbName + "`.* TO '" + username + "'@'%'")
	if err != nil {
		return fmt.Errorf("error granting privileges: %w", err)
	}
	return nil
}

// RevokeAll revokes all privileges for the given database and user
func (c *MySQLController) RevokeAll(dbName, username string) error {
	err := validateDBName(dbName)
	if err != nil {
		return fmt.Errorf("error validating database name: %w", err)
	}

	err = validateUsername(username)
	if err != nil {
		return fmt.Errorf("error validating username: %w", err)
	}

	_, err = c.db.Exec("REVOKE ALL PRIVILEGES ON `" + dbName + "`.* FROM '" + username + "'@'%'")
	if err != nil {
		return fmt.Errorf("error revoking privileges: %w", err)
	}
	return nil
}

// Grant grants the given grant to the given database and user
func (c *MySQLController) Grant(grantName, dbName, username string) error {
	err := validateDBName(dbName)
	if err != nil {
		return fmt.Errorf("error validating database name: %w", err)
	}

	err = validateUsername(username)
	if err != nil {
		return fmt.Errorf("error validating username: %w", err)
	}

	grantName = strings.ToUpper(grantName)
	err = validateGrant(grantName)
	if err != nil {
		return fmt.Errorf("error validating grant: %w", err)
	}

	q := fmt.Sprintf("GRANT %s ON `%s`.* TO '%s'@'%%'", grantName, dbName, username)
	_, err = c.db.Exec(q)
	if err != nil {
		return fmt.Errorf("error granting privileges: %w", err)
	}
	return nil
}

// Revoke revokes the given grant from the given database and user
func (c *MySQLController) Revoke(grantName, dbName, username string) error {
	err := validateDBName(dbName)
	if err != nil {
		return fmt.Errorf("error validating database name: %w", err)
	}

	err = validateUsername(username)
	if err != nil {
		return fmt.Errorf("error validating username: %w", err)
	}

	grantName = strings.ToUpper(grantName)
	err = validateGrant(grantName)
	if err != nil {
		return fmt.Errorf("error validating grant: %w", err)
	}

	q := fmt.Sprintf("REVOKE %s ON `%s`.* FROM '%s'@'%%'", grantName, dbName, username)
	_, err = c.db.Exec(q)
	if err != nil {
		return fmt.Errorf("error revoking privileges: %w", err)
	}

	return nil
}

// GrantExists returns true if the given grant exists for the given database and user
func (c *MySQLController) GrantExists(grantName, dbName, username string) (bool, error) {
	err := validateDBName(dbName)
	if err != nil {
		return false, fmt.Errorf("error validating database name: %w", err)
	}

	err = validateUsername(username)
	if err != nil {
		return false, fmt.Errorf("error validating username: %w", err)
	}

	grantName = strings.ToUpper(grantName)
	err = validateGrant(grantName)
	if err != nil {
		return false, fmt.Errorf("error validating grant: %w", err)
	}

	grantColumn := grants[grantName]

	q := fmt.Sprintf("SELECT COUNT(*) FROM mysql.db WHERE Db = '%s' AND User = '%s' AND %s = 'Y'", dbName, username, grantColumn)
	var count int
	err = c.db.QueryRow(q).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if grant exists: %w", err)
	}

	return count > 0, nil
}

// validateGrant checks if the given grant is valid
func validateGrant(grantName string) error {
	if grantName == "" {
		return ErrInvalidGrant
	}

	if _, ok := grants[grantName]; !ok {
		return ErrInvalidGrant
	}

	return nil
}
