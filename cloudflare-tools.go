package main

// check the dnscords have useing proxy with shopify!
// author honux.micheal
import (
	"database/sql"
	"flag"
	"fmt"
	cloudflare "github.com/cloudflare/cloudflare-go"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type singleton map[string]string

type Datas struct {
	domains string
	jump    string
}

var (
	email   string
	api_key string
)

var (
	Old_data = New()
)

var (
	New_data = make(map[string]string)
	Up_data  = make(map[string]string)
)

var (
	once     sync.Once
	instance singleton
)

const (
	host     = "x.x.x.x"
	port     = xxxx
	user     = "xxxx"
	password = "xxxx"
	dbname   = "xxxx"
)

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

	stmt, err := db.Prepare("update public.dns set jump=$1  where domains=$2 and \"Status\" = 'A'")
	checkErr(err)

	stmt.Exec(domains, cf)
	defer stmt.Close()
}

func New() singleton {
	once.Do(func() {
		instance = make(singleton)
	})

	return instance
}

func getDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("successfull connected!")
	return db
}

func GetDatas(db *sql.DB) map[string]string {
	rows, err := db.Query("select domains, jump from public.dns where jump <> 'NULL' and \"Status\" = 'A'")
	if err != nil {
		log.Fatal(err)
	}
	datas := Datas{}
	for rows.Next() {
		err := rows.Scan(&datas.domains, &datas.jump)
		if err != nil {
			log.Fatal(err)
		}
		Old_data[datas.domains] = datas.jump
	}
	return Old_data
}

func GetCurFilename() string {
	_, fulleFilename, _, _ := runtime.Caller(0)
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fulleFilename)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix)

	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)

	return filenameOnly
}

func main() {
	db := getDB()
	GetDatas(db)
	//for k, v := range Old_data {
	//	fmt.Println(k, v)
	//}
	flag.StringVar(&email, "e", "", "email, default is nil")
	flag.StringVar(&api_key, "a", "", "api_key, default is nil")
	flag.Parse()
	//fmt.Println(api_key, email)
	var filenameOnly string = GetCurFilename()
	var logFilename string = filenameOnly + ".log"
	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE, 0777)
	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		fmt.Printf("open file error=%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logFile.Close()
	api, err := cloudflare.New(api_key, email)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		// Fetch all zones available to this user.
		zones, err := api.ListZones()
		if err != nil {
			log.Fatal(err)
		} else {
			for _, z := range zones {
				//fmt.Println(z.Name)
				zoneID, err := api.ZoneIDByName(z.Name)
				if err != nil {
					fmt.Println(err)
				}
				localhost := cloudflare.DNSRecord{Content: "23.227.38.32", Name: ("letter." + z.Name)}
				recs, err := api.DNSRecords(zoneID, localhost)
				if err != nil {
					log.Fatal(err)
					logger.Printf("%+v not shopify website", z.Name)
				} else {
					for _, r := range recs {
						if strconv.FormatBool(r.Proxied) == "false" {
							//fmt.Println("no proxy!")
							logger.Printf("%+v has no proxy", z.Name)
						} else {
							logger.Printf("%+v\n having proxy", z.Name)
							recsrule, err := api.ListPageRules(zoneID)
							if err != nil {
								log.Fatal(err)
								return
							} else {
								for _, v := range recsrule {
									if v.Status == "active" {
										m := v.Actions[0].Value.(map[string]interface{})
										logger.Println(z.Name, m["status_code"], m["url"], v.Status)
										ms := m["url"].(string)
										//fmt.Printf("%+v %+v %+v %+v\n", z.Name, m["status_code"], ms, v.Status)
										fmt.Println(z.Name, ms)
										New_data[z.Name] = ms
										//go_cf_postgres.Upda(ms, z.Name)
										//go_cf_postgres.Get_data()
									}
								}
							}
						}
					}
				}
			}
			//fmt.Println(Old_data)
			//fmt.Println(New_data)
			for k, _ := range New_data {
				//fmt.Println(k)
				if _, ok := Old_data[k]; ok {
					if strings.EqualFold(New_data[k], Old_data[k]) {
						fmt.Printf("%v, %v is in\n", k, New_data[k])
					} else {
						fmt.Printf("%v, %v has changed, I will update it!\n", k, New_data[k])
						Upda(New_data[k], k)
					}
				} else {
					fmt.Printf("%v %v not in, I will insert right now!\n", k, New_data[k])
					Upda(New_data[k], k)
				}
			}
		}
	}
}
