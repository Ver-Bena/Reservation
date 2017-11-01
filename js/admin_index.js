window.onload = init;

function init() {
    document.getElementById("manageOrder").addEventListener("click", manageOrder);
    document.getElementById("manageMember").addEventListener("click", manageMember);
    document.getElementById("logout").addEventListener("click", logout);
}

function manageOrder() {
    location.href = "/manage_order";
}

function manageMember() {
    location.href = "/manage_member";
}

function logout() {
    location.href = "/logout";
}