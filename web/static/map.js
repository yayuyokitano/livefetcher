function getLive(liveID) {
  return document.querySelector(`input[value="${liveID}"]`).closest("li");
}

function createMapDataRetriever() {
  const leaflet = L.map("map").setView([34.671662, 135.497672], 13);
  L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution:
      '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
  }).addTo(leaflet);
  const layerGroup = L.layerGroup().addTo(leaflet);
	const now = new Date();
	if (now.getHours() > 18) {
		now.setDate(now.getDate() + 1);
	}
	const dateString = `${now.getFullYear()}-${(now.getMonth() + 1).toString().padStart(2, "0")}-${now.getDate().toString().padStart(2, "0")}`

  return {
    isLoading: false,
    map: {
      lives: [],
      geoJson: [],
    },
    filteredLives: [],
    liveMarkers: {},
    date: dateString,
    filterLives() {
      const bounds = leaflet.getBounds();
      const minLat = bounds.getSouth();
      const maxLat = bounds.getNorth();
      const minLng = bounds.getWest();
      const maxLng = bounds.getEast();
      this.filteredLives = this.map.lives.filter((live) => {
        if (live.venue.latitude < minLat || live.venue.latitude > maxLat) {
          return false;
        }
        if (live.venue.longitude < minLng || live.venue.longitude > maxLng) {
          return false;
        }
        return true;
      });
    },
    initMapData() {
      leaflet.addEventListener("moveend", this.filterLives.bind(this));
      leaflet.addEventListener("zoomend", this.filterLives.bind(this));
      this.getMapData();
    },
    getMapData() {
      this.isLoading = true;
      const urlDate = this.date.split("-").join("/");
      fetch("/api/dailylives/" + urlDate)
        .then((res) => res.json())
        .then((map) => {
          this.map = map;
          layerGroup.clearLayers();
          L.geoJSON(map.geoJson, {
            onEachFeature: this.onEachFeature.bind(this),
          }).addTo(layerGroup);
          layerGroup.addTo(leaflet);
          this.filterLives();
          this.isLoading = false;
        });
    },
    onEachFeature(feature, layer) {
      if (
        feature.properties &&
        feature.properties.name &&
        feature.properties.popupContent
      ) {
        layer.bindPopup(
          `<b>${feature.properties.name}</b><br/><br/>${feature.properties.popupContent}`,
        );
      }
      if (feature.properties.id) {
        this.liveMarkers[feature.properties.id] = layer;
        layer.on("click", () => {
          const liveEl = getLive(feature.properties.id);
          liveEl.focus();
          liveEl.scrollIntoView({
            block: "center",
            inline: "center",
          });
        });
      }
    },
  };
}

function getLiveCounter(num) {
  if (document.documentElement.lang == "en") {
    return `${num} lives`;
  }
  return `${num}ライブ`;
}

