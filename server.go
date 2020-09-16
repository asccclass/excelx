package main

import (
   "os"
   "fmt"
   "github.com/asccclass/staticfileserver"
)

func main() {
   port := os.Getenv("PORT")
   if port == "" {
      port = "80"
   }
   Production := os.Getenv("Production")
   if Production == "" {  // 預設為開發區
      Production = "Development"
   }
   // 取得 web 目錄
   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www"
   }
   templateRoot := os.Getenv("TemplateRoot")
   if templateRoot == "" {
      templateRoot = "www/template"
   }

   server, err := SherryServer.NewServer(":" + port, documentRoot, templateRoot)
   if err != nil {
      panic(fmt.Errorf("start server error:%s", err.Error()))
   }
   server.Server.Handler, err = NewRouter(server, documentRoot)
   if err != nil {
      panic(fmt.Errorf("start error:%s", err.Error()))
   } 
   server.Start()
}
