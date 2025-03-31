package queries

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func GetUserLiveLists(ctx context.Context, userID int, loggedInUser datastructures.AuthUser) (liveLists []datastructures.LiveList, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	rows, err := tx.Query(ctx, `SELECT livelist.id, title, list_description, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		WHERE users_id=$1`, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ll datastructures.LiveList
		err = rows.Scan(&ll.ID, &ll.Title, &ll.Desc, &ll.User.ID, &ll.User.Username, &ll.User.Nickname, &ll.User.Location)
		if err != nil {
			return
		}
		liveLists = append(liveLists, ll)
	}
	for i, ll := range liveLists {
		isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, loggedInUser.ID, ll.ID)
		if err == nil {
			liveLists[i].FavoriteCount = int(favoriteCount)
			liveLists[i].IsFavorited = isFavorited
		}
		err = nil
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func GetLiveLiveLists(ctx context.Context, liveID int, loggedInUser datastructures.AuthUser) (liveLists []datastructures.LiveList, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	rows, err := tx.Query(ctx, `SELECT livelist.id, title, list_description, live_description, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		INNER JOIN livelistlives ON (livelist.id = livelistlives.livelists_id)
		WHERE livelistlives.lives_id=$1`, liveID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ll datastructures.LiveList
		err = rows.Scan(&ll.ID, &ll.Title, &ll.Desc, &ll.LiveDesc, &ll.User.ID, &ll.User.Username, &ll.User.Nickname, &ll.User.Location)
		if err != nil {
			return
		}
		liveLists = append(liveLists, ll)
	}
	for i, ll := range liveLists {
		isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, loggedInUser.ID, ll.ID)
		if err == nil {
			liveLists[i].FavoriteCount = int(favoriteCount)
			liveLists[i].IsFavorited = isFavorited
		}
		err = nil
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func UserOwnsLiveList(ctx context.Context, r *http.Request, liveListID int, loggedInUser datastructures.AuthUser) *logging.StatusError {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}
	defer counters.RollbackTransaction(ctx, tx)

	var userID int
	err = tx.QueryRow(ctx, "SELECT users_id FROM livelists WHERE id=$1", liveListID).Scan(&userID)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	if userID != loggedInUser.ID {
		return logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.action-not-permitted"))
	}
	return nil
}

func UserOwnsLiveListLive(ctx context.Context, r *http.Request, liveListLiveID int, loggedInUser datastructures.AuthUser) *logging.StatusError {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}
	defer counters.RollbackTransaction(ctx, tx)

	var userID int
	err = tx.QueryRow(ctx, "SELECT users_id FROM livelistlives INNER JOIN livelists ON (livelists.id = livelistlives.livelists_id) WHERE livelistlives.id=$1", liveListLiveID).Scan(&userID)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	if userID != loggedInUser.ID {
		return logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.action-not-permitted"))
	}
	return nil
}

func GetLiveList(ctx context.Context, liveListID int, loggedInUser datastructures.AuthUser) (liveList datastructures.LiveList, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	err = tx.QueryRow(ctx, `SELECT livelist.id, title, list_description, livelist.created_at, livelist.updated_at, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		WHERE livelist.id=$1`, liveListID).Scan(&liveList.ID, &liveList.Title, &liveList.Desc, &liveList.CreatedAt, &liveList.UpdatedAt, &liveList.User.ID, &liveList.User.Username, &liveList.User.Nickname, &liveList.User.Location)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, loggedInUser.ID, liveList.ID)
	if err == nil {
		liveList.FavoriteCount = int(favoriteCount)
		liveList.IsFavorited = isFavorited
	}
	err = nil

	err = counters.CommitTransaction(ctx, tx)
	if err != nil {
		return
	}

	liveListLives, err := GetLives(ctx, LiveQuery{
		LiveListId: liveListID,
	}, loggedInUser)
	if err != nil {
		return
	}

	liveList.Lives = liveListLives.Lives
	return
}

func PostLiveList(ctx context.Context, livelist datastructures.LiveListWriteRequest) (id int, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	err = tx.QueryRow(ctx, "INSERT INTO livelists (users_id, title, list_description) VALUES ($1, $2, $3) RETURNING id", livelist.UserID, livelist.Title, livelist.Desc).Scan(&id)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func PutLiveList(ctx context.Context, livelist datastructures.LiveListWriteRequest) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "UPDATE livelists SET title=$1, list_description=$2 WHERE id=$3", livelist.Title, livelist.Desc, livelist.ID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func DeleteLiveList(ctx context.Context, liveListID int) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "DELETE FROM livelists WHERE id=$1", liveListID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func PostLiveListLive(ctx context.Context, liveListID int, liveID int, desc string) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "INSERT INTO livelistlives (livelists_id, lives_id, live_description) VALUES ($1, $2, $3)", liveListID, liveID, desc)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func DeleteLiveListLive(ctx context.Context, liveListLiveID int) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "DELETE FROM livelistlives WHERE id=$1", liveListLiveID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func FavoriteLiveList(ctx context.Context, userid int, livelistid int) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "INSERT INTO livelistfavorites (users_id, livelists_id) VALUES ($1, $2)", userid, livelistid)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, userid, livelistid)
	if err != nil {
		return
	}
	favoriteButtonInfo = datastructures.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		ID:            int(livelistid),
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func UnfavoriteLiveLiveList(ctx context.Context, userid int, livelistid int) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "DELETE FROM livelistfavorites WHERE users_id=$1 AND livelists_id=$2", userid, livelistid)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, userid, livelistid)
	if err != nil {
		return
	}
	favoriteButtonInfo = datastructures.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		ID:            int(livelistid),
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func GetUserFavoriteLiveLists(ctx context.Context, user datastructures.AuthUser) (liveLists []datastructures.LiveList, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	// TODO: only use one transaction
	favoriteTx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, favoriteTx)

	rows, err := tx.Query(ctx, `SELECT livelist.id, title, list_description, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		INNER JOIN livelistfavorites lf ON (livelist.users_id = lf.users_id)
		WHERE users_id=$1`, user.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ll datastructures.LiveList
		err = rows.Scan(&ll.ID, &ll.Title, &ll.Desc, &ll.User.ID, &ll.User.Username, &ll.User.Nickname, &ll.User.Location)
		if err != nil {
			return
		}
		liveLists = append(liveLists, ll)
	}

	for i, ll := range liveLists {
		isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, favoriteTx, user.ID, ll.ID)
		if err == nil {
			liveLists[i].FavoriteCount = int(favoriteCount)
			liveLists[i].IsFavorited = isFavorited
		}
		err = nil
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func getLiveListFavoriteAndCount(ctx context.Context, tx pgx.Tx, userid int, livelistid int) (isFavorited bool, favoriteCount int, err error) {
	row := tx.QueryRow(ctx, "SELECT count(*) FROM livelistfavorites WHERE livelists_id=$1", livelistid)
	err = row.Scan(&favoriteCount)
	if err != nil || favoriteCount == 0 || userid == 0 {
		return
	}

	var selfCount int
	selfRow := tx.QueryRow(ctx, "SELECT count(*) FROM livelistfavorites WHERE livelists_id=$1 AND users_id=$2", livelistid, userid)
	err = selfRow.Scan(&selfCount)
	if err != nil {
		return
	}
	if selfCount > 0 {
		isFavorited = true
	}
	return
}
