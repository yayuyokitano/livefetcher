package counters

import "context"

func GetLiveCount(ctx context.Context) (n int64, err error) {
	tx, err := FetchTransaction()
	defer RollbackTransaction(tx)
	if err != nil {
		return
	}
	err = tx.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM lives",
	).Scan(&n)
	if err != nil {
		return
	}
	err = CommitTransaction(tx)
	return
}

func GetLiveHouseCount(ctx context.Context) (n int64, err error) {
	tx, err := FetchTransaction()
	defer RollbackTransaction(tx)
	if err != nil {
		return
	}
	err = tx.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM livehouses",
	).Scan(&n)
	if err != nil {
		return
	}
	err = CommitTransaction(tx)
	return
}

func GetArtistCount(ctx context.Context) (n int64, err error) {
	tx, err := FetchTransaction()
	defer RollbackTransaction(tx)
	if err != nil {
		return
	}
	err = tx.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM artists",
	).Scan(&n)
	if err != nil {
		return
	}
	err = CommitTransaction(tx)
	return
}

func GetAreaCount(ctx context.Context) (n int64, err error) {
	tx, err := FetchTransaction()
	defer RollbackTransaction(tx)
	if err != nil {
		return
	}
	err = tx.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM areas",
	).Scan(&n)
	if err != nil {
		return
	}
	err = CommitTransaction(tx)
	return
}
