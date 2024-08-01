package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func getFavoriteAndCount(ctx context.Context, tx pgx.Tx, userid int64, liveid int64) (isFavorited bool, favoriteCount int64, err error) {
	row := tx.QueryRow(ctx, "SELECT count(*) FROM userfavorites WHERE lives_id=$1", liveid)
	err = row.Scan(&favoriteCount)
	if err != nil || favoriteCount == 0 || userid == 0 {
		return
	}

	var selfCount int64
	selfRow := tx.QueryRow(ctx, "SELECT count(*) FROM userfavorites WHERE lives_id=$1 AND users_id=$2", liveid, userid)
	err = selfRow.Scan(&selfCount)
	if err != nil {
		return
	}
	if selfCount > 0 {
		isFavorited = true
	}
	return
}

func FavoriteLive(ctx context.Context, userid int64, liveid int64) (favoriteButtonInfo util.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO userfavorites (users_id, lives_id) VALUES ($1, $2)", userid, liveid)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, tx, userid, liveid)
	if err != nil {
		return
	}
	favoriteButtonInfo = util.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		Id:            int(liveid),
	}

	err = counters.CommitTransaction(tx)
	return
}

func UnfavoriteLive(ctx context.Context, userid int64, liveid int64) (favoriteButtonInfo util.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM userfavorites WHERE users_id=$1 AND lives_id=$2", userid, liveid)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, tx, userid, liveid)
	if err != nil {
		return
	}
	favoriteButtonInfo = util.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		Id:            int(liveid),
	}

	err = counters.CommitTransaction(tx)
	return
}
