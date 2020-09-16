package wordSrv

import (
   // "io"
   "fmt"
   "time"
   "bytes"
   // "strings"
   "strconv"
   "net/http"
   "io/ioutil"
   "encoding/json"
   "github.com/gorilla/mux"
   "github.com/nguyenthenguyen/docx"
   // "github.com/asccclass/sherrytime"
   "github.com/asccclass/staticfileserver"
)

type WordTemplate struct {
   Docx			*docx.ReplaceDocx
   FilePath		string
   FromFile		bool			// true)從檔案 false)從記憶體
   Srv			*SherryServer.ShryServer
}

// 變數名稱=>取代的值
func(word *WordTemplate) ReplaceParams(params map[string]string)(*docx.Docx, error) {
   docx := word.Docx.Editable()
   if len(params) <= 0 {
      return docx, nil
   }
   for key, val := range params {
      if err := docx.Replace(key, val, -1); err != nil {
         return nil, err
      }
   }
   // docx1.ReplaceLink("http://example.com/", "https://github.com/nguyenthenguyen/docx")
   // docx1.ReplaceFooter("Change This Footer", "new footer")
   return docx, nil
}

// 從檔案開啟word template
func(word *WordTemplate) Open(filepath string, data []byte)(error) {
   var err error
   if word.FromFile  {
      word.Docx, err = docx.ReadDocxFile(filepath)
      word.FilePath = filepath
      word.FromFile = true 
   } else {
      word.Docx, err = docx.ReadDocxFromMemory(bytes.NewReader(data), int64(len(data)))  // *ReplaceDocx
      word.FilePath = ""
      word.FromFile = false
   }
   if err != nil {
      return err
   }
   return nil
}

// 將 word 轉為 pdf
func(word *WordTemplate) ConvertWord2PdfFromWeb(w http.ResponseWriter, r *http.Request) {
   r.ParseMultipartForm(32 << 20)
   if err := r.ParseForm(); err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   // file, handler, err := r.FormFile("templatefile")   // handler 內有檔案資訊
   file, _, err := r.FormFile("templatefile")
   if err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   data, err := ioutil.ReadAll(file)  // 讀取檔案到 buffer
   if err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   defer file.Close()
   ps := r.FormValue("params")   // 解出web參數
   params := make(map[string]string)
   if err := json.Unmarshal([]byte(ps), &params); err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   if err := word.Open("", data); err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   docx, err := word.ReplaceParams(params)
   if err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   // 資料取代完，寫回docx檔
   filepath := "/tmp/word/"
   // st := sherrytime.NewSherryTime("Asia/Taipei", "-")  // Initial
   filename := "test"  // st.NewUUID()   // 設定檔案名稱(亂數)
   docx.WriteToFile(filepath + filename + ".docx")
   word.Docx.Close()

   downloadBytes, err := ioutil.ReadFile(filepath + filename + ".docx")
   if err != nil {
      word.Srv.Error.Error2Web(w, err)
      return
   }
   topdf := r.FormValue("topdf")  // 是否要轉成pdf輸出 0)不須 1)轉成pdf
   if topdf == "" || topdf == "0" {
      mime := http.DetectContentType(downloadBytes)
      fileSize := len(string(downloadBytes))
      w.Header().Set("Content-Type", mime)
      w.Header().Set("Content-Disposition", "attachment; filename=a.docx")
      w.Header().Set("Content-Length", strconv.Itoa(fileSize))
      http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(downloadBytes))
      return
   } 
   word.Srv.Error.Error2Web(w, fmt.Errorf("目前尚未支援PDF輸出"))
}

// Web Router
func(word *WordTemplate) AddRouter(router *mux.Router) {
   router.HandleFunc("/word2pdf", word.ConvertWord2PdfFromWeb).Methods("POST")
}

func NewWordTemplate(srv *SherryServer.ShryServer, filepath string)(*WordTemplate, error) {
   var err error
   var r WordTemplate

   r.Srv = srv

   if filepath != "" {
      r.FromFile = true   // 從檔案
      err = r.Open(filepath, nil) 
      if err != nil {
         return nil, err
      }
   } else {
      r.FromFile = false    // 從記憶體
   }
   return &r, nil
}

/*
func main() {
   r, err := NewWordTemplate("./V1.docx")
   if err != nil {
      panic(err)
   }

    st := sherrytime.NewSherryTime("Asia/Taipei", "-")  // Initial
    td := st.Today()
    s := strings.Split(td, st.Delimiter)

    params := make(map[string]string)
    params["{SerialNumber}"] = "20200102-33333"
    params["{name}"] = "劉智漢"
    params["{coursey}"] = "109"
    params["{coursem}"] = "09"
    params["{coursed}"] = "21"
    params["{mxourse}"] = "誠信提升計畫課程研究倫理教育課程"
    params["{Hours}"] = "2.9"
    t := st.Year2Chinese(s[0])
    params["{Year}"] = t
    params["{Month}"] = s[1]
    params["{day}"] = s[2]

   docx1, err := r.ReplaceParams(params)

   docx1.WriteToFile("./new.docx")

   // Or write to ioWriter
   // docx2.Write(ioWriter io.Writer)

   r.Docx.Close()
}
*/
