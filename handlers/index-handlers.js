const elements = {
    users: document.querySelector("#users"),
    admin: document.querySelector("#admin"),
    all: document.querySelector("#all"),
    premium: document.querySelector("#premium")
};

function updateDisplay(selector) {
    // Скрываем все элементы
    Object.values(elements).forEach(element => {
        element.classList.remove("link-open");
        element.classList.add("link-closed");
    });

    // Показываем выбранный элемент
    elements[selector].classList.remove("link-closed");
    elements[selector].classList.add("link-open");
}