package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
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
}
