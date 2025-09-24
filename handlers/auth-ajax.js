const errors = document.getElementById("errors");
const errorsT = document.getElementById("error-title");
const errorsP = document.getElementById("error-p");

function clearError() {
    errors.style.display = "none";
    errorsT.innerHTML = "";
    errorsP.innerHTML = "";
}

function reg() {
    clearError();
    switch (validatorReg()) {
        case "-1.x":
            errors.style.display = "block";
            errorsT.innerHTML = "Заполните все поля";
            errorsP.innerHTML = "Пожалуйста, заполните поля перед авторизацией";
            break;
        case "success":
            backEndReg();
            break;
        default:
            errors.style.display = "block";
            errorsT.innerHTML = "Неизвестная ошибка";
            errorsP.innerHTML = "Мы незнаем, почему это произошло";
    }
}

function log(type) {
    clearError();
    switch (validatorLog(type)) {
        case "-10":
            errors.style.display = "block";
            errorsT.innerHTML = "Неизвестная ошибка";
            errorsP.innerHTML = "JS получил невалидный type в значениии oncklick";
            break;
        case "-2.x":
            errors.style.display = "block";
            errorsT.innerHTML = "Заполните все поля";
            errorsP.innerHTML = "Пожалуйста, заполните поля перед авторизацией";
            break;
        case "success":
            backEndLog();
            break;
        default:
            errors.style.display = "block";
            errorsT.innerHTML = "Неизвестная ошибка";
            errorsP.innerHTML = "Мы незнаем, почему это произошло";
    }
}