## Excel 讀取程式
主要用來處理Excel file 轉成 JSON 格式輸出。標題為輸出欄位，須盡量採用英文，以免發生不可預期錯誤。

#### 使用函式庫
```
go get github.com/360EntSecGroup-Skylar/excelize
```

### 建立環境
```
mkdir tmp
```

### RESRful API
* /excel2json
  - 功能：將上傳的Excel檔案轉成 json 檔

* /json2excel
  - 功能：將json檔案轉成excel檔

* /healthz
  - 功能：健康檢查

### 參考資料
* [go reflect](https://stackoverflow.com/questions/47187680/how-do-i-change-fields-a-slice-of-structs-using-reflect)
