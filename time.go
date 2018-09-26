package main

import (
	"microSocket-master/util"
	"log"
)

var timerStruct = util.NewTimer()

func cron1(a,b,c int, d string)  {
	log.Print(a,b,c,d)
}

func cron2() {
	log.Printf("定时器2\r\n")
}

func main()  {
	log.Printf("1\r\n")
	timerStruct.RegisterTimer(cron1, "1s", 2,3,4, "tuzisir")
	timerStruct.RegisterTimer(cron2, "3s")
	log.Printf("8\r\n")
	timerStruct.ExecTimer()
	// 测试使用，等待所有线程退出，测试代码永不退出
	<-timerStruct.WaitTimerFinsh
	log.Print("2")
}
