package queries

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

func isSameLive(live util.Live, oldLive util.Live, oldLives []util.Live, lives []util.Live) bool {
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

func addLive(tx pgx.Tx, ctx context.Context, live util.Live, artists *map[string]bool, liveartists *[][]interface{}) (added int64, err error) {
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
	liveartistmap := make(map[string]bool)
	for _, artist := range live.Artists {
		liveartistmap[artist] = true
		(*artists)[artist] = true
	}
	for k := range liveartistmap {
		*liveartists = append(*liveartists, []interface{}{liveid, k})
	}
	return
}

func arraysHaveSameItems(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	slices.Sort(a)
	slices.Sort(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func shouldUpdateLive(live util.Live, oldLive util.Live) bool {
	if live.Title != oldLive.Title {
		return true
	}
	if live.OpenTime != oldLive.OpenTime {
		return true
	}
	if live.StartTime != oldLive.StartTime {
		return true
	}
	if !arraysHaveSameItems(live.Artists, oldLive.Artists) {
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
	return false
}

func tryUpdateLive(tx pgx.Tx, ctx context.Context, live util.Live, oldLive util.Live, artists *map[string]bool, liveartists *[][]interface{}) (modified int64, err error) {
	if !shouldUpdateLive(live, oldLive) {
		return
	}

	cmd, err := tx.Exec(ctx, "UPDATE lives SET (title, opentime, starttime, url, price, price_en, livehouses_id) = ($1, $2, $3, $4, $5, $6, $7) WHERE id=$8", live.Title, live.OpenTime, live.StartTime, live.URL, live.Price, live.PriceEnglish, live.Venue.ID, oldLive.ID)
	if err != nil {
		return
	}
	modified = cmd.RowsAffected()

	liveartistmap := make(map[string]bool)
	for _, artist := range live.Artists {
		liveartistmap[artist] = true
		(*artists)[artist] = true
	}
	for k := range liveartistmap {
		*liveartists = append(*liveartists, []interface{}{oldLive.ID, k})
	}
	return
}

func updateAndAddLives(tx pgx.Tx, ctx context.Context, lives []util.Live, oldLives []util.Live) (artists map[string]bool, liveartists [][]interface{}, added int64, modified int64, deleted int64, err error) {
	oldLiveFoundIndices := make(map[int]bool)
	liveartists = make([][]interface{}, 0)
	artists = make(map[string]bool)

	for _, live := range lives {
		foundLive := false
		for oldLiveIndex, oldLive := range oldLives {
			if !isSameLive(live, oldLive, oldLives, lives) {
				continue
			}
			foundLive = true
			oldLiveFoundIndices[oldLiveIndex] = true
			m, err := tryUpdateLive(tx, ctx, live, oldLive, &artists, &liveartists)
			if err == nil {
				modified += m
			}
			break
		}
		if foundLive {
			continue
		}
		a, err := addLive(tx, ctx, live, &artists, &liveartists)
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

func PostLives(ctx context.Context, lives []util.Live) (deleted int64, added int64, modified int64, addedArtists int64, err error) {

	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	venues := make([]util.LiveHouse, 0)
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
	cmd, err := tx.Exec(ctx, "DELETE FROM lives WHERE livehouses_id=ANY($1) AND starttime NOT BETWEEN $2 AND $3", util.GetUniqueVenueIDs(venues), firstLive, lastLive)
	if err != nil {
		return
	}
	deleted += cmd.RowsAffected()

	oldLives, err := getLiveHouseLives(tx, livehouses)
	if err != nil {
		return
	}

	artists, liveartists, added, modified, d, err := updateAndAddLives(tx, ctx, lives, oldLives)
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

	liveIds := make([]int64, 0)
	for _, l := range lives {
		liveIds = append(liveIds, l.ID)
	}
	_, err = tx.Exec(ctx, "DELETE FROM liveartists WHERE lives_id=ANY($1)", liveIds)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO liveartists SELECT * FROM tmp_liveartists ON CONFLICT DO NOTHING")
	if err != nil {
		return
	}

	err = counters.CommitTransaction(tx)
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

func GetLives(query LiveQuery, user util.AuthUser) (lives []util.Live, err error) {
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
		args = append(args, query.Artist+"%")
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
	ctx := context.Background()
	rows, err := tx.Query(ctx, queryStr, args...)
	if err != nil {
		return
	}
	for rows.Next() {
		var l util.Live
		err = rows.Scan(&l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL, &l.Venue.Longitude, &l.Venue.Latitude)
		if err != nil {
			return
		}
		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, favoriteTx, user.ID, l.ID)
		if err == nil {
			l.FavoriteCount = int(favoriteCount)
			l.IsFavorited = isFavorited
		}
		err = nil
		lives = append(lives, l)
	}
	err = counters.CommitTransaction(tx)
	return
}

func getLiveHouseLives(tx pgx.Tx, livehouses []string) (lives []util.Live, err error) {
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
	rows, err := tx.Query(context.Background(), queryStr, livehouses)
	if err != nil {
		return
	}
	for rows.Next() {
		var l util.Live
		err = rows.Scan(&l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.PriceEnglish, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL)
		if err != nil {
			return
		}
		lives = append(lives, l)
	}
	return
}

func GetUserFavoriteLives(ctx context.Context, userid int64) (lives []util.Live, err error) {
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
	rows, err := tx.Query(context.Background(), queryStr, userid)
	if err != nil {
		return
	}
	for rows.Next() {
		var l util.Live
		err = rows.Scan(&l.ID, &l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.PriceEnglish, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL)
		if err != nil {
			return
		}

		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, favoriteTx, userid, l.ID)
		if err == nil {
			l.FavoriteCount = int(favoriteCount)
			l.IsFavorited = isFavorited
		}
		err = nil
		lives = append(lives, l)
	}
	err = counters.CommitTransaction(tx)
	return
}
