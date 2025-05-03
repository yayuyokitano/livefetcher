package runner

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	coreconnectors "github.com/yayuyokitano/livefetcher/internal/core/connectors"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
)

func RunConnector(ctx context.Context, connectorID string) (err error) {
	if _, ok := coreconnectors.Connectors[connectorID]; !ok {
		err = errors.New("Connector not found: " + connectorID)
		return
	}
	fetcher := coreconnectors.Connectors[connectorID]
	err = fetcher.Fetch()
	if len(fetcher.Lives) == 0 {
		fmt.Println(err)
		return
	}

	deleted, added, modified, addedArtists, err := queries.PostLives(ctx, fetcher.Lives, &http.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Lives: deleted %d, added %d, modified %d, added %d artists\n", deleted, added, modified, addedArtists)
	return
}

func RunConnectorTest(connectorID string) (err error) {
	if _, ok := coreconnectors.Connectors[connectorID]; !ok {
		err = errors.New("Connector not found: " + connectorID)
		return
	}
	fetcher := coreconnectors.Connectors[connectorID]
	err = fetcher.Fetch()
	if len(fetcher.Lives) == 0 {
		fmt.Println(err)
		return
	}
	fmt.Println(fetcher.Lives)
	return
}
