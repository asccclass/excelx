package main

import(
   "os"
   "github.com/gorilla/mux"
   "github.com/asccclass/staticfileserver"
   "github.com/asccclass/serverstatus"
   "github.com/asccclass/excelx/libs/excel"
   "github.com/asccclass/excelx/libs/word"
)

// Create your Router function
func NewRouter(srv *SherryServer.ShryServer, documentRoot string)(*mux.Router, error) {
   router := mux.NewRouter()

   // excel Service
   excel, _ := excelSrv.NewExcelsrv(srv)
   excel.AddRouter(router)
   // word Service
   word, err := wordSrv.NewWordTemplate(srv,"")  // "" is file path
   if err != nil {
      return nil, err
   }
   word.AddRouter(router)
   
   //logger
   router.Use(SherryServer.ZapLogger(srv.Logger))

   // health check
   systemName := os.Getenv("SystemName")
   m := serverstatus.NewServerStatus(systemName)
   router.HandleFunc("/healthz", m.Healthz).Methods("GET")

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   router.PathPrefix("/").Handler(staticfileserver)

   return router, nil
}
