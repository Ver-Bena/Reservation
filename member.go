// member
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Join(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"join.html",
		gin.H{},
	)
}

func PostJoin(c *gin.Context) {
	id := c.PostForm("id")
	password := c.PostForm("password")
	name := c.PostForm("name")
	birthday := c.PostForm("birthday")
	tel := c.PostForm("tel")
	email := c.PostForm("email")

	stmt, err := db.Prepare("insert into members (id, name, tel, email, password, birthday) values(?, ?, ?, ?, ?, ?);")

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	_, err = stmt.Exec(id, name, tel, email, password, birthday)

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	defer stmt.Close()

	c.JSON(
		http.StatusOK,
		gin.H{
			"complete": true,
		},
	)
}

func SendEmail(c *gin.Context) {
	email, _ = c.GetQuery("email")

	emailUser := &EmailUser{"dlekrud0503@gmail.com", "good0503", "smtp.gmail.com", 587}

	auth := smtp.PlainAuth(
		"",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)

	vrf_code := rand.Intn(9999-1000) + 1000
	msg := "인증번호는 " + strconv.Itoa(vrf_code) + " 입니다!"
	var err error

	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
		auth,
		emailUser.Username,
		[]string{email},
		[]byte(msg))

	if err != nil {
		log.Print("ERROR: attempting to send a mail ", err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"vrf_code": vrf_code,
		},
	)
}

func OverlapID(c *gin.Context) {
	id, _ := c.GetQuery("id")
	var overlap_id string

	check := db.QueryRow("select id from members where id = ?;", id) // 똑같은 아이디가 있는지 검색
	isOverlap := check.Scan(&overlap_id)                             // 검색 결과가 존재하는지 확인 (중복 여부 저장)

	if isOverlap == sql.ErrNoRows { // 중복이 아니라면
		c.JSON(
			http.StatusOK,
			gin.H{
				"isValid": true,
			},
		)
	} else { // 중복이라면
		c.JSON(
			http.StatusOK,
			gin.H{
				"isValid": false,
			},
		)
	}
}

func Login(c *gin.Context) {
	session := sessions.Default(c)

	getId := session.Get("id")

	if getId != nil { // 로그인하고 있을 때
		c.Redirect(http.StatusMovedPermanently, "/member_index")
	} else {
		c.HTML(
			http.StatusOK,
			"login.html",
			gin.H{},
		)
	}
}

func MemberIndex(c *gin.Context) {
	session := sessions.Default(c)
	var name string

	getId := session.Get("id")

	if getId == nil { // 로그인하지 않았을 때
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else { // 로그인했을 때
		id := getId.(string)

		accountName := db.QueryRow("select name from members where id = ?;", id)
		_ = accountName.Scan(&name)

		c.HTML(
			http.StatusOK,
			"member_index.html",
			gin.H{
				"id":      id,
				"name":    name,
				"isLogin": true,
			},
		)
	}
}

func PostLogin(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)

	var id string
	var password string
	var name string

	getId := session.Get("id")
	getPw := session.Get("password")

	if getId == nil && getPw == nil {
		id = c.PostForm("id")
		password = c.PostForm("password")
	} else {
		id = getId.(string)
		id = c.PostForm("id")
		password = getPw.(string)
		password = c.PostForm("password")
	}

	//selct하여 존재하는 ID, 비밀번호인지 검사한다

	accountName := db.QueryRow("select name from members where id = ? and password = ?;", id, password)
	isMember := accountName.Scan(&name)

	if id == admin.Id && password == admin.Password {
		session.Set("id", id)
		session.Set("password", password)
		session.Save()
		c.JSON(
			200,
			gin.H{
				"isAdmin": true,
			},
		)
	} else if isMember != sql.ErrNoRows { // 회원이 맞다면
		session.Set("id", id)
		session.Set("password", password)
		session.Save()
		c.JSON(
			200,
			gin.H{
				"isMember": true,
			},
		)
	} else { // 회원이 아니라면
		c.JSON(
			200,
			gin.H{
				"isMember": false,
			},
		)
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)

	getId := session.Get("id")

	if getId != nil { // 로그인하고 있을 때
		session.Delete("id")
		session.Delete("password")
		session.Save()
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func MemberRsv(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)

	getId := session.Get("id")

	if getId != nil { // 로그인하고 있을 때
		param_id := getId.(string)

		var member Member

		name := db.QueryRow("select name from members where id = ?;", param_id)
		_ = name.Scan(&member.Name)

		tel := db.QueryRow("select tel from members where id = ?;", param_id)
		_ = tel.Scan(&member.Tel)

		email := db.QueryRow("select email from members where id = ?;", param_id)
		_ = email.Scan(&member.Email)

		c.HTML(
			http.StatusOK,
			"member_rsv.html",
			gin.H{
				"name":  member.Name,
				"tel":   member.Tel,
				"email": member.Email,
			},
		)
	} else { // 비회원일 때
		c.HTML(
			http.StatusOK,
			"alien_rsv.html",
			gin.H{},
		)
	}
}

func MemberDiscount(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)

	getId := session.Get("id")
	people_cnt, _ := c.GetQuery("people_cnt")

	if getId != nil { // 회원이 로그인 중이라면
		param_id := getId.(string)

		var member_grade string
		var member_birth string
		var discount = 0.0

		rsv_date, _ := c.GetQuery("rsv_date")

		// 총 예약 횟수.생일을 DB에서 가져온다

		birthday := db.QueryRow("select birthday from members where id = ?;", param_id)
		_ = birthday.Scan(&member_birth)

		// 예약 횟수에 따라 등급 update

		_, err := db.Query("update members set grade = case when rsv_cnt > 60 then '다이아' when rsv_cnt > 30 then '골드' when rsv_cnt > 10 then '실버' else grade end")

		if err != nil {
			fmt.Print(err.Error())
		}

		grade := db.QueryRow("select grade from members where id = ?;", param_id)
		_ = grade.Scan(&member_grade)

		switch member_grade {
		case "실버":
			discount = 5

		case "골드":
			discount = 10

		case "다이아":
			discount = 25
		}

		sub_rsv_date := rsv_date[5:len(rsv_date)]
		sub_member_birth := member_birth[5:len(member_birth)]

		if sub_rsv_date == sub_member_birth { // 예약날짜가 생일인가?
			discount += 10
		}

		atoi_pc, _ := strconv.Atoi(people_cnt)
		sum := float64(atoi_pc) * 29900 * (float64(100-discount) / float64(100))

		c.JSON(
			http.StatusOK,
			gin.H{
				"sum": sum,
			},
		)
	} else { // 비회원이라면
		atoi_pc, _ := strconv.Atoi(people_cnt)
		sum := float64(atoi_pc) * 29900

		c.JSON(
			http.StatusOK,
			gin.H{
				"sum": sum,
			},
		)
	}
}

func PostMemberRsv(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")

	if getId != nil { // 회원일 때
		param_id := getId.(string)
		var member_rsv Member_Rsv

		id := db.QueryRow("select id from members where id = ?;", param_id)
		_ = id.Scan(&member_rsv.Id)

		name := db.QueryRow("select name from members where id = ?;", param_id)
		_ = name.Scan(&member_rsv.Name)

		tel := db.QueryRow("select tel from members where id = ?;", param_id)
		_ = tel.Scan(&member_rsv.Tel)

		grade := db.QueryRow("select grade from members where id = ?;", param_id)
		_ = grade.Scan(&member_rsv.Grade)

		member_rsv.Rsv_Date = c.PostForm("rsv_date")
		member_rsv.Rsv_Time = c.PostForm("rsv_time") + ":00"
		member_rsv.People_Num = c.PostForm("people_cnt")
		member_rsv.Requests = c.PostForm("requests")
		member_rsv.Sum = c.PostForm("sum")

		fmt.Println(member_rsv.Rsv_Date)

		if member_rsv.Rsv_Date != "" {
			stmt, err := db.Prepare("insert into members_order (id, name, tel, grade, rsv_date, rsv_time, people_cnt, requests, sum) values(?, ?, ?, ?, ?, ?, ?, ?, ?);")

			if err != nil {
				fmt.Print(err.Error())
			}

			_, err = stmt.Exec(member_rsv.Id, member_rsv.Name, member_rsv.Tel, member_rsv.Grade, member_rsv.Rsv_Date, member_rsv.Rsv_Time, member_rsv.People_Num, member_rsv.Requests, member_rsv.Sum)

			if err != nil {
				fmt.Print(err.Error())
			}

			//예약횟수 1 증가

			_, err = db.Query("update members set rsv_cnt = rsv_cnt+1 where id = ?", param_id)

			//이메일로 예약 내역 전송

			email := c.PostForm("email")

			emailUser := &EmailUser{"dlekrud0503@gmail.com", "good0503", "smtp.gmail.com", 587}

			auth := smtp.PlainAuth(
				"",
				emailUser.Username,
				emailUser.Password,
				emailUser.EmailServer,
			)

			msg := `*** 주문내역 ***` + "\n" +
				`이름 : ` + member_rsv.Name + "\n" +
				`ID : ` + member_rsv.Id + "\n" +
				`전화번호 : ` + member_rsv.Tel + "\n" +
				`등급 : ` + member_rsv.Grade + "\n" +
				`예약 날짜 : ` + member_rsv.Rsv_Date + "\n" +
				`예약 시간 : ` + member_rsv.Rsv_Time + "\n" +
				`인원수 : ` + member_rsv.People_Num + "\n" +
				`요청사항 : ` + member_rsv.Requests + "\n" +
				`예약비용 : ` + member_rsv.Sum + "\n"

			err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
				auth,
				emailUser.Username,
				[]string{email},
				[]byte(msg))

			if err != nil {
				log.Print("ERROR!", err)
				return
			}

			defer stmt.Close()

			c.JSON(
				http.StatusOK,
				gin.H{},
			)
		} else { // 회원이 아닐 때

		}
	}
}

func Member_InquireRsv_CfmPw(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"inquire_confirm_pw.html",
		gin.H{},
	)
}

func PostMember_InquireRsv_CfmPw(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	getPw := session.Get("password")

	if getId != nil { // 회원일 때
		param_pw := getPw.(string)

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

func InquireMemberOrder(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")

	var rsvs []Member_Rsv

	if getId != nil { // 회원일 때
		param_id := getId.(string)
		rows, err := db.Query("select * from members_order where id = ?;", param_id)

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}

		for rows.Next() {
			var rsv Member_Rsv
			err = rows.Scan(&rsv.Num, &rsv.Id, &rsv.Name, &rsv.Tel, &rsv.Grade, &rsv.Rsv_Date, &rsv.Rsv_Time, &rsv.People_Num, &rsv.Requests, &rsv.Sum)
			rsvs = append(rsvs, rsv)

			fmt.Println(rsv)

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
	} else { // 비회원일 때

	}
}

func InquireMemberOrderList(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"member_rsv_list.html",
		gin.H{},
	)
}

func DeleteMemberOrder(c *gin.Context) {
	num, _ := c.GetQuery("num")
	atoi_num, _ := strconv.Atoi(num)

	fmt.Println(num)

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

var chg_num string
var chg_id string
var chg_name string
var chg_tel string
var chg_grade string
var chg_rsv_date string
var chg_rsv_time string
var chg_people_cnt string
var chg_requests string

func Setting_ChangeMemberOrder(c *gin.Context) {
	chg_num, _ = c.GetQuery("num")
	chg_id, _ = c.GetQuery("id")
	chg_name, _ = c.GetQuery("name")
	chg_tel, _ = c.GetQuery("tel")
	chg_grade, _ = c.GetQuery("grade")
	chg_rsv_date, _ = c.GetQuery("rsv_date")
	chg_rsv_time, _ = c.GetQuery("rsv_time")
	chg_people_cnt, _ = c.GetQuery("people_cnt")
	chg_requests, _ = c.GetQuery("requests")
}

func ChangeMemberOrder(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"change_order.html",
		gin.H{
			"num":        chg_num,
			"id":         chg_id,
			"name":       chg_name,
			"tel":        chg_tel,
			"grade":      chg_grade,
			"rsv_date":   chg_rsv_date,
			"rsv_time":   chg_rsv_time,
			"people_cnt": chg_people_cnt,
			"requests":   chg_requests,
		},
	)
}

func PostChangeMemberOrder(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")

	if getId != nil { // 회원일 때
		rsv_date := c.PostForm("rsv_date")
		rsv_time := c.PostForm("rsv_time")
		people_cnt := c.PostForm("people_cnt")
		requests := c.PostForm("requests")
		sum := c.PostForm("sum")

		if rsv_date != "" {
			atoi_sum, _ := strconv.Atoi(sum)

			_, err := db.Query("update members_order set rsv_date = ?, rsv_time = ?, people_cnt = ?, requests = ?, sum = ? where num = ?;",
				rsv_date, rsv_time, people_cnt, requests, atoi_sum, chg_num)

			if err != nil {
				fmt.Print(err.Error() + "\n")
			}

			email := c.PostForm("email")

			emailUser := &EmailUser{"dlekrud0503@gmail.com", "good0503", "smtp.gmail.com", 587}

			auth := smtp.PlainAuth(
				"",
				emailUser.Username,
				emailUser.Password,
				emailUser.EmailServer,
			)

			msg := `*** 주문내역 ***` + "\n" +
				`이름 : ` + chg_name + "\n" +
				`ID : ` + chg_id + "\n" +
				`전화번호 : ` + chg_tel + "\n" +
				`등급 : ` + chg_grade + "\n" +
				`예약 날짜 : ` + rsv_date + "\n" +
				`예약 시간 : ` + rsv_time + "\n" +
				`인원수 : ` + people_cnt + "\n" +
				`요청사항 : ` + requests + "\n" +
				`예약비용 : ` + sum + "\n"

			err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
				auth,
				emailUser.Username,
				[]string{email},
				[]byte(msg))

			if err != nil {
				log.Print("ERROR!", err)
				return
			}

			c.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
				},
			)
		}
	} else {

	}
}

func Member_ChangeMemberInfo_CfmPw(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"change_memberinfo_confirm_pw.html",
		gin.H{},
	)
}

func PostMember_ChangeMemberInfo_CfmPw(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	getPw := session.Get("password")

	if getId != nil { // 회원일 때
		param_pw := getPw.(string)

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

func Member_ChangeMemberInfo(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	var member Member

	if getId != nil { // 회원일 때
		param_id := getId.(string)

		id := db.QueryRow("select id from members where id = ?;", param_id)
		_ = id.Scan(&member.Id)

		name := db.QueryRow("select name from members where id = ?;", param_id)
		_ = name.Scan(&member.Name)

		tel := db.QueryRow("select tel from members where id = ?;", param_id)
		_ = tel.Scan(&member.Tel)

		email := db.QueryRow("select email from members where id = ?;", param_id)
		_ = email.Scan(&member.Email)

		birthday := db.QueryRow("select birthday from members where id = ?;", param_id)
		_ = birthday.Scan(&member.Birthday)

		c.HTML(
			http.StatusOK,
			"change_memberInfo.html",
			gin.H{
				"id":       member.Id,
				"name":     member.Name,
				"tel":      member.Tel,
				"email":    member.Email,
				"birthday": member.Birthday,
			},
		)
	}
}

func PostMember_ChangeMemberInfo(c *gin.Context) {
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")

	if getId != nil { // 회원일 때
		//param_id := getId.(string)
		password := c.PostForm("password")
		name := c.PostForm("name")
		tel := c.PostForm("tel")
		email := c.PostForm("email")

		fmt.Println(password + " " + name + " " + tel + " " + email)

		//		_, err := db.Query("update members set name = ?, tel = ?, email = ?, password = ? where id = ?;",
		//			password, name, tel, email, param_id)

		//		if err != nil {
		//			fmt.Print(err.Error() + "\n")
		//		}
	}
}

func Member_SecessionRsv_CfmPw(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"secession_confirm_pw.html",
		gin.H{},
	)
}

func PostMember_SecessionRsv_CfmPw(c *gin.Context) {
	// 세션의 ID를 가져와 회원을 탈퇴함
	session := sessions.Default(c) // default 세션을 갖는다 (Session 리턴)
	getId := session.Get("id")
	getPw := session.Get("password")

	if getId != nil { // 회원일 때
		param_id := getId.(string)
		param_pw := getPw.(string)
		password := c.PostForm("password")

		if password == param_pw { // 일치할 경우
			stmt, err := db.Prepare("delete from members where id = ?;")

			if err != nil {
				fmt.Print(err.Error())
			}

			_, err = stmt.Exec(param_id)

			if err != nil {
				fmt.Print(err.Error())
			}

			stmt.Close()

			session.Delete("id")
			session.Delete("password")
			session.Save()

			c.JSON(
				http.StatusOK,
				gin.H{
					"complete": true,
				},
			)
		} else { // 불일치
			c.JSON(
				http.StatusOK,
				gin.H{
					"complete": false,
				},
			)
		}
	} else { // 회원이 아닐 때
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}
