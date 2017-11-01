window.onload = init;
var server_vrf_code; // 서버에서 발생시킨 인증번호
var isVCAgree = false; // 인증번호가 일치하는가?
var curEmail; // 인증 당시의 이메일

function init() {
    document.getElementById("changeForm").onsubmit = Update;
    document.getElementById("send_vc").addEventListener("click", SendVrfCode);
    document.getElementById("cfm_vc").addEventListener("click", ConfirmVrfCode);
}
    
function SendVrfCode() { // 인증번호 전송 버튼을 누를 시 호출되는 함수
    curEmail = $("#email").val();
        
    if (curEmail == "") {
        alert("이메일을 입력해주세요!");
    }

    else {
        var json_email = `{ 
            "email": "`+curEmail+`"
        }`;

        var parse_email = JSON.parse(json_email);
            
        //서버에게 사용자가 입력한 이메일로 인증번호를 보낼 것을 요청함
            
        $.ajax({
            url: "/send_email",
            type: "GET",
            data: parse_email,
            contentType: "application/json",
            success: function(result) {
                if (result) {
                    alert("입력하신 이메일로 인증번호를 전송하였습니다! ");

                    $.each(result, function(key, value) {
                        server_vrf_code = value; // 서버에서 생성한 인증번호를 전달한다
                    });
                }

                else {
                    alert("이메일이 전송되지 않았습니다.");
                }
            }
        });

        return false;
    }
}

function ConfirmVrfCode() { // 인증번호 일치 여부 확인
    var vrf_code = $("#vrf_code").val();

    if (vrf_code == "") {
        alert("인증번호를 입력해주세요!");
    }

    else {
        if (server_vrf_code == vrf_code) { // 인증번호가 일치한다면
            alert("인증이 완료되었습니다!");
            isVCAgree = true;
        }

        else {
            alert("인증번호가 일치하지 않습니다!");
            isVCAgree = false;
        }
    }
}

function Update() { // 폼을 제출할 때 실행되는 함수
    //사용자에게 이름, 생년월일, 이메일, 비밀번호, 비밀번호 확인, 전화번호, 인증번호를 입력받는다.
    var password = $("#password").val();
    var password_cfm = $("#password_cfm").val();
    var name = $("#name").val();
    var email;
    var tel = $("#tel").val();
    var vrf_code = $("#vrf_code").val();
    var agree = $("#agree").val();
    
    //서버에 아이디를 보내어 DB에 사용자가 입력한 ID가 존재하는지 여부 판별

    if (password == password_cfm) { // 비밀번호와 비밀번호 확인이 일치하는가?
        document.getElementById("checkPW").innerHTML = "<font color = 'green'>비밀번호와 비밀번호 확인이 일치합니다.</font>";
    }

    else {
        document.getElementById("checkPW").innerHTML = "<font color = 'red'>비밀번호와 비밀번호 확인이 일치하지 않습니다.</font>";
    }
    
    if (isVCAgree) { // 인증번호를 입력하고 올바른 아이디를 입력했을 때
        id = $("#id").val();
        email = $("#email").val();

        //인증 당시의 ID, 이메일(curID, curEmail)과 현재 입력되어 있는 ID, 이메일(id, email)이 일치하는지 판별
        if (curID == id && curEmail == email) { // 같을 경우
            //서버에게 이름, 생년월일, 이메일, ID, 비밀번호, 전화번호를 서버로 전달하여 회원 개인정보 테이블에 추가할 것을 요청함

            var json_memberInfo = `{
                "password": "`+password+`",
                "name": "`+name+`",
                "tel": "`+tel+`",
                "email": "`+email+`"
            }`;
            
            var parse_memberInfo = JSON.parse(json_memberInfo);

            $.ajax({
                url: "/member_info/change",
                type: "POST",
                data: parse_memberInfo,
                contentType: "application/x-www-form-urlencoded",
                success: function(result) {
                    if (result) {
                        alert("개인정보 수정이 완료되었습니다!");
                        location.href = "/login"
                    }

                    else {
                        alert("개인정보 수정이 완료되지 않았습니다.")
                    }
                }
            })
        }

        else if (curEmail != email) { // email만 다른가?
            alert("이메일을 재인증해주세요!");
        }
    }

    else {
        alert("아이디가 중복이거나 인증번호가 올바르지 않습니다.");
    }

    return false;
}