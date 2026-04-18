document.addEventListener("DOMContentLoaded", async () => {
  const slug = window.location.pathname.split("/").pop();
  if (!slug) {
    document.getElementById("step-container").innerHTML =
      '<p class="error">Неверный адрес квеста</p>';
    return;
  }

  const state = {
    quest: null,
    currentIndex: 0,
    selectedAnswer: null,
    builderSelection: [],
  };

  const els = {
    title: document.getElementById("quest-title"),
    desc: document.getElementById("quest-desc"),
    container: document.getElementById("step-container"),
    progress: document.getElementById("progress-fill"),
    currentNum: document.getElementById("current-step-num"),
    totalNum: document.getElementById("total-steps"),
    prevBtn: document.getElementById("prev-btn"),
    nextBtn: document.getElementById("next-btn"),
  };

  try {
    const res = await fetch(`/api/quests/${slug}`);
    if (!res.ok) throw new Error("Квест не найден");
    const { data: quest } = await res.json();

    state.quest = quest;
    renderQuest(quest);
  } catch (err) {
    console.error(err);
    els.container.innerHTML = `<p class="error">Ошибка загрузки: ${err.message}</p>`;
    els.nextBtn.disabled = true;
  }

  function renderQuest(quest) {
    els.title.textContent = quest.title;
    els.desc.textContent = quest.description;
    els.totalNum.textContent = quest.steps.length;
    updateProgress();
    renderStep(quest.steps[0]);
  }

  function renderStep(step) {
    els.currentNum.textContent = state.currentIndex + 1;
    updateProgress();
    updateButtons();

    state.selectedAnswer = null;
    state.builderSelection = [];

    switch (step.type) {
      case "intro":
        renderIntro(step);
        break;
      case "quiz":
        renderQuiz(step);
        break;
      case "builder":
        renderBuilder(step);
      case "map_match":
        renderMapMatch(step);
        break;
      default:
        els.container.innerHTML = `<p>Неизвестный тип шага: ${step.type}</p>`;
    }
  }

  function renderIntro(step) {
    const content = step.content;
    els.container.innerHTML = `
      <div class="step-intro">
        <h2>${step.title}</h2>
        ${content.image ? `<img src="${content.image}" alt="">` : ""}
        <p>${content.text}</p>
      </div>
    `;
    els.nextBtn.onclick = goToNext;
  }

  function renderQuiz(step) {
    const { question, options } = step.content;
    const optionsHtml = options
      .map(
        (opt) => `
      <label class="quiz-option" data-id="${opt.id}">
        <input type="radio" name="quiz-answer" value="${opt.id}">
        <span>${opt.text}</span>
      </label>
    `,
      )
      .join("");

    els.container.innerHTML = `
      <div class="step-quiz">
        <h2>${step.title}</h2>
        <p class="quiz-question">${question}</p>
        <div class="quiz-options">${optionsHtml}</div>
        <div id="quiz-result" class="quiz-result" hidden></div>
      </div>
    `;

    document.querySelectorAll(".quiz-option").forEach((label) => {
      label.onclick = (e) => {
        if (e.target.tagName === "INPUT") return;
        const input = label.querySelector("input");
        input.checked = true;
        state.selectedAnswer = input.value;
      };
    });

    els.nextBtn.onclick = () => {
      const resultEl = document.getElementById("quiz-result");
      const selected = state.selectedAnswer;

      if (!selected) {
        resultEl.textContent = "Выберите вариант ответа";
        resultEl.className = "quiz-result error";
        resultEl.hidden = false;
        return;
      }

      const correct = step.content.options.find(
        (o) => o.id === selected,
      )?.correct;

      if (correct) {
        resultEl.textContent = "Правильно!";
        resultEl.className = "quiz-result success";
        document
          .querySelector(`[data-id="${selected}"]`)
          .classList.add("correct");
        els.nextBtn.onclick = goToNext;
      } else {
        resultEl.textContent = "Попробуйте ещё раз";
        resultEl.className = "quiz-result error";
        document
          .querySelector(`[data-id="${selected}"]`)
          .classList.add("incorrect");
      }
      resultEl.hidden = false;
    };
  }

  function renderBuilder(step) {
    console.log("RENDER BUILDER CALLED");
    console.log("Step:", step);
    console.log("Step type:", step.type);
    console.log("Step content:", step.content);
    console.log("Patterns:", step.content?.patterns);

    const { base, patterns, goal } = step.content;

    const patternsHtml = patterns
      .map(
        (p) => `
      <button class="pattern-btn" data-id="${p.id}" data-file="${p.file}" data-symbol="${p.symbol || "🎨"}">
        ${p.name}
      </button>
    `,
      )
      .join("");

    els.container.innerHTML = `
      <div class="step-builder">
        <h2>${step.title}</h2>
        <p>${goal}</p>
        <div id="builder-canvas" class="builder-canvas">
          ${base && base !== "empty" ? `<span style="font-size:3rem;">${base}</span>` : `<span style="color:#94a3b8;">Нажмите на элементы ниже</span>`}
        </div>
        <div class="patterns-palette">${patternsHtml}</div>
        <div id="builder-hint" class="hint" hidden>Нажмите "Далее" для проверки</div>
      </div>
    `;

    document.querySelectorAll(".pattern-btn").forEach((btn) => {
      btn.onclick = () => {
        btn.classList.toggle("active");
        const id = btn.dataset.id;

        if (state.builderSelection.includes(id)) {
          state.builderSelection = state.builderSelection.filter(
            (p) => p !== id,
          );
        } else {
          state.builderSelection.push(id);
        }

        const canvas = document.getElementById("builder-canvas");
        canvas.innerHTML = state.builderSelection
          .map((selId) => {
            const item = patterns.find((p) => p.id === selId);
            return `<img src="${item.file}" alt="${item.name}" style="width:80px;height:80px;object-fit:contain;margin:4px;" onerror="this.style.display='none'; this.nextElementSibling.style.display='inline'"><span style="display:none;font-size:48px;">${item.symbol}</span>`;
          })
          .join("");
      };
    });

    els.nextBtn.onclick = () => {
      if (state.builderSelection.includes("solnce")) {
        document.getElementById("builder-canvas").style.background = "#f0fdf4";
        document.getElementById("builder-canvas").innerHTML +=
          "<div style='color:#22c55e; font-weight:bold; margin-top:8px;'>Отлично!</div>";
        els.nextBtn.onclick = goToNext;
      } else {
        document.getElementById("builder-hint").textContent =
          "💡 Подсказка: обязательно добавьте узор «Солнце»";
        document.getElementById("builder-hint").hidden = false;
        document.getElementById("builder-canvas").style.background = "#fef2f2";
      }
    };
  }

  function goToNext() {
    if (state.currentIndex < state.quest.steps.length - 1) {
      state.currentIndex++;
      renderStep(state.quest.steps[state.currentIndex]);
    } else {
      els.container.innerHTML = `
        <div class="step-complete">
          <h2>Квест завершён!</h2>
          <p>Вы прошли «${state.quest.title}»</p>
          <button class="btn btn-primary" onclick="window.location.href='/'">На главную</button>
        </div>
      `;
      els.prevBtn.disabled = true;
      els.nextBtn.disabled = true;
    }
  }

  function goToPrev() {
    if (state.currentIndex > 0) {
      state.currentIndex--;
      renderStep(state.quest.steps[state.currentIndex]);
    }
  }

  function updateButtons() {
    els.prevBtn.disabled = state.currentIndex === 0;
    els.prevBtn.onclick = goToPrev;

    els.nextBtn.disabled = false;
    els.nextBtn.textContent =
      state.currentIndex === state.quest.steps.length - 1
        ? "Завершить"
        : "Далее";
  }

  function updateProgress() {
    const percent = ((state.currentIndex + 1) / state.quest.steps.length) * 100;
    els.progress.style.width = `${percent}%`;
  }

  function renderMapMatch(step) {
    console.log("renderMapMatch START");
    const { goal, toys, regions } = step.content;
    state.matchedToys = {};
    state.selectedToy = null;

    els.container.innerHTML = `
    <div class="step-map-match">
      <h2>${step.title}</h2>
      <p>${goal}</p>
      
      <div class="map-match-layout">
        <!-- Панель с игрушками -->
        <div class="toys-panel">
          <h3>Игрушки:</h3>
          <div class="toys-list">
            ${toys
              .map(
                (toy) => `
              <div class="toy-card" draggable="true" data-id="${toy.id}">
                <img src="${toy.image}" alt="${toy.name}" />
                <div class="toy-info">
                  <strong>${toy.name}</strong>
                  <p>${toy.description}</p>
                </div>
              </div>
            `,
              )
              .join("")}
          </div>
        </div>

        <!-- Карта -->
        <div class="map-container" id="match-map"></div>
      </div>

      <div id="match-result" class="match-result" hidden></div>
    </div>
  `;

    const map = L.map("match-map").setView([55.5, 40.0], 6);
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 8,
      minZoom: 5,
    }).addTo(map);

    const regionMarkers = {};
    regions.forEach((region) => {
      const marker = L.circleMarker(region.center, {
        radius: 30,
        fillColor: "#e2e8f0",
        color: "#94a3b8",
        weight: 2,
        fillOpacity: 0.3,
      }).addTo(map);

      marker.bindTooltip(region.name, { permanent: true, direction: "center" });
      marker.bindPopup(`<strong>${region.name}</strong>`);
      regionMarkers[region.id] = marker;

      marker.on("dragover", (e) => e.preventDefault());
      marker.on("drop", (e) => {
        e.preventDefault();
        const toyId = e.dataTransfer.getData("text/plain");
        handleMatch(toyId, region.id);
      });

      marker.on("click", () => {
        if (state.selectedToy) {
          handleMatch(state.selectedToy, region.id);
        }
      });
    });

    document.querySelectorAll(".toy-card").forEach((card) => {
      const toyId = card.dataset.id;

      card.addEventListener("dragstart", (e) => {
        e.dataTransfer.setData("text/plain", toyId);
        card.classList.add("dragging");
      });
      card.addEventListener("dragend", () => card.classList.remove("dragging"));

      card.addEventListener("click", () => {
        document
          .querySelectorAll(".toy-card")
          .forEach((c) => c.classList.remove("selected"));
        card.classList.add("selected");
        state.selectedToy = toyId;
      });
    });

    function handleMatch(toyId, regionId) {
      const toy = toys.find((t) => t.id === toyId);
      const region = regions.find((r) => r.id === regionId);
      const isCorrect = region.correct_toys.includes(toyId);

      if (isCorrect) {
        state.matchedToys[toyId] = regionId;

        const marker = regionMarkers[regionId];
        marker.setStyle({ fillColor: "#22c55e", fillOpacity: 0.6 });
        marker.setPopupContent(
          `<strong>${region.name}</strong><br>✓ ${toy.name}`,
        );
        marker.openPopup();

        const toyCard = document.querySelector(`[data-id="${toyId}"]`);
        toyCard.classList.add("matched");
        toyCard.draggable = false;

        if (Object.keys(state.matchedToys).length === toys.length) {
          const resultEl = document.getElementById("match-result");
          resultEl.textContent = "🎉 Отлично! Все игрушки нашли свои регионы!";
          resultEl.className = "match-result success";
          resultEl.hidden = false;
          els.nextBtn.disabled = false;
          els.nextBtn.onclick = goToNext;
        }
      } else {
        const marker = regionMarkers[regionId];
        marker.setStyle({ fillColor: "#ef4444", fillOpacity: 0.6 });
        setTimeout(() => {
          marker.setStyle({ fillColor: "#e2e8f0", fillOpacity: 0.3 });
        }, 1000);

        const resultEl = document.getElementById("match-result");
        resultEl.textContent = `✗ ${toy.name} не из ${region.name}. Попробуй другой регион!`;
        resultEl.className = "match-result error";
        resultEl.hidden = false;
        setTimeout(() => {
          resultEl.hidden = true;
        }, 2000);
      }
    }

    els.nextBtn.disabled = true;
  }
});
