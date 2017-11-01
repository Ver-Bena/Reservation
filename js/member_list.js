window.onload = init;

function init() { // 서버에서 해당 계정의 모든 예약 조회, 그에 맞게 테이블 태그 생성
    $.ajax({
        url: "/manage_member/inquire",
        type: "GET",
        contentType: "application/json",
        complete: function(result) {
            if (result) {
                if (result.responseJSON.result != null) {
                        var table = $('<table></table>').addClass('table table-striped table-hover');

                        var head = $('<tr></tr>');
                        var row0 = $('<td></td>').text("ID");
                        var row1 = $('<td></td>').text("이름");
                        var row2 = $('<td></td>').text("전화번호");
                        var row3 = $('<td></td>').text("이메일");
                        var row4 = $('<td></td>').text("비밀번호");
                        var row5 = $('<td></td>').text("등급");
                        var row6 = $('<td></td>').text("생년월일");
                        var row7 = $('<td></td>').text("예약횟수");
                        var row8 = $('<td></td>').text("추방");

                        head.append(row0);
                        head.append(row1);
                        head.append(row2);
                        head.append(row3);
                        head.append(row4);
                        head.append(row5);
                        head.append(row6);
                        head.append(row7);
                        head.append(row8);

                        table.append(head);

                    for (var i = 0; i < result.responseJSON.result.length; i++) {
                        var col = $('<tr></tr>');
                        row0 = $('<td></td>').attr("id", "id"+i).text(result.responseJSON.result[i].Id);
                        row1 = $('<td></td>').attr("id", "name"+i).text(result.responseJSON.result[i].Name);
                        row2 = $('<td></td>').attr("id", "tel"+i).text(result.responseJSON.result[i].Tel);
                        row3 = $('<td></td>').attr("id", "email"+i).text(result.responseJSON.result[i].Email);
                        row4 = $('<td></td>').attr("id", "password"+i).text(result.responseJSON.result[i].Password);
                        row5 = $('<td></td>').attr("id", "grade"+i).text(result.responseJSON.result[i].Grade);
                        row6 = $('<td></td>').attr("id", "birthday"+i).text(result.responseJSON.result[i].Birthday);
                        row7 = $('<td></td>').attr("id", "birthday"+i).text(result.responseJSON.result[i].Rsv_Cnt);
                        row8 = $('<td></td>');

                        var deleteBtn = $('<button></button>').addClass('btn btn-danger').attr("id", "delete"+i).text("추방").click(DeleteOrder);

                        row7.append(deleteBtn);
                                
                        col.append(row0);
                        col.append(row1);
                        col.append(row2);
                        col.append(row3);
                        col.append(row4);
                        col.append(row5);
                        col.append(row6);
                        col.append(row7);
                        col.append(row8);

                        table.append(col);
                        $('#table').append(table);
                    }
                }

                else {
                    alert("예약내역이 존재하지 않습니다!");
                    location.href = "/admin_index";
                }
            }

            else {
                
            }
        }
    });

    return false;
}

function DeleteOrder(e) {
    var btn = e.target.id;
    var index = btn.substring(7, 6); // 라인넘버 추출
    var del_id = document.getElementById("id"+index).innerHTML;

    var json_delete = `{ 
        "id": "`+del_id+`"
    }`;

    var parse_delete = JSON.parse(json_delete);

    $.ajax({
        url: "/manage_order/delete",
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