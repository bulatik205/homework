const buttonR = document.getElementById("buttonR");
const buttonL = document.getElementById("buttonL");
const regElemet = document.getElementById("reg");
const logElemet = document.getElementById("log");

const buttonLP = document.getElementById("buttonLP");
const buttonLB = document.getElementById("buttonLB");
const logElemetP = document.getElementById("p");
const logElemetB = document.getElementById("b");

function navClick(navValue) {
    switch (navValue) {
        case "register":
            regElemet.style.display = "block";
            logElemet.style.display = "none";

            buttonR.classList.add("active-type");
            buttonR.classList.remove("close-type");

            buttonL.classList.remove("active-type");
            buttonL.classList.add("close-type");
            break;

        case "login":
            regElemet.style.display = "none";
            logElemet.style.display = "block";

            buttonL.classList.add("active-type");
            buttonL.classList.remove("close-type");

            buttonR.classList.remove("active-type");
            buttonR.classList.add("close-type");
            break;
    }
}

function logNavClick(logNavValue) {
    switch (logNavValue) {
        case "p":
            logElemetP.style.display = "flex";
            logElemetB.style.display = "none";

            buttonLP.classList.add("active-log-type");
            buttonLP.classList.remove("close-log-type");

            buttonLB.classList.remove("active-log-type");
            buttonLB.classList.add("close-log-type");
            break;

        case "b":
            logElemetB.style.display = "block";
            logElemetP.style.display = "none";

            buttonLP.classList.remove("active-log-type");
            buttonLP.classList.add("close-log-type");

            buttonLB.classList.add("active-log-type");
            buttonLB.classList.remove("close-log-type");
            break;
    }
}