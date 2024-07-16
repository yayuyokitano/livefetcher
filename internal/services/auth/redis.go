package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/yayuyokitano/livefetcher/internal/services"
)

func registerRefreshToken(ctx context.Context, rtid string) (err error) {
	err = services.RDB.Set(ctx, rtid, "0", refreshTokenDuration).Err()
	if err != nil {
		return
	}
	return
}

func deleteRefreshToken(ctx context.Context, rtid string) (err error) {
	err = services.RDB.Del(ctx, rtid).Err()
	return
}

func checkRefreshToken(ctx context.Context, rtid string, useCount int) bool {
	val, err := services.RDB.Get(ctx, rtid).Result()
	if err != nil {
		return false
	}
	expectedUseCount, err := strconv.Atoi(val)
	if err != nil {
		services.RDB.Del(ctx, rtid).Err()
		return false
	}

	// disable family of refresh tokens including newer ones completely if there is any attempted duplicated use
	if expectedUseCount != useCount {
		services.RDB.Del(ctx, rtid).Err()
		return false
	}

	err = services.RDB.Set(ctx, rtid, fmt.Sprintf("%b", useCount+1), refreshTokenDuration).Err()
	if err != nil {
		services.RDB.Del(ctx, rtid).Err()
		return false
	}

	return true
}
