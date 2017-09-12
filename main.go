package main

import (
    "github.com/julienschmidt/httprouter"
    "flag"
    "net/http"
)

var (
    iniFile = flag.String("conf", "./config.ini", "配置文件路径")
    httpAddr = flag.String("addr", ":8080", "监听地址")
    config  *Config
)

func main() {
    flag.Parse()
    config, _ = NewConfig(*iniFile)

    router := httprouter.New()
    router.GET("/", func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        http.ServeFile(w, r, "home.html")
    })
    router.GET("/branches", decorate(apiBranches))
    router.GET("/pull", decorate(apiPull))
    router.GET("/checkout", decorate(apiCheckout))
    router.GET("/status", decorate(apiStatus))
    router.GET("/command", decorate(apiCommand))

    println("http server started on " + *httpAddr)
    http.ListenAndServe(*httpAddr, router)
}
