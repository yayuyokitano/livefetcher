package queries

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
)

func GetUserLiveLists(ctx context.Context, userID int64, loggedInUser datastructures.AuthUser) (liveLists []datastructures.LiveList, err error) {
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

func GetLiveLiveLists(ctx context.Context, liveID int64, loggedInUser datastructures.AuthUser) (liveLists []datastructures.LiveList, err error) {
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

func UserOwnsLiveList(ctx context.Context, liveListID int64, loggedInUser datastructures.AuthUser) *logging.StatusError {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	defer counters.RollbackTransaction(ctx, tx)

	var userID int64
	err = tx.QueryRow(ctx, "SELECT users_id FROM livelists WHERE id=$1", liveListID).Scan(&userID)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	if userID != loggedInUser.ID {
		return logging.SE(http.StatusUnauthorized, errors.New("not owner of live list"))
	}
	return nil
}

func UserOwnsLiveListLive(ctx context.Context, liveListLiveID int64, loggedInUser datastructures.AuthUser) *logging.StatusError {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	defer counters.RollbackTransaction(ctx, tx)

	var userID int64
	err = tx.QueryRow(ctx, "SELECT users_id FROM livelistlives INNER JOIN livelists ON (livelists.id = livelistlives.livelists_id) WHERE livelistlives.id=$1", liveListLiveID).Scan(&userID)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	if userID != loggedInUser.ID {
		return logging.SE(http.StatusUnauthorized, errors.New("not owner of live list live"))
	}
	return nil
}

func GetLiveList(ctx context.Context, liveListID int64, loggedInUser datastructures.AuthUser) (liveList datastructures.LiveList, err error) {
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

	queryStr := `WITH queriedlives AS (
		SELECT livelistlive.id AS id, livelist.users_id AS livelist_owner_id, live_description, live.id AS live_id, live.title AS title, opentime, starttime, COALESCE(live.price,'') AS price, COALESCE(live.price_en,'') AS price_en, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url
		FROM lives AS live
		INNER JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		INNER JOIN livelistlives livelistlive ON (livelistlive.lives_id = live.id AND livelistlive.livelists_id=$1)
		INNER JOIN livelists livelist ON (livelistlive.livelists_id = livelist.id)
	)
	SELECT id, livelist_owner_id, live_description, live_id, array_agg(DISTINCT liveartists.artists_name), title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	FROM queriedlives
	INNER JOIN liveartists ON (liveartists.lives_id = queriedlives.live_id)
	GROUP BY id, livelist_owner_id, live_description, live_id, title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	ORDER BY starttime`

	rows, err := tx.Query(ctx, queryStr, liveList.ID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l datastructures.Live
		err = rows.Scan(&l.LiveListLiveID, &l.LiveListOwnerID, &l.Desc, &l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.PriceEnglish, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL)
		if err != nil {
			return
		}
		liveList.Lives = append(liveList.Lives, l)
	}
	for i, l := range liveList.Lives {
		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, tx, loggedInUser.ID, l.ID)
		if err == nil {
			liveList.Lives[i].FavoriteCount = int(favoriteCount)
			liveList.Lives[i].IsFavorited = isFavorited
		}
		err = nil
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func PostLiveList(ctx context.Context, livelist datastructures.LiveListWriteRequest) (id int64, err error) {
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

func DeleteLiveList(ctx context.Context, liveListID int64) (err error) {
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

func PostLiveListLive(ctx context.Context, liveListID int64, liveID int64, desc string) (err error) {
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

func DeleteLiveListLive(ctx context.Context, liveListLiveID int64) (err error) {
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

func FavoriteLiveList(ctx context.Context, userid int64, livelistid int64) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
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

func UnfavoriteLiveLiveList(ctx context.Context, userid int64, livelistid int64) (favoriteButtonInfo datastructures.FavoriteButtonInfo, err error) {
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

func getLiveListFavoriteAndCount(ctx context.Context, tx pgx.Tx, userid int64, livelistid int64) (isFavorited bool, favoriteCount int64, err error) {
	row := tx.QueryRow(ctx, "SELECT count(*) FROM livelistfavorites WHERE livelists_id=$1", livelistid)
	err = row.Scan(&favoriteCount)
	if err != nil || favoriteCount == 0 || userid == 0 {
		return
	}

	var selfCount int64
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
