window.onload = init;

function init() {
    document.getElementById("loginForm").onsubmit = Login;
}

function Login() {
    var id = $("#id").val();
    var password = $("#password").val();

    //서버에게 id, password를 전달하여 해당 ID와 비밀번호를 가진 계정으로 로그인할 것을 요청한다.

    var json_account = `{ 
            "id": "`+id+`",
            "password": "`+password+`"
        }`

    var parse_account = JSON.parse(json_account);

    $.ajax({
        url: "/login",
        type: "POST",
        data: parse_account,
        contentType: "application/x-www-form-urlencoded",
        complete: function(result) {
            parse_result = JSON.parse(result.responseText);

            if (parse_result.isAdmin) {
                location.href = "/admin_index";
            }
            
            else if (parse_result.isMember) {
                location.href = "/member_index";
            }

            else {
                alert("ID나 비밀번호가 일치하지 않습니다!");
                return;
            }
        }
    })

    return false;
}