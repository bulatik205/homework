(() => {
    const msg = document.getElementById("msg");
    const dateMain = document.getElementById("date-main");
    const dateSub = document.getElementById("date-sub");
    const prevBtn = document.getElementById("prev-day");
    const nextBtn = document.getElementById("next-day");
    const tomorrowBtn = document.getElementById("jump-tomorrow");
    const cardsEl = document.getElementById("cards");

    let currentDate = new Date();
    let subjectsMap = {};

    function fmtDate(d) {
        const y = d.getFullYear();
        const m = String(d.getMonth() + 1).padStart(2, "0");
        const day = String(d.getDate()).padStart(2, "0");
        return `${y}-${m}-${day}`;
    }

    function addDays(d, n) {
        const x = new Date(d);
        x.setDate(x.getDate() + n);
        return x;
    }

    function humanDate(d) {
        const days = ["вс", "пн", "вт", "ср", "чт", "пт", "сб"];
        const months = ["янв", "фев", "мар", "апр", "мая", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"];
        return `${days[d.getDay()]}, ${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
    }

    function escapeHtml(s) {
        return String(s ?? "")
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll('"', "&quot;");
    }

    function displayName(code) {
        return subjectsMap[code] || code;
    }

    function showMsg(text, ok) {
        msg.textContent = text;
        msg.className = "msg " + (ok ? "ok" : "err");
        if (ok) setTimeout(() => msg.classList.add("hidden"), 2500);
    }

    async function api(url) {
        const r = await fetch(url, { credentials: "same-origin" });
        let body = null;
        try { body = await r.json(); } catch { }
        return { status: r.status, body };
    }

    async function loadSubjects() {
        const { status, body } = await api("/api/v1/subjects");
        if (status === 200 && body?.success) {
            subjectsMap = {};
            for (const s of body.body || []) {
                subjectsMap[s.code] = s.display;
            }
        } else {
            console.warn("subjects load failed", status, body);
        }
    }

    function updateDateUI() {
        dateMain.textContent = humanDate(currentDate);

        const today = new Date();
        const diff = Math.round(
            (currentDate - new Date(today.getFullYear(), today.getMonth(), today.getDate())) / 86400000
        );

        if (diff === 0) dateSub.textContent = "сегодня";
        else if (diff === 1) dateSub.textContent = "завтра";
        else if (diff === -1) dateSub.textContent = "вчера";
        else if (diff > 1) dateSub.textContent = `через ${diff} дн.`;
        else dateSub.textContent = `${-diff} дн. назад`;
    }

    function setDate(d) {
        currentDate = d;
        updateDateUI();
        loadTasks();
    }

    prevBtn.addEventListener("click", () => setDate(addDays(currentDate, -1)));
    nextBtn.addEventListener("click", () => setDate(addDays(currentDate, 1)));
    tomorrowBtn.addEventListener("click", () => {
        const t = new Date();
        setDate(new Date(t.getFullYear(), t.getMonth(), t.getDate() + 1));
    });

    async function loadTasks() {
        cardsEl.innerHTML = `<div class="empty">Загрузка…</div>`;
        const dateStr = fmtDate(currentDate);
        const { status, body } = await api(`/api/v1/tasks?date=${dateStr}`);

        if (status !== 200 || !body?.success) {
            cardsEl.innerHTML = `<div class="empty">Не удалось загрузить</div>`;
            showMsg(body?.error || `Ошибка (${status})`, false);
            return;
        }
        renderCards(body.body || []);
    }

    function renderCards(tasks) {
        if (!tasks.length) {
            cardsEl.innerHTML = `<div class="empty">На этот день задач нет</div>`;
            return;
        }

        cardsEl.innerHTML = tasks.map(t => {
            const instead = t.instead_of
                ? `<span class="badge instead">вместо ${escapeHtml(displayName(t.instead_of))}</span>`
                : "";
            const status = t.status
                ? `<span class="badge">${escapeHtml(t.status)}</span>`
                : "";

            return `
        <div class="task-card">
          <div class="body">
            <div class="subject">
              <span class="badge subject">${escapeHtml(displayName(t.subject))}</span>
              ${instead}
              ${status}
            </div>
            <div class="task-text">${escapeHtml(t.task)}</div>
          </div>
        </div>
      `;
        }).join("");
    }

    async function init() {
        await loadSubjects();
        const today = new Date();
        setDate(new Date(today.getFullYear(), today.getMonth(), today.getDate() + 1));
    }

    init();
})();