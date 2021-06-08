package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/colinrs/pkgx/httpclient"
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache"
)

type Req struct {
	Key string `form:"key" binding:"required"`
}

const (
	DataSource = "http://127.0.0.1:9005/get/"
)


var cacheGroup *groupcache.Group
var httpClient httpclient.Client


func GetData(c *gin.Context) {
	req := new(Req)
	err := c.ShouldBind(req)
	if err!=nil{
		c.String(http.StatusOK, err.Error())
		return
	}
	var b []byte
	cacheGroup.Get(c.Request.Context(), req.Key, groupcache.AllocatingByteSliceSink(&b))
	result := map[string]interface{}{
		"key": req.Key,
		"value": string(b),
	}
	c.JSON(http.StatusOK, result)
}

func InitHttpClient(){
	httpClient = httpclient.InitClient()
}

func main() {
	var port = flag.String("port", "9001", "api port")
	flag.Parse()
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/get", GetData)

	// Listen and serve on defined port
	log.Printf("Listening on port %s", *port)
	InitHttpClient()
	go InitCache(*port)

	r.Run(":" + *port)
}


func InitCache(port string) {
	i, err := strconv.Atoi(port)
	if err != nil {
		// handle error
		log.Fatal(err.Error())
	}
	var cachePort = ":" + strconv.Itoa(i-1000)
	peers := groupcache.NewHTTPPool("http://127.0.0.1:" + cachePort)

	cacheGroup = groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			// get data from source
			param := map[string]string{
				"key":key,
			}
			_, ok := peers.PickPeer(key)
			log.Printf("PickPeer:%+v\n", ok)
			var result string
			log.Printf("get key:%s data from:%s\n",key, DataSource)

			if err := httpClient.Get(ctx, DataSource, param, &result);err!=nil{
				log.Printf("get key:%s data err:%s", key, err.Error())
				return err
			}
			log.Printf("get key:%s data from result:%s\n",key, result)
			dest.SetBytes([]byte(result))
			return nil
		}))
	peers.Set("http://127.0.0.1:8001", "http://127.0.0.1:8002", "http://127.0.0.1:8003")

	log.Printf("cachegroup:%s slave starting on:%s\n",cacheGroup.Name(), cachePort)
	http.ListenAndServe("127.0.0.1:"+cachePort, http.HandlerFunc(peers.ServeHTTP))
}