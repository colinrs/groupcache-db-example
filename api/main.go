package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/colinrs/pkgx/httpclient"
	"github.com/colinrs/pkgx/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache"
	"github.com/golang/groupcache/consistenthash"
)

type Req struct {
	Key string `form:"key" binding:"required"`
}

const (
	DataSource = "http://127.0.0.1:9005/get/"
)

type CacheGroup struct {
	group *groupcache.Group
	peerMap *consistenthash.Map
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

func WhereKey(c *gin.Context) {

	req := new(Req)
	err := c.ShouldBind(req)
	if err!=nil{
		c.String(http.StatusOK, err.Error())
		return
	}
	result := map[string]interface{}{
		"key": req.Key,
		"where": apiCacheGroup.peerMap.Get(req.Key),
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
	r.GET("/where", WhereKey)

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
	opt := &groupcache.HTTPPoolOptions{
		Replicas: 1,
		BasePath: "/gouache/",
	}
	cacheGroupHosts := []string{"http://127.0.0.1:8001", "http://127.0.0.1:8002", "http://127.0.0.1:8003"}
	peers := groupcache.NewHTTPPoolOpts("http://127.0.0.1:" + cachePort, opt)
	apiCacheGroup.peerMap = consistenthash.New(opt.Replicas, opt.HashFn)
	apiCacheGroup.peerMap.Add(cacheGroupHosts...)
	cacheGroup := groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			// get data from source
			param := map[string]string{
				"key":key,
			}
			var result string
			logger.Info("get key:%s data from:%s",key, DataSource)

			if err := httpClient.Get(ctx, DataSource, param, &result);err!=nil{
				logger.Info("get key:%s data err:%s", key, err.Error())
				return err
			}
			logger.Info("get key:%s data from result:%s",key, result)
			dest.SetBytes([]byte(result))
			return nil
		}))
	apiCacheGroup.group = cacheGroup
	peers.Set(cacheGroupHosts...)
	logger.Info("cachegroup:%s slave starting on:127.0.0.1:%s",cacheGroup.Name(), cachePort)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s",cachePort),http.HandlerFunc(peers.ServeHTTP)))
}