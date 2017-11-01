window.onload = init;

function init() {
    document.getElementById("addOrder").addEventListener("click", addOrder);
    document.getElementById("manageMyOrders").addEventListener("click", manageMyOrders);
}

function addOrder() {
    location.href = "/alien_rsv";
}

function manageMyOrders() {
    location.href = "/alien_rsv/inquire_confirm_on";
}