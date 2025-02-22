package coreconnectors

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

type Translation struct {
	Prefectures map[string]string            `toml:"prefecture"`
	Areas       map[string]map[string]string `toml:"area"`
	Livehouses  map[string]string            `toml:"livehouse"`
}

var translations map[string]Translation

func TestConnectors(t *testing.T) {
	initTranslations(t)
	testResults := make(chan util.ConnectorTestResult, 1)

	// test only one connector if specified, otherwise test all
	connectorID := os.Getenv("CONNECTOR_ID")
	if _, ok := Connectors[connectorID]; ok {
		errorCreator := util.ConnectorTestErrorCreator{
			T:    t,
			Chan: testResults,
			Name: Connectors[connectorID].VenueID,
		}
		var wg sync.WaitGroup
		wg.Add(1)
		executeConnectorTest(t, Connectors[connectorID], &wg, &errorCreator)
		return
	}

	var wg sync.WaitGroup
	for _, connector := range Connectors {
		errorCreator := util.ConnectorTestErrorCreator{
			T:    t,
			Chan: testResults,
			Name: connector.VenueID,
		}
		wg.Add(1)
		go executeConnectorTest(t, connector, &wg, &errorCreator)
	}
	go func() {
		wg.Wait()
		close(testResults)
	}()

	res := make(util.ConnectorTestResults, 0)
	for {
		t, ok := <-testResults
		if !ok {
			break
		}
		res = append(res, t)
	}
	fmt.Println(res.String())
}

func executeConnectorTest(t *testing.T, connector fetchers.Simple, wg *sync.WaitGroup, errorCreator *util.ConnectorTestErrorCreator) {
	for lang, translation := range translations {
		if translation.Prefectures[connector.PrefectureName] == "" {
			errorCreator.Errorf("No %s translation found for prefecture %s.", lang, connector.PrefectureName)
		}

		if translation.Areas[connector.PrefectureName][connector.AreaName] == "" {
			errorCreator.Errorf("No %s translation found for area %s in prefecture %s.", lang, connector.AreaName, connector.PrefectureName)
		}

		if translation.Livehouses[connector.VenueID] == "" {
			errorCreator.Errorf("No %s translation found for venue %s.", lang, connector.VenueID)
		}
	}

	t.Run(connector.VenueID, func(t *testing.T) {
		defer wg.Done()

		connector.Test(t, errorCreator)
	})
}

func initTranslations(t *testing.T) {
	translations = make(map[string]Translation)
	files, err := os.ReadDir("../../i18n/locales")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		var translation Translation
		_, err = toml.DecodeFile("../../i18n/locales/"+f.Name(), &translation)
		if err != nil {
			t.Fatal(err)
		}
		translations[f.Name()] = translation
	}
}
