package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
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

func deleteLives(tx pgx.Tx, ctx context.Context, lives []datastructures.Live) (deleted int, err error) {
	liveIds := make([]int, 0)
	for _, live := range lives {
		liveIds = append(liveIds, live.ID)
		err = notifyDeletedLive(ctx, tx, live)
		if err != nil {
			return
		}
	}

	cmd, err := tx.Exec(ctx, "DELETE FROM lives WHERE id=ANY($1)", liveIds)
	if err != nil {
		return
	}
	deleted = int(cmd.RowsAffected())
	return
}

func addLive(tx pgx.Tx, ctx context.Context, live datastructures.Live, artists *map[string]bool, liveartists *[][]interface{}) (added int, err error) {
	var liveid int
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
	live.ID = liveid

	added++
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

func tryUpdateLive(tx pgx.Tx, ctx context.Context, live datastructures.Live, oldLive datastructures.Live, artists *map[string]bool, liveartists *[][]interface{}) (modified int, err error) {
	if !shouldUpdateLive(live, oldLive) {
		return
	}

	_, err = tx.Exec(ctx, "UPDATE lives SET (title, opentime, starttime, url, price, price_en, livehouses_id) = ($1, $2, $3, $4, $5, $6, $7) WHERE id=$8", live.Title, live.OpenTime, live.StartTime, live.URL, live.Price, live.PriceEnglish, live.Venue.ID, oldLive.ID)
	if err != nil {
		return
	}
	modified = 1

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

func createNotification(ctx context.Context, tx pgx.Tx, liveId *int, nt datastructures.NotificationType) (notificationID int, err error) {
	err = tx.QueryRow(ctx, "INSERT INTO notifications (lives_id, notification_type) VALUES ($1, $2) RETURNING id", liveId, nt).Scan(&notificationID)
	return
}

func pushNotificationToUsers(ctx context.Context, tx pgx.Tx, notificationId int, userIds []int, notificationFields []datastructures.NotificationField) (err error) {
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

	notificationId, err := createNotification(ctx, tx, &oldLive.ID, datastructures.NotificationTypeEdited)
	if err != nil {
		return
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, userIds, notificationFields)
	return
}

func updateAndAddLives(tx pgx.Tx, ctx context.Context, lives []datastructures.Live, oldLives []datastructures.Live) (artists map[string]bool, liveartists [][]interface{}, added int, modified int, deleted int, err error) {
	oldLiveFoundIds := make(map[int]bool)
	liveartists = make([][]interface{}, 0)
	artists = make(map[string]bool)

	for _, live := range lives {
		foundLive := false
		for _, oldLive := range oldLives {
			if !isSameLive(live, oldLive, oldLives, lives) {
				continue
			}
			foundLive = true
			oldLiveFoundIds[oldLive.ID] = true
			m, _ := tryUpdateLive(tx, ctx, live, oldLive, &artists, &liveartists)
			modified += m
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
	oldLivesToDelete := make([]datastructures.Live, 0)
	for _, oldLive := range oldLives {
		if !oldLiveFoundIds[oldLive.ID] {
			oldLivesToDelete = append(oldLivesToDelete, oldLive)
		}
	}
	deleted, err = deleteLives(tx, ctx, oldLivesToDelete)

	return
}

func PostLives(ctx context.Context, lives []datastructures.Live, r *http.Request) (deleted int, added int, modified int, addedArtists int, err error) {
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

	oldLives, err := GetLives(ctx, LiveQuery{
		IncludeOldLives: true,
		LiveHouses:      livehouses,
	}, datastructures.AuthUser{}, r)
	if err != nil {
		return
	}

	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	artists, liveartists, added, modified, d, err := updateAndAddLives(tx, ctx, lives, oldLives.Lives)
	if err != nil {
		fmt.Println("updateandaddlives: ", err)
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

	liveids := make([]int, 0)
	for _, a := range liveartists {
		liveids = append(liveids, a[0].(int))
	}

	_, err = tx.Exec(ctx, "DELETE FROM liveartists WHERE lives_id=ANY($1)", liveids)
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
	Areas             map[int]bool    `form:"areas"`
	Artist            string          `form:"artist"`
	From              time.Time       `form:"from"`
	To                time.Time       `form:"to"`
	Id                int             `form:"id"`
	IncludeOldLives   bool            `form:"includeOldLives"`
	LiveHouses        []string        `form:"livehouses"`
	UserFavoritesId   int             `form:"userFavoritesId"`
	LiveListId        int             `form:"liveListId"`
	SavedSearchUserId int             `form:"savedSearchUserId"`
	Limit             int             `form:"limit"`
	Offset            int             `form:"offset"`
	AdditionalArtists map[string]bool `form:"additionalArtists"`
	AllowAllLocations bool            `form:"allowAllLocations"`
}

func GetLives(ctx context.Context, query LiveQuery, user datastructures.AuthUser, r *http.Request) (lives datastructures.Lives, err error) {
	localizer := i18nloader.GetLocalizer(r)
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	index := 0
	whereUsed := false
	incIndex := func() int {
		index++
		return index
	}
	args := []any{}

	queryStr := `WITH queriedlives AS ( SELECT live.id, array_agg(DISTINCT liveartists.artists_name) AS matching_artists, live.title AS live_title, opentime, starttime, COALESCE(live.price,'') AS price, COALESCE(live.price_en,'') AS price_en, livehouses_id, COALESCE(livehouse.url,'') AS livehouse_url, COALESCE(livehouse.description,'') AS livehouse_description, livehouse.areas_id AS areas_id, area.prefecture AS prefecture, area.name AS name, COALESCE(live.url,'') AS live_url, longitude, latitude, COALESCE(event.open_id, '') AS open_id, COALESCE(event.start_id, '') AS start_id, COUNT(*) OVER() AS count`
	if query.LiveListId != 0 {
		queryStr += ", livelistlive.id AS livelistlive_id, livelist.users_id AS livelist_owner_id, live_description"
	}

	queryStr += `
		FROM lives AS live
		LEFT JOIN liveartists ON (liveartists.lives_id = live.id)
		INNER JOIN livehouses livehouse ON (livehouse.id = live.livehouses_id)
		INNER JOIN areas area ON (area.id = livehouse.areas_id)
		LEFT JOIN artistaliases alias ON (alias.artists_name = liveartists.artists_name)
		`

	if query.UserFavoritesId != 0 {
		queryStr += `INNER JOIN userfavorites uf ON (uf.lives_id = live.id)
		`
	}

	if query.LiveListId != 0 {
		queryStr += fmt.Sprintf(`INNER JOIN livelistlives livelistlive ON (livelistlive.lives_id = live.id AND livelistlive.livelists_id=$%d)
			INNER JOIN livelists livelist ON (livelistlive.livelists_id = livelist.id)
		`, incIndex())
		args = append(args, query.LiveListId)
	}

	if query.SavedSearchUserId != 0 {
		queryStr += fmt.Sprintf(`
		INNER JOIN saved_searches ss ON (
			ss.users_id = $%d AND
			alias.alias ILIKE ss.keyword
		)
		`, incIndex())
		args = append(args, query.SavedSearchUserId)
	}
	queryStr += `LEFT JOIN calendarevents event ON (event.lives_id = live.id)`

	addCondition := func(condition string, conditionArgs ...any) {
		argInts := make([]any, 0)
		args = append(args, conditionArgs...)
		for i := 0; i < len(conditionArgs); i++ {
			argInts = append(argInts, incIndex())
		}
		if !whereUsed {
			queryStr += " WHERE " + fmt.Sprintf(condition, argInts...)
			whereUsed = true
		} else {
			queryStr += " AND " + fmt.Sprintf(condition, argInts...)
		}
	}

	if !query.IncludeOldLives {
		addCondition("starttime > NOW()")
	}

	if len(query.Areas) != 0 {
		var areaArray []int
		for area, isActive := range query.Areas {
			if isActive {
				areaArray = append(areaArray, area)
			}
		}
		addCondition("livehouse.areas_id = ANY($%d)", areaArray)
	}

	if query.Artist != "" {
		if query.Artist[0] == '"' && query.Artist[len(query.Artist)-1] == '"' {
			addCondition("alias.alias ILIKE $%d", query.Artist[1:len(query.Artist)-1])
		} else {
			addCondition("alias.alias ILIKE $%d", query.Artist+"%")
		}
	}

	if !query.From.IsZero() {
		addCondition("live.starttime >= $%d", query.From)
	}

	if !query.To.IsZero() {
		addCondition("live.starttime <= $%d", query.To)
	}

	if query.Id != 0 {
		addCondition("live.id = $%d", query.Id)
	}

	if len(query.LiveHouses) != 0 {
		addCondition("livehouses_id=ANY($%d)", query.LiveHouses)
	}

	if query.UserFavoritesId != 0 {
		addCondition("uf.users_id=$%d", query.UserFavoritesId)
	}

	queryStr += `
		GROUP BY live.id, live_title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url, latitude, longitude, open_id, start_id
		ORDER BY starttime, id`

	if query.Offset != 0 {
		queryStr += fmt.Sprintf(" OFFSET $%d", incIndex())
		args = append(args, query.Offset)
	}
	if query.Limit != 0 {
		queryStr += fmt.Sprintf(" LIMIT $%d", incIndex())
		args = append(args, query.Limit)
	}

	queryStr += `)
		SELECT live.id, matching_artists, array_agg(DISTINCT liveartists.artists_name) AS artists, live_title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url, longitude, latitude, open_id, start_id, count`
	if query.LiveListId != 0 {
		queryStr += ", livelistlive.id AS livelistlive_id, livelist.users_id AS livelist_owner_id, live_description"
	}
	queryStr += `
		FROM queriedlives AS live
		LEFT JOIN liveartists ON (liveartists.lives_id = live.id)
		LEFT JOIN artistaliases alias ON (alias.artists_name = liveartists.artists_name)
		GROUP BY live.id, matching_artists, live_title, opentime, starttime, price, price_en, livehouses_id, livehouse_url, livehouse_description, areas_id, prefecture, name, live_url, latitude, longitude, open_id, start_id, count`

	if query.LiveListId != 0 {
		queryStr += `, livelistlive_id, livelist_owner_id, live_description`
	}

	rows, err := tx.Query(ctx, queryStr, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l datastructures.Live
		var artists []*string
		var matchingArtists []*string
		scans := make([]any, 0)
		scans = append(scans, &l.ID, &matchingArtists, &artists, &l.Title, &l.OpenTime, &l.StartTime, &l.Price, &l.PriceEnglish, &l.Venue.ID, &l.Venue.Url, &l.Venue.Description, &l.Venue.Area.ID, &l.Venue.Area.Prefecture, &l.Venue.Area.Area, &l.URL, &l.Venue.Longitude, &l.Venue.Latitude, &l.CalendarOpenEventId, &l.CalendarStartEventId, &lives.Paginator.Total)
		if query.LiveListId != 0 {
			scans = append(scans, &l.LiveListLiveID, &l.LiveListOwnerID, &l.Desc)
		}
		err = rows.Scan(scans...)
		if err != nil {
			return
		}
		for _, a := range artists {
			if a != nil {
				l.Artists = append(l.Artists, *a)
			}
		}
		for _, ma := range matchingArtists {
			if ma != nil {
				l.MatchingArtists = append(l.MatchingArtists, *ma)
			}
		}
		lives.Lives = append(lives.Lives, l)
	}
	lives.Paginator.Offset = query.Offset
	lives.Paginator.Limit = query.Limit
	if query.Limit != 0 {
		lives.Paginator.Page = (query.Offset / query.Limit) + 1
	} else {
		lives.Paginator.Page = 1
	}

	if lives.Paginator.Total == 0 || lives.Paginator.Limit == 0 {
		lives.Paginator.TotalPages = 1
	} else {
		lives.Paginator.TotalPages = ((lives.Paginator.Total - 1) / lives.Paginator.Limit) + 1
	}

	for i, l := range lives.Lives {
		isFavorited, favoriteCount, err := getFavoriteAndCount(ctx, tx, user.ID, l.ID)
		if err == nil {
			lives.Lives[i].FavoriteCount = int(favoriteCount)
			lives.Lives[i].IsFavorited = isFavorited
		}
		err = nil

		lives.Lives[i].Venue.Name = localizer.Localize("livehouse." + lives.Lives[i].Venue.ID)
		lives.Lives[i].LocalizedTime = i18nloader.FormatOpenStartTime(lives.Lives[i].OpenTime, lives.Lives[i].StartTime, i18nloader.GetLanguages(r))
		lives.Lives[i].LocalizedPrice = lives.Lives[i].PriceEnglish
		for _, lang := range i18nloader.GetLanguages(r) {
			if strings.HasPrefix(lang, "ja") {
				lives.Lives[i].LocalizedPrice = lives.Lives[i].Price
				break
			}
			if strings.HasPrefix(lang, "en") {
				break
			}
		}
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

func GetLiveFavoritedUsers(ctx context.Context, tx pgx.Tx, liveID int) (userIDs []int, err error) {
	rows, err := tx.Query(ctx, "SELECT users_id FROM userfavorites WHERE lives_id = $1", liveID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userID int
		err = rows.Scan(&userID)
		if err != nil {
			return
		}
		userIDs = append(userIDs, userID)
	}
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

func createOldLiveNotificationFields(live datastructures.Live) (fields []datastructures.NotificationField, err error) {
	artists, err := json.Marshal(live.Artists)
	if err != nil {
		return
	}

	return []datastructures.NotificationField{{
		Type:     datastructures.NotificationFieldTitle,
		OldValue: live.Title,
	}, {
		Type:     datastructures.NotificationFieldOpenTime,
		OldValue: live.OpenTime.Format(time.RFC3339),
	}, {
		Type:     datastructures.NotificationFieldStartTime,
		OldValue: live.StartTime.Format(time.RFC3339),
	}, {
		Type:     datastructures.NotificationFieldPrice,
		OldValue: live.Price,
	}, {
		Type:     datastructures.NotificationFieldPriceEnglish,
		OldValue: live.PriceEnglish,
	}, {
		Type:     datastructures.NotificationFieldURL,
		OldValue: live.URL,
	}, {
		Type:     datastructures.NotificationFieldVenue,
		OldValue: live.Venue.ID,
	}, {
		Type:     datastructures.NotificationFieldArtists,
		OldValue: string(artists),
	}}, nil
}

func getUnnotifiedUsers(ctx context.Context, tx pgx.Tx, userIds []int, live datastructures.Live) (unnotifiedUserIds []int, err error) {
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
		var uid int
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
		fmt.Println("createnewlivenotificationfields:", err)
		return
	}

	notificationId, err := createNotification(ctx, tx, &live.ID, datastructures.NotificationTypeAdded)
	if err != nil {
		fmt.Println("createnotification:", err)
		return
	}

	userIds := make([]int, 0)
	for _, s := range ss {
		userIds = append(userIds, s.UserId)
	}

	unnotifiedUsers, err := getUnnotifiedUsers(ctx, tx, userIds, live)
	if err != nil {
		fmt.Println("getunnotifiedusers:", err)
		return
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, unnotifiedUsers, nf)
	fmt.Println("pushnotificationtousers:", err)
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

	notificationId, err := createNotification(ctx, tx, &live.ID, datastructures.NotificationTypeAdded)
	if err != nil {
		return
	}

	userIds := make([]int, 0)
	for _, s := range ss {
		userIds = append(userIds, s.UserId)
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, userIds, nf)
	return
}

func notifyDeletedLive(ctx context.Context, tx pgx.Tx, live datastructures.Live) (err error) {

	nf, err := createOldLiveNotificationFields(live)
	if err != nil {
		return
	}

	userIds, err := GetLiveFavoritedUsers(ctx, tx, live.ID)
	if err != nil {
		return
	}

	notificationId, err := createNotification(ctx, tx, nil, datastructures.NotificationTypeDeleted)
	if err != nil {
		return
	}

	err = pushNotificationToUsers(ctx, tx, notificationId, userIds, nf)
	return
}

func getMatchingSavedSearches(ctx context.Context, tx pgx.Tx, live datastructures.Live) (savedSearches []datastructures.SavedSearch, err error) {
	for _, artist := range live.Artists {
		var rows pgx.Rows
		rows, err = tx.Query(ctx, `
			SELECT users_id, keyword FROM saved_searches s
			LEFT JOIN user_saved_search_areas a ON s.users_id = a.users_id
			WHERE $1 ILIKE s.keyword AND (a.areas_id IS NULL OR a.areas_id = $2 OR s.allow_all_locations IS TRUE)
		`, artist, live.Venue.Area.ID)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var ss datastructures.SavedSearch
			err = rows.Scan(&ss.UserId, &ss.TextSearch)
			if err != nil {
				return
			}

			savedSearches = append(savedSearches, ss)
		}
	}
	return
}
