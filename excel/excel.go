package excelSrv

import(
   "io"
   "os"
   "fmt"
   "reflect"
   "strings"
   "strconv"
   "path/filepath"
   "net/http"
   "encoding/json"
   log "github.com/sirupsen/logrus"
   "github.com/360EntSecGroup-Skylar/excelize"
)

type Excelsrv struct {
   Rows map[string]string
}

//將Json轉為Excel
func (xls *Excelsrv) Json2Excel(rows []map[string]string)(error) {
   return nil
}

// 將Excel檔案轉為Json輸出
func (xsl *Excelsrv) Rows2JsonSingleLine(row interface{}, title interface{}) (map[string]string, error) {
   cellValues := reflect.ValueOf(row)
   titleFields := reflect.ValueOf(title)

   if cellValues.Kind() != reflect.Slice {  // 只能是 Slice
      return nil, fmt.Errorf("Row2Json only accept array or slice. But got %v.", cellValues.Kind())
   }
   var len int = cellValues.Len()
   var flen int = titleFields.Len()
   if len != flen { // Check the length
      return nil, fmt.Errorf("Should be %d, but got %d.", flen, len)
   }

   cells := make(map[string]string, len)
   for i := 0; i < cellValues.Len(); i++ {
      cells[titleFields.Index(i).String()] = cellValues.Index(i).String()
   }
   return cells, nil
}

func (xsl *Excelsrv) Rows2Json(rows [][]string)([]map[string]string, error) {
   if len(rows) == 0 {
      return nil, fmt.Errorf("There has no data in Excel.")
   }
   title := rows[0]
   var errorStr string = ""

   vals := make([]map[string]string, len(rows)-1)
   for i, row := range rows {
      if i == 0 { continue }
      val, err := xsl.Rows2JsonSingleLine(row, title)
      if err != nil {
         errorStr += "Line " + strconv.Itoa(i) + ": " + err.Error() + "\n"
      }
      // vals = append(vals, val) 
      vals[i-1] = val 
   }
   if errorStr != "" {
      return nil, fmt.Errorf("Row2Json:%v", errorStr)
   }
   return vals, nil
}

func (xsl *Excelsrv) NewExcelSrv(f string)([]map[string]string, error) {
   xf, err := excelize.OpenFile(f)
   if err != nil { return nil, err }
   sheetName := xf.GetSheetName(1) // start from 1
   if sheetName == "" {
      return nil, fmt.Errorf("Can not get " + f + " sheet name.")
   }
   rows, err := xf.GetRows(sheetName)
   // rows := xf.GetRows(sheetName)
   if err != nil { return nil, err }
   if len(rows) == 0 {
      return nil, fmt.Errorf("No data in file.")
   }
   vals, err := xsl.Rows2Json(rows)
   if err != nil { return nil, err }
   return vals, nil
}

func (xls *Excelsrv) ParsefileFromWeb(w http.ResponseWriter, r *http.Request) {
   savedPath := ""

   defer func() {
      if savedPath != "" {
         _ = os.Remove(savedPath)
      }
   }()
   // 10 << 20 specifies a maximum upload of 10 MB files.
   r.ParseMultipartForm(10 << 20)
   file, handler, err := r.FormFile("xlsx")
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      log.Printf("FormFile error): %v", err)
      return
   }
   defer file.Close()

   nameParts := strings.Split(handler.Filename, ".")
   filename := nameParts[1]
   savedPath = filepath.Join("/tmp/excel/", filename)
   f, err := os.OpenFile(savedPath, os.O_WRONLY|os.O_CREATE, 0666)
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      log.Printf("os.OpenFile() error:%v(%v)", err, savedPath)
      return
   }
   defer f.Close()
   _, err = io.Copy(f, file)
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      log.Printf("io.Copy() Error:%v", err)
      return
   }

   vals, err := xls.NewExcelSrv(savedPath)
   if err != nil { 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      log.Printf("NewExcelSrc: %v", err)
      return 
   }

   valj, err := json.Marshal(vals)
   if err != nil { 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      log.Printf("Marshal error:%v", err)
      return 
   }
   w.Header().Set("Content-Type", "application/json")
   w.Write(valj)
}

// check health
func (xls *Excelsrv) Healthz(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
}

