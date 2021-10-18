package main

import (
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/app/routers"
	_ "github.com/LittleBenx86/Benlog/internal/bootstrap"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"log"
)

func main() {
	server := routers.InitAEHttpServer()
	addr := fmt.Sprintf("%s:%s", "", variables.YmlAppConfig.GetString("App.Server.Port"))
	err := server.Listen(addr)
	if err != nil {
		log.Fatal("admin server startup failed")
	}
}
