package mysqlctl

import (
	"fmt"
	"strings"
)

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

var _ UserController = &MySQLController{}

var baseUsers = []string{"root", "mysql.sys", "mysql.session", "mysql.infoschema"}

var (
	ErrUserExists       = fmt.Errorf("user exists")
	ErrUserDoesNotExist = fmt.Errorf("user does not exist")
)

func (c *MySQLController) CreateUser(username, password string) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

	_, err = c.db.Exec("CREATE USER `" + username + "` IDENTIFIED BY '" + password + "'")
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			return ErrUserExists
		}
	}
	return err
}

func (c *MySQLController) CreateUserWithMaxConn(username, password string, maxConn int) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

	_, err = c.db.Exec("CREATE USER `" + username + "` IDENTIFIED BY '" + password + "' WITH MAX_USER_CONNECTIONS " + fmt.Sprintf("%d", maxConn))
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			return ErrUserExists
		}
	}
	return err
}

func (c *MySQLController) GetUserMaxConn(username string) (int, error) {
	err := validateUsername(username)
	if err != nil {
		return 0, err
	}

	rows, err := c.db.Query("SELECT MAX_USER_CONNECTIONS FROM mysql.user WHERE User = '" + username + "'")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var maxConn int
	for rows.Next() {
		err = rows.Scan(&maxConn)
		if err != nil {
			return 0, err
		}
	}

	return maxConn, nil
}

func (c *MySQLController) UpdateUserMaxConn(username string, maxConn int) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

	_, err = c.db.Exec("ALTER USER `" + username + "` WITH MAX_USER_CONNECTIONS " + fmt.Sprintf("%d", maxConn))
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			return ErrUserDoesNotExist
		}
	}
	return err
}

func (c *MySQLController) UpdateUserPassword(username, password string) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

	err = validatePassword(password)
	if err != nil {
		return err
	}

	_, err = c.db.Exec("SET PASSWORD FOR `" + username + "` = '" + password + "'")
	if err != nil {
		if strings.Contains(err.Error(), "Error 1133") {
			return ErrUserDoesNotExist
		}
	}

	return err
}

func (c *MySQLController) DeleteUser(username string) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(fmt.Sprintf("DROP USER `%s`", username))
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			return ErrUserDoesNotExist
		}
	}
	return err
}

func (c *MySQLController) ListUsers() ([]string, error) {
	rows, err := c.db.Query("SELECT user FROM mysql.user WHERE host = '%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return filterUsers(users), nil
}

func (c *MySQLController) UserExists(username string) (bool, error) {
	err := validateUsername(username)
	if err != nil {
		return false, err
	}

	var exists bool
	err = c.db.QueryRow("SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = ? AND host = '%')", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func filterUsers(users []string) []string {
	var filtered []string
	for _, user := range users {
		if !contains(baseUsers, user) {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if contains(baseUsers, username) {
		return fmt.Errorf("username %s is reserved", username)
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	return nil
}
