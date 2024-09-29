package adapters

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/zhikh23/pgutils"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type pgUsersRepository struct {
	db *sqlx.DB
}

func NewPGUsersRepository() (sm.UsersRepository, func() error) {
	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		panic("DATABASE_URL environment variable not set")
	}
	db := sqlx.MustConnect("postgres", uri)

	return &pgUsersRepository{db: db}, db.Close
}

func (r *pgUsersRepository) Save(ctx context.Context, user sm.User) error {
	if err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		return r.save(ctx, tx, user)
	}); pgutils.IsUniqueViolationError(err) {
		return sm.ErrUserAlreadyExists
	} else if err != nil {
		return err
	}
	return nil
}

func (r *pgUsersRepository) User(ctx context.Context, username string) (sm.User, error) {
	var user sm.User
	var err error
	if err = pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		user, err = r.user(ctx, tx, username)
		return err
	}); errors.Is(err, sql.ErrNoRows) {
		return sm.User{}, sm.ErrUserNotFound
	} else if err != nil {
		return sm.User{}, err
	}
	return user, nil
}

func (r *pgUsersRepository) save(ctx context.Context, ex sqlx.ExtContext, user sm.User) error {
	return r.requireExecResult(sqlx.NamedExecContext(ctx, ex,
		`INSERT INTO users (username, role) VALUES (:username, :role)`, marshallUserToRow(user),
	))
}

func (r *pgUsersRepository) user(ctx context.Context, qx sqlx.QueryerContext, username string) (sm.User, error) {
	var userRow userRow
	if err := sqlx.GetContext(ctx, qx, &userRow,
		`SELECT username, role FROM users WHERE username = $1`, username,
	); err != nil {
		return sm.User{}, err
	}
	return unmarshallUserFromRow(userRow)
}

func (r *pgUsersRepository) requireExecResult(res sql.Result, err error) error {
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type userRow struct {
	Username string `db:"username"`
	Role     string `db:"role"`
}

func marshallUserToRow(u sm.User) userRow {
	return userRow{
		Username: u.Username,
		Role:     u.Role.String(),
	}
}

func unmarshallUserFromRow(u userRow) (sm.User, error) {
	return sm.UnmarshallUserFromDB(u.Username, u.Role)
}
