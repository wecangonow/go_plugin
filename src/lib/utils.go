package lib

import(
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"github.com/astaxie/beego/context"
)


func IpToISOCode(ip_str string, db_path string) (string, error) {

	db, err := geoip2.Open(db_path)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer db.Close()
	ip := net.ParseIP(ip_str)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return record.Country.IsoCode, nil
}

func SetHeader(ctx *context.Context) {
	ctx.Output.Header("Content-Type","application/javascript; charset=UTF-8")
	ctx.Output.Header("Access-Control-Allow-Origin",AppConfig.AccessControllAllowOrigin)
	ctx.Output.Header("Access-Control-Allow-Methods","POST, GET, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept")
}