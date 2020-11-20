package main

import (
	"os"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	db := sqlx.MustConnect("mysql", "fss:password@tcp(db:3306)/fss_db?parseTime=true")
	defer db.Close()

	m := model.NewModel(db)
	h := handler.NewHandler(m)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/signup", h.Signup)
	e.GET("/login", h.Login)

	g := e.Group("/user")
	config := middleware.JWTConfig{
		Claims: &handler.JWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	g.Use(middleware.JWTWithConfig(config))
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
	g.PUT("/folder", h.PutRenameFolder)
	g.DELETE("/folder", h.DeleteFolder)

	e.Logger.Fatal(e.Start(":8057"))
}