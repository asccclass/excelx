### JSON 格式轉為 xlsx Excel 檔案


### Usage

* POST https://devwteamapi.test5.sinica.edu.tw/excel/json2excel

```
[
	{"sheetName": "資訊服務處",
	 "rows":[
                  {"chName":"坂本和","instName":"天文及天文物理研究所","namez":"一般同仁","hours":"3"},
                  {"chName":"Alexandre","instName":"天文及天文物理研究所","namez":"一般同仁","hours":"3"}
         ]
        },
   
	{"sheetName": "總務處",
	 "rows":[
                  {"chName":"坂本和","instName":"天文及天文物理研究所","namez":"一般同仁","hours":"3"},
                  {"chName":"Alexandre","instName":"天文及天文物理研究所","namez":"一般同仁","hours":"3"}
         ]
        },
]
```

* Result
```
[
    {
        "sheetname": "資訊服務處",
        "rows": [
            {
                "chName": "坂本和",
                "email": "ksakamoto@asiaa.sinica.edu.tw",
                "hours": "3",
                "hr": "0",
                "instName": "天文及天文物理研究所",
                "namez": "一般同仁",
                "phone": "02-23665472",
                "title": "副研究員"
            },
            {
                "chName": "Alexandre",
                "email": "",
                "hours": "3",
                "hr": "0",
                "instName": "天文及天文物理研究所",
                "namez": "一般同仁",
                "phone": "0613266641",
                "title": "工讀生"
            }
        ],
        "Headers": [
            "chName",
            "email",
            "hours",
            "hr",
            "instName",
            "namez",
            "phone",
            "title"
        ]
    },
    {
        "sheetname": "總務處",
        "rows": [
            {
                "chName": "坂本和AAAC",
                "email": "ksakamoto@asiaa.sinica.edu.tw",
                "hours": "3",
                "hr": "0",
                "instName": "天文及天文物理研究所",
                "namez": "一般同仁",
                "phone": "02-23665472",
                "title": "副研究員"
            },
            {
                "chName": "憨憨憨",
                "email": "",
                "hours": "3",
                "hr": "0",
                "instName": "天文及天文物理研究所",
                "namez": "一般同仁",
                "phone": "0613266641",
                "title": "工讀生"
            }
        ],
        "Headers": [
            "chName",
            "email",
            "hours",
            "hr",
            "instName",
            "namez",
            "phone",
            "title"
        ]
    }
]
```
