package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	for i := 0; i < 1; i++ {
		//如果是北京时间的晚上8点到第二天早上8点，就不要请求了
		//获取当前时间
		t := time.Now()
		fmt.Println(t)
		//获取当前时间的小时

		now := time.Now().In(time.FixedZone("CST", 8*3600)) // Beijing time
		hour := now.Hour()

		if (hour >= 20) || (hour < 8) {
			fmt.Println("It's between 8pm and 8am Beijing time.")
			return
		}

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
			f.WriteString(err.Error() + "\n")
		} else if gotTicket {
			fmt.Println("有票啦，赶紧抢票")
		}
		resp = timeReserveListApiFamily()
		handleRespFamily(resp)
		time.Sleep(1 * time.Minute)
	}
}
