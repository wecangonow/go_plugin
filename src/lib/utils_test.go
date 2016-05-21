package lib

import(
	"testing"
	"github.com/astaxie/beego/config"
	"os"
	"math/rand"
	"fmt"
)

func TestJson(t *testing.T) {
	var jsoncontext = `{
"appname": "beeapi",
"testnames": "foo;bar",
"httpport": 8080,
"mysqlport": 3600,
"PI": 3.1415976,
"runmode": "dev",
"autorender": false,
"copyrequestbody": true,
"database": {
        "host": "host",
        "port": "port",
        "database": "database",
        "username": "username",
        "password": "password",
		"conns":{
			"maxconnection":12,
			"autoconnect":true,
			"connectioninfo":"info"
		}
    }
}`
	f, err := os.Create("testjson.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testjson.conf")
	jsonconf, err := config.NewConfig("json", "testjson.conf")
	jsonconf2, _ := config.NewConfig("json", "/Users/og/gocode/learning/src/adapi/conf/ad.conf")
	if err != nil {
		t.Fatal(err)
	}
	if jsonconf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}
	if jsonconf2.String("test") != "test" {
		t.Fatal("test not equal to test")
	}
	if jsonconf2.String("country::BR") != "巴西" {
		t.Fatal("get country::BR error")
	}
	if port, err := jsonconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if port, err := jsonconf.Int64("mysqlport"); err != nil || port != 3600 {
		t.Error(port)
		t.Fatal(err)
	}
	if pi, err := jsonconf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}
	if jsonconf.String("runmode") != "dev" {
		t.Fatal("runmode not equal to dev")
	}
	if v := jsonconf.Strings("unknown"); len(v) > 0 {
		t.Fatal("unknown strings, the length should be 0")
	}
	if v := jsonconf.Strings("testnames"); len(v) != 2 {
		t.Fatal("testnames length should be 2")
	}
	if v, err := jsonconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := jsonconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	if err = jsonconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if jsonconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
	if jsonconf.String("database::host") != "host" {
		t.Fatal("get database::host error")
	}
	if jsonconf.String("database::conns::connectioninfo") != "info" {
		t.Fatal("get database::conns::connectioninfo error")
	}
	if maxconnection, err := jsonconf.Int("database::conns::maxconnection"); err != nil || maxconnection != 12 {
		t.Fatal("get database::conns::maxconnection error")
	}
	if db, err := jsonconf.DIY("database"); err != nil {
		t.Fatal(err)
	} else if m, ok := db.(map[string]interface{}); !ok {
		t.Log(db)
		t.Fatal("db not map[string]interface{}")
	} else {
		if m["host"].(string) != "host" {
			t.Fatal("get host err")
		}
	}

	if _, err := jsonconf.Int("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int")
	}

	if _, err := jsonconf.Int64("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int64")
	}

	if _, err := jsonconf.Float("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Float")
	}

	if _, err := jsonconf.DIY("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an interface{}")
	}

	if val := jsonconf.String("unknown"); val != "" {
		t.Error("unknown keys should return an empty string when expecting a String")
	}

	if _, err := jsonconf.Bool("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Bool")
	}

	if !jsonconf.DefaultBool("unknow", true) {
		t.Error("unknown keys with default value wrong")
	}
}


func random(min, max int) int {

	return rand.Intn(max - min) + min
}
func Test_randrange(t *testing.T) {
	//rand.Seed(time.Now().Unix())

	for i := 1; i < 100; i++ {
		myrand := random(1,6)

		if 1 <= myrand && myrand <= 6 {
			t.Log("right")
		} else {
			t.Error("wrong")
		}
	}

}

func Test_adCount(t *testing.T) {
	countInfo := AdCountIndex{
		Uuid:"adbddsdsd",
		Ad_type:1,
		Ad_id:12 }
	num := GetAdCount(countInfo, "user")
	num2 := GetAdCount(countInfo, "ad")
	fmt.Println(" num is ", num)
	fmt.Println(" num2 is ", num)
	IncrementUserAdCountByOne(countInfo)
	IncrementUserAdCountByOne(countInfo)
	IncrementUserAdCountByOne(countInfo)
	IncrementUserAdCountByOne(countInfo)
	num = GetAdCount(countInfo, "user")
	num2 = GetAdCount(countInfo, "ad")

	fmt.Println(" num is ", num)
	fmt.Println(" num2 is ", num2)
}