package user

import "database/sql"

type UserRepository interface {
	Save(user User) error
	FindByEmail(email string) (User, error)
	FindById(userId int) (User, error)
}

type serRepository struct {
	db *sql.DB
}

func NewRepositoryUser(db *sql.DB) UserRepository {
	return &serRepository{
		db: db,
	}
}

func (r *serRepository) Save(user User) error {
	query := `
		INSERT INTO users (email, password) VALUES ($1, $2)
	`

	if _, err := r.db.Exec(query, user.Email, user.Password); err != nil {
		return err
	}

	return nil

}

func (r *serRepository) FindByEmail(email string) (User, error) {
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

func (r *serRepository) FindById(userId int) (User, error) {
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