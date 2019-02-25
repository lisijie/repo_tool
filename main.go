package main

import (
    "flag"
    "github.com/julienschmidt/httprouter"
    "log"
    "net/http"
)

var (
    Version   = "unknown"
    BuildTime = ""
    iniFile   = flag.String("conf", "./config.ini", "配置文件路径")
    httpAddr  = flag.String("addr", ":8080", "监听地址")
    config    *Config
)

func main() {
    flag.Parse()
    config, _ = NewConfig(*iniFile)
    log.SetFlags(log.Ldate | log.Ltime)

    router := httprouter.New()
    router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        http.ServeFile(w, r, "home.html")
    })
    router.GET("/branches", decorate(apiBranches))
    router.GET("/pull", decorate(apiPull))
    router.GET("/checkout", decorate(apiCheckout))
    router.GET("/clean", decorate(apiClean))
    router.GET("/status", decorate(apiStatus))
    router.GET("/command", decorate(apiCommand))

    log.Println("http server started on " + *httpAddr)
    if err := http.ListenAndServe(*httpAddr, router); err != nil {
        log.Fatalln(err)
    }
}
