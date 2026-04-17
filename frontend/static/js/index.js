document.addEventListener("DOMContentLoaded", () => {
  const ALLOWED_REGIONS = [
    "Москва",
    "Санкт-Петербург",
    "Карачаево-Черкесская Республика",
  ];

  const map = L.map("map", {
    zoomControl: true,
  }).setView([60, 90], 3);

  map.attributionControl.setPrefix("");

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution:
      '<a href="https://leafletjs.com" title="A JS library for interactive maps">Leaflet</a> | © <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    maxZoom: 18,
  }).addTo(map);

  const infoPanel = document.getElementById("info-panel");
  const regionNameEl = document.getElementById("region-name");
  const regionInfoEl = document.getElementById("region-info");
  const closePanelBtn = document.getElementById("close-panel");

  closePanelBtn.addEventListener("click", () => {
    infoPanel.style.display = "none";
  });

  const defaultStyle = () => ({
    weight: 1,
    color: "#cbd5e1",
    fillColor: "#64748b",
    fillOpacity: 0.5,
  });
  const highlightStyle = () => ({
    weight: 2,
    color: "#eb2577",
    fillColor: "rgb(247, 119, 198)",
    fillOpacity: 0.75,
  });

  function updatePanel(feature) {
    regionNameEl.textContent =
      feature.properties.name || feature.properties.NAME_1 || "Регион";
    regionInfoEl.textContent =
      feature.properties.info ||
      feature.properties.description ||
      "Информация не указана.";
    infoPanel.style.display = "block";
  }

  function onEachFeature(feature, layer) {
    layer.bindTooltip(feature.properties.name || "Регион", { sticky: true });
    layer.on({
      mouseover: (e) => e.target.setStyle(highlightStyle()),
      mouseout: (e) => e.target.setStyle(defaultStyle()),
      click: (e) => {
        updatePanel(feature);
        map.fitBounds(e.target.getBounds(), { padding: [40, 40], maxZoom: 7 });
      },
    });
  }

  const GEOJSON_URL =
    "https://raw.githubusercontent.com/codeforamerica/click_that_hood/master/public/data/russia.geojson";

  fetch(GEOJSON_URL)
    .then((res) => {
      if (!res.ok) throw new Error("Не удалось загрузить карту");
      return res.json();
    })
    .then((data) => {
      const geoLayer = L.geoJSON(data, {
        style: defaultStyle,
        onEachFeature: onEachFeature,
      }).addTo(map);

      setTimeout(() => {
        map.invalidateSize();
        const allLayers = geoLayer.getLayers();
        const validLayers = allLayers.filter(
          (l) => l.getBounds && l.getBounds().isValid(),
        );

        if (ALLOWED_REGIONS.length > 0) {
          console.log(
            "📦 Реальные названия в GeoJSON:",
            validLayers.map(
              (l) => l.feature.properties.name || l.feature.properties.NAME_1,
            ),
          );
        }

        const pool =
          ALLOWED_REGIONS.length > 0
            ? validLayers.filter((layer) => {
                const name = (
                  layer.feature.properties.name ||
                  layer.feature.properties.NAME_1 ||
                  ""
                )
                  .trim()
                  .toLowerCase();
                return ALLOWED_REGIONS.some(
                  (allowed) => allowed.toLowerCase() === name,
                );
              })
            : validLayers;

        const targetPool = pool.length > 0 ? pool : validLayers;
        if (targetPool.length === 0) return;

        const randomLayer =
          targetPool[Math.floor(Math.random() * targetPool.length)];
        randomLayer.setStyle(highlightStyle());
        updatePanel(randomLayer.feature);
        map.fitBounds(randomLayer.getBounds(), {
          padding: [50, 50],
          animate: true,
          duration: 1,
        });
      }, 400);
    })
    .catch((err) => {
      document.getElementById("map").innerHTML = `
        <div style="display:flex;align-items:center;justify-content:center;height:100%;color:#ef4444;padding:20px;">
          ${err.message}<br><small>Для локального теста используйте Live Server</small>
        </div>`;
    });
});
