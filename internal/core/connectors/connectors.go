package coreconnectors

import (
	"github.com/yayuyokitano/livefetcher/internal/connectors"
	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
)

type ConnectorsType map[string]fetchers.Simple

var Connectors = ConnectorsType{
	"ShibuyaEggmanDay":         connectors.ShibuyaEggmanDayFetcher,
	"ShibuyaEggmanNight":       connectors.ShibuyaEggmanNightFetcher,
	"ShibuyaOCrest":            connectors.ShibuyaOCrestFetcher,
	"ShibuyaOEast":             connectors.ShibuyaOEastFetcher,
	"ShibuyaONest":             connectors.ShibuyaONestFetcher,
	"ShibuyaOWest":             connectors.ShibuyaOWestFetcher,
	"ShibuyaWWW":               connectors.ShibuyaWWWFetcher,
	"ShibuyaWWWX":              connectors.ShibuyaWWWXFetcher,
	"ShibuyaWWWBeta":           connectors.ShibuyaWWWBetaFetcher,
	"ShimokitazawaArtist":      connectors.ShimokitazawaArtistFetcher,
	"ShimokitazawaBasementBar": connectors.ShimokitazawaBasementBarFetcher,
	"ShimokitazawaChikamatsu":  connectors.ShimokitazawaChikamatsuFetcher,
	"ShimokitazawaChikamichi":  connectors.ShimokitazawaChikamichiFetcher,
	"ShimokitazawaClub251":     connectors.ShimokitazawaClub251Fetcher,
	"ShimokitazawaClubQue":     connectors.ShimokitazawaClubQueFetcher,
	"ShimokitazawaDaisyBar":    connectors.ShimokitazawaDaisyBarFetcher,
	"ShimokitazawaDyCube":      connectors.ShimokitazawaDyCubeFetcher,
	"ShimokitazawaEra":         connectors.ShimokitazawaEraFetcher,
	"ShimokitazawaFlowersLoft": connectors.ShimokitazawaFlowersLoftFetcher,
	"ShimokitazawaLaguna":      connectors.ShimokitazawaLagunaFetcher,
	"ShimokitazawaLivehaus":    connectors.ShimokitazawaLiveHausFetcher,
	"ShimokitazawaLiveHolic":   connectors.ShimokitazawaLiveHolicFetcher,
	"ShimokitazawaMonaRecords": connectors.ShimokitazawaMonaRecordsFetcher,
	"ShimokitazawaMosaic":      connectors.ShimokitazawaMosaicFetcher,
	"ShimokitazawaOtemae":      connectors.ShimokitazawaOtemaeFetcher,
	"ShimokitazawaReg":         connectors.ShimokitazawaRegFetcher,
	"ShimokitazawaShangrila":   connectors.ShimokitazawaShangrilaFetcher,
	"ShimokitazawaShelter":     connectors.ShimokitazawaShelterFetcher,
	"ShimokitazawaSpread":      connectors.ShimokitazawaSpreadFetcher,
	"ShimokitazawaThree":       connectors.ShimokitazawaThreeFetcher,
	"ShimokitazawaWaver":       connectors.ShimokitazawaWaverFetcher,
	"ShindaitaFever":           connectors.ShindaitaFeverFetcher,
	"ShinjukuLoft":             connectors.ShinjukuLoftFetcher,
	"ShinsaibashiAnima":        connectors.ShinsaibashiAnimaFetcher,
	"ShinsaibashiBeyond":       connectors.ShinsaibashiBeyondFetcher,
	"ShinsaibashiBigcat":       connectors.ShinsaibashiBigcatFetcher,
	"ShinsaibashiBronze":       connectors.ShinsaibashiBronzeFetcher,
	"ShinsaibashiFanjtwice":    connectors.ShinsaibashiFanjtwiceFetcher,
	"ShinsaibashiHokage":       connectors.ShinsaibashiHokageFetcher,
	"ShinsaibashiJanus":        connectors.ShinsaibashiJanusFetcher,
	"ShinsaibashiKurage":       connectors.ShinsaibashiKurageFetcher,
	"ShinsaibashiMuse":         connectors.ShinsaibashiMuseFetcher,
	"ShinsaibashiPangea":       connectors.ShinsaibashiPangeaFetcher,
	"ShinsaibashiUtausakana":   connectors.ShinsaibashiUtausakanaFetcher,
}
