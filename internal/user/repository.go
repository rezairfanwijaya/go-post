package user

import (
	"database/sql"
)

type UserRepository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindById(userId int) (User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewRepositoryUser(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(user User) (User, error) {
	res := User{}
	query := `
		INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *
	`
	row := r.db.QueryRow(query, user.Email, user.Password)
	if err := row.Scan(&res.Id, &res.Email, &res.Password); err != nil {
		return res, err
	}

	return res, nil
}

func (r *userRepository) FindByEmail(email string) (User, error) {
	var user User

	query := `
		SELECT * FROM users WHERE email = $1
	`

	row := r.db.QueryRow(query, email)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindById(userId int) (User, error) {
	var user User

	query := `
		SELECT * FROM users WHERE id = $1
	`

	row := r.db.QueryRow(query, userId)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
