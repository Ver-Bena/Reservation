window.onload = init;

function init() {
    document.getElementById("searchForm").onsubmit = Search;
}

function Search() { // 폼을 제출할 때 실행되는 함수
    //사용자에게 이름, 생년월일, 이메일, 비밀번호, 비밀번호 확인, 전화번호, 인증번호를 입력받는다
    var id = $("#id").val();
    var name = $("#name").val();
    var birthday = $("#birthday").val();

    var search_json = `{
         "id": "`+id+`",
         "name": "`+name+`",
         "birthday": "`+birthday+`"
    }`;

    var parse_search = JSON.parse(search_json);

    $.ajax({
        url: "/manage_member/filter",
        type: "POST",
        data: parse_search,
        contentType: "application/x-www-form-urlencoded",
        success: function(result) {
            location.href = "/manage_member/inquire_list";
        }
    })

    return false;
}