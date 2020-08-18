package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/Ehco1996/v2scar"

	"github.com/tidwall/gjson"
)

var SYNC_TIME int

func main() {

	app := cli.NewApp()
	app.Name = "v2scar"
	app.Usage = "sidecar for V2ray"
	app.Version = "0.0.11"
	app.Author = "Ehco1996"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "grpc-endpoint, gp",
			Value:       "127.0.0.1:8080",
			Usage:       "V2ray开放的GRPC地址",
			EnvVar:      "V2SCAR_GRPC_ENDPOINT",
			Destination: &v2scar.GRPC_ENDPOINT,
		},
		cli.StringFlag{
			Name:        "api-endpoint, ap",
			Value:       "http://fun.com/api",
			Usage:       "django-sspanel开放的Vemss Node Api地址",
			EnvVar:      "V2SCAR_API_ENDPOINT",
			Destination: &v2scar.API_ENDPOINT,
		},
		cli.IntFlag{
			Name:        "sync-time, st",
			Value:       60,
			Usage:       "与django-sspanel同步的时间间隔",
			EnvVar:      "V2SCAR_SYNC_TIME",
			Destination: &SYNC_TIME,
		},
	}

	app.Action = func(c *cli.Context) error {
		getMyIp(&v2scar.IP)
		log.Println("本机的IP是: ", v2scar.IP)
		up := v2scar.NewUserPool()
		log.Println("Waitting v2ray start...")
		time.Sleep(time.Second * 3)
		tick := time.Tick(time.Duration(SYNC_TIME) * time.Second)
		for {
			go v2scar.SyncTask(up)
			<-tick
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func getMyIp(ip *string) {
	resp, err := http.Get("https://api.ip.sb/geoip")
	if err != nil {
		log.Println(err)
		*ip = ""
		return
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	if !gjson.Valid(string(res)) {
		log.Println("invalid json")
	}
	*ip = gjson.Get(string(res), `ip`).String()
}

