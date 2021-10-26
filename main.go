package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "SMG"
	app.Version = "1.0.0"
	app.Email = "shanquan54@gmail.com"
	app.Usage = "什么鬼,将请求参数返回给请求者,用于检查网关修改数据过程。"
	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "启动SMG",
			Action: run,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "host",
					Value: "0.0.0.0",
					Usage: "host",
				},
				&cli.StringFlag{
					Name:  "port",
					Value: "11111",
					Usage: "port",
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(c *cli.Context) {
	engine := Engine{}

	address := fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))

	fmt.Printf("SMG Listen: %s\n", address)
	if err := http.ListenAndServe(address, &engine); err != nil {
		panic(err)
	}
}

type Engine struct {
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	reqBodyByte, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var body interface{}
	if req.Header.Get("Content-Type") == "application/json" {
		json.Unmarshal(reqBodyByte, &body)
	} else {
		body = string(reqBodyByte)
	}
	json.Unmarshal(reqBodyByte, &body)

	obj := map[string]interface{}{
		"Method": req.Method,
		"Header": req.Header,
		"Router": req.URL.Path,
		"Params": req.URL.Query(),
		"Body":   body,
	}
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(obj); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
