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
        url: "/manage_member_order/confirm_pw",
        type: "POST",
        data: parse_pw,
        contentType: "application/x-www-form-urlencoded",
        success: function(result) {
            $.each(result, function(key, value) {
                if (value) {
                    location.href = "/manage_order/member_filter"
                }

                else {
                    alert("비밀번호가 일치하지 않습니다!");
                }
            })
        }
    })

    return false;
}