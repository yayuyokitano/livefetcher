package coreconnectors

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/yayuyokitano/livefetcher/lib/core/fetchers"
)

type Translation struct {
	Prefectures map[string]string            `toml:"prefecture"`
	Areas       map[string]map[string]string `toml:"area"`
	Livehouses  map[string]string            `toml:"livehouse"`
}

var translations map[string]Translation

func TestConnectors(t *testing.T) {
	initTranslations(t)

	var wg sync.WaitGroup
	for _, connector := range Connectors {
		wg.Add(1)
		go executeConnectorTest(t, connector, &wg)
	}
	wg.Wait()
}

func executeConnectorTest(t *testing.T, connector fetchers.Simple, wg *sync.WaitGroup) {
	for lang, translation := range translations {
		if translation.Prefectures[connector.PrefectureName] == "" {
			t.Errorf("No %s translation found for prefecture %s.", lang, connector.PrefectureName)
		}

		if translation.Areas[connector.PrefectureName][connector.AreaName] == "" {
			t.Errorf("No %s translation found for area %s in prefecture %s.", lang, connector.AreaName, connector.PrefectureName)
		}

		if translation.Livehouses[connector.VenueID] == "" {
			t.Errorf("No %s translation found for venue %s.", lang, connector.VenueID)
		}
	}

	t.Run(connector.VenueID, func(t *testing.T) {
		defer wg.Done()

		err := connector.Test()
		if err != nil {
			t.Error(err)
		}
	})
}

func initTranslations(t *testing.T) {
	translations = make(map[string]Translation)
	files, err := ioutil.ReadDir("../../../i18nloader/locales")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		var translation Translation
		_, err = toml.DecodeFile("../../../i18nloader/locales/"+f.Name(), &translation)
		if err != nil {
			t.Fatal(err)
		}
		translations[f.Name()] = translation
	}
}
