package main

// check the dnscords have useing proxy with shopify!
// author honux.micheal
import (
	"flag"
	"fmt"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"log"
	"strconv"
)

func main() {
	var email string
	var api_key string
	flag.StringVar(&email, "e", "", "default is nil")
	flag.StringVar(&api_key, "a", "", "default is nil")
	flag.Parse()
	fmt.Println(api_key, email)
	api, err := cloudflare.New(api_key, email)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Fetch all zones available to this user.
	zones, err := api.ListZones()
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, z := range zones {
		//fmt.Println(z.Name)
		zoneID, err := api.ZoneIDByName(z.Name)
		//ct := "www." + z.Name
		//fmt.Println(ct)
		localhost := cloudflare.DNSRecord{Content: "35.203.158.174", Name: ("www." + z.Name)}
		if err != nil {
			fmt.Println(err)
			return
		}
		recs, err := api.DNSRecords(zoneID, localhost)
		if err != nil {
			log.Fatal(err)
			return
		} else {
			for _, r := range recs {
				if strconv.FormatBool(r.Proxied) == "false" {
					return
				} else {
					recsrule, err := api.ListPageRules(zoneID)
					if err != nil {
						log.Fatal(err)
						return
					} else {
						for _, v := range recsrule {
							if v.Status == "active" {
								m := v.Actions[0].Value.(map[string]interface{})
								fmt.Printf("%+v %+v %+v %+v\n", z.Name, m["status_code"], m["url"], v.Status)
								//fmt.Println(v.Actions[0].Value)
								//for _, vasd := range m {
								//	fmt.Printf("%+v %+v %+v\n", z.Name, vasd, m["url"])
								//}
							}
						}
					}
				}
			}
		}
	}
}
