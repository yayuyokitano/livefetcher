package queries

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-playground/form"
	"github.com/jackc/pgx/v4"
	"github.com/yayuyokitano/livefetcher/lib/core/counters"
	"github.com/yayuyokitano/livefetcher/lib/core/logging"
	"github.com/yayuyokitano/livefetcher/lib/core/util"
)

func PostLives(ctx context.Context, lives []util.Live) (n int64, err error) {

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

	cmd, err := tx.Exec(ctx, "DELETE FROM lives WHERE livehouses_id=ANY($1)", util.GetUniqueVenueIDs(venues))
	if err != nil {
		return
	}
	n -= cmd.RowsAffected()

	liveartists := make([][]interface{}, 0)
	artists := make(map[string]bool)
	// this is slow, but I can't think of a better way to do it while getting the generated ids
	for _, live := range lives {
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
		n++
		liveartistmap := make(map[string]bool)
		for _, artist := range live.Artists {
			liveartistmap[artist] = true
			artists[artist] = true
		}
		for k := range liveartistmap {
			liveartists = append(liveartists, []interface{}{liveid, k})
		}
	}

	newartists := 0
	artistSlice := make([]string, 0)
	for artist := range artists {
		artistSlice = append(artistSlice, artist)
	}
	n, err = PostArtists(ctx, artistSlice)
	if err != nil {
		return
	}
	logging.AddArtists(newartists)

	fmt.Println(liveartists)

	o, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"liveartists"},
		[]string{"lives_id", "artists_name"},
		pgx.CopyFromRows(liveartists),
	)
	if err != nil {
		return
	}
	n += o
	logging.AddLives(int(n))

	err = counters.CommitTransaction(tx)
	return
}

type LiveQuery struct {
	Areas  map[int]bool
	Artist string
}

func GetLives(values url.Values) (lives []util.Live, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	queryStr := `WITH queriedlives AS (
		SELECT live.id AS id, title, opentime, starttime, COALESCE(live.price,'') AS price, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url
		FROM lives AS live
		INNER JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		INNER JOIN artistaliases alias ON (alias.artists_name = liveartists.artists_name)
		WHERE starttime > NOW()`
	args := []any{}

	decoder := form.NewDecoder()
	var query LiveQuery
	err = decoder.Decode(&query, values)
	if err != nil {
		return
	}

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

	queryStr += `)
	SELECT array_agg(DISTINCT liveartists.artists_name), title, opentime, starttime, price, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	FROM queriedlives
	INNER JOIN liveartists ON (liveartists.lives_id = queriedlives.id)
	GROUP BY title, opentime, starttime, price, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url
	ORDER BY starttime`
	rows, err := tx.Query(context.Background(), queryStr, args...)
	if err != nil {
		return
	}
	for rows.Next() {
		var l util.Live
		err = rows.Scan(&l.Artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL)
		if err != nil {
			return
		}
		lives = append(lives, l)
	}
	err = counters.CommitTransaction(tx)
	return
}
