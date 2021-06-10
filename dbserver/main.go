package main


import (
	"flag"
	"fmt"
	"net/http"

	"github.com/colinrs/groupcache-db-example/db"

	"github.com/colinrs/pkgx/logger"
	"github.com/gin-gonic/gin"
)

type DataSourceServer struct {
	db *db.DB
}

type Req struct {
	Key string `form:"key" binding:"required"`
}

type SetReq struct {
	Key string `form:"key" binding:"required"`
	Value string `form:"value" binding:"required"`
}

func GetData(c *gin.Context) {
	req := new(Req)
	err := c.ShouldBind(req)
	if err!=nil{
		c.String(http.StatusOK, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataSource.db.Get(req.Key))
}

func DelData(c *gin.Context) {
	req := new(Req)
	err := c.ShouldBind(req)
	if err!=nil{
		c.String(http.StatusOK, err.Error())
		return
	}
	dataSource.db.Del(req.Key)
	c.JSON(http.StatusOK, "OK")
}

func SetData(c *gin.Context) {
	setReq := new(SetReq)
	err := c.ShouldBind(setReq)
	if err!=nil{
		c.String(http.StatusOK, err.Error())
		return
	}
	dataSource.db.Set(setReq.Key, setReq.Value)
	c.JSON(http.StatusOK, setReq)
}

func CleanData(c *gin.Context) {
	dataSource.db = db.NewDB()
	c.JSON(http.StatusOK, "ok")
}


func LookDataset(c *gin.Context) {
	c.JSON(http.StatusOK, dataSource.db.Data())
}



var dataSource *DataSourceServer

func main() {
	var port = flag.String("port", "9005", "data source port")
	flag.Parse()
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/get", GetData)
	r.GET("/set", SetData)
	r.GET("/del", DelData)
	r.GET("/clean", CleanData)
	r.GET("/look", LookDataset)

	// Listen and serve on defined port
	InitServer()
	logger.Info("Listening on port %s", *port)
	r.Run(":" + *port)
}

func InitServer() {
	dataSource = &DataSourceServer{
		db: db.NewDB(),
	}
	for i:=0;i<100;i++{
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		dataSource.db.Set(key, value)
	}
}