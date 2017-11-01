window.onload = init;

function init() { // 서버에서 해당 계정의 모든 예약 조회, 그에 맞게 테이블 태그 생성
    $.ajax({
        url: "/alien_rsv/inquire_json",
        type: "GET",
        contentType: "application/json",
        complete: function(result) {
            if (JSON.parse(result.responseText).result.Order_Num != "") {
                parse_result = JSON.parse(result.responseText);

                var table = $('<table></table>').addClass('table table-striped table-hover');

                var head = $('<tr></tr>');
                var row0 = $('<td></td>').text("주문번호");
                var row1 = $('<td></td>').text("이름");
                var row2 = $('<td></td>').text("전화번호");
                var row3 = $('<td></td>').text("이메일");
                var row4 = $('<td></td>').text("예약날짜");
                var row5 = $('<td></td>').text("예약시간");
                var row6 = $('<td></td>').text("인원수");
                var row7 = $('<td></td>').text("요청사항");
                var row8 = $('<td></td>').text("합계");
                var row9 = $('<td></td>').text(" 변경 / 삭제 ");

                head.append(row0);
                head.append(row1);
                head.append(row2);
                head.append(row3);
                head.append(row4);
                head.append(row5);
                head.append(row6);
                head.append(row7);
                head.append(row8);
                head.append(row9);

                table.append(head);

                if (parse_result) {
                    var col = $('<tr></tr>');
                    row0 = $('<td></td>').attr("id", "order_num").text(parse_result.result.Order_Num);
                    row1 = $('<td></td>').attr("id", "name").text(parse_result.result.Name);
                    row2 = $('<td></td>').attr("id", "tel").text(parse_result.result.Tel);
                    row3 = $('<td></td>').attr("id", "email").text(parse_result.result.Email);
                    row4 = $('<td></td>').attr("id", "rsv_date").text(parse_result.result.Rsv_Date);
                    row5 = $('<td></td>').attr("id", "rsv_time").text(parse_result.result.Rsv_Time);
                    row6 = $('<td></td>').attr("id", "people_cnt").text(parse_result.result.People_Num);
                    row7 = $('<td></td>').attr("id", "requests").text(parse_result.result.Requests);
                    row8 = $('<td></td>').attr("id", "sum").text(parse_result.result.Sum);
                    row9 = $('<td></td>');

                    var changeBtn = $('<button></button>').addClass('btn btn-success').attr("id", "change").text("변경").click(ChangeOrder);
                    var deleteBtn = $('<button></button>').addClass('btn btn-danger').attr("id", "delete").text("삭제").click(DeleteOrder);

                    row9.append(changeBtn);
                    row9.append(deleteBtn);
                        
                    col.append(row0);
                    col.append(row1);
                    col.append(row2);
                    col.append(row3);
                    col.append(row4);
                    col.append(row5);
                    col.append(row6);
                    col.append(row7);
                    col.append(row8);
                    col.append(row9);

                    table.append(col);
                }

                $('#table').append(table);
            }

            else {
                alert("예약내역이 존재하지 않습니다!");
                location.href = "/alien_index";
            }
        }
    });

    return false;
}

function ChangeOrder() {
    var chg_order_num = document.getElementById("order_num").innerHTML;
    var chg_name = document.getElementById("name").innerHTML;
    var chg_tel = document.getElementById("tel").innerHTML;
    var chg_email = document.getElementById("email").innerHTML;
    var chg_rsv_date = document.getElementById("rsv_date").innerHTML;
    var chg_rsv_time = document.getElementById("rsv_time").innerHTML;
    var chg_people_cnt = document.getElementById("people_cnt").innerHTML;
    var chg_requests = document.getElementById("requests").innerHTML;
    var chg_sum = document.getElementById("sum").innerHTML;

    var json_change = `{ 
        "order_num": "`+chg_order_num+`",
        "name": "`+chg_name+`",
        "tel": "`+chg_tel+`",
        "email": "`+chg_email+`",
        "rsv_date": "`+chg_rsv_date+`",
        "rsv_time": "`+chg_rsv_time+`",
        "people_cnt": "`+chg_people_cnt+`",
        "requests": "`+chg_requests+`",
        "sum": "`+chg_sum+`"
    }`;

    var parse_change = JSON.parse(json_change);

    $.ajax({
        url: "/alien_rsv/change_setting",
        type: "GET",
        data: parse_change,
        contentType: "application/json",
        success: function(result) {
            location.href = "/alien_rsv/change";
        }
    })
}

function DeleteOrder() {
    var del_order_num = document.getElementById("order_num").innerHTML;

    var json_delete = `{ 
        "order_num": "`+del_order_num+`"
    }`;

    var parse_delete = JSON.parse(json_delete);

    $.ajax({
        url: "/alien_rsv/delete",
        type: "GET",
        data: parse_delete,
        contentType: "application/json",
        success: function(result) {
            if (result) {
                alert("삭제가 완료되었습니다!");
                location.reload();
            }

            else {
                alert("에러가 발생하였습니다. 잠시 후에 다시 시도하여 주세요.");
            }
        }
    })
}