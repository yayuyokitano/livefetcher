package queries

import (
	"context"
	"strings"

	"github.com/gojp/kana"
	"github.com/jackc/pgconn"
	"github.com/shogo82148/go-mecab"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
)

func PostArtists(ctx context.Context, artists []string) (n int64, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}
	for _, artist := range artists {
		var cmd pgconn.CommandTag
		cmd, err = tx.Exec(ctx, "INSERT INTO artists (name) VALUES ($1) ON CONFLICT DO NOTHING", artist)
		if err != nil {
			return
		}
		var tagger mecab.MeCab
		tagger, err = mecab.New(map[string]string{"output-format-type": "yomi"})
		if err != nil {
			return
		}
		defer tagger.Destroy()
		var katakana string
		katakana, err = tagger.ParseToString(artist)
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
		n += cmd.RowsAffected()
	}
	err = counters.CommitTransaction(tx)
	return
}
