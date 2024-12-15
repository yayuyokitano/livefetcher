package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
)

func isSameLive(live datastructures.Live, oldLive datastructures.Live, oldLives []datastructures.Live, lives []datastructures.Live) bool {
	if live.StartTime == oldLive.StartTime {
		return true
	}
	if live.Title == oldLive.Title {
		closest := live.StartTime.Sub(oldLive.StartTime).Abs()
		closestIndex := 0
		for i, l := range oldLives {
			cur := live.StartTime.Sub(l.StartTime).Abs()
			if closest < cur {
				continue
			}
			closestIndex = i
			closest = cur
		}
		if oldLives[closestIndex].ID != oldLive.ID {
			return false
		}

		closest = oldLive.StartTime.Sub(live.StartTime).Abs()
		closestIndex = 0
		for i, l := range lives {
			cur := oldLive.StartTime.Sub(l.StartTime).Abs()
			if closest < cur {
				continue
			}
			closestIndex = i
			closest = cur
		}
		if lives[closestIndex].ID == live.ID {
			return true
		}
	}

	return false
}

func deleteLives(tx pgx.Tx, ctx context.Context, liveIDs []int64) (deleted int64, err error) {
	cmd, err := tx.Exec(ctx, "DELETE FROM lives WHERE id=ANY($1)", liveIDs)
	if err != nil {
		return
	}
	deleted = cmd.RowsAffected()
	return
}

func addLive(tx pgx.Tx, ctx context.Context, live datastructures.Live, artists *map[string]bool, liveartists *[][]interface{}, editedLives *[]int64) (added int64, err error) {
	var liveid int64
	err = tx.QueryRow(
		ctx,
		"INSERT INTO lives (title, opentime, starttime, url, price, price_en, livehouses_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		live.Title,
		live.OpenTime,
		live.StartTime,
		live.URL,
		live.Price,
		live.PriceEnglish,
		live.Venue.ID,
	).Scan(&liveid)
	if err != nil {
		fmt.Println(err)
		return
	}

	added++
	*editedLives = append(*editedLives, liveid)
	liveartistmap := make(map[string]bool)
	for _, artist := range live.Artists {
		liveartistmap[artist] = true
		(*artists)[artist] = true
	}
	for k := range liveartistmap {
		*liveartists = append(*liveartists, []interface{}{liveid, k})
	}

	err = notifyNewLive(ctx, tx, live)
	return
}

func isArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	slices.Sort(a)
	slices.Sort(b)
	for i := range b {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func shouldUpdateLive(live datastructures.Live, oldLive datastructures.Live) bool {
	if live.Title != oldLive.Title {
		return true
	}
	if live.OpenTime != oldLive.OpenTime {
		return true
	}
	if live.StartTime != oldLive.StartTime {
		return true
	}
	if live.Price != oldLive.Price {
		return true
	}
	if live.PriceEnglish != oldLive.PriceEnglish {
		return true
	}
	if live.URL != oldLive.URL {
		return true
	}
	if live.Venue.ID != oldLive.Venue.ID {
		return true
	}
	if !isArrayEqual(live.Artists, oldLive.Artists) {
		return true
	}

	return false
}

func getNotificationFields(live, oldLive datastructures.Live) (fields []datastructures.NotificationField, err error) {
	oldLiveArtists, err := json.Marshal(oldLive.Artists)
	if err != nil {
		return
	}
	liveArtists, err := json.Marshal(live.Artists)
	if err != nil {
		return
	}

	return []datastructures.NotificationField{{
		Type:     datastructures.NotificationFieldTitle,
		OldValue: oldLive.Title,
		NewValue: live.Title,
	}, {
		Type:     datastructures.NotificationFieldOpenTime,
		OldValue: oldLive.OpenTime.Format(time.RFC3339),
		NewValue: live.OpenTime.Format(time.RFC3339),
	}, {
		Type:     datastructures.NotificationFieldStartTime,
		OldValue: oldLive.StartTime.Format(time.RFC3339),
		NewValue: live.StartTime.Format(time.RFC3339),
	}, {
		Type:     datastructures.NotificationFieldPrice,
		OldValue: oldLive.Price,
		NewValue: live.Price,
	}, {
		Type:     datastructures.NotificationFieldPriceEnglish,
		OldValue: oldLive.PriceEnglish,
		NewValue: live.PriceEnglish,
	}, {
		Type:     datastructures.NotificationFieldURL,
		OldValue: oldLive.URL,
		NewValue: live.URL,
	}, {
		Type:     datastructures.NotificationFieldVenue,
		OldValue: oldLive.Venue.ID,
		NewValue: live.Venue.ID,
	}, {
		Type:     datastructures.NotificationFieldArtists,
		OldValue: string(oldLiveArtists),
		NewValue: string(liveArtists),
	}}, nil
}

func tryUpdateLive(tx pgx.Tx, ctx context.Context, live datastructures.Live, oldLive datastructures.Live, artists *map[string]bool, liveartists *[][]interface{}) (modified int64, err error) {
	if !shouldUpdateLive(live, oldLive) {
		return
	}

	cmd, err := tx.Exec(ctx, "UPDATE lives SET (title, opentime, starttime, url, price, price_en, livehouses_id) = ($1, $2, $3, $4, $5, $6, $7) WHERE id=$8", live.Title, live.OpenTime, live.StartTime, live.URL, live.Price, live.PriceEnglish, live.Venue.ID, oldLive.ID)
	if err != nil {
		return
	}
	modified = cmd.RowsAffected()

	err = notifyUpdates(ctx, tx, oldLive, live)
	if err != nil {
		// ignore and log
		fmt.Println("notifyUpdates", err)
		err = nil
	}

	liveartistmap := make(map[string]bool)
	for _, artist := range live.Artists {
		liveartistmap[artist] = true
		(*artists)[artist] = true
	}
	for k := range liveartistmap {
		*liveartists = append(*liveartists, []interface{}{oldLive.ID, k})
	}

	err = notifyChangedLive(ctx, tx, live)
	return
}

func createNotification(ctx context.Context, tx pgx.Tx, liveId int64) (notificationID int64, err error) {
	err = tx.QueryRow(ctx, "INSERT INTO notifications (lives_id) VALUES ($1) RETURNING id", liveId).Scan(&notificationID)
	return
}

func pushNotificationToUsers(ctx context.Context, tx pgx.Tx, notificationId int64, userIds []int64, notificationFields []datastructures.NotificationField) (err error) {
	// this is massively inefficient, but for now i do not care
	// TODO: make this more efficient
	for _, userId := range userIds {
		_, err = tx.Exec(ctx, "INSERT INTO usernotifications (notifications_id, users_id) VALUES ($1, $2)", notificationId, userId)
		if err != nil {
			return
		}

		for _, notificationField := range notificationFields {
			_, err = tx.Exec(ctx, "INSERT INTO notification_fields (notifications_id, field_type, old_value, new_value) VALUES ($1, $2, $3, $4)", notificationId, notificationField.Type, notificationField.OldValue, notificationField.NewValue)
			if err != nil {
				return
			}
		}
	}
	return
}

func notifyUpdates(ctx context.Context, tx pgx.Tx, oldLive, live datastructures.Live) (err error) {
	userIds, err := GetLiveFavoritedUsers(ctx, tx, oldLive.ID)
	if err != nil || len(userIds) == 0 {
		return
	}

	notificationFields, err := getNotificationFields(live, oldLive)
	if err != nil {
		return
	}

	notificationId, err := createNotification(ctx, tx, oldLive.ID)
	if err != nil {
		return
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, userIds, notificationFields)
	return
}

func updateAndAddLives(tx pgx.Tx, ctx context.Context, lives []datastructures.Live, oldLives []datastructures.Live) (artists map[string]bool, liveartists [][]interface{}, editedLives []int64, added int64, modified int64, deleted int64, err error) {
	oldLiveFoundIndices := make(map[int]bool)
	liveartists = make([][]interface{}, 0)
	artists = make(map[string]bool)
	editedLives = make([]int64, 0)

	for _, live := range lives {
		foundLive := false
		for oldLiveIndex, oldLive := range oldLives {
			if !isSameLive(live, oldLive, oldLives, lives) {
				continue
			}
			foundLive = true
			oldLiveFoundIndices[oldLiveIndex] = true
			m, err := tryUpdateLive(tx, ctx, live, oldLive, &artists, &liveartists)
			if err == nil && m != 0 {
				editedLives = append(editedLives, oldLive.ID)
				modified += m
			}
			break
		}
		if foundLive {
			continue
		}
		a, err := addLive(tx, ctx, live, &artists, &liveartists, &editedLives)
		if err == nil {
			added += a
		}
	}
	oldLivesToDelete := make([]int64, 0)
	for i, oldLive := range oldLives {
		if !oldLiveFoundIndices[i] {
			oldLivesToDelete = append(oldLivesToDelete, oldLive.ID)
		}
	}
	deleted, err = deleteLives(tx, ctx, oldLivesToDelete)

	return
}

func PostLives(ctx context.Context, lives []datastructures.Live) (deleted int64, added int64, modified int64, addedArtists int64, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	venues := make([]datastructures.LiveHouse, 0)
	for _, live := range lives {
		venues = append(venues, live.Venue)
	}

	newlivehouses, err := PostLiveHouses(ctx, util.GetUniqueVenues(venues))
	if err != nil {
		return
	}
	logging.AddLiveHouses(newlivehouses)

	livehouses := util.GetUniqueVenueIDs(venues)
	var firstLive, lastLive time.Time
	for _, live := range lives {
		if firstLive.IsZero() {
			firstLive = live.StartTime
			lastLive = live.StartTime
			continue
		}
		if live.StartTime.Before(firstLive) {
			firstLive = live.StartTime
		}
		if live.StartTime.After(lastLive) {
			lastLive = live.StartTime
		}
	}
	cmd, err := tx.Exec(ctx, "DELETE FROM lives WHERE livehouses_id=ANY($1) AND starttime < NOW()", util.GetUniqueVenueIDs(venues))
	if err != nil {
		return
	}
	deleted += cmd.RowsAffected()

	oldLives, err := getLiveHouseLives(ctx, tx, livehouses)
	if err != nil {
		return
	}

	artists, liveartists, editedLives, added, modified, d, err := updateAndAddLives(tx, ctx, lives, oldLives)
	if err != nil {
		return
	}
	deleted += d

	newartists := 0
	artistSlice := make([]string, 0)
	for artist := range artists {
		artistSlice = append(artistSlice, artist)
	}
	addedArtists, err = PostArtists(ctx, artistSlice)
	if err != nil {
		return
	}
	logging.AddArtists(newartists)

	fmt.Println(liveartists)

	_, err = tx.Exec(ctx, "CREATE TEMP TABLE tmp_liveartists (LIKE liveartists INCLUDING DEFAULTS) ON COMMIT DROP")
	if err != nil {
		return
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_liveartists"},
		[]string{"lives_id", "artists_name"},
		pgx.CopyFromRows(liveartists),
	)
	if err != nil {
		return
	}

	fmt.Println(editedLives)
	_, err = tx.Exec(ctx, "DELETE FROM liveartists WHERE lives_id=ANY($1)", editedLives)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO liveartists SELECT * FROM tmp_liveartists ON CONFLICT DO NOTHING")
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	if err != nil {
		return
	}

	logging.AddLives(int(added - deleted))
	return
}

type LiveQuery struct {
	Areas  map[int]bool
	Artist string
	From   time.Time
	To     time.Time
}

func GetLives(ctx context.Context, query LiveQuery, user datastructures.AuthUser) (lives []datastructures.Live, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	queryStr := `WITH queriedlives AS (
		SELECT live.id AS id, title, opentime, starttime, COALESCE(live.price,'') AS price, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url, ST_X(location::geometry) AS longitude, ST_Y(location::geometry) AS latitude
		FROM lives AS live
		INNER JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		INNER JOIN artistaliases alias ON (alias.artists_name = liveartists.artists_name)
		WHERE starttime > NOW()`
	args := []any{}

	index := 0
	incIndex := func() int {
		index++
		return index
	}

	if len(query.Areas) != 0 {
		var areaArray []int
		for area, isActive := range query.Areas {
			if isActive {
				areaArray = append(areaArray, area)
			}
		}
		args = append(args, areaArray)
		queryStr += fmt.Sprintf(" AND livehouse.areas_id = ANY($%d)", incIndex())
	}

	if query.Artist != "" {
		if query.Artist[0] == '"' && query.Artist[len(query.Artist)-1] == '"' {
			args = append(args, query.Artist[1:len(query.Artist)-1])
		} else {
			args = append(args, query.Artist+"%")
		}

		queryStr += fmt.Sprintf(" AND alias.alias ILIKE $%d", incIndex())
	}

	if !query.From.IsZero() {
		args = append(args, query.From)
		queryStr += fmt.Sprintf(" AND live.starttime >= $%d", incIndex())
	}

	if !query.To.IsZero() {
		args = append(args, query.To)
		queryStr += fmt.Sprintf(" AND live.starttime <= $%d", incIndex())
	}

	queryStr += `)
	SELECT id, array_agg(DISTINCT liveartists.artists_name), title, opentime, starttime, price, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url, longitude, latitude
	FROM queriedlives
	INNER JOIN liveartists ON (liveartists.lives_id = queriedlives.id)
	GROUP BY id, title, opentime, starttime, price, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url, latitude, longitude
	ORDER BY starttime`
	rows, err := tx.Query(ctx, queryStr, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l datastructures.Live
		err = rows.Scan(&l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL, &l.Venue.Longitude, &l.Venue.Latitude)
		if err != nil {
			return
		}
		lives = append(lives, l)
	}

	for i, l := range lives {
		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, tx, user.ID, l.ID)
		if err == nil {
			lives[i].FavoriteCount = int(favoriteCount)
			lives[i].IsFavorited = isFavorited
		}
		err = nil
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func getLiveHouseLives(ctx context.Context, tx pgx.Tx, livehouses []string) (lives []datastructures.Live, err error) {
	queryStr := `WITH queriedlives AS (
		SELECT live.id AS id, title, opentime, starttime, COALESCE(live.price,'') AS price, COALESCE(live.price_en,'') AS price_en, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url
		FROM lives AS live
		INNER JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		WHERE livehouses_id=ANY($1)
	)
	SELECT id, array_agg(DISTINCT liveartists.artists_name), title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	FROM queriedlives
	INNER JOIN liveartists ON (liveartists.lives_id = queriedlives.id)
	GROUP BY id, title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	ORDER BY starttime`
	rows, err := tx.Query(ctx, queryStr, livehouses)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l datastructures.Live
		err = rows.Scan(&l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.PriceEnglish, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL)
		if err != nil {
			return
		}
		lives = append(lives, l)
	}
	return
}

func GetLiveFavoritedUsers(ctx context.Context, tx pgx.Tx, liveID int64) (userIDs []int64, err error) {
	rows, err := tx.Query(ctx, "SELECT users_id FROM userfavorites WHERE lives_id = $1", liveID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userID int64
		err = rows.Scan(&userID)
		if err != nil {
			return
		}
		userIDs = append(userIDs, userID)
	}
	return
}

func GetUserFavoriteLives(ctx context.Context, userid int64) (lives []datastructures.Live, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	queryStr := `WITH queriedlives AS (
		SELECT live.id AS id, title, opentime, starttime, COALESCE(live.price,'') AS price, COALESCE(live.price_en,'') AS price_en, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url
		FROM lives AS live
		INNER JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		INNER JOIN userfavorites uf ON (uf.lives_id = live.id)
		WHERE uf.users_id=$1
	)
	SELECT id, array_agg(DISTINCT liveartists.artists_name), title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	FROM queriedlives
	INNER JOIN liveartists ON (liveartists.lives_id = queriedlives.id)
	GROUP BY id, title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	ORDER BY starttime`
	rows, err := tx.Query(ctx, queryStr, userid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l datastructures.Live
		err = rows.Scan(&l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.PriceEnglish, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL)
		if err != nil {
			return
		}
		lives = append(lives, l)
	}
	for i, l := range lives {
		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, tx, userid, l.ID)
		if err == nil {
			lives[i].FavoriteCount = int(favoriteCount)
			lives[i].IsFavorited = isFavorited
		}
		err = nil
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func createNewLiveNotificationFields(live datastructures.Live) (fields []datastructures.NotificationField, err error) {
	artists, err := json.Marshal(live.Artists)
	if err != nil {
		return
	}

	return []datastructures.NotificationField{{
		Type:     datastructures.NotificationFieldTitle,
		NewValue: live.Title,
	}, {
		Type:     datastructures.NotificationFieldOpenTime,
		NewValue: live.OpenTime.Format(time.RFC3339),
	}, {
		Type:     datastructures.NotificationFieldStartTime,
		NewValue: live.StartTime.Format(time.RFC3339),
	}, {
		Type:     datastructures.NotificationFieldPrice,
		NewValue: live.Price,
	}, {
		Type:     datastructures.NotificationFieldPriceEnglish,
		NewValue: live.PriceEnglish,
	}, {
		Type:     datastructures.NotificationFieldURL,
		NewValue: live.URL,
	}, {
		Type:     datastructures.NotificationFieldVenue,
		NewValue: live.Venue.ID,
	}, {
		Type:     datastructures.NotificationFieldArtists,
		NewValue: string(artists),
	}}, nil
}

func getUnnotifiedUsers(ctx context.Context, tx pgx.Tx, userIds []int64, live datastructures.Live) (unnotifiedUserIds []int64, err error) {
	rows, err := tx.Query(ctx, `
	SELECT un.users_id FROM usernotifications un
	INNER JOIN notifications n ON n.id = un.notifications_id AND n.lives_id = $1
	WHERE un.users_id = ANY($2)
	`, live.ID, userIds)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid int64
		err = rows.Scan(&uid)
		if err != nil {
			return
		}
		unnotifiedUserIds = append(unnotifiedUserIds, uid)
	}
	return
}

func notifyChangedLive(ctx context.Context, tx pgx.Tx, live datastructures.Live) (err error) {
	ss, err := getMatchingSavedSearches(ctx, tx, live)
	if err != nil || len(ss) == 0 {
		return
	}

	nf, err := createNewLiveNotificationFields(live)
	if err != nil {
		return
	}

	notificationId, err := createNotification(ctx, tx, live.ID)
	if err != nil {
		return
	}

	userIds := make([]int64, 0)
	for _, s := range ss {
		userIds = append(userIds, s.UserId)
	}

	unnotifiedUsers, err := getUnnotifiedUsers(ctx, tx, userIds, live)
	if err != nil {
		return
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, unnotifiedUsers, nf)
	return
}

func notifyNewLive(ctx context.Context, tx pgx.Tx, live datastructures.Live) (err error) {
	ss, err := getMatchingSavedSearches(ctx, tx, live)
	if err != nil || len(ss) == 0 {
		return
	}

	nf, err := createNewLiveNotificationFields(live)
	if err != nil {
		return
	}

	notificationId, err := createNotification(ctx, tx, live.ID)
	if err != nil {
		return
	}

	userIds := make([]int64, 0)
	for _, s := range ss {
		userIds = append(userIds, s.UserId)
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, userIds, nf)
	return
}

func getMatchingSavedSearches(ctx context.Context, tx pgx.Tx, live datastructures.Live) (savedSearches []datastructures.SavedSearch, err error) {
	for _, artist := range live.Artists {
		var rows pgx.Rows
		rows, err = tx.Query(ctx, `
			SELECT id, users_id, text_search FROM saved_searches s
			LEFT JOIN saved_search_areas a ON s.id = a.saved_searches_id
			WHERE $1 ILIKE s.text_search AND (a.id IS NULL OR a.id = $2)
		`, artist, live.Venue.Area.ID)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var ss datastructures.SavedSearch
			err = rows.Scan(&ss.Id, ss.UserId, ss.TextSearch)
			if err != nil {
				return
			}

			savedSearches = append(savedSearches, ss)
		}
	}
	return
}
