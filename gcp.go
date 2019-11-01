package go_cf_postgres

import (
"database/sql"
"fmt"

_ "github.com/lib/pq"
)

const (
	host     = "172.21.4.118"
	port     = 55432
	user     = "postgres"
	password = "123.com"
	dbname   = "website"
)

type Teacher struct {
	domain   string
	location string
}

func Inst(domains, cf string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	sqlStatement := `
INSERT INTO site.info (domain , location)
VALUES ($1, $2)
	RETURNING domain `
	var domain string
	//domain := "dcits.app"
	err = db.QueryRow(sqlStatement, domains, cf).Scan(&domain)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record is:", domain)

}
