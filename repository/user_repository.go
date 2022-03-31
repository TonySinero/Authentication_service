package repository

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"time"
)

type UserPostgres struct {
	db     *sql.DB
	logger logging.Logger
}

func NewUserPostgres(db *sql.DB, logger logging.Logger) *UserPostgres {
	return &UserPostgres{db: db, logger: logger}
}

// GetUserByID ...
func (u UserPostgres) GetUserByID(id int) (*model.ResponseUser, error) {
	var user model.ResponseUser
	result := u.db.QueryRow("SELECT id, email, role, created_at FROM users WHERE id = $1", id)
	if err := result.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt); err != nil {
		u.logger.Errorf("GetUserByID: error while scanning for user:%s", err)
		return nil, fmt.Errorf("getUserByID: repository error:%w", err)
	}
	return &user, nil
}

// GetUserPasswordByID ...
func (u UserPostgres) GetUserPasswordByID(id int) (string, error) {
	var password string
	result := u.db.QueryRow("SELECT password FROM users WHERE id = $1", id)
	if err := result.Scan(&password); err != nil {
		u.logger.Errorf("GetUserPasswordByID: error while scanning for user:%s", err)
		return "", fmt.Errorf("getUserPasswordByID: repository error:%w", err)
	}
	return password, nil
}

// GetUserAll ...
func (u *UserPostgres) GetUserAll(page int, limit int) ([]model.ResponseUser, int, error) {
	transaction, err := u.db.Begin()
	if err != nil {
		u.logger.Errorf("GetUserAll: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getUserAll: can not starts transaction:%w", err)
	}
	var Users []model.ResponseUser
	var query string
	var pages int
	var rows *sql.Rows
	if page == 0 || limit == 0 {
		query = "SELECT id, email, role, created_at FROM users WHERE deleted = false ORDER BY id"
		rows, err = transaction.Query(query)
		if err != nil {
			u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
		}
		pages = 1
	} else {
		query = "SELECT id, email, role, created_at FROM users WHERE deleted = false ORDER BY id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, limit, (page-1)*limit)
		if err != nil {
			u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
		}
	}
	for rows.Next() {
		var User model.ResponseUser
		if err := rows.Scan(&User.ID, &User.Email, &User.Role, &User.CreatedAt); err != nil {
			u.logger.Errorf("Error while scanning for user:%s", err)
			return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
		}
		Users = append(Users, User)
	}
	if pages == 0 {
		query = "SELECT CEILING(COUNT(id)/$1::float) FROM users WHERE deleted = false"
		row := transaction.QueryRow(query, limit)
		if err := row.Scan(&pages); err != nil {
			u.logger.Errorf("Error while scanning for pages:%s", err)
		}
	}
	return Users, pages, transaction.Commit()
}
func (u *UserPostgres) GetUserByRoleFilter(page int, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error) {
	transaction, err := u.db.Begin()
	if err != nil {
		u.logger.Errorf("GetUserAll: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getUserAll: can not starts transaction:%w", err)
	}
	var Users []model.ResponseUser
	var pages int
	var rows *sql.Rows
	if page == 0 || limit == 0 {
		if filters.ShowDeleted {
			query := "SELECT id, email, role, created_at FROM users WHERE role = $1 ORDER BY id"
			rows, err = transaction.Query(query, filters.Role)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		} else {
			query := "SELECT id, email, role, created_at FROM users WHERE deleted = false AND role = $1 ORDER BY id"
			rows, err = transaction.Query(query, filters.Role)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		}
		pages = 1
	} else {
		if filters.ShowDeleted {
			query := "SELECT CEILING(COUNT(id)/$1::float) FROM users WHERE role = $2"
			row := transaction.QueryRow(query, limit, filters.Role)
			if err := row.Scan(&pages); err != nil {
				u.logger.Errorf("Error while scanning for pages:%s", err)
			}
			query2 := "SELECT id, email, role, created_at FROM users WHERE role = $1 ORDER BY id LIMIT $2 OFFSET $3"
			rows, err = transaction.Query(query2, filters.Role, limit, (page-1)*limit)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		} else {
			query := "SELECT CEILING(COUNT(id)/$1::float) FROM users WHERE role = $2 AND deleted = false"
			row := transaction.QueryRow(query, limit, filters.Role)
			if err := row.Scan(&pages); err != nil {
				u.logger.Errorf("Error while scanning for pages:%s", err)
			}
			query2 := "SELECT id, email, role, created_at FROM users WHERE role = $1 AND deleted = false ORDER BY id LIMIT $2 OFFSET $3"
			rows, err = transaction.Query(query2, filters.Role, limit, (page-1)*limit)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		}
	}
	for rows.Next() {
		var User model.ResponseUser
		if err := rows.Scan(&User.ID, &User.Email, &User.Role, &User.CreatedAt); err != nil {
			u.logger.Errorf("Error while scanning for user:%s", err)
			return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
		}
		Users = append(Users, User)
	}
	return Users, pages, transaction.Commit()
}

func (u *UserPostgres) GetUserByDataFilter(page int, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error) {
	transaction, err := u.db.Begin()
	if err != nil {
		u.logger.Errorf("GetUserAll: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("getUserAll: can not starts transaction:%w", err)
	}
	var Users []model.ResponseUser
	var pages int
	var rows *sql.Rows
	if page == 0 || limit == 0 {
		if filters.ShowDeleted {
			query := "SELECT id, email, role, created_at FROM users WHERE created_at >= $1 AND created_at <= $2 ORDER BY id"
			rows, err = transaction.Query(query, filters.StartTime, filters.EndTime)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		} else {
			query := "SELECT id, email, role, created_at FROM users WHERE created_at >= $1 AND created_at <= $2 AND deleted = false ORDER BY id"
			rows, err = transaction.Query(query, filters.StartTime, filters.EndTime)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		}
		pages = 1
	} else {
		if filters.ShowDeleted {
			query := "SELECT CEILING(COUNT(id)/$1::float) FROM users WHERE created_at >= $2 AND created_at <= $3"
			row := transaction.QueryRow(query, limit, filters.StartTime, filters.EndTime)
			if err := row.Scan(&pages); err != nil {
				u.logger.Errorf("Error while scanning for pages:%s", err)
			}
			query2 := "SELECT id, email, role, created_at FROM users WHERE created_at >= $1 AND created_at <= $2 ORDER BY id LIMIT $3 OFFSET $4"
			rows, err = transaction.Query(query2, filters.StartTime, filters.EndTime, limit, (page-1)*limit)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		} else {
			query2 := "SELECT CEILING(COUNT(id)/$1::float) FROM users WHERE created_at >= $2 AND created_at <= $3 AND deleted = false"
			row := transaction.QueryRow(query2, limit, filters.StartTime, filters.EndTime)
			if err := row.Scan(&pages); err != nil {
				u.logger.Errorf("Error while scanning for pages:%s", err)
			}
			query := "SELECT id, email, role, created_at FROM users WHERE created_at >= $1 AND created_at <= $2 AND deleted = false ORDER BY id LIMIT $3 OFFSET $4"
			rows, err = transaction.Query(query, filters.StartTime, filters.EndTime, limit, (page-1)*limit)
			if err != nil {
				u.logger.Errorf("GetUserAll: can not executes a query:%s", err)
				return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
			}
		}
	}
	for rows.Next() {
		var User model.ResponseUser
		if err := rows.Scan(&User.ID, &User.Email, &User.Role, &User.CreatedAt); err != nil {
			u.logger.Errorf("Error while scanning for user:%s", err)
			return nil, 0, fmt.Errorf("getUserAll:repository error:%w", err)
		}
		Users = append(Users, User)
	}
	return Users, pages, transaction.Commit()
}

// CreateStaff ...
func (u *UserPostgres) CreateStaff(user *model.CreateStaff) (int, error) {
	var id int
	row := u.db.QueryRow("INSERT INTO users (email, password, role, created_at, deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id", user.Email, user.Password, user.Role, time.Now().Format(model.Layout), false)
	if err := row.Scan(&id); err != nil {
		u.logger.Errorf("CreateStaff: error while scanning for user:%s", err)
		return 0, fmt.Errorf("CreateStaff: error while scanning for user:%w", err)
	}
	return id, nil
}

// CreateCustomer ...
func (u *UserPostgres) CreateCustomer(user *model.CreateCustomer) (int, error) {
	var id int
	row := u.db.QueryRow("INSERT INTO users (email, password, role, created_at, deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id", user.Email, user.Password, "Authorized Customer", time.Now().Format(model.Layout), false)
	if err := row.Scan(&id); err != nil {
		u.logger.Errorf("CreateCustomer: error while scanning for user:%s", err)
		return 0, fmt.Errorf("CreateCustomer: error while scanning for user:%w", err)
	}
	return id, nil
}

// UpdateUser ...
func (u *UserPostgres) UpdateUser(user *model.UpdateUser) error {
	_, err := u.db.Exec("UPDATE users SET password = $1 WHERE email = $2", user.NewPassword, user.Email)
	if err != nil {
		u.logger.Errorf("UpdateUser: error while updating user:%s", err)
		return fmt.Errorf("updateUser: error while updating user:%w", err)
	}
	return nil
}

// DeleteUserByID ...
func (u *UserPostgres) DeleteUserByID(id int) (int, error) {
	var userId int
	row := u.db.QueryRow("UPDATE users SET deleted = true WHERE id=$1 RETURNING id", id)
	if err := row.Scan(&userId); err != nil {
		u.logger.Errorf("DeleteUserByID: error while scanning for userId:%s", err)
		return 0, fmt.Errorf("deleteUserByID: error while scanning for userId:%w", err)
	}
	return userId, nil
}

// GetUserByEmail ...
func (u *UserPostgres) GetUserByEmail(email string) (*model.User, error) {
	var User model.User
	query := "SELECT id, email, password, role, deleted FROM users WHERE email = $1"
	row := u.db.QueryRow(query, email)
	if err := row.Scan(&User.ID, &User.Email, &User.Password, &User.Role, &User.Deleted); err != nil {
		u.logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("getUserByEmail: repository error:%w", err)

	}
	return &User, nil
}

// CheckEmail ...
func (u *UserPostgres) CheckEmail(email string) error {
	var exist bool
	query := "SELECT EXISTS (select 1 from users where email = $1)"
	row := u.db.QueryRow(query, email)
	if err := row.Scan(&exist); err != nil {
		u.logger.Errorf("Error while scanning for issued email:%s", err)
		return err
	}
	if !exist {
		u.logger.Error("user with this email does not exist")
		return pkg.ErrorEmailDoesNotExist
	}
	return nil
}
func (u *UserPostgres) RestorePassword(restore *model.RestorePassword) error {
	query := "UPDATE users SET password = $1 WHERE email=$2"
	_, err := u.db.Exec(query, restore.Password, restore.Email)
	if err != nil {
		return err
	}
	return nil
}
