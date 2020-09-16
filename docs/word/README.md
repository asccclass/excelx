### Word 公用程式

### /word2md
* Convert Microsoft Word Document to Markdown
* web 參數：
   - templatefile：樣板檔案 docx
   - params：要取代的內容（json格式），例如：
   - topdf: 0)不需要 or 1)轉成pdf下載

```
{"{name}": "劉智漢", "{SerialNumber}": "00123133-2313","{coursey}":"2020","{coursem}":"08","{coursed}":"31","{mxourse}":"誠信提升計畫轉換測試功能","{Hours}":"3","{Year}":"109","{Month}":"08","{day}":"31"}
```

### 嘗試過函數庫或工具
* 使用libreoffice
會有中文字掉字的問題，解決後，產生出來的pdf也不符合需求，浮水印沒法使用，因此無法使用。 

```
libreoffice --invisible --convert-to pdf test.docx.docx
```

* [Gotenberg](https://thecodingmachine.github.io/gotenberg/)
一樣會有word浮水印問題，沒法解決。


### 參考文件
* [Gotenberg](https://thecodingmachine.github.io/gotenberg/)
