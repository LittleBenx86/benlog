package main

import (
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/app/routers"
	_ "github.com/LittleBenx86/Benlog/internal/bootstrap"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"log"
)

func main() {
	server := routers.InitUEHttpServer()
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", variables.YmlAppConfig.GetString("App.Server.Port"))
	err := server.Listen(addr)
	if err != nil {
		log.Fatal("user server startup failed")
	}
}
