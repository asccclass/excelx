package main

import(
   "time"
   "net/http"
   "os"
   "os/signal"
   "syscall"
   "context"
   log "github.com/sirupsen/logrus"
   "github.com/gorilla/mux"
   "github.com/asccclass/excelx/excel"
)

func main() {
   var err error
   defer func() {
      if err != nil {
         log.Printf("error in myHandler - error: %v", err)
         // w.WriteHeader(http.StatusInternalServerErrror)  // for web clien
      }
   }()

   port := os.Getenv("PORT")
   if port == "" { port = "80" }

   xls := excelSrv.Excelsrv{}

   // API
   router := mux.NewRouter()
   router.HandleFunc("/excel2json", xls.ParsefileFromWeb).Methods("POST")
   // router.HandleFunc("/csv2json", xls.ParseCSVfileFromWeb).Methods("POST")
   router.HandleFunc("/json2excel", xls.DownloadFileFromWeb).Methods("POST")
   router.HandleFunc("/healthz", xls.Healthz).Methods("GET")

   interrupt := make(chan os.Signal, 1)
   signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

   srv := &http.Server{
      Addr:    ":" + port,
      Handler: router,
      WriteTimeout: 15 * time.Second,
      ReadTimeout:  15 * time.Second,
   }

   log.Println("Server running at " + port)
   go func() {
      if err := srv.ListenAndServe(); err != nil {
         log.Printf("listen Error %v\n", err)
      }
   }()

   c := make(chan os.Signal, 1)
   signal.Notify(c, os.Interrupt)
   <-c
   ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
   defer cancel()
   srv.Shutdown(ctx)
   log.Println("\nshutting down...")
   os.Exit(0)
}
