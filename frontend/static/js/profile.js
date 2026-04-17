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
      window.location.href = "/login.html";
      return;
    }

    const user = await res.json();
    userNameEl.textContent = user.name || user.email?.split("@")[0];

    await loadQuestRecommendations();
  } catch (err) {
    console.error("Ошибка загрузки кабинета:", err);
    window.location.href = "/login.html";
  }

  logoutBtn?.addEventListener("click", async () => {
    await fetch("/api/logout", {
      method: "POST",
      credentials: "include",
    });
    window.location.href = "/login.html";
  });

  async function loadQuestRecommendations() {
    try {
      const mockQuests = [
        {
          id: 1,
          title: "Тайны северных народов",
          description: "Исследуй традиции и обряды коренных народов Севера",
          region: "Ямал",
          difficulty: "Средний",
        },
        {
          id: 2,
          title: "Путь кочевника",
          description: "Погрузись в культуру степных народов и их миграции",
          region: "Алтай",
          difficulty: "Лёгкий",
        },
        {
          id: 3,
          title: "Голос тайги",
          description: "Узнай о верованиях и ремёслах сибирских народов",
          region: "Бурятия",
          difficulty: "Сложный",
        },
      ];

      renderQuests(mockQuests);
    } catch (err) {
      questsGrid.innerHTML = `<div class="quest-card placeholder">Не удалось загрузить рекомендации</div>`;
      console.error("Ошибка загрузки квестов:", err);
    }
  }

  function renderQuests(quests) {
    if (!quests.length) {
      questsGrid.innerHTML = `<div class="quest-card placeholder">Пока нет рекомендаций</div>`;
      return;
    }

    questsGrid.innerHTML = quests
      .map(
        (quest) => `
        <div class="quest-card" data-id="${quest.id}">
          <h3>${escapeHtml(quest.title)}</h3>
          <p>${escapeHtml(quest.description)}</p>
          <div class="quest-meta">
            <span class="tag">📍 ${escapeHtml(quest.region)}</span>
            <span class="tag">⚡ ${escapeHtml(quest.difficulty)}</span>
          </div>
        </div>
      `,
      )
      .join("");

    document.querySelectorAll(".quest-card").forEach((card) => {
      card.addEventListener("click", () => {
        const id = card.dataset.id;
        console.log(`Переход к квесту #${id}`);
      });
    });
  }

  function escapeHtml(str) {
    if (!str) return "";
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
  }
});
