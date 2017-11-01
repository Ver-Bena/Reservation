window.onload = init;

function init() {
    document.getElementById("pwForm").onsubmit = orderList;
}

function orderList() {
    //서버로 비밀번호를 보내어 값이 일치하는지 확인함
    var password = document.getElementById("password").value;

    var json_pw = `{ 
        "password": "`+password+`"
    }`;

    var parse_pw = JSON.parse(json_pw);

    $.ajax({
        url: "/member_rsv/secession_confirm_pw",
        type: "POST",
        data: parse_pw,
        contentType: "application/x-www-form-urlencoded",
        complete: function(result) {
            if (result.complete) {
                alert("탈퇴가 완료되었습니다!");
                location.href = "/";
            }

            else {
                alert("비밀번호가 일치하지 않습니다!");
            }
        }
    })

    return false;
}