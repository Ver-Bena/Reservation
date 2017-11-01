window.onload = init;

function init() { // 서버에서 해당 계정의 모든 예약 조회, 그에 맞게 테이블 태그 생성
    $.ajax({
        url: "/manage_order/member_inquire",
        type: "GET",
        contentType: "application/json",
        success: function(result) {
            console.log(result);
            if (result) {
                var table = $('<table></table>').addClass('table table-striped table-hover');

                var head = $('<tr></tr>');
                var row0 = $('<td></td>').text("순번");
                var row1 = $('<td></td>').text("ID");
                var row2 = $('<td></td>').text("이름");
                var row3 = $('<td></td>').text("전화번호");
                var row4 = $('<td></td>').text("등급");
                var row5 = $('<td></td>').text("예약날짜");
                var row6 = $('<td></td>').text("예약시간");
                var row7 = $('<td></td>').text("인원수");
                var row8 = $('<td></td>').text("요청사항");
                var row9 = $('<td></td>').text("합계");
                var row10 = $('<td></td>').text("삭제");

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
                head.append(row10);

                table.append(head);

                if (result.result != null) {
                    for (var i = 0; i < result.result.length; i++) {
                        //console.log(result.result[0]);
                        var col = $('<tr></tr>');
                        row0 = $('<td></td>').attr("id", "num"+i).text(result.result[i].Num);
                        row1 = $('<td></td>').attr("id", "id"+i).text(result.result[i].Id);
                        row2 = $('<td></td>').attr("id", "name"+i).text(result.result[i].Name);
                        row3 = $('<td></td>').attr("id", "tel"+i).text(result.result[i].Tel);
                        row4 = $('<td></td>').attr("id", "grade"+i).text(result.result[i].Grade);
                        row5 = $('<td></td>').attr("id", "rsv_date"+i).text(result.result[i].Rsv_Date);
                        row6 = $('<td></td>').attr("id", "rsv_time"+i).text(result.result[i].Rsv_Time);
                        row7 = $('<td></td>').attr("id", "people_cnt"+i).text(result.result[i].People_Num);
                        row8 = $('<td></td>').attr("id", "requests"+i).text(result.result[i].Requests);
                        row9 = $('<td></td>').attr("id", "sum"+i).text(result.result[i].Sum);
                        row10 = $('<td></td>');

                        var deleteBtn = $('<button></button>').addClass('btn btn-danger').attr("id", "delete"+i).text("삭제").click(DeleteOrder);

                        row10.append(deleteBtn);
                        
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
                        col.append(row10);

                        table.append(col);
                    }

                    $('#table').append(table);
                }

                else {
                    alert("예약내역이 존재하지 않습니다!");
                    location.href = "/admin_index";
                }
            }
        }
    });

    return false;
}

function ChangeOrder(e) {
    var btn = e.target.id;
    var index = btn.substring(7, 6); // 라인넘버 추출

    // 해당 라인의 정보 추출
    
    var chg_num = document.getElementById("num"+index).innerHTML;
    var chg_id = document.getElementById("id"+index).innerHTML;
    var chg_name = document.getElementById("name"+index).innerHTML;
    var chg_tel = document.getElementById("tel"+index).innerHTML;
    var chg_grade = document.getElementById("grade"+index).innerHTML;
    var chg_rsv_date = document.getElementById("rsv_date"+index).innerHTML;
    var chg_rsv_time = document.getElementById("rsv_time"+index).innerHTML;
    var chg_people_cnt = document.getElementById("people_cnt"+index).innerHTML;
    var chg_requests = document.getElementById("requests"+index).innerHTML;
    var chg_sum = document.getElementById("sum"+index).innerHTML;

    var json_change = `{ 
        "num": "`+chg_num+`",
        "id": "`+chg_id+`",
        "name": "`+chg_name+`",
        "tel": "`+chg_tel+`",
        "grade": "`+chg_grade+`",
        "rsv_date": "`+chg_rsv_date+`",
        "rsv_time": "`+chg_rsv_time+`",
        "people_cnt": "`+chg_people_cnt+`",
        "requests": "`+chg_requests+`",
        "sum": "`+chg_sum+`"
    }`;

    var parse_change = JSON.parse(json_change);

    $.ajax({
        url: "/member_rsv/change_setting",
        type: "GET",
        data: parse_change,
        contentType: "application/json",
        success: function(result) {
            location.href = "/member_rsv/change";
        }
    })
}

function DeleteOrder(e) {
    var btn = e.target.id;
    var index = btn.substring(7, 6); // 라인넘버 추출
    var del_num = document.getElementById("num"+index).innerHTML;

    var json_delete = `{ 
        "num": "`+del_num+`"
    }`;

    var parse_delete = JSON.parse(json_delete);

    $.ajax({
        url: "/manage_order/member_delete",
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