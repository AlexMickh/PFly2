package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AlexMickh/PFly2/internal/models"
	"github.com/AlexMickh/PFly2/internal/storage"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(
	host string,
	port int,
	user string,
	password string,
	dbname string,
) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(user models.User) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (name, email, password, image_url, description, interests) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(user.Name, user.Email, user.Password, user.ImageUrl, user.Description, pq.Array(user.Interests))
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return 1, nil
}

func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	const op = "storage.postgres.GetUserByEmail"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	err = stmt.QueryRow(email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.ImageUrl, &user.Description, pq.Array(&user.Interests))
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, storage.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
