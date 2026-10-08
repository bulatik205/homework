(() => {
    const meEl = document.getElementById("me");
    const logoutBtn = document.getElementById("logout");
    const reloadBtn = document.getElementById("reload");
    const msg = document.getElementById("msg");
    const tbody = document.getElementById("tasks-body");

    function showMsg(text, ok) {
        msg.textContent = text;
        msg.className = "msg " + (ok ? "ok" : "err");
        if (ok) setTimeout(() => msg.classList.add("hidden"), 2000);
    }

    async function api(url, opts = {}) {
        const r = await fetch(url, {
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            ...opts,
        });
        let body = null;
        try { body = await r.json(); } catch { }
        if (r.status === 401) {
            location.href = "/auth";
            throw new Error("unauthorized");
        }
        return { status: r.status, body };
    }

    async function loadMe() {
        const { status, body } = await api("/api/auth/me");
        if (status === 200 && body?.success) {
            meEl.textContent = `${body.body.username} · ${body.body.role}`;
        } else {
            location.href = "/auth";
        }
    }

    async function loadTasks() {
        tbody.innerHTML = `<tr><td colspan="6" class="empty">Загрузка…</td></tr>`;
        const { status, body } = await api("/api/admin/tasks");
        if (status !== 200 || !body?.success) {
            showMsg(body?.error || `Ошибка (${status})`, false);
            tbody.innerHTML = `<tr><td colspan="6" class="empty">Не удалось загрузить</td></tr>`;
            return;
        }
        const tasks = body.body || [];
        if (!tasks.length) {
            tbody.innerHTML = `<tr><td colspan="6" class="empty">Пусто</td></tr>`;
            return;
        }
        tbody.innerHTML = tasks.map(t => `
      <tr>
        <td>${t.id}</td>
        <td>${escapeHtml(t.subject)}</td>
        <td>${escapeHtml(t.task)}</td>
        <td>${t.date_to ? escapeHtml(t.date_to) : ""}</td>
        <td>${t.instead_of ? `<span class="badge warn">${escapeHtml(t.instead_of)}</span>` : ""}</td>
        <td>${t.status ? escapeHtml(t.status) : ""}</td>
      </tr>
    `).join("");
    }

    function escapeHtml(s) {
        return String(s)
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll('"', "&quot;");
    }

    logoutBtn.addEventListener("click", async () => {
        await api("/api/auth/logout", { method: "POST" });
        location.href = "/auth";
    });

    reloadBtn.addEventListener("click", loadTasks);

    loadMe();
    loadTasks();
})();