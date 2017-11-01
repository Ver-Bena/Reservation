window.onload = init;

function init() {
    document.getElementById("searchForm").onsubmit = Search;
}

function Search() { // 폼을 제출할 때 실행되는 함수
    //사용자에게 이름, 생년월일, 이메일, 비밀번호, 비밀번호 확인, 전화번호, 인증번호를 입력받는다
    var name = $("#name").val();
    var rsv_date = $("#rsv_date").val();
    var rsv_time = $("#rsv_time").val();

    var search_json = `{
         "name": "`+name+`",
         "rsv_date": "`+rsv_date+`",
         "rsv_time": "`+rsv_time+`"
    }`;

    var parse_search = JSON.parse(search_json);

    $.ajax({
        url: "/manage_order/member_filter",
        type: "POST",
        data: parse_search,
        contentType: "application/x-www-form-urlencoded",
        success: function(result) {
            location.href = "/manage_order/alien";
        }
    })

    return false;
}