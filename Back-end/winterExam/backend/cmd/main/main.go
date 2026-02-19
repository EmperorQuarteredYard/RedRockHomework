package main

import (
	"fmt"
	"homeworkSystem/backend/configs"
	"homeworkSystem/backend/internal/cli"
	"homeworkSystem/backend/internal/router"
	"log"
)

func main() {
	cfg := configs.LoadConfig()
	db := configs.InitDB(cfg)

	r := router.SetupRouter(db)

	cliExit := false
	go func() {
		if err := r.Run(":51443"); err != nil {
			log.Fatal("服务器启动失败:", err)
		}
		fmt.Print("webServe exit")
		cliExit = true
	}()
	ctl := cli.NewCLIController(db, &cliExit)
	ctl.Serve()
}
