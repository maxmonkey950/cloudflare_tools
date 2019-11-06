package main

// check the dnscords have useing proxy with shopify!
// author honux.micheal
import (
	"flag"
	"fmt"
	cloudflare "github.com/cloudflare/cloudflare-go"
	go_cf_postgres "github.com/maxmonkey950/cloudflare_tools/ext"

	//go_cf_postgres "github.com/maxmonkey950/cloudflare_tools/ext"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

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

var email string
var api_key string

func main() {
	flag.StringVar(&email, "e", "", "email, default is nil")
	flag.StringVar(&api_key, "a", "", "api_key, default is nil")
	flag.Parse()
	fmt.Println(api_key, email)
	var filenameOnly string = GetCurFilename()
	var logFilename string = filenameOnly + ".log"
	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE, 0777)
	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		fmt.Printf("open file error=%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logFile.Close()
	//api, err := cloudflare.New("xxx", "xxx@gmail.com")
	//api, err := cloudflare.New("xxx", "xxx@xxx.com")
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
				fmt.Println(z.Name)
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
										fmt.Printf("%+v %+v %+v %+v\n", z.Name, m["status_code"], m["url"], v.Status)
										go_cf_postgres.Upda(ms, z.Name)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
