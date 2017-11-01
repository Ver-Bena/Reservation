// Admin
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func AdminIndex(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"admin_index.html",
		gin.H{},
	)
}

func ManageOrder(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"manage_order.html",
		gin.H{},
	)
}

func ManageMemberRsv_CfmPw(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"member_order_confirm_pw.html",
		gin.H{},
	)
}

func PostManageMemberRsv_CfmPw(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	getPw := session.Get("password")
	param_pw := getPw.(string)

	if getId != nil { // 회원일 때
		//로그인하고 있는 계정의 비밀번호와 사용자가 입력한 비밀번호가 일치하는지 확인
		password := c.PostForm("password")

		if password == param_pw { // 일치할 경우
			c.JSON(
				http.StatusOK,
				gin.H{
					"isMember": true,
				},
			)
		} else { // 불일치
			c.JSON(
				http.StatusOK,
				gin.H{
					"isMember": false,
				},
			)
		}
	} else { // 비회원일 경우

	}
}

func Member_Rsv_Filter(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"member_order_filter.html",
		gin.H{},
	)
}

var name string
var rsv_date string
var rsv_time string

func PostMember_Rsv_Filter(c *gin.Context) {
	name = c.PostForm("name")
	rsv_date = c.PostForm("rsv_date")
	rsv_time = c.PostForm("rsv_time")

	fmt.Println(rsv_time)
}

func ManageMemberRsv_Inquire(c *gin.Context) {
	count := 0

	query := "select * from members_order"

	if name != "" || rsv_date != "" || rsv_time != "" {
		query += " where "

		if name != "" {
			query += "name = '" + name + "'"
			count++
		}

		if rsv_date != "" {
			if count == 0 {
				query += "rsv_date = '" + rsv_date + ":00" + "'"
			} else {
				query += " and rsv_date = '" + rsv_date + ":00" + "'"
			}

			count++
		}

		if rsv_time != "" {
			if count == 0 {
				query += "rsv_time = '" + rsv_time + "'"
			} else {
				query += " and rsv_time = '" + rsv_time + "'"
			}

			count++
		}
	}

	query += ";"

	fmt.Println(query)

	rows, err := db.Query(query)

	var rsvs []Member_Rsv

	if err != nil {
		fmt.Print(err.Error() + "\n")
	}

	for rows.Next() {
		var rsv Member_Rsv
		err = rows.Scan(&rsv.Num, &rsv.Id, &rsv.Name, &rsv.Tel, &rsv.Grade, &rsv.Rsv_Date, &rsv.Rsv_Time, &rsv.People_Num, &rsv.Requests, &rsv.Sum)
		rsvs = append(rsvs, rsv)

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": rsvs,
		},
	)

	defer rows.Close()
}

func ManageMemberRsv(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"admin_member_rsv_list.html",
		gin.H{},
	)
}

func ManageMemberRsv_Delete(c *gin.Context) {
	num, _ := c.GetQuery("num")
	atoi_num, _ := strconv.Atoi(num)

	stmt, err := db.Prepare("delete from members_order where num = ?;")

	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(atoi_num)

	if err != nil {
		fmt.Print(err.Error())
	}

	stmt.Close()

	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}

func ManageAlienRsv_CfmPw(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"alien_order_confirm_pw.html",
		gin.H{},
	)
}

func PostManageAlienRsv_CfmPw(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	getPw := session.Get("password")
	param_pw := getPw.(string)

	if getId != nil { // 회원일 때
		//로그인하고 있는 계정의 비밀번호와 사용자가 입력한 비밀번호가 일치하는지 확인
		password := c.PostForm("password")

		if password == param_pw { // 일치할 경우
			c.JSON(
				http.StatusOK,
				gin.H{
					"isMember": true,
				},
			)
		} else { // 불일치
			c.JSON(
				http.StatusOK,
				gin.H{
					"isMember": false,
				},
			)
		}
	} else { // 비회원일 경우

	}
}

func Alien_Rsv_Filter(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"alien_order_filter.html",
		gin.H{},
	)
}

var a_name string
var a_rsv_date string
var a_rsv_time string

func PostAlien_Rsv_Filter(c *gin.Context) {
	a_name = c.PostForm("name")
	a_rsv_date = c.PostForm("rsv_date")
	a_rsv_time = c.PostForm("rsv_time")
}

func ManageAlienRsv_Inquire(c *gin.Context) {
	count := 0

	query := "select * from aliens_order"

	if a_name != "" || a_rsv_date != "" || a_rsv_time != "" {
		query += " where "

		if a_name != "" {
			query += "name = '" + a_name + "'"
			count++
		}

		if a_rsv_date != "" {
			if count == 0 {
				query += "rsv_date = '" + a_rsv_date + ":00" + "'"
			} else {
				query += " and rsv_date = '" + a_rsv_date + ":00" + "'"
			}

			count++
		}

		if a_rsv_time != "" {
			if count == 0 {
				query += "rsv_time = '" + a_rsv_time + "'"
			} else {
				query += " and rsv_time = '" + a_rsv_time + "'"
			}

			count++
		}
	}

	query += ";"

	rows, err := db.Query(query)

	var rsvs []Alien_Rsv

	if err != nil {
		fmt.Print(err.Error() + "\n")
	}

	for rows.Next() {
		var rsv Alien_Rsv
		err = rows.Scan(&rsv.Order_Num, &rsv.Name, &rsv.Tel, &rsv.Email, &rsv.Rsv_Date, &rsv.Rsv_Time, &rsv.People_Num, &rsv.Requests, &rsv.Sum)
		rsvs = append(rsvs, rsv)

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}
	}

	fmt.Println(rsvs)

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": rsvs,
		},
	)

	defer rows.Close()
}

func ManageAlienRsv(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"admin_alien_rsv_list.html",
		gin.H{},
	)
}

func ManageAlienRsv_Delete(c *gin.Context) {
	order_num, _ := c.GetQuery("order_num")
	atoi_order_num, _ := strconv.Atoi(order_num)

	stmt, err := db.Prepare("delete from aliens_order where order_num = ?;")

	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(atoi_order_num)

	if err != nil {
		fmt.Print(err.Error())
	}

	stmt.Close()

	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}

func ManageMember_CmfPW(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"member_confirm_pw.html",
		gin.H{},
	)
}

func PostManageMember_CmfPW(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	getPw := session.Get("password")
	param_pw := getPw.(string)

	if getId != nil { // 회원일 때
		//로그인하고 있는 계정의 비밀번호와 사용자가 입력한 비밀번호가 일치하는지 확인
		password := c.PostForm("password")

		if password == param_pw { // 일치할 경우
			c.JSON(
				http.StatusOK,
				gin.H{
					"isMember": true,
				},
			)
		} else { // 불일치
			c.JSON(
				http.StatusOK,
				gin.H{
					"isMember": false,
				},
			)
		}
	} else { // 비회원일 경우

	}
}

func MemberFilter(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"member_filter.html",
		gin.H{},
	)
}

var m_ID string
var m_name string
var m_birthday string

func PostMemberFilter(c *gin.Context) {
	m_ID = c.PostForm("id")
	m_name = c.PostForm("name")
	m_birthday = c.PostForm("birthday")
}

func ManageMember_Inquire(c *gin.Context) {
	count := 0

	query := "select * from members "

	if m_ID != "" || m_name != "" || m_birthday != "" {
		query += " where "

		if m_ID != "" {
			query += "id = '" + m_ID + "'"
			count++
		}

		if m_name != "" {
			if count == 0 {
				query += "name = '" + m_name + ":00" + "'"
			} else {
				query += " and name = '" + m_name + ":00" + "'"
			}

			count++
		}

		if m_birthday != "" {
			if count == 0 {
				query += "birthday = '" + m_birthday + "'"
			} else {
				query += " and birthday = '" + m_birthday + "'"
			}

			count++
		}
	}

	query += ";"

	rows, err := db.Query(query)

	var members []Member

	if err != nil {
		fmt.Print(err.Error() + "\n")
	}

	for rows.Next() {
		var member Member
		err = rows.Scan(&member.Id, &member.Name, &member.Tel, &member.Email, &member.Password, &member.Grade, &member.Birthday, &member.Rsv_Cnt)
		members = append(members, member)

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}
	}

	fmt.Println(members)

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": members,
		},
	)

	defer rows.Close()
}

func ManageMember(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"member_list.html",
		gin.H{},
	)
}

func ManageMember_Delete(c *gin.Context) {
	id, _ := c.GetQuery("id")

	stmt, err := db.Prepare("delete from members where id = ?;")

	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Print(err.Error())
	}

	stmt.Close()

	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}
