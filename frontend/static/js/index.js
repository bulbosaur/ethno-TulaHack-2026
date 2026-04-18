document.addEventListener("DOMContentLoaded", async () => {
  let allowedRegionNames = [];
  let allowedRegionsMap = {};

  try {
    const res = await fetch("/api/regions?limit=1");
    if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
    const data = await res.json();

    allowedRegionNames = data.map((r) => r.name);

    allowedRegionsMap = data.reduce((acc, r) => {
      acc[r.name.toLowerCase()] = r;
      return acc;
    }, {});

    console.log("Loaded regions from backend:", allowedRegionNames);
  } catch (err) {
    console.error("Failed to load regions:", err);
    allowedRegionNames = [];
  }

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

  function updatePanel(feature, dbData = null) {
    const name =
      feature.properties.name || feature.properties.NAME_1 || "Region";
    regionNameEl.textContent = name;

    if (dbData) {
      regionNameEl.textContent = dbData.title || name;
      regionInfoEl.textContent =
        dbData.summary ||
        feature.properties.description ||
        "Information not available.";
    } else {
      regionInfoEl.textContent =
        feature.properties.info ||
        feature.properties.description ||
        "Information not available.";
    }
    infoPanel.style.display = "block";
  }

  function onEachFeature(feature, layer) {
    layer.bindTooltip(feature.properties.name || "Region", { sticky: true });
    layer.on({
      mouseover: (e) => e.target.setStyle(highlightStyle()),
      mouseout: (e) => e.target.setStyle(defaultStyle()),
      click: (e) => {
        const name = feature.properties.name || feature.properties.NAME_1;
        const dbData = allowedRegionsMap[name.toLowerCase()];
        updatePanel(feature, dbData);
        map.fitBounds(e.target.getBounds(), { padding: [40, 40], maxZoom: 7 });
      },
    });
  }

  const GEOJSON_URL =
    "https://raw.githubusercontent.com/codeforamerica/click_that_hood/master/public/data/russia.geojson";

  fetch(GEOJSON_URL)
    .then((res) => {
      if (!res.ok) throw new Error("Failed to load map");
      return res.json();
    })
    .then((data) => {
      const geoLayer = L.geoJSON(data, {
        style: defaultStyle,
        onEachFeature: onEachFeature,
        // filter: function (feature) {
        //   if (allowedRegionNames.length === 0) return true;

        //   const name = (
        //     feature.properties.name ||
        //     feature.properties.NAME_1 ||
        //     ""
        //   )
        //     .toLowerCase()
        //     .trim();

        //   return allowedRegionNames.some(
        //     (allowed) => allowed.toLowerCase() === name,
        //   );
        // },
      }).addTo(map);

      setTimeout(() => {
        map.invalidateSize();
        const allLayers = geoLayer.getLayers();
        const validLayers = allLayers.filter(
          (l) => l.getBounds && l.getBounds().isValid(),
        );

        const targetPool =
          allowedRegionNames.length > 0
            ? validLayers.filter((layer) => {
                const name = (
                  layer.feature.properties.name ||
                  layer.feature.properties.NAME_1 ||
                  ""
                )
                  .toLowerCase()
                  .trim();
                return allowedRegionNames.some(
                  (allowed) => allowed.toLowerCase() === name,
                );
              })
            : validLayers;

        if (targetPool.length === 0) return;

        const randomLayer =
          targetPool[Math.floor(Math.random() * targetPool.length)];
        randomLayer.setStyle(highlightStyle());

        const name =
          randomLayer.feature.properties.name ||
          randomLayer.feature.properties.NAME_1;
        const dbData = allowedRegionsMap[name.toLowerCase()];

        updatePanel(randomLayer.feature, dbData);

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
          ${err.message}<br><small>For local testing, use Live Server</small>
        </div>`;
    });
});

async function loadQuests() {
  const grid = document.getElementById("quests-grid");

  try {
    const res = await fetch("/api/quests");
    if (!res.ok) throw new Error("Failed to load quests");
    const { data: quests } = await res.json();

    if (quests.length === 0) {
      grid.innerHTML = '<div class="quest-card">Квесты пока недоступны</div>';
      return;
    }

    grid.innerHTML = quests
      .map(
        (quest) => `
      <a href="/quests/${quest.slug}" class="quest-card" data-slug="${quest.slug}">
        ${
          quest.cover
            ? `<img src="${quest.cover}" alt="${quest.title}" class="quest-cover" onerror="this.classList.add('missing'); this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22240%22 height=%22120%22><text y=%22.6em%22 font-size=%2240%22>🎨</text></svg>'">`
            : `<div class="quest-cover missing">🎨</div>`
        }
        <div class="quest-info">
          <h4>${quest.title}</h4>
          <p>${quest.description ? quest.description.slice(0, 80) + (quest.description.length > 80 ? "..." : "") : ""}</p>
        </div>
        <div class="quest-meta">
          <span class="quest-badge">Начать</span>
          <span>→</span>
        </div>
      </a>
    `,
      )
      .join("");
  } catch (err) {
    console.error("Quests load error:", err);
    grid.innerHTML =
      '<div class="quest-card" style="border-color:#ef4444">Не удалось загрузить квесты</div>';
  }
}

document.addEventListener("DOMContentLoaded", () => {
  setTimeout(loadQuests, 500);
});
