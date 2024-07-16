package queries

import (
	"context"

	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

func GetUserByID(ctx context.Context, id int) (user util.User, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}
	row := tx.QueryRow(ctx, "SELECT id, email, username, nickname, password_hash, bio, location, is_verified FROM users WHERE id=$1", id)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Nickname, &user.PasswordHash, &user.Bio, &user.Location, &user.IsVerified)
	return
}

func GetUserByUsername(ctx context.Context, username string) (user util.User, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}
	row := tx.QueryRow(ctx, "SELECT id, email, username, nickname, password_hash, bio, location, is_verified FROM users WHERE username=$1", username)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Nickname, &user.PasswordHash, &user.Bio, &user.Location, &user.IsVerified)
	return
}

func GetUserByUsernameOrEmail(ctx context.Context, query string) (user util.User, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}
	row := tx.QueryRow(ctx, "SELECT id, email, username, nickname, password_hash, bio, location, is_verified FROM users WHERE username=$1 OR email=$1", query)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Nickname, &user.PasswordHash, &user.Bio, &user.Location, &user.IsVerified)
	return
}

func PostUser(ctx context.Context, user util.User) (err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	cmd, err := tx.Exec(ctx, "INSERT INTO users (email, username, nickname, password_hash, bio, location) VALUES ($1, $2, $3, $4, $5, $6)", user.Email, user.Username, user.Nickname, user.PasswordHash, user.Bio, user.Location)
	if err != nil {
		return
	}
	if cmd.RowsAffected() == 1 {
		logging.IncrementUsers()
	}

	err = counters.CommitTransaction(tx)
	return
}
