package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Result struct {
	Data   []DataItem `json:"data"`
	Status int        `json:"status"`
}

type DataItem struct {
	StockCode string `json:"stockCode"`
	Num       int    `json:"num"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	OccDate   string `json:"occDate"`
	Name      string `json:"name"`
	FsStatus  string `json:"fsStatus"`
}

func executeEvery(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			resp := timeReserveListApi()
			gotTicket, err := handleResp(resp)
			if err != nil {
				fmt.Println("错误信息:", err)
				//写入文件
				fileName := "error.log"
				f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					fmt.Println("打开文件错误:", err)
				}
				defer f.Close()
				_, err = f.WriteString(err.Error() + "\n")
			} else if gotTicket {
				fmt.Println("有票啦，赶紧抢票")
				time.Sleep(5 * time.Minute)

			}
			resp = timeReserveListApiFamily()
			gotTicket, err = handleRespFamily(resp)
		}
	}
}
func handleRespFamily(resp string) (bool, error) {
	//解析json
	var result Result
	err := json.Unmarshal([]byte(resp), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	for _, item := range result.Data {
		fmt.Println("Family ticket!!! OccDate:", item.OccDate, "Num:", item.Num)
		// if item.OccDate == target_date || item.OccDate == test_date ||
		// 	strings.Contains(item.OccDate, "2-17") ||
		// 	strings.Contains(item.OccDate, "1-12") ||
		// 	strings.Contains(item.OccDate, "1-08") {
		if item.Num > 0 {
			for i := 0; i < 10; i++ {
				sendMsg(item.OccDate+"臻享家庭票", item.Num)
				return true, nil
			}
		}
		// }
	}
	return false, nil
}

var target_date string
var test_date string

func main() {
	// fmt.Println("Hello world")
	//每个固定时间执行一次
	//定时任务
	//打印当前时间

	target_date = "2024-02-17"
	// 参数获取一个日期

	if len(os.Args) > 1 {
		test_date = os.Args[1]
	} else {
		test_date = "2024-01-07"
	}
	sendMsg(test_date, 10)
	executeEvery(1 * time.Minute)
	// executeEvery(5 * time.Second)
}

// https://www.cnblogs.com/wushaoyu/p/16884766.html
// https://help.aliyun.com/zh/ssl-certificate/user-guide/obtain-the-webhook-url-of-a-dingtalk-chatbot
// https://juejin.cn/post/7143491849330622478
func sendMsg(date string, num int) {
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=b346583700f0adae8939458016ad9922a41daf04867bc41ae95704e5fcbe4755`
	// content := `{"msgtype": "text",
	//         "text": {"content": "` + "告警" + msg + `"}
	//     }`
	msg := fmt.Sprintf("有%s的票啦, 还有%d张票  , 小奇哥 ,小马哥起床抢票啦,麻溜的 ", date, num)
	content := `{"msgtype": "text",
    "text": {"content": "` + "告警test" + msg + `"},
         "at": {
            "atMobiles": [
              "15071244227"
            ],
            "isAtAll": false
         }
       }`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func handleResp(resp string) (bool, error) {
	//解析json
	var result Result
	err := json.Unmarshal([]byte(resp), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	for _, item := range result.Data {
		fmt.Println("OccDate:", item.OccDate, "Num:", item.Num)
		if item.OccDate == target_date || item.OccDate == test_date ||
			strings.Contains(item.OccDate, "2-17") ||
			strings.Contains(item.OccDate, "1-12") ||
			strings.Contains(item.OccDate, "1-08") {
			if item.Num > 0 {
				sendMsg(item.OccDate, item.Num)
				return true, nil
			}
		}
	}
	return false, nil
}

func PostRequest(url string, headers map[string]string, body string) (string, error) {
	// Create a Transport that ignores SSL certificate verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a Client with the custom Transport
	client := &http.Client{Transport: transport}

	// Create a new request
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	// Set the headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func timeReserveListApi() string {
	//构造下面注释的包进行发送,https
	/*
	   POST /lotsapi/product/api/product/timeReserveList HTTP/1.1
	   Host: wap.lotsmall.cn
	   Accept: application/json, text/plain,
	   Trace_device_id:
	   User-Agent: Mozilla/5.0 (Linux; Android 11; Nexus 6 Build/RQ3A.211001.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3262 MMWEBSDK/201201 Mobile Safari/537.36 MMWEBID/8212 MicroMessenger/8.0.1840(0x2800003A) Process/tools WeChat/arm32 Weixin NetType/WIFI Language/zh_CN ABI/arm32
	   Content-Type: application/x-www-form-urlencoded;charset=UTF-8
	   Origin: https://wap.lotsmall.cn
	   X-Requested-With: com.tencent.mm
	   Sec-Fetch-Site: same-origin
	   Sec-Fetch-Mode: cors
	   Sec-Fetch-Dest: empty
	   Referer: https://wap.lotsmall.cn/vue/order/ticket?scenicId=1566&ticketId=133428&m_id=381
	   Accept-Encoding: gzip, deflate
	   Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7
	   Connection: close
	   Content-Length: 175

	   endTime=2024-01-09&externalCode=PST20231227825706&merchantId=381&merchantInfoId=381&modelCode=MP2023122717171632116&startTime=2024-01-06&xj_time_stamp_2019_11_28=
	*/
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Linux; Android 11; Nexus 6 Build/RQ3A.211001.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3262 MMWEBSDK/201201 Mobile Safari/537.36 MMWEBID/8212 MicroMessenger/8.0.1840(0x2800003A) Process/tools WeChat/arm32 Weixin NetType/WIFI Language/zh_CN ABI/arm32",
		"Referer":      "https://wap.lotsmall.cn/vue/order/ticket?scenicId=1566&ticketId=133428&m_id=381",
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}
	host := "https://wap.lotsmall.cn"
	url := "/lotsapi/product/api/product/timeReserveList"

	target := host + url
	starttime := "2024-01-06"
	endtime := "2024-01-09"
	//starttime 设置为当天
	//endtime 设置为当天+3

	t := time.Now()
	starttime = t.Format("2006-01-02")

	tmp := time.Now()
	tmp = tmp.AddDate(0, 0, 7)
	endtime = tmp.Format("2006-01-02")

	body := fmt.Sprintf("endTime=%s&externalCode=PST20231227825706&merchantId=381&merchantInfoId=381&modelCode=MP2023122717171632116&startTime=%s&xj_time_stamp_2019_11_28=",
		endtime, starttime)
	resp, err := PostRequest(target, headers, body)

	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	fmt.Println(resp)

	//帮我解析如下结果{"data":[{"stockCode":"SC221210111820072","num":0,"startTime":"11:00","endTime":"21:30","occDate":"2024-01-06","name":"","fsStatus":"T"},{"stockCode":"SC221210111820072","num":0,"startTime":"11:00","endTime":"21:30","occDate":"2024-01-07","name":"","fsStatus":"T"},{"stockCode":"SC221210111820072","num":0,"startTime":"11:00","endTime":"21:30","occDate":"2024-01-08","name":"","fsStatus":"T"},{"stockCode":"SC221210111820072","num":0,"startTime":"11:00","endTime":"21:30","occDate":"2024-01-09","name":"","fsStatus":"T"}],"status":200}

	return resp
}

func timeReserveListApiFamily() string {
	//构造下面注释的包进行发送,https
	/*
	   POST /lotsapi/product/api/product/timeReserveList HTTP/1.1
	   Host: wap.lotsmall.cn
	   Accept: application/json, text/plain,
	   Trace_device_id:
	   User-Agent: Mozilla/5.0 (Linux; Android 11; Nexus 6 Build/RQ3A.211001.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3262 MMWEBSDK/201201 Mobile Safari/537.36 MMWEBID/8212 MicroMessenger/8.0.1840(0x2800003A) Process/tools WeChat/arm32 Weixin NetType/WIFI Language/zh_CN ABI/arm32
	   Content-Type: application/x-www-form-urlencoded;charset=UTF-8
	   Origin: https://wap.lotsmall.cn
	   X-Requested-With: com.tencent.mm
	   Sec-Fetch-Site: same-origin
	   Sec-Fetch-Mode: cors
	   Sec-Fetch-Dest: empty
	   Referer: https://wap.lotsmall.cn/vue/order/ticket?scenicId=1566&ticketId=133428&m_id=381
	   Accept-Encoding: gzip, deflate
	   Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7
	   Connection: close
	   Content-Length: 175

	   endTime=2024-01-09&externalCode=PST20231227825706&merchantId=381&merchantInfoId=381&modelCode=MP2023122717171632116&startTime=2024-01-06&xj_time_stamp_2019_11_28=
	*/
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Linux; Android 11; Nexus 6 Build/RQ3A.211001.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3262 MMWEBSDK/201201 Mobile Safari/537.36 MMWEBID/8212 MicroMessenger/8.0.1840(0x2800003A) Process/tools WeChat/arm32 Weixin NetType/WIFI Language/zh_CN ABI/arm32",
		"Referer":      "https://wap.lotsmall.cn/vue/order/ticket?scenicId=1566&ticketId=133428&m_id=381",
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}
	host := "https://wap.lotsmall.cn"
	url := "/lotsapi/product/api/product/timeReserveList"

	target := host + url
	starttime := "2024-01-06"
	endtime := "2024-01-09"
	//starttime 设置为当天
	//endtime 设置为当天+3

	t := time.Now()
	starttime = t.Format("2006-01-02")

	tmp := time.Now()
	tmp = tmp.AddDate(0, 0, 7)
	endtime = tmp.Format("2006-01-02")

	body := fmt.Sprintf("endTime=%s&externalCode=PST20231227825706&merchantId=381&merchantInfoId=381&modelCode=MP2023122717171632116&startTime=%s&xj_time_stamp_2019_11_28=",
		endtime, starttime)

	body = "endTime=2024-02-17&externalCode=PST20231202813329&merchantId=381&merchantInfoId=381&modelCode=MP2023121320020822826&startTime=2024-02-17&xj_time_stamp_2019_11_28="

	resp, err := PostRequest(target, headers, body)
	if err != nil {
		fmt.Println("臻享家庭票 Error:", err)
		return err.Error()
	}
	fmt.Println(resp)

	return resp
}
