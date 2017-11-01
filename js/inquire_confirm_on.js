window.onload = init;

function init() {
    document.getElementById("onForm").onsubmit = orderList;
}

function orderList() {
    //서버로 비밀번호를 보내어 값이 일치하는지 확인함
    var ordernum = document.getElementById("ordernum").value;

    var json_on = `{ 
        "ordernum": "`+ordernum+`"
    }`;

    var parse_on = JSON.parse(json_on);

    $.ajax({
        url: "/alien_rsv/inquire_confirm_on",
        type: "POST",
        data: parse_on,
        contentType: "application/x-www-form-urlencoded",
        success: function(result) {
            location.href = "/alien_rsv/inquire_list";
        } 
    })

    return false;
}