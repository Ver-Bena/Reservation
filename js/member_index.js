window.onload = init;

function init() {
    document.getElementById("addOrder").addEventListener("click", addOrder);
    document.getElementById("manageMyOrders").addEventListener("click", manageMyOrders);
    document.getElementById("changeInfo").addEventListener("click", changeInfo);
    document.getElementById("logout").onsubmit = logout;
    document.getElementById("secession").addEventListener("click", secession);
}

function addOrder() {
    location.href = "/member_rsv";
}

function manageMyOrders() {
    location.href = "/member_rsv/inquire_confirm_pw";
}

function changeInfo() {
    location.href = "/member_info/change_confirm_pw";
}

function logout() {
    $.ajax({
        url: "/logout",
        type: "POST",
        contentType: "application/JSON",
        success: function(result) { }
    })
}

function secession() {
    location.href = "/member_rsv/secession_confirm_pw";
}