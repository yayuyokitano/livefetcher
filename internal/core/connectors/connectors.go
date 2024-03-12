package coreconnectors

import (
	"github.com/yayuyokitano/livefetcher/internal/connectors"
	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
)

type ConnectorsType map[string]fetchers.Simple

var Connectors = ConnectorsType{
	"ShibuyaCyclone":              connectors.ShibuyaCycloneFetcher,
	"ShibuyaDive":                 connectors.ShibuyaDiveFetcher,
	"ShibuyaEggmanDay":            connectors.ShibuyaEggmanDayFetcher,
	"ShibuyaEggmanNight":          connectors.ShibuyaEggmanNightFetcher,
	"ShibuyaGarret":               connectors.ShibuyaGarretFetcher,
	"ShibuyaLaDonna":              connectors.ShibuyaLaDonnaFetcher,
	"ShibuyaOCrest":               connectors.ShibuyaOCrestFetcher,
	"ShibuyaOEast":                connectors.ShibuyaOEastFetcher,
	"ShibuyaONest":                connectors.ShibuyaONestFetcher,
	"ShibuyaOWest":                connectors.ShibuyaOWestFetcher,
	"ShibuyaStrobe":               connectors.ShibuyaStrobeFetcher,
	"ShibuyaTokioTokyo":           connectors.ShibuyaTokioTokyoFetcher,
	"ShibuyaVeats":                connectors.ShibuyaVeatsFetcher,
	"ShibuyaWWW":                  connectors.ShibuyaWWWFetcher,
	"ShibuyaWWWX":                 connectors.ShibuyaWWWXFetcher,
	"ShibuyaWWWBeta":              connectors.ShibuyaWWWBetaFetcher,
	"ShimokitazawaArtist":         connectors.ShimokitazawaArtistFetcher,
	"ShimokitazawaBasementBar":    connectors.ShimokitazawaBasementBarFetcher,
	"ShimokitazawaChikamatsu":     connectors.ShimokitazawaChikamatsuFetcher,
	"ShimokitazawaChikamichi":     connectors.ShimokitazawaChikamichiFetcher,
	"ShimokitazawaClub251":        connectors.ShimokitazawaClub251Fetcher,
	"ShimokitazawaClubQue":        connectors.ShimokitazawaClubQueFetcher,
	"ShimokitazawaDaisyBar":       connectors.ShimokitazawaDaisyBarFetcher,
	"ShimokitazawaDyCube":         connectors.ShimokitazawaDyCubeFetcher,
	"ShimokitazawaEra":            connectors.ShimokitazawaEraFetcher,
	"ShimokitazawaFlowersLoft":    connectors.ShimokitazawaFlowersLoftFetcher,
	"ShimokitazawaLaguna":         connectors.ShimokitazawaLagunaFetcher,
	"ShimokitazawaLivehaus":       connectors.ShimokitazawaLiveHausFetcher,
	"ShimokitazawaLiveHolic":      connectors.ShimokitazawaLiveHolicFetcher,
	"ShimokitazawaMonaRecords":    connectors.ShimokitazawaMonaRecordsFetcher,
	"ShimokitazawaMosaic":         connectors.ShimokitazawaMosaicFetcher,
	"ShimokitazawaOtemae":         connectors.ShimokitazawaOtemaeFetcher,
	"ShimokitazawaReg":            connectors.ShimokitazawaRegFetcher,
	"ShimokitazawaShangrila":      connectors.ShimokitazawaShangrilaFetcher,
	"ShimokitazawaShelter":        connectors.ShimokitazawaShelterFetcher,
	"ShimokitazawaSpread":         connectors.ShimokitazawaSpreadFetcher,
	"ShimokitazawaThree":          connectors.ShimokitazawaThreeFetcher,
	"ShimokitazawaWaver":          connectors.ShimokitazawaWaverFetcher,
	"ShindaitaFever":              connectors.ShindaitaFeverFetcher,
	"ShinjukuLoft":                connectors.ShinjukuLoftFetcher,
	"ShinjukuZircoTokyo":          connectors.ShinjukuZircoTokyoFetcher,
	"ShinsaibashiAnima":           connectors.ShinsaibashiAnimaFetcher,
	"ShinsaibashiBeyond":          connectors.ShinsaibashiBeyondFetcher,
	"ShinsaibashiBigcat":          connectors.ShinsaibashiBigcatFetcher,
	"ShinsaibashiBronze":          connectors.ShinsaibashiBronzeFetcher,
	"ShinsaibashiClapper":         connectors.ShinsaibashiClapperFetcher,
	"ShinsaibashiClubVijon":       connectors.ShinsaibashiClubVijonFetcher,
	"ShinsaibashiConpass":         connectors.ShinsaibashiConpassFetcher,
	"ShinsaibashiDrop":            connectors.ShinsaibashiDropFetcher,
	"ShinsaibashiFanjtwice":       connectors.ShinsaibashiFanjtwiceFetcher,
	"ShinsaibashiHillsPan":        connectors.ShinsaibashiHillsPanFetcher,
	"ShinsaibashiHokage":          connectors.ShinsaibashiHokageFetcher,
	"ShinsaibashiJanus":           connectors.ShinsaibashiJanusFetcher,
	"ShinsaibashiKanon":           connectors.ShinsaibashiKanonFetcher,
	"ShinsaibashiKingCobra":       connectors.ShinsaibashiKingCobraFetcher,
	"ShinsaibashiKnave":           connectors.ShinsaibashiKnaveFetcher,
	"ShinsaibashiKurage":          connectors.ShinsaibashiKurageFetcher,
	"ShinsaibashiLoftPlusOneWest": connectors.ShinsaibashiLoftPlusOneWestFetcher,
	"ShinsaibashiMuse":            connectors.ShinsaibashiMuseFetcher,
	"ShinsaibashiPangea":          connectors.ShinsaibashiPangeaFetcher,
	"ShinsaibashiSocoreFactory":   connectors.ShinsaibashiSocoreFactoryFetcher,
	"ShinsaibashiSoma":            connectors.ShinsaibashiSomaFetcher,
	"ShinsaibashiQupe":            connectors.ShinsaibashiQupeFetcher,
	"ShinsaibashiUtausakana":      connectors.ShinsaibashiUtausakanaFetcher,
	"ShinsaibashiVaron":           connectors.ShinsaibashiVaronFetcher,
}
