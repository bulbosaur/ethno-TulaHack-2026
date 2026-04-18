function renderMarkdown(text) {
  if (!text) return "";

  let html = text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

  html = html
    .replace(/^### (.*$)/gim, "<h3>$1</h3>")
    .replace(/^## (.*$)/gim, "<h2>$1</h2>")
    .replace(/^# (.*$)/gim, "<h1>$1</h1>");

  html = html
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.*?)\*/g, "<em>$1</em>");

  if (html.includes("\n- ") || html.startsWith("- ")) {
    html = html
      .replace(/^\- (.*$)/gim, "<li>$1</li>")
      .replace(/(<li>.*<\/li>)/gs, "<ul>$1</ul>");
  }

  html = html.replace(/\n/g, "<br>");

  return html;
}

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

async function updateAuthUI() {
  const nav = document.querySelector(".sidebar-nav");
  if (!nav) return;

  try {
    const res = await fetch("/api/me", { credentials: "include" });

    if (res.ok) {
      const user = await res.json();

      nav.innerHTML = `
        <a href="/profile" class="sidebar-link">👤 ${user.username || "Личный кабинет"}</a>
        <button id="logout-btn" class="sidebar-link" style="background:none;border:none;cursor:pointer;text-align:left;width:100%;font:inherit;color:var(--text);padding:12px 16px;border-radius:8px;">Выйти</button>
      `;

      document
        .getElementById("logout-btn")
        ?.addEventListener("click", async (e) => {
          e.preventDefault();

          await fetch("/api/logout", {
            method: "POST",
            credentials: "include",
          });

          localStorage.removeItem("auth_token");
          sessionStorage.removeItem("auth_token");

          await new Promise((resolve) => setTimeout(resolve, 150));

          updateAuthUI();
          window.location.href = "/";
        });
    } else {
      nav.innerHTML = `
        <a href="/login" class="sidebar-link">Войти</a>
        <a href="/register" class="sidebar-link">Зарегистрироваться</a>
      `;
    }
  } catch (err) {
    console.error("Ошибка updateAuthUI:", err);
    nav.innerHTML = `
      <a href="/login" class="sidebar-link">Войти</a>
      <a href="/register" class="sidebar-link">Зарегистрироваться</a>
    `;
  }
}

document.addEventListener("DOMContentLoaded", async () => {
  let allowedRegionNames = [];
  let allowedRegionsMap = {};

  try {
    const res = await fetch("/api/regions");
    if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
    const data = await res.json();

    console.log("📦 API Response:", data);
    if (data[0]) {
      console.log("🔑 Keys:", Object.keys(data[0]));
      console.log("📝 First item:", JSON.stringify(data[0], null, 2));
    }

    allowedRegionNames = data.map((r) => r.name || r.Name || r.title || "");

    allowedRegionNames = data.map((r) => r.name);
    allowedRegionsMap = data.reduce((acc, r) => {
      acc[r.name.toLowerCase()] = r;
      return acc;
    }, {});

    console.log("Loaded all regions from backend:", allowedRegionNames.length);
  } catch (err) {
    console.error("Failed to load regions:", err);
  }

  const map = L.map("map", { zoomControl: true }).setView([60, 90], 3);
  map.attributionControl.setPrefix("");

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "© OpenStreetMap contributors",
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

    if (dbData) {
      regionNameEl.textContent = dbData.title || name;

      let summaryText = "Information not available.";
      if (dbData.summary?.text) {
        summaryText = dbData.summary.text;
      } else if (typeof dbData.summary === "string") {
        summaryText = dbData.summary;
      }
      regionInfoEl.innerHTML = renderMarkdown(summaryText);
    } else {
      regionNameEl.textContent = name;
      regionInfoEl.textContent =
        feature.properties.info ||
        feature.properties.description ||
        "Информация о регионе пока не добавлена.";
    }

    infoPanel.style.display = "block";
  }

  function onEachFeature(feature, layer) {
    layer.bindTooltip(feature.properties.name || "Region", { sticky: true });

    layer.on({
      mouseover: (e) => {
        if (!e.target._isRandomlyHighlighted) {
          e.target.setStyle(highlightStyle());
        }
      },
      mouseout: (e) => {
        if (!e.target._isRandomlyHighlighted) {
          e.target.setStyle(defaultStyle());
        }
      },
      click: (e) => {
        map.eachLayer((l) => {
          if (l._isRandomlyHighlighted && l.setStyle) {
            l.setStyle(defaultStyle());
            l._isRandomlyHighlighted = false;
          }
        });

        const name = feature.properties.name || feature.properties.NAME_1;
        const dbData = allowedRegionsMap[name.toLowerCase()];

        updatePanel(feature, dbData);

        e.target.setStyle(highlightStyle());

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
      }).addTo(map);

      setTimeout(() => {
        map.invalidateSize();

        const allLayers = geoLayer.getLayers();
        const validLayers = allLayers.filter(
          (l) => l.getBounds && l.getBounds().isValid(),
        );

        const layersWithDbData = validLayers.filter((layer) => {
          const name = (
            layer.feature.properties.name ||
            layer.feature.properties.NAME_1 ||
            ""
          )
            .toLowerCase()
            .trim();
          return allowedRegionsMap[name] !== undefined;
        });

        const targetPool =
          layersWithDbData.length > 0 ? layersWithDbData : validLayers;

        if (targetPool.length === 0) return;

        const randomLayer =
          targetPool[Math.floor(Math.random() * targetPool.length)];

        randomLayer.setStyle(highlightStyle());
        randomLayer._isRandomlyHighlighted = true;

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

  await updateAuthUI();
  setTimeout(loadQuests, 500);
});
