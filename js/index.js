window.onload = init;

function init() {
    document.getElementById("member_mode").addEventListener("click", memberMode);
    document.getElementById("alien_mode").addEventListener("click", alienMode);
    document.getElementById("login").addEventListener("click", login);
    document.getElementById("join").addEventListener("click", join);
}

function memberMode() {
    location.href = "/member_index";
}

function alienMode() {
    location.href = "/alien_index";
}

function login() {
    location.href = "/login";
}

function join() {
    location.href = "/join";
}