package main

import (
	"github.com/Azure/storage-blobs-go-quickstart/file_handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func init() {
	godotenv.Load()
}

func main() {
	e := echo.New()
	e.POST("/api/v1/file", file_handler.Upload)
	e.GET("/api/v1/fileurl/:cliente_id/:bc_id/:paciente_id/:filetype/:filename", file_handler.GetURL)
	e.GET("/api/v1/file/:cliente_id/:bc_id/:paciente_id/:file_type/:filename/", file_handler.Download)
	e.Logger.Fatal(e.Start(":8080"))
}
