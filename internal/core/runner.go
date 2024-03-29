package runner

import (
	"context"
	"errors"
	"fmt"

	coreconnectors "github.com/yayuyokitano/livefetcher/internal/core/connectors"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
)

func RunConnector(connectorID string) (err error) {
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

	n, err := queries.PostLives(context.Background(), fetcher.Lives)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Inserted", n, "lives")
	return
}

func RunConnectorTest(connectorID string) (err error) {
	if _, ok := coreconnectors.Connectors[connectorID]; !ok {
		err = errors.New("Connector not found: " + connectorID)
		return
	}
	fetcher := coreconnectors.Connectors[connectorID]
	err = fetcher.Fetch()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fetcher.Lives)
	return
}
