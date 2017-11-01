package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Member struct {
	Id       string
	Name     string
	Tel      string
	Email    string
	Password string
	Grade    string
	Birthday string
	Rsv_Cnt  string
}

type Member_Rsv struct {
	Num        string
	Id         string
	Name       string
	Tel        string
	Grade      string
	Rsv_Date   string
	Rsv_Time   string
	People_Num string
	Requests   string
	Sum        string
}

type Alien_Rsv struct {
	Order_Num  string
	Name       string
	Tel        string
	Email      string
	Rsv_Date   string
	Rsv_Time   string
	People_Num string
	Requests   string
	Sum        string
}

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

type Admin struct {
	Id       string
	Password string
}

var Rsvs = []*Member_Rsv{}
var email = ""
var db *sql.DB
var admin Admin

func main() {
	db, _ = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/Restaurant")

	defer db.Close()

	// make sure connection is available
	err := db.Ping()

	if err != nil {
		fmt.Print(err.Error())
	}

	rand.Seed(time.Now().UTC().UnixNano())

	router := gin.Default()
	store := sessions.NewCookieStore([]byte("secret")) // string형의 secret을 byte 배열로 변환하여 매개변수로 넘긴다 (secret이 key)
	router.Use(sessions.Sessions("mysession", store))  // store을 mysession이라는 이름으로 사용

	router.LoadHTMLGlob("template/*")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/images", "./images")

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	admin = Admin{}
	err = decoder.Decode(&admin)
	if err != nil {
		fmt.Println("error:", err)
	}

	//메인화면 - 회원모드/비회원모드 선택
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})

	//회원모드 - 회원가입 페이지
	router.GET("/join", Join)

	//회원가입 폼을 제출했을 때 (members 테이블에 insert)
	router.POST("/join", PostJoin)

	//회원가입 : 이메일 인증
	router.GET("/send_email", SendEmail)

	//회원모드 - 회원가입 : ID 중복 체크
	router.GET("/overlap_id", OverlapID)

	//회원모드 - 로그인
	router.GET("/login", Login)

	//회원모드 메인화면
	router.GET("/member_index", MemberIndex)

	//회원모드 - 로그인 (폼 제출 시)
	router.POST("/login", PostLogin)

	//회원모드 - 로그아웃
	router.POST("/logout", Logout)

	//회원모드 - 예약하기 : 회원의 이름, 전화번호 세팅
	router.GET("/member_rsv", MemberRsv)

	//회원모드 - 예약하기 : 총 합계 계산
	router.GET("/member_rsv/discount", MemberDiscount)

	//회원모드 - 예약하기 (폼 제출할 때)
	router.POST("/member_rsv", PostMemberRsv)

	//회원모드 - 예약조회 : 비밀번호 입력
	router.GET("/member_rsv/inquire_confirm_pw", Member_InquireRsv_CfmPw)

	//회원모드 - 예약조회 (폼 제출 시) : 비밀번호 체크
	router.POST("/member_rsv/inquire_confirm_pw", PostMember_InquireRsv_CfmPw)

	//회원모드 - 예약내역 준비
	router.GET("/member_rsv/inquire_json", InquireMemberOrder)

	//회원모드 - 예약 조회창
	router.GET("/member_rsv/inquire_list", InquireMemberOrderList)

	//회원모드 - 예약 취소
	router.GET("/member_rsv/delete", DeleteMemberOrder)

	//회원모드 - 예약 변경 : 기존 정보 전부 세팅
	router.GET("/member_rsv/change_setting", Setting_ChangeMemberOrder)

	//회원모드 - 예약 변경창
	router.GET("/member_rsv/change", ChangeMemberOrder)

	//회원모드 - 예약 변경창 (폼 제출 시)
	router.POST("/member_rsv/change", PostChangeMemberOrder)

	//회원모드 - 개인정보 수정 : 본인확인
	router.GET("/member_info/change_confirm_pw", Member_ChangeMemberInfo_CfmPw)

	//회원모드 - 개인정보 수정 : 본인확인 (폼 제출 시)
	router.POST("/member_info/change_confirm_pw", PostMember_ChangeMemberInfo_CfmPw)

	//회원모드 - 개인정보 수정 : 개인정보 수정 페이지
	router.GET("/member_info/change", Member_ChangeMemberInfo)

	//회원모드 - 개인정보 수정 : 개인정보 수정 페이지
	router.POST("/member_info/change", PostMember_ChangeMemberInfo)

	//회원모드 - 탈퇴 : 본인확인
	router.GET("/member_rsv/secession_confirm_pw", Member_SecessionRsv_CfmPw)

	//회원모드 - 탈퇴 : 본인확인 (폼 제출 후)
	router.POST("/member_rsv/secession_confirm_pw", PostMember_SecessionRsv_CfmPw)

	//비회원모드 메인화면
	router.GET("/alien_index", AlienIndex)

	//비회원모드 - 예약하기
	router.GET("/alien_rsv", AlienRsv)

	//비회원모드 - 예약하기 : 총 합계 계산
	router.GET("/alien_rsv/discount", AlienDiscount)

	//비회원모드 - 예약하기 (폼 제출할 때)
	router.POST("/alien_rsv", PostAlienRsv)

	//비회원모드 - 예약조회 : 비밀번호 입력
	router.GET("/alien_rsv/inquire_confirm_on", Alien_InquireRsv_CfmPw)

	//비회원 예약조회 (폼 제출 시) : 주문번호 체크
	router.POST("/alien_rsv/inquire_confirm_on", PostAlien_InquireRsv_CfmPw)

	//비회원 예약내역 준비
	router.GET("/alien_rsv/inquire_json", InquireAlienOrder)

	//비회원 예약 조회창
	router.GET("/alien_rsv/inquire_list", InquireAlienOrderList)

	//비회원 예약 취소
	router.GET("/alien_rsv/delete", DeleteAlienOrder)

	//비회원 예약 변경 : 기존 정보 전부 세팅
	router.GET("/alien_rsv/change_setting", Setting_ChangeAlienOrder)

	//비회원 예약 변경
	router.GET("/alien_rsv/change", ChangeAlienOrder)

	router.POST("/alien_rsv/change", PostChangeAlienOrder)

	//관리자모드 : 메인화면
	router.GET("/admin_index", AdminIndex)

	//관리자모드 : 예약 관리 - 회원/비회원 선택
	router.GET("/manage_order", ManageOrder)

	//관리자모드 : 예약 관리 - 회원 예약 관리를 위한 비밀번호 인증창
	router.GET("/manage_member_order/confirm_pw", ManageMemberRsv_CfmPw)

	//관리자모드 : 예약 관리 - 회원 예약 관리를 위한 비밀번호 인증창 (폼 제출 시)
	router.POST("/manage_member_order/confirm_pw", PostManageMemberRsv_CfmPw)

	//관리자모드 : 예약 관리 - 회원 예약 데이터 조회를 위한 조건 입력 폼
	router.GET("/manage_order/member_filter", Member_Rsv_Filter)

	//관리자모드 : 예약 관리 - 회원 예약 데이터 조회를 위한 조건 입력 폼 (폼 제출시)
	router.POST("/manage_order/member_filter", PostMember_Rsv_Filter)

	//관리자모드 : 예약 관리 - 회원 예약 데이터 조회
	router.GET("/manage_order/member_inquire", ManageMemberRsv_Inquire)

	//관리자모드 : 예약 관리 - 회원 예약 관리 페이지
	router.GET("/manage_order/member", ManageMemberRsv)

	//관리자모드 : 회원 예약 삭제
	router.GET("/manage_order/member_delete", ManageMemberRsv_Delete)

	//관리자모드 : 예약 관리 - 비회원 예약 관리를 위한 비밀번호 인증창
	router.GET("/manage_alien_order/confirm_pw", ManageAlienRsv_CfmPw)

	//관리자모드 : 예약 관리 - 비회원 예약 관리를 위한 비밀번호 인증창 (폼 제출 시)
	router.POST("/manage_alien_order/confirm_pw", PostManageAlienRsv_CfmPw)

	//관리자모드 : 예약 관리 - 비회원 예약 데이터 조회를 위한 조건 입력 폼
	router.GET("/manage_order/alien_filter", Alien_Rsv_Filter)

	//관리자모드 : 예약 관리 - 비회원 예약 데이터 조회를 위한 조건 입력 폼 (폼 제출시)
	router.POST("/manage_order/alien_filter", PostAlien_Rsv_Filter)

	//관리자모드 : 예약 관리 - 비회원 예약 데이터 조회
	router.GET("/manage_order/alien_inquire", ManageAlienRsv_Inquire)

	//관리자모드 : 예약 관리 - 비회원 예약 관리 페이지
	router.GET("/manage_order/alien", ManageAlienRsv)

	//관리자모드 : 비회원 예약 삭제
	router.GET("/manage_order/alien_delete", ManageAlienRsv_Delete)

	//관리자모드 : 회원관리 - 비밀번호 확인
	router.GET("/manage_member", ManageMember_CmfPW)

	//관리자모드 : 회원관리 - 비밀번호 확인 (폼 제출시)
	router.POST("/manage_member/confirm_pw", PostManageMember_CmfPW)

	//관리자모드 : 회원관리 - 회원 정보 조회를 위한 조건 입력폼
	router.GET("/manage_member/filter", MemberFilter)

	//관리자모드 : 회원관리 - 회원 정보 조회를 위한 조건 입력폼 (폼 제출시)
	router.POST("/manage_member/filter", PostMemberFilter)

	//관리자모드 : 회원관리 - 회원 정보 조회
	router.GET("/manage_member/inquire", ManageMember_Inquire)

	//관리자모드 : 회원관리 - 회원 관리 페이지
	router.GET("/manage_member/inquire_list", ManageMember)

	//관리자모드 : 회원 관리 - 회원 추방
	router.GET("/manage_order/delete", ManageMember_Delete)

	router.Run(":8000")
}
