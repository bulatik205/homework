function validatorReg() {
    const ur = document.getElementById("usernameR").value;
    const ure = document.getElementById("usernameR");
    const pr = document.getElementById("passwordR").value;
    const pre = document.getElementById("passwordR");
    const rr = document.getElementById("rules").checked;
    const rre = document.getElementById("rules");
    const rret = document.getElementById("rulesText");

    // Сброс стилей
    ure.style.borderColor = "#9368b6";
    pre.style.borderColor = "#9368b6";
    rre.style.boxShadow = "";
    rret.style.color = "#fff";

    let hasError = false;

    if (ur == "") {
        ure.style.borderColor = "rgb(255, 0, 0)";
        hasError = true;
    }

    if (pr == "") {
        pre.style.borderColor = "rgb(255, 0, 0)";
        hasError = true;
    }

    if (rr == false) {
        rre.style.boxShadow = "0px 0px 5px rgb(255, 0, 0)";
        rret.style.color = "rgb(255, 0, 0)";
        hasError = true;
    }

    return hasError ? "-1.x" : "success"; // возвращаем успех если ошибок нет
}

function validatorLog(type) {
    const ul = document.getElementById("usernameL").value;
    const ule = document.getElementById("usernameL");
    const pl = document.getElementById("passwordL").value;
    const ple = document.getElementById("passwordL");
    const tl = document.getElementById("telegramL").value;
    const tle = document.getElementById("telegramL");

    // Сброс стилей для всех полей
    ule.style.borderColor = "#9368b6";
    ple.style.borderColor = "#9368b6";
    tle.style.borderColor = "#9368b6";

    let hasError = false;

    switch (type) {
        case "password":
            if (ul == "") {
                ule.style.borderColor = "rgb(255, 0, 0)";
                hasError = true;
            }

            if (pl == "") {
                ple.style.borderColor = "rgb(255, 0, 0)";
                hasError = true;
            }
            break;
        case "telegram":
            if (tl == "") {
                tle.style.borderColor = "rgb(255, 0, 0)";
                hasError = true;
            }
            break;
        default:
            return "-10";
    }

    return hasError ? "-2.x" : "success";
}