package main

import (
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	cache "github.com/patrickmn/go-cache"

	"github.com/Tayu0404/file-sync-system-server/api/handler"
	"github.com/Tayu0404/file-sync-system-server/api/model"
)

func main() {
	db := sqlx.MustConnect("mysql", "fss:password@tcp(db:3306)/fss_db?parseTime=true")
	defer db.Close()

	snowflakeNode, _ := strconv.ParseInt(os.Getenv("SnowflakeNode"), 10, 64)

	m := model.NewModel(db)
	n, _ := snowflake.NewNode(snowflakeNode)
	c := cache.New(5*time.Minute, 30*time.Second)
	h := handler.NewHandler(m, n, c)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)

	g := e.Group("/user")
	config := middleware.JWTConfig{
		Claims: &handler.JWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	g.Use(middleware.JWTWithConfig(config))
	
	/*
	g.GET("/detail", h.Detail)
	g.GET("/sync", h.GetSync)
	g.PUT("/changepass", h.PutChangePass)

	g.GET("/list", h.GetList)

	g,GET("/file/:id", h.GetFileDetaile)
	g.PUT("/file/:id", h.PutRenameFile)
	g,GET("/file/download/:id", h.GetDownloadFile)
	g.POST("/file/upload/:id", h.PostUploadFile)
	g.PUT("/file/move/:id", h.PutMoveFile)
	g.DELETE("/file/delete/:id", h.DeleteFile)

	g.POST("/folder", h.PostCreateFolder)
	g.PUT("/folder/:id", h.PutRenameFolder)
	g.DELETE("/folder/:id", h.DeleteFolder)
	*/
	e.Logger.Fatal(e.Start(":8057"))
}