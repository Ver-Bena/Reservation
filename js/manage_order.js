window.onload = init;

function init() {
    document.getElementById("member").addEventListener("click", memberOrder);
    document.getElementById("alien").addEventListener("click", alienOrder);
}

function memberOrder() {
    location.href = "/manage_member_order/confirm_pw";
}

function alienOrder() {
    location.href = "/manage_alien_order/confirm_pw";
}