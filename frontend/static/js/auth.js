document.getElementById("login-form")?.addEventListener("submit", async (e) => {
  e.preventDefault();
  const messageEl = document.getElementById("message");

  const data = {
    email: document.getElementById("login-email").value,
    password: document.getElementById("login-password").value,
  };

  try {
    const res = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
      credentials: "include",
    });

    const result = await res.json();

    if (res.ok) {
      messageEl.textContent = "Успешный вход!";
      messageEl.className = "message success";
      messageEl.classList.remove("hidden");
      setTimeout(() => (window.location.href = "/profile"), 1500);
    } else {
      messageEl.textContent = result.error || "Ошибка входа";
      messageEl.className = "message error";
      messageEl.classList.remove("hidden");
    }
  } catch (err) {
    messageEl.textContent = "Ошибка сети";
    messageEl.className = "message error";
    messageEl.classList.remove("hidden");
  }
});

document
  .getElementById("register-form")
  ?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const messageEl = document.getElementById("message");

    const data = {
      name: document.getElementById("reg-name").value,
      email: document.getElementById("reg-email").value,
      password: document.getElementById("reg-password").value,
    };

    try {
      const res = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      });

      const result = await res.json();

      if (res.ok) {
        messageEl.textContent = "Регистрация успешна!";
        messageEl.className = "message success";
        messageEl.classList.remove("hidden");
        setTimeout(() => (window.location.href = "/login.html"), 1500);
      } else {
        messageEl.textContent = result.error || "Проверьте введенные данные";
        messageEl.className = "message error";
        messageEl.classList.remove("hidden");
      }
    } catch (err) {
      messageEl.textContent = "Проверьте введенные данные";
      messageEl.className = "message error";
      messageEl.classList.remove("hidden");
    }
  });
