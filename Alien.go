// Alien
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	//"time"

	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func AlienIndex(c *gin.Context) {
	c.HTML(
		200,
		"alien_index.html",
		gin.H{},
	)
}

func AlienRsv(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"alien_rsv.html",
		gin.H{},
	)
}

func AlienDiscount(c *gin.Context) {
	people_cnt, _ := c.GetQuery("people_cnt")

	atoi_pc, _ := strconv.Atoi(people_cnt)
	sum := atoi_pc * 29900

	c.JSON(
		http.StatusOK,
		gin.H{
			"sum": sum,
		},
	)
}

func PostAlienRsv(c *gin.Context) {
	var alien_rsv Alien_Rsv

	alien_rsv.Name = c.PostForm("name")
	alien_rsv.Tel = c.PostForm("tel")
	alien_rsv.Email = c.PostForm("email")
	alien_rsv.Rsv_Date = c.PostForm("rsv_date")
	alien_rsv.Rsv_Time = c.PostForm("rsv_time")
	alien_rsv.People_Num = c.PostForm("people_cnt")
	alien_rsv.Requests = c.PostForm("requests")
	alien_rsv.Sum = c.PostForm("sum")

	if alien_rsv.Name != "" {
		atoi_sum, _ := strconv.Atoi(alien_rsv.Sum)

		for {
			alien_rsv.Order_Num = strconv.Itoa(rand.Intn(99999999-10000000) + 10000000)
			order_num := db.QueryRow("select order_num from aliens_order where order_num = ?;", alien_rsv.Order_Num)
			isExist := order_num.Scan(&alien_rsv.Order_Num)

			if isExist == sql.ErrNoRows {
				break
			}
		}

		stmt, err := db.Prepare("insert into aliens_order values(?, ?, ?, ?, ?, ?, ?, ?, ?);")

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}

		_, err = stmt.Exec(alien_rsv.Order_Num, alien_rsv.Name, alien_rsv.Tel, alien_rsv.Email, alien_rsv.Rsv_Date, alien_rsv.Rsv_Time, alien_rsv.People_Num, alien_rsv.Requests, atoi_sum)

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}

		emailUser := &EmailUser{"dlekrud0503@gmail.com", "good0503", "smtp.gmail.com", 587}

		auth := smtp.PlainAuth(
			"",
			emailUser.Username,
			emailUser.Password,
			emailUser.EmailServer,
		)

		msg := `*** 주문내역 ***` + "\n" +
			`주문번호 : ` + alien_rsv.Order_Num + "\n" +
			`이름 : ` + alien_rsv.Name + "\n" +
			`전화번호 : ` + alien_rsv.Tel + "\n" +
			`예약 날짜 : ` + alien_rsv.Rsv_Date + "\n" +
			`예약 시간 : ` + alien_rsv.Rsv_Time + "\n" +
			`인원수 : ` + alien_rsv.People_Num + "\n" +
			`요청사항 : ` + alien_rsv.Requests + "\n" +
			`예약비용 : ` + alien_rsv.Sum + "\n"

		err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
			auth,
			emailUser.Username,
			[]string{alien_rsv.Email},
			[]byte(msg))

		if err != nil {
			log.Print("ERROR!", err)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{},
		)

		defer stmt.Close()
	}
}

func Alien_InquireRsv_CfmPw(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"inquire_confirm_on.html",
		gin.H{},
	)
}

var param_order_num string

func PostAlien_InquireRsv_CfmPw(c *gin.Context) {
	param_order_num = c.PostForm("ordernum")
	var res string

	order_num := db.QueryRow("select order_num from aliens_order where order_num = ?;", param_order_num)
	isExist := order_num.Scan(&res)

	if isExist != sql.ErrNoRows { // 결과가 있다면
		c.JSON(
			http.StatusOK,
			gin.H{
				"isExist": true,
			},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"isExist": false,
			},
		)
	}

}

func InquireAlienOrder(c *gin.Context) {
	var alien_rsv Alien_Rsv

	rows := db.QueryRow("select * from aliens_order where order_num = ?;", param_order_num)
	_ = rows.Scan(&alien_rsv.Order_Num, &alien_rsv.Name, &alien_rsv.Tel, &alien_rsv.Email, &alien_rsv.Rsv_Date, &alien_rsv.Rsv_Time, &alien_rsv.People_Num, &alien_rsv.Requests, &alien_rsv.Sum)

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": alien_rsv,
		},
	)
}

func InquireAlienOrderList(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"alien_rsv_list.html",
		gin.H{},
	)
}

func DeleteAlienOrder(c *gin.Context) {
	order_num, _ := c.GetQuery("order_num")

	stmt, err := db.Prepare("delete from aliens_order where order_num = ?;")

	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(order_num)

	if err != nil {
		fmt.Print(err.Error())
	}

	stmt.Close()

	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}

var a_chg_order_num string
var a_chg_name string
var a_chg_tel string
var a_chg_email string
var a_chg_rsv_date string
var a_chg_rsv_time string
var a_chg_people_cnt string
var a_chg_requests string

func Setting_ChangeAlienOrder(c *gin.Context) {
	a_chg_order_num, _ = c.GetQuery("order_num")
	a_chg_name, _ = c.GetQuery("name")
	a_chg_tel, _ = c.GetQuery("tel")
	a_chg_email, _ = c.GetQuery("email")
	a_chg_rsv_date, _ = c.GetQuery("rsv_date")
	a_chg_rsv_time, _ = c.GetQuery("rsv_time")
	a_chg_people_cnt, _ = c.GetQuery("people_cnt")
	a_chg_requests, _ = c.GetQuery("requests")
}

func ChangeAlienOrder(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"alien_change_order.html",
		gin.H{
			"order_num":  a_chg_order_num,
			"name":       a_chg_name,
			"tel":        a_chg_tel,
			"email":      a_chg_email,
			"rsv_date":   a_chg_rsv_date,
			"rsv_time":   a_chg_rsv_time,
			"people_cnt": a_chg_people_cnt,
			"requests":   a_chg_requests,
		},
	)
}

func PostChangeAlienOrder(c *gin.Context) {
	email := c.PostForm("email")
	rsv_date := c.PostForm("rsv_date")
	rsv_time := c.PostForm("rsv_time")
	people_cnt := c.PostForm("people_cnt")
	requests := c.PostForm("requests")
	sum := c.PostForm("sum")

	if email != "" {
		atoi_sum, _ := strconv.Atoi(sum)

		_, err := db.Query("update aliens_order set email = ?, rsv_date = ?, rsv_time = ?, people_cnt = ?, requests = ?, sum = ? where order_num = ?;",
			email, rsv_date, rsv_time, people_cnt, requests, atoi_sum, a_chg_order_num)

		if err != nil {
			fmt.Print(err.Error() + "\n")
		}

		emailUser := &EmailUser{"dlekrud0503@gmail.com", "good0503", "smtp.gmail.com", 587}

		auth := smtp.PlainAuth(
			"",
			emailUser.Username,
			emailUser.Password,
			emailUser.EmailServer,
		)

		msg := `*** 주문내역 ***` + "\n" +
			`주문번호 : ` + a_chg_order_num + "\n" +
			`이름 : ` + a_chg_name + "\n" +
			`전화번호 : ` + a_chg_tel + "\n" +
			`이메일 : ` + email + "\n" +
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
	}
}
