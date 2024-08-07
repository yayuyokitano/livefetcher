package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

func GetUserLiveLists(ctx context.Context, userID int64, loggedInUser util.AuthUser) (liveLists []util.LiveList, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	favoriteTx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(favoriteTx)
	if err != nil {
		return
	}

	rows, err := tx.Query(ctx, `SELECT livelist.id, title, list_description, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		WHERE users_id=$1`, userID)
	if err != nil {
		return
	}

	for rows.Next() {
		var ll util.LiveList
		err = rows.Scan(&ll.ID, &ll.Title, &ll.Desc, &ll.User.ID, &ll.User.Username, &ll.User.Nickname, &ll.User.Location)
		if err != nil {
			return
		}
		isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, favoriteTx, loggedInUser.ID, ll.ID)
		if err == nil {
			ll.FavoriteCount = int(favoriteCount)
			ll.IsFavorited = isFavorited
		}
		err = nil
		liveLists = append(liveLists, ll)
	}

	err = counters.CommitTransaction(tx)
	return
}

func GetLiveLiveLists(ctx context.Context, liveID int64, loggedInUser util.AuthUser) (liveLists []util.LiveList, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	favoriteTx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(favoriteTx)
	if err != nil {
		return
	}

	rows, err := tx.Query(ctx, `SELECT livelist.id, title, list_description, live_description, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		INNER JOIN livelistlives ON (livelist.id = livelistlives.livelists_id)
		WHERE livelistlives.lives_id=$1`, liveID)
	if err != nil {
		return
	}

	for rows.Next() {
		var ll util.LiveList
		err = rows.Scan(&ll.ID, &ll.Title, &ll.Desc, &ll.LiveDesc, &ll.User.ID, &ll.User.Username, &ll.User.Nickname, &ll.User.Location)
		if err != nil {
			return
		}
		isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, favoriteTx, loggedInUser.ID, ll.ID)
		if err == nil {
			ll.FavoriteCount = int(favoriteCount)
			ll.IsFavorited = isFavorited
		}
		err = nil
		liveLists = append(liveLists, ll)
	}

	err = counters.CommitTransaction(tx)
	return
}

func GetLiveList(ctx context.Context, liveListID int64, loggedInUser util.AuthUser) (liveList util.LiveList, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	favoriteTx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(favoriteTx)
	if err != nil {
		return
	}

	err = tx.QueryRow(ctx, `SELECT livelist.id, title, list_description, livelist.created_at, livelist.updated_at, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		WHERE livelist.id=$1`, liveListID).Scan(&liveList.ID, &liveList.Title, &liveList.Desc, &liveList.CreatedAt, &liveList.UpdatedAt, &liveList.User.ID, &liveList.User.Username, &liveList.User.Nickname, &liveList.User.Location)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, favoriteTx, loggedInUser.ID, liveList.ID)
	if err == nil {
		liveList.FavoriteCount = int(favoriteCount)
		liveList.IsFavorited = isFavorited
	}
	err = nil

	queryStr := `WITH queriedlives AS (
		SELECT livelistlive.id AS id, live_description, live.id AS live_id, title, opentime, starttime, COALESCE(live.price,'') AS price, COALESCE(live.price_en,'') AS price_en, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url
		FROM lives AS live
		INNER JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		INNER JOIN livelistlives livelistlive ON (livelistlive.lives_id = live.id)
		WHERE livelistlive.livelists_id=$1
	)
	SELECT id, live_description, live_id, array_agg(DISTINCT liveartists.artists_name), title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	FROM queriedlives
	INNER JOIN liveartists ON (liveartists.lives_id = queriedlives.id)
	GROUP BY id, live_description, live_id, title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	ORDER BY starttime`

	rows, err := tx.Query(ctx, queryStr, liveList.ID)
	for rows.Next() {
		var l util.LiveListLive
		err = rows.Scan(&l.ID, &l.Desc, &l.Live.Artists, &l.Live.Title, &l.Live.OpenTime, &l.Live.StartTime, &l.Live.Price, &l.Live.PriceEnglish, &l.Live.Venue.ID, &l.Live.Venue.Url, &l.Live.Venue.Description, &l.Live.Venue.Area.ID, &l.Live.Venue.Area.Prefecture, &l.Live.Venue.Area.Area, &l.Live.URL)
		if err != nil {
			return
		}
		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, favoriteTx, loggedInUser.ID, l.Live.ID)
		if err == nil {
			l.Live.FavoriteCount = int(favoriteCount)
			l.Live.IsFavorited = isFavorited
		}
		err = nil
		liveList.Lives = append(liveList.Lives, l)
	}

	err = counters.CommitTransaction(tx)
	return
}

func PostLiveList(ctx context.Context, livelist util.LiveListWriteRequest) (id int64, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	err = tx.QueryRow(ctx, "INSERT INTO livelists (users_id, title, list_description) VALUES ($1, $2, $3) RETURNING id", livelist.UserID, livelist.Title, livelist.Desc).Scan(&id)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(tx)
	return
}

func PutLiveList(ctx context.Context, livelist util.LiveListWriteRequest) (err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "UPDATE livelists SET title=$1, list_description=$2 WHERE id=$3", livelist.Title, livelist.Desc, livelist.ID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(tx)
	return
}

func DeleteLiveList(ctx context.Context, liveListID int64) (err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM livelists WHERE id=$1", liveListID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(tx)
	return
}

func PostLiveListLive(ctx context.Context, liveListID int64, liveID int64, desc string) (err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO livelistlives (livelists_id, lives_id, live_description) VALUES ($1, $2, $3)", liveListID, liveID, desc)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(tx)
	return
}

func DeleteLiveListLive(ctx context.Context, liveListLiveID int64) (err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM livelistlives WHERE id=$1", liveListLiveID)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(tx)
	return
}

func FavoriteLiveList(ctx context.Context, userid int64, livelistid int64) (favoriteButtonInfo util.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO livelistfavorites (users_id, livelists_id) VALUES ($1, $2)", userid, livelistid)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, userid, livelistid)
	if err != nil {
		return
	}
	favoriteButtonInfo = util.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		ID:            int(livelistid),
	}

	err = counters.CommitTransaction(tx)
	return
}

func UnfavoriteLiveLiveList(ctx context.Context, userid int64, livelistid int64) (favoriteButtonInfo util.FavoriteButtonInfo, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM livelistfavorites WHERE users_id=$1 AND livelists_id=$2", userid, livelistid)
	if err != nil {
		return
	}
	isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, tx, userid, livelistid)
	if err != nil {
		return
	}
	favoriteButtonInfo = util.FavoriteButtonInfo{
		IsFavorited:   isFavorited,
		FavoriteCount: int(favoriteCount),
		ID:            int(livelistid),
	}

	err = counters.CommitTransaction(tx)
	return
}

func GetUserFavoriteLiveLists(ctx context.Context, user util.AuthUser) (liveLists []util.LiveList, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	favoriteTx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(favoriteTx)
	if err != nil {
		return
	}

	rows, err := tx.Query(ctx, `SELECT livelist.id, title, list_description, users.id, users.username, users.nickname, users.location
		FROM livelists as livelist
		INNER JOIN users ON (livelist.users_id = users.id)
		INNER JOIN livelistfavorites lf ON (livelist.users_id = lf.users_id)
		WHERE users_id=$1`, user.ID)
	if err != nil {
		return
	}

	for rows.Next() {
		var ll util.LiveList
		err = rows.Scan(&ll.ID, &ll.Title, &ll.Desc, &ll.User.ID, &ll.User.Username, &ll.User.Nickname, &ll.User.Location)
		if err != nil {
			return
		}
		isFavorited, favoriteCount, err := getLiveListFavoriteAndCount(ctx, favoriteTx, user.ID, ll.ID)
		if err == nil {
			ll.FavoriteCount = int(favoriteCount)
			ll.IsFavorited = isFavorited
		}
		err = nil
		liveLists = append(liveLists, ll)
	}

	err = counters.CommitTransaction(tx)
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
