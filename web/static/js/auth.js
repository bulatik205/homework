(() => {
    const tabLogin = document.getElementById("tab-login");
    const tabReg = document.getElementById("tab-register");
    const formLogin = document.getElementById("form-login");
    const formReg = document.getElementById("form-register");
    const title = document.getElementById("title");
    const msg = document.getElementById("msg");

    function showLogin() {
        tabLogin.classList.add("active");
        tabReg.classList.remove("active");
        formLogin.classList.remove("hidden");
        formReg.classList.add("hidden");
        title.textContent = "Вход";
        clearMsg();
    }

    function showRegister() {
        tabReg.classList.add("active");
        tabLogin.classList.remove("active");
        formReg.classList.remove("hidden");
        formLogin.classList.add("hidden");
        title.textContent = "Регистрация";
        clearMsg();
    }

    function showMsg(text, ok) {
        msg.textContent = text;
        msg.className = "msg " + (ok ? "ok" : "err");
    }

    function clearMsg() {
        msg.textContent = "";
        msg.className = "msg hidden";
    }

    async function postJSON(url, data) {
        const r = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "same-origin",
            body: JSON.stringify(data),
        });
        let body;
        try { body = await r.json(); } catch { body = { success: false, error: "bad response" }; }
        return { status: r.status, body };
    }

    tabLogin.addEventListener("click", showLogin);
    tabReg.addEventListener("click", showRegister);

    formLogin.addEventListener("submit", async (e) => {
        e.preventDefault();
        const fd = new FormData(formLogin);
        const { status, body } = await postJSON("/api/auth/login", {
            username: fd.get("username"),
            password: fd.get("password"),
        });
        if (body.success) {
            showMsg("Успех! Сейчас перенаправим…", true);
            setTimeout(() => location.href = "/", 800);
        } else {
            showMsg(body.error || `Ошибка (${status})`, false);
        }
    });

    formReg.addEventListener("submit", async (e) => {
        e.preventDefault();
        const fd = new FormData(formReg);
        const { status, body } = await postJSON("/api/auth/register", {
            username: fd.get("username"),
            password: fd.get("password"),
            key: fd.get("key"),
        });
        if (body.success) {
            showMsg("Аккаунт создан! Сейчас перенаправим…", true);
            setTimeout(() => location.href = "/", 800);
        } else {
            showMsg(body.error || `Ошибка (${status})`, false);
        }
    });
})();