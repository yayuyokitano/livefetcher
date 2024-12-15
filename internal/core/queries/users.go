package queries

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
)

func GetUserByID(ctx context.Context, id int) (user datastructures.User, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}
	return getUserByIDQuery(ctx, tx, id)
}

func getUserByIDQuery(ctx context.Context, tx pgx.Tx, id int) (user datastructures.User, err error) {
	row := tx.QueryRow(ctx, "SELECT id, email, username, nickname, password_hash, bio, location, is_verified FROM users WHERE id=$1", id)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Nickname, &user.PasswordHash, &user.Bio, &user.Location, &user.IsVerified)
	return
}

func GetUserByUsername(ctx context.Context, username string) (user datastructures.User, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}
	row := tx.QueryRow(ctx, "SELECT id, email, username, nickname, password_hash, bio, location, is_verified FROM users WHERE username=$1", username)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Nickname, &user.PasswordHash, &user.Bio, &user.Location, &user.IsVerified)
	return
}

func GetUserByUsernameOrEmail(ctx context.Context, query string) (user datastructures.User, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}
	row := tx.QueryRow(ctx, "SELECT id, email, username, nickname, password_hash, bio, location, is_verified FROM users WHERE username=$1 OR email=$1", query)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Nickname, &user.PasswordHash, &user.Bio, &user.Location, &user.IsVerified)
	return
}

func PostUser(ctx context.Context, user datastructures.User) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
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

	err = counters.CommitTransaction(ctx, tx)
	return
}

func PatchUser(ctx context.Context, patchInfo datastructures.User) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	user, err := getUserByIDQuery(ctx, tx, int(patchInfo.ID))
	if err != nil {
		return
	}
	if patchInfo.ID != user.ID {
		return errors.New("user not found")
	}

	if patchInfo.Username != "" {
		user.Username = patchInfo.Username
	}
	if patchInfo.Nickname != "" {
		user.Nickname = patchInfo.Nickname
	}
	if patchInfo.Email != "" && patchInfo.Email != user.Email {
		user.Email = patchInfo.Email
		user.IsVerified = false
	}
	if patchInfo.PasswordHash != "" {
		user.PasswordHash = patchInfo.PasswordHash
	}
	if patchInfo.Bio != "" {
		user.Bio = patchInfo.Bio
	}
	if patchInfo.Location != "" {
		user.Location = patchInfo.Location
	}

	_, err = tx.Exec(ctx, "UPDATE users SET email=$1, username=$2, nickname=$3, password_hash=$4, bio=$5, location=$6, is_verified=$7 WHERE id=$8", user.Email, user.Username, user.Nickname, user.PasswordHash, user.Bio, user.Location, user.IsVerified, user.ID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
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

func FavoriteLive(ctx context.Context, userid int64, liveid int64) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
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
	favoriteButtonInfo = datastructures.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		ID:            int(liveid),
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func UnfavoriteLive(ctx context.Context, userid int64, liveid int64) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
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
	favoriteButtonInfo = datastructures.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		ID:            int(liveid),
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func PostSavedSearch(ctx context.Context, userid int64, search string, areaIds []int64) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	var searchId int64
	if search[0] == '"' && search[len(search)-1] == '"' {
		err = tx.QueryRow(ctx, "INSERT INTO saved_searches (users_id, text_search) VALUES ($1, $2) RETURNING id", userid, search[1:len(search)-1]).Scan(&searchId)
		if err != nil {
			return
		}
	} else {
		err = tx.QueryRow(ctx, "INSERT INTO saved_searches (users_id, text_search) VALUES ($1, $2) RETURNING id", userid, search+"%").Scan(&searchId)
		if err != nil {
			return
		}
	}

	for _, areaId := range areaIds {
		_, err = tx.Exec(ctx, "INSERT INTO saved_search_areas (saved_searches_id, areas_id) VALUES ($1, $2)", searchId, areaId)
		if err != nil {
			return
		}
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}
