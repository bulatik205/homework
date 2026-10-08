(() => {
    const meEl = document.getElementById("me");
    const logoutBtn = document.getElementById("logout");
    const msg = document.getElementById("msg");

    const dateMain = document.getElementById("date-main");
    const dateSub = document.getElementById("date-sub");
    const prevBtn = document.getElementById("prev-day");
    const nextBtn = document.getElementById("next-day");
    const tomorrowBtn = document.getElementById("jump-tomorrow");
    const cardsEl = document.getElementById("cards");
    const openCreateBtn = document.getElementById("open-create");

    const backdrop = document.getElementById("modal-backdrop");
    const modalTitle = document.getElementById("modal-title");
    const modalClose = document.getElementById("modal-close");
    const modalCancel = document.getElementById("modal-cancel");
    const form = document.getElementById("task-form");
    const fSubject = document.getElementById("f-subject");
    const fDateTo = document.getElementById("f-date-to");
    const fInsteadEnable = document.getElementById("f-instead-of-enable");
    const fInsteadWrap = document.getElementById("instead-of-wrap");
    const fInstead = document.getElementById("f-instead-of");
    const fTask = document.getElementById("f-task");
    const submitBtn = document.getElementById("submit-btn");
    const formHint = document.getElementById("form-hint");

    let subjects = [];
    let currentDate = new Date();
    let editingId = null;
    let editingTask = null;

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
        const s = subjects.find(x => x.code === code);
        return s ? s.display : code;
    }
    function showMsg(text, ok) {
        msg.textContent = text;
        msg.className = "msg " + (ok ? "ok" : "err");
        if (ok) setTimeout(() => msg.classList.add("hidden"), 2500);
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
    logoutBtn.addEventListener("click", async () => {
        await api("/api/auth/logout", { method: "POST" });
        location.href = "/auth";
    });

    async function loadSubjects() {
        const { status, body } = await api("/api/admin/subjects");
        if (status === 200 && body?.success) {
            subjects = body.body || [];
            fSubject.innerHTML = subjects
                .map(s => `<option value="${escapeHtml(s.code)}">${escapeHtml(s.display)}</option>`)
                .join("");
        }
    }

    async function fetchNextLessonDate(subject) {
        const { status, body } = await api(`/api/admin/schedule/next?subject=${encodeURIComponent(subject)}`);
        if (status === 200 && body?.success) return body.body.date;
        return null;
    }

    async function refreshInsteadOptions(dateStr) {
        const { status, body } = await api(`/api/admin/schedule/day?date=${dateStr}`);
        if (status === 200 && body?.success) {
            const lessons = body.body.lessons || [];
            fInstead.innerHTML = `<option value="">— не выбрано —</option>` +
                lessons.map(l => `<option value="${escapeHtml(l.code)}">${escapeHtml(l.display)}</option>`).join("");
        } else {
            fInstead.innerHTML = `<option value="">— нет уроков —</option>`;
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
        const { status, body } = await api(`/api/admin/tasks?date=${dateStr}`);
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
                ? `<span class="badge">${escapeHtml(t.status)}</span>` : "";
            return `
        <div class="task-card" data-id="${t.id}">
          <div class="body">
            <div class="subject">
              <span class="badge subject">${escapeHtml(displayName(t.subject))}</span>
              ${instead}
              ${status}
            </div>
            <div class="task-text">${escapeHtml(t.task)}</div>
            <div class="meta">
              <span>#${t.id}</span>
              ${t.date_to ? `<span>сдать: ${escapeHtml(t.date_to)}</span>` : ""}
              ${t.date_from ? `<span>выдано: ${escapeHtml(t.date_from)}</span>` : ""}
            </div>
          </div>
          <div class="actions">
            <button type="button" data-action="edit" data-id="${t.id}">Изменить</button>
            <button type="button" class="danger" data-action="delete" data-id="${t.id}">Удалить</button>
          </div>
        </div>
      `;
        }).join("");
    }

    cardsEl.addEventListener("click", async (e) => {
        const btn = e.target.closest("button[data-action]");
        if (!btn) return;
        const id = Number(btn.dataset.id);
        const action = btn.dataset.action;

        if (action === "delete") {
            if (!confirm(`Удалить задачу #${id}?`)) return;
            const { status, body } = await api(`/api/admin/tasks/${id}`, { method: "DELETE" });
            if (status === 200 && body?.success) {
                showMsg("Удалено", true);
                loadTasks();
            } else {
                showMsg(body?.error || `Ошибка (${status})`, false);
            }
            return;
        }

        if (action === "edit") {
            const dateStr = fmtDate(currentDate);
            const { status, body } = await api(`/api/admin/tasks?date=${dateStr}`);
            if (status !== 200 || !body?.success) return;
            const t = (body.body || []).find(x => x.id === id);
            if (!t) return;
            openEdit(t);
        }
    });

    function openModal() {
        backdrop.classList.remove("hidden");
        document.body.style.overflow = "hidden";
    }
    function closeModal() {
        backdrop.classList.add("hidden");
        document.body.style.overflow = "";
    }

    function openCreate() {
        editingId = null;
        editingTask = null;
        modalTitle.textContent = "Добавить ДЗ";
        submitBtn.textContent = "Добавить";
        form.reset();
        fInsteadEnable.checked = false;
        fInsteadWrap.classList.add("hidden");
        fInstead.innerHTML = "";

        const defaultSubject = subjects[0]?.code || "";
        fSubject.value = defaultSubject;
        fSubject.disabled = false;
        fDateTo.disabled = false;

        const fallbackDate = fmtDate(currentDate);
        fDateTo.value = fallbackDate;

        openModal();
        if (defaultSubject) {
            fetchNextLessonDate(defaultSubject).then(d => {
                if (d && !fDateTo.dataset.userModified) fDateTo.value = d;
            });
        }
    }

    function openCreate() {
        editingId = null;
        editingTask = null;
        modalTitle.textContent = "Добавить ДЗ";
        submitBtn.textContent = "Добавить";
        form.reset();
        fInsteadEnable.checked = false;
        fInsteadWrap.classList.add("hidden");
        fInstead.innerHTML = "";

        fSubject.disabled = false;
        fDateTo.disabled = false;
        fInsteadEnable.disabled = false;
        fInstead.disabled = false;

        delete fDateTo.dataset.userModified;

        const defaultSubject = subjects[0]?.code || "";
        fSubject.value = defaultSubject;

        const fallbackDate = fmtDate(currentDate);
        fDateTo.value = fallbackDate;

        openModal();
        if (defaultSubject) {
            fetchNextLessonDate(defaultSubject).then(d => {
                if (d && !fDateTo.dataset.userModified) fDateTo.value = d;
            });
        }
    }

    openCreateBtn.addEventListener("click", openCreate);
    modalClose.addEventListener("click", closeModal);
    modalCancel.addEventListener("click", closeModal);
    backdrop.addEventListener("click", (e) => {
        if (e.target === backdrop) closeModal();
    });

    fSubject.addEventListener("change", async () => {
        if (editingId) return;
        const subject = fSubject.value;
        if (!subject) return;
        const d = await fetchNextLessonDate(subject);
        if (d) fDateTo.value = d;
        if (fInsteadEnable.checked) refreshInsteadOptions(fDateTo.value);
    });

    fDateTo.addEventListener("change", () => {
        fDateTo.dataset.userModified = "1";
        if (fInsteadEnable.checked) refreshInsteadOptions(fDateTo.value);
    });

    fInsteadEnable.addEventListener("change", () => {
        if (fInsteadEnable.checked) {
            fInsteadWrap.classList.remove("hidden");
            refreshInsteadOptions(fDateTo.value);
        } else {
            fInsteadWrap.classList.add("hidden");
            fInstead.innerHTML = "";
        }
    });

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        submitBtn.disabled = true;

        let url = "/api/admin/tasks";
        let method = "POST";
        let payload;

        if (editingId) {
            url = `/api/admin/tasks/${editingId}`;
            method = "PATCH";
            payload = {
                subject: editingTask.subject,
                task: fTask.value.trim(),
                date_from: editingTask.date_from || null,
                date_to: editingTask.date_to || null,
                status: editingTask.status || null,
                instead_of: editingTask.instead_of || null,
            };
        } else {
            payload = {
                subject: fSubject.value,
                task: fTask.value.trim(),
                date_to: fDateTo.value || null,
                date_from: null,
                status: null,
                instead_of: fInsteadEnable.checked && fInstead.value ? fInstead.value : null,
            };
        }

        const { status, body } = await api(url, {
            method,
            body: JSON.stringify(payload),
        });

        submitBtn.disabled = false;

        if (status === 200 && body?.success) {
            showMsg(editingId ? "Обновлено" : "Добавлено", true);
            closeModal();
            if (body.body?.date_to) {
                const [y, m, d] = body.body.date_to.split("-").map(Number);
                setDate(new Date(y, m - 1, d));
            } else {
                loadTasks();
            }
        } else {
            showMsg(body?.error || `Ошибка (${status})`, false);
        }
    });

    async function init() {
        await loadMe();
        await loadSubjects();
        const today = new Date();
        setDate(new Date(today.getFullYear(), today.getMonth(), today.getDate() + 1));
    }

    init();
})();