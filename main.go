package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Ipsave struct {
	id int
	ip string
}

var tpl *template.Template

var db *sql.DB

var ip string

func getMacAddr() ([]string, error) {

	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func main() {

	tpl, _ = template.ParseGlob("template/*.html")
	var err error

	as, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range as {
		fmt.Println("Mac address", a)
	}

	ip = as[4]
	fmt.Println("Mac 4th one address", ip)

	db, err = sql.Open("mysql", "root:ijse@tcp(localhost:3306)/ipsave")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/insert", insertHandler)
	http.ListenAndServe("localhost:8080", nil)

}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****insertHandler running******")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}
	r.ParseForm()
	ip := ip

	fmt.Println("Mac 4th one another one address :", ip)

	var err error
	if ip == "" {
		fmt.Println("Error Insert row:", err)
		tpl.ExecuteTemplate(w, "insert.html", "Error insert data, please check all fields")
		return
	}
	var ins *sql.Stmt
	var exists bool

	row := db.QueryRow("select * from `ipsave`.`ip2` where `ip`=?;", ip)

	if err := row.Scan(&ip, &exists); err != nil {

		var ip bool
		ip, err := strconv.ParseBool("true")
		ip = exists

		ins, err = db.Prepare("insert into `ipsave`.`ip2`(`ip`)values(?);")
		if err != nil {
			panic(err)
		}
		defer ins.Close()

		fmt.Println("checking exists 1 : ", exists)
		fmt.Println("checking ip 1 : ", ip)

		panic(err)

	} else if !exists {

		fmt.Println("===================== : ")
		fmt.Println("checking exists 2 : ", exists)
		fmt.Println("checking ip 2 : ", ip)

		ins, err = db.Prepare("insert into `ipsave`.`ip2`(`ip`)values(?);")
		if err != nil {
			panic(err)
		}
		defer ins.Close()

	}

	res, err := ins.Exec(ip)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		tpl.ExecuteTemplate(w, "insert.html", "mac address Successfully Inserted")
	}
}
