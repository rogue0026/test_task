package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rogue0026/test_/internal/models"
	"github.com/rogue0026/test_/internal/storage"
)

type UsersRepository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (UsersRepository, error) {
	const fn = "storage.users.postgres.New"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return UsersRepository{}, fmt.Errorf("%s: %w", fn, err)
	}
	r := UsersRepository{
		pool: pool,
	}
	return r, nil
}

func (ur UsersRepository) GetUnregisteredUserByID(ctx context.Context, id string) (models.User, error) {
	const fn = "storage.users.postgres.GetUnregisteredUserByID"
	query := `SELECT * FROM get_unregistered_user_by_id(@id)`
	usr := models.User{}
	if err := ur.pool.QueryRow(ctx, query, pgx.NamedArgs{"id": id}).Scan(
		&usr.ID,
		&usr.Login,
		&usr.Name,
		&usr.Email,
		&usr.Password,
		&usr.IsVerified,
		&usr.VerificationCode,
		&usr.VerificationCodeExpires); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}
	return usr, nil
}

func (ur UsersRepository) SaveUnregisteredUser(ctx context.Context, usr models.User) (string, error) {
	const fn = "storage.users.postgres.SaveUnregisteredUser"
	if err := usr.HashPassword(); err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	args := pgx.NamedArgs{
		"login":                     usr.Login,
		"name":                      usr.Name,
		"email":                     usr.Email,
		"password":                  usr.Password,
		"verification_code":         usr.VerificationCode,
		"verification_code_expires": usr.VerificationCodeExpires,
	}
	query := `SELECT * FROM create_unregistered_user(@login, @name, @email, @password, @verification_code, @verification_code_expires)`
	row := ur.pool.QueryRow(ctx, query, args)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return userID, nil
}

func (ur UsersRepository) GetRegisteredUserByLogin(ctx context.Context, login string) (models.User, error) {
	const fn = "storage.users.postgres.GetRegisteredUserByLogin"
	query := `select * get_registered_user_by_login(@login)`
	user := models.User{}
	if err := ur.pool.QueryRow(ctx, query, pgx.NamedArgs{"login": login}).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.Email,
		&user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}
	return user, nil
}

func (ur UsersRepository) GetRegisteredUserByEmail(ctx context.Context, email string) (models.User, error) {
	const fn = "storage.users.postgres.GetRegisteredUserByEmail"
	query := `select * get_registered_user_by_email(@email)`
	user := models.User{}
	if err := ur.pool.QueryRow(ctx, query, pgx.NamedArgs{"email": email}).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.Email,
		&user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}
	return user, nil
}

func (ur UsersRepository) SaveRegisteredUser(ctx context.Context, u models.User) error {
	const fn = "storage.users.postgres.SaveRegisteredUser"
	query := `call register_user(@id, @login, @name, @email, @password)`
	args := pgx.NamedArgs{
		"id":       u.ID,
		"login":    u.Login,
		"name":     u.Name,
		"email":    u.Email,
		"password": u.Password,
	}
	_, err := ur.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}
