package excelSrv

import(
   "io"
   "os"
   "fmt"
   "sort"
   "time"
   // "bytes"
   "reflect"
   "strings"
   "strconv"
   "math/rand"
   "path/filepath"
   "net/http"
   "io/ioutil"
   "encoding/csv"
   "encoding/json"
   log "github.com/sirupsen/logrus"
   "github.com/asccclass/sherrytime"
   // _ "github.com/paulrosania/go-charset/data"
   // "github.com/paulrosania/go-charset/charset"
   "github.com/saintfish/chardet"
   "github.com/360EntSecGroup-Skylar/excelize"
)

var (
   Chars = map[int]string{ 0:"A", 1:"B", 2:"C", 3:"D", 4:"E", 5:"F", 6:"G", 7:"H", 8:"I", 9:"J", 10:"K", 11:"L", 12:"M",  13:"N", 14:"O", 15:"P", 16:"Q", 17:"R", 18:"S", 19:"T", 20:"U", 21:"V", 22:"W", 23:"X", 24:"Y", 25:"Z", 26:"AA", 27:"AB",}
)

type ExcelJson struct {
   SheetName	string			`json:"sheetname"`
   Rows		[]map[string]string	`json:"rows"`
   Headers	[]string
}

type Excelsrv struct {
   Rows map[string]string
}

func(xls *Excelsrv) Error(w http.ResponseWriter, err error) {
   w.Header().Set("Content-Type", "application/json;charset=UTF-8")
   w.WriteHeader(http.StatusOK)
   fmt.Fprintf(w, "{\"errMsg\": \"%s\"}", err.Error())
}

// 取得tab數量及名稱
func(xls *Excelsrv) GetSheetTabs(f *excelize.File)(int, []string) {
   var sheetCnt int = 0
   var sheetName []string
   for index, name := range f.GetSheetMap() {
      sheetCnt = index
      sheetName = append(sheetName, name)
   }
   return sheetCnt, sheetName
}

// 從associate arrray 取得Headers
func(xls *Excelsrv) SetHeader(rows []map[string]string)([]string) {
   var header []string
   if len(rows) <= 0 {
      return header
   }
   for key, _ := range rows[0] {
      header = append(header, key)
   }
   sort.Strings(header)
   return header
}

//將Json轉為Excel
func (xls *Excelsrv) Json2Excel(f *excelize.File, sheetName string, rows ExcelJson)(error) {
   index := f.NewSheet(sheetName)
   line := 1
   for i, row := range rows.Rows {		// row
      if(i == 0)  {  // write header
         for k, head := range rows.Headers {
            f.SetCellValue(sheetName, Chars[k] + fmt.Sprint(line), head) 
         }
         line++
      }
      for j, head := range rows.Headers {	// cell
         f.SetCellValue(sheetName, Chars[j] + fmt.Sprint(line), row[head])
      } 
      line++
   }
   f.SetActiveSheet(index)
/*
   if err := f.SaveAs("/tmp/excel/tmp1.xlsx"); err != nil {
      return nil, fmt.Errorf("Write file error(%v).", err)
   }
*/
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
   st := sherrytime.NewSherryTime("Asia/Taipei", "-")  // Initial
   for i := 0; i < cellValues.Len(); i++ {
      // 檢查日期格式，將其他分隔符號都轉為-分隔
      if titleFields.Index(i).String() == "certified_date" {	// 日期: - or /
         val := cellValues.Index(i).String()
         idx := strings.Index(val, "/")
         if idx != -1 {  // 分隔為 /
            val = strings.ReplaceAll(val, "/", "-") 
         } else {   // 若分隔為 .
            if idx = strings.Index(val, "."); idx != -1 {
               val = strings.ReplaceAll(val, ".", "-")
            }
         }
         certdate, err := st.TransferFormat(val)
         if err != nil {
            cells[titleFields.Index(i).String()] = "日期格式錯誤:" + err.Error()
         } else {
            cells[titleFields.Index(i).String()] = certdate
         }
      }  else  {
         cells[titleFields.Index(i).String()] = cellValues.Index(i).String()
      }
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
         errorStr += "Line " + strconv.Itoa(i) + ": " + err.Error()
      }
      // vals = append(vals, val) 
      vals[i-1] = val 
   }
   if errorStr != "" {
      return nil, fmt.Errorf("Row2Json:%s", errorStr)
   }
   return vals, nil
}

func (xsl *Excelsrv) NewExcelSrv(f string)([]map[string]string, error) {
   xf, err := excelize.OpenFile(f)
   if err != nil { return nil, err }

   // 檢查sheet數量
   if xf.SheetCount > 1 {
      return nil, fmt.Errorf("上傳之Excel檔案表單(sheet)數量大於一個.")
   }
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

// CharsetDetect 偵測檔案內編碼
func (xls *Excelsrv) CharsetDetect(str []byte)(*chardet.Result, error) {
   detector := chardet.NewTextDetector()
   result, err := detector.DetectBest(str)
   if err != nil {
      return nil, err
   }
   return result, nil
   // Detected charset is %s, language is %s", result.Charset, result.Language)
   // Detected charset is GB-18030, language is zh
}

func (xls *Excelsrv) DownloadFileFromWeb(w http.ResponseWriter, r *http.Request) {
   var err error
   defer func() {
      if err != nil {
         w.Header().Set("Content-Type", "application/json;charset=UTF-8")
         w.WriteHeader(http.StatusOK)
         fmt.Fprintf(w, "{\"errMsg\": \"%v\"}", err)
      }
   }()
   b, err := ioutil.ReadAll(r.Body)  // read json file from web
   if err != nil {
      return
   }
   jsonfile := []ExcelJson{} 
   if !json.Valid(b)  {  // check json file
      err = fmt.Errorf("json data is invalid(%v).", string(b))
      return;
   }
   defer r.Body.Close()
   if err = json.Unmarshal(b, &jsonfile); err != nil {
      err = fmt.Errorf("Unmarshal body error(%v).", err)
      return
   }
   f := excelize.NewFile()
   for i, jsf := range jsonfile {
      jsonfile[i].Headers = xls.SetHeader(jsf.Rows)
      if err = xls.Json2Excel(f, jsf.SheetName, jsonfile[i]); err != nil {
         err = fmt.Errorf("Json to Excel error(%v).", err)
         return
      }
   }
   rand.Seed(time.Now().UnixNano())
   x := rand.Intn(10000)   // filename
   fname := fmt.Sprintf("%v.xlsx", x)
   if err := f.SaveAs("/tmp/excel/" + fname); err != nil {
      err = fmt.Errorf("Write file error(%v).", err)
      return
   }

   // 輸出下載檔案
   Openfile, err := os.Open("/tmp/excel/" + fname)
   defer Openfile.Close() //Close after function return
   if err != nil {
      err = fmt.Errorf("File not found.")
      return
   }
   FileHeader := make([]byte, 512)
   Openfile.Read(FileHeader)
   FileStat, _ := Openfile.Stat()                     //Get info from file
   FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

   w.Header().Set("Content-Disposition", "attachment; filename="+fname)
   w.Header().Set("Content-Type",  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
   w.Header().Set("Content-Length", FileSize)
   w.Header().Set("Content-Transfer-Encoding", "binary")
   w.Header().Set("Expires", "0")
   Openfile.Seek(0, 0)
   io.Copy(w, Openfile) //'Copy' the file to the client
   _ = os.Remove("/tmp/excel/" + fname)
}

// ParsefileFromWeb 將excel/csv轉為CSV
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
   fileExtension := nameParts[len(nameParts)-1]
   var vals []map[string]string
   var xcharset *chardet.Result

   if fileExtension == "csv"  {
      r := csv.NewReader(file)
      line := 0  
      var header []string
      for {
         rows, err := r.Read()   // []string
         if err == io.EOF || err != nil {
            break
         }
         if line == 0 {
            xcharset, err = xls.CharsetDetect([]byte(rows[0]))
            if err != nil {
               xls.Error(w, fmt.Errorf("encoding detect error"))
               return
            }
            log.Printf("charset=%v", xcharset)
            for _, h := range rows {
               header = append(header, h)
            }
            line++
            continue
         }
         cell := make(map[string]string, len(rows))
         for j, data := range rows {
            cell[header[j]] = data
         }
         vals = append(vals, cell) 
      }
      if len(vals) == 0 {
         xls.Error(w, fmt.Errorf("no data in csv file"))
         return
      }
   } else {
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

      vals, err = xls.NewExcelSrv(savedPath)
      if err != nil { 
         w.Header().Set("Content-Type", "application/json; charset=UTF-8")
         w.WriteHeader(http.StatusOK)
         fmt.Fprintf(w, "{\"errMsg\": \"%s\"}", err.Error())
         return 
      }
   }

   valj, err := json.Marshal(vals)
   if err != nil { 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      log.Printf("Marshal error:%v", err)
      return 
   }
   w.Header().Set("Content-Type", "application/json; charset=UTF-8")
   if fileExtension == "csv" {  // 要作格式轉換
   }
   w.Write(valj)
}

// check health
func (xls *Excelsrv) Healthz(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
}

