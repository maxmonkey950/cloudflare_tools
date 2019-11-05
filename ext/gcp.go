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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func Upda(domains, cf string) {
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
	stmt, err := db.Prepare("update site.info set status=$1 where domain=$2")
	checkErr(err)

	stmt.Exec(domains, cf)
	defer stmt.Close()

}
