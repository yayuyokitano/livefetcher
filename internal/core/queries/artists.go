package queries

import (
	"context"
	"strings"

	"github.com/gojp/kana"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util/mecab"
)

func PostArtists(ctx context.Context, artists []string) (n int, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	for _, artist := range artists {
		var cmd pgconn.CommandTag
		cmd, err = tx.Exec(ctx, "INSERT INTO artists (name) VALUES ($1) ON CONFLICT DO NOTHING", artist)
		if err != nil {
			return
		}
		var katakana string
		katakana, err = mecab.Mecab(artist)
		if err != nil {
			return
		}
		romaji := kana.KanaToRomaji(katakana)
		romajiSingleN := strings.ReplaceAll(romaji, "nn", "n")
		hiragana := kana.RomajiToHiragana(romaji)
		aliases := make(map[string]bool, 0)
		aliases[artist] = true
		aliases[katakana] = true
		aliases[romaji] = true
		aliases[romajiSingleN] = true
		aliases[hiragana] = true
		aliasSlice := make([]string, 0)
		for alias := range aliases {
			aliasSlice = append(aliasSlice, alias)
		}
		_, err = PushAliases(ctx, tx, artist, aliasSlice)
		if err != nil {
			return
		}
		n += int(cmd.RowsAffected())
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}
