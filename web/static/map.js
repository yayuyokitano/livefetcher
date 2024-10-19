function onEachFeature(feature, layer) {
  if (
    feature.properties &&
    feature.properties.name &&
    feature.properties.popupContent
  ) {
    layer.bindPopup(
      `<b>${feature.properties.name}</b><br/><br/>${feature.properties.popupContent}`,
    );
  }
}

function createMapDataRetriever() {
  const leaflet = L.map("map").setView([34.671662, 135.497672], 13);
  L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution:
      '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
  }).addTo(leaflet);
  const layerGroup = L.layerGroup().addTo(leaflet);

  return {
    isLoading: false,
    map: {
      lives: [],
      geoJson: [],
    },
    date: "2024-10-30",
    getMapData() {
      this.isLoading = true;
      console.log("hi");
      const urlDate = this.date.split("-").join("/");
      fetch("/api/dailylives/" + urlDate)
        .then((res) => res.json())
        .then((map) => {
          this.map = map;
          layerGroup.clearLayers();
          L.geoJSON(map.geoJson, {
            onEachFeature: onEachFeature,
          }).addTo(layerGroup);
          this.isLoading = false;
        });
    },
  };
}
