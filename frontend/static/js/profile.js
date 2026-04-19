document.addEventListener("DOMContentLoaded", async () => {
  const userNameEl = document.getElementById("user-name");
  const questsGrid = document.getElementById("quests-grid");
  const logoutBtn = document.getElementById("logout-btn");

  try {
    const res = await fetch("/api/me", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    });

    if (!res.ok) {
      window.location.href = "/login";
      return;
    }

    const user = await res.json();
    userNameEl.textContent =
      user.username || user.email?.split("@")[0] || "Пользователь";
    await loadQuestRecommendations();
  } catch (err) {
    console.error("Ошибка загрузки кабинета:", err);
    userNameEl.textContent = "Ошибка";
    questsGrid.innerHTML = `<div class="quest-card" style="border-color:#ef4444">Не удалось загрузить рекомендации</div>`;
  }

  logoutBtn?.addEventListener("click", async (e) => {
    e.preventDefault();

    try {
      await fetch("/api/logout", {
        method: "POST",
        credentials: "include",
      });
    } catch (err) {
      console.error("Ошибка при выходе:", err);
    } finally {
      localStorage.removeItem("auth_token");
      sessionStorage.removeItem("auth_token");
      window.location.href = "/";
    }
  });

  async function loadQuestRecommendations() {
    try {
      const mockQuests = [
        {
          id: 1,
          title: "Тропа мастеров Ханты",
          description:
            "Познакомься с древними ремёслами и традициями хантов. Узнай о тамгах, узорах и создай свой оберег.",
          region: "Ямал",
          difficulty: "Средний",
          slug: "khanty-crafts",
          cover: null,
        },
        {
          id: 2,
          title: "Путь кочевника",
          description: "Погрузись в культуру степных народов и их миграции",
          region: "Алтай",
          difficulty: "Лёгкий",
          slug: "nomad-path",
          cover: null,
        },
        {
          id: 3,
          title: "Голос тайги",
          description: "Узнай о верованиях и ремёслах сибирских народов",
          region: "Бурятия",
          difficulty: "Сложный",
          slug: "taiga-voice",
          cover: null,
        },
      ];

      renderQuests(mockQuests);
    } catch (err) {
      questsGrid.innerHTML = `<div class="quest-card" style="border-color:#ef4444">Не удалось загрузить рекомендации</div>`;
      console.error("Ошибка загрузки квестов:", err);
    }
  }

  function renderQuests(quests) {
    if (!quests?.length) {
      questsGrid.innerHTML = `<div class="quest-card">Пока нет рекомендаций</div>`;
      return;
    }

    questsGrid.innerHTML = quests
      .map(
        (quest) => `
        <a href="/quests/${quest.slug}" class="quest-card">
          ${
            quest.cover
              ? `<img src="${quest.cover}" alt="${quest.title}" class="quest-cover" onerror="this.classList.add('missing'); this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22240%22 height=%22120%22><text y=%22.6em%22 font-size=%2240%22>🎨</text></svg>'">`
              : `<div class="quest-cover missing">🎨</div>`
          }
          <div class="quest-info">
            <h4>${escapeHtml(quest.title)}</h4>
            <p>${escapeHtml(quest.description?.slice(0, 80) + (quest.description?.length > 80 ? "..." : ""))}</p>
          </div>
          <div class="quest-meta">
            <span class="quest-badge">📍 ${escapeHtml(quest.region)}</span>
            <span>⚡ ${escapeHtml(quest.difficulty)}</span>
            <span>→</span>
          </div>
        </a>
      `,
      )
      .join("");
  }

  function escapeHtml(str) {
    if (!str) return "";
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
  }
});
