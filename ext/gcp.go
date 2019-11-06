package go_cf_postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "x.x.x.x"
	port     = 5432
	user     = "postgres"
	password = "xxxx"
	dbname   = "cmdb"
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
	checkErr(err)

	defer db.Close()

	err = db.Ping()
	checkErr(err)

	stmt, err := db.Prepare("update public.dns set jump=$1, sitetype = 'shopify', location='cloudflare'  where domains=$2")
	checkErr(err)

	stmt.Exec(domains, cf)
	defer stmt.Close()

}
