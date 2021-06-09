package main

import (
	"flag"
	"fmt"
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

type CacheGroup struct {
	group *groupcache.Group
	peer *groupcache.HTTPPool
}


var apiCacheGroup *CacheGroup
var httpClient httpclient.Client


func GetData(c *gin.Context) {

	req := new(Req)
	err := c.ShouldBind(req)
	if err!=nil{
		c.String(http.StatusOK, err.Error())
		return
	}
	var b []byte
	apiCacheGroup.group.Get(c.Request.Context(), req.Key, groupcache.AllocatingByteSliceSink(&b))
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
	apiCacheGroup = new(CacheGroup)
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
	var cachePort = strconv.Itoa(i-1000)
	peers := groupcache.NewHTTPPoolOpts("http://127.0.0.1:" + cachePort, &groupcache.HTTPPoolOptions{
		Replicas: 1,
		BasePath: "/gouache/",
	})
	apiCacheGroup.peer = peers
	cacheGroup := groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			// get data from source
			param := map[string]string{
				"key":key,
			}
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
	apiCacheGroup.group = cacheGroup
	peers.Set("http://127.0.0.1:8001", "http://127.0.0.1:8002", "http://127.0.0.1:8003")
	log.Printf("cachegroup:%s slave starting on:127.0.0.1:%s\n",cacheGroup.Name(), cachePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s",cachePort),http.HandlerFunc(peers.ServeHTTP)))
}