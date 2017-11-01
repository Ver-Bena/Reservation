window.onload = init;

    function init() {
        document.getElementById("member_rsvForm").onsubmit = Reservation;
        document.getElementById("calculate").addEventListener("click", Calculate);
    }

    var curRsv_Date;
    var curPeople_Cnt;
    var isCalcul = false;

    function Calculate() {
        //할인율 계산 (해당 회원의 ID를 기준으로 등급, 생일 db에서 가져오기)
        curRsv_Date = $("#rsv_date").val();
        curPeople_Cnt = $("#people_cnt").val();

        if (curRsv_Date == "" || curPeople_Cnt == "") {
            if (curRsv_Date == "") {
                alert("예약날짜를 입력해주세요!");
            }

            if (curPeople_Cnt == "") {
                alert("인원수를 입력해주세요!");
            }
        }

        else {
            var json_discount = `{ 
                "rsv_date": "`+curRsv_Date+`",
                "people_cnt": "`+curPeople_Cnt+`"
            }`;

            var parse_discount = JSON.parse(json_discount);

            $.ajax({
                url: "/member_rsv/discount",
                type: "GET",
                data: parse_discount,
                contentType: "application/JSON",
                success: function(result) {
                    if (result) {
                        $.each(result, function(key, value) {
                            document.getElementById("sum").value = value;
                            isCalcul = true;
                        });
                    }

                    else {
                        alert("에러가 발생하였습니다. 잠시 후에 다시 시도하여 주세요.");
                    }
                }
            });

            return false;
        }

        //sum을 서버로 보내기
    }

    function Reservation() {
        //이메일, 인원수, 예약날짜, 예약시간, 요청사항 폼에서 받아오기
        var email = $("#email").val();
        var people_cnt = $("#people_cnt").val();
        var rsv_date = $("#rsv_date").val();
        var rsv_time = $("#rsv_time").val();
        var requests = $("#requests").val();
        var sum = $("#sum").val();

        if (curPeople_Cnt == people_cnt && curRsv_Date == rsv_date && isCalcul) {
            //당일로부터 최소 7일 이상인가?

            var d_rsv_date = new Date(rsv_date);
            var today = new Date();
            
            if (((d_rsv_date.getTime()-today.getTime())/1000/60/60/24) > 6) { // 당일로부터 7일 뒤인가?
                var json_rsv = `{ 
                    "email": "`+email+`",
                    "people_cnt": "`+people_cnt+`",
                    "rsv_date": "`+rsv_date+`",
                    "rsv_time": "`+rsv_time+`",
                    "requests": "`+requests+`",
                    "sum": "`+sum+`"
                }`;

                var parse_rsv = JSON.parse(json_rsv);

                $.ajax({
                    url: "/member_rsv",
                    type: "POST",
                    data: parse_rsv,
                    contentType: "application/x-www-form-urlencoded",
                    success: function(result) {
                        if (result) {
                            alert("예약이 완료되었습니다! 예약내역이 이메일로 전송되었습니다:)");
                            location.href = "/member_index";
                        }

                        else {
                            alert("에러가 발생하였습니다. 잠시 후에 다시 시도하여 주세요.");
                        }
                    }
                })

                return false;
            }

            else {
                $.ajax({
                    url: "/member_rsv",
                    type: "POST",
                    contentType: "application/x-www-form-urlencoded",
                    success: function(result) {
                        alert("예약은 당일로부터 최소 7일 후의 날짜로 예약해주세요!");
                    }
                })

                return false;
            }
        }

        else {
            $.ajax({
                url: "/member_rsv",
                type: "POST",
                contentType: "application/x-www-form-urlencoded",
                success: function(result) {
                    alert("비용 계산 버튼을 눌러주세요!");
                    isCalcul = false;
                }
            })

            return false;
        }
    }