package main

import (
	"app/internal/application"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	//fmt.Println(products)
	godotenv.Load()
	app := application.NewDefaultHTTP("localhost:8080")
	if err := app.Run(); err != nil {
		fmt.Println("Error al iniciar el servidor", err)
		return
	}

}
