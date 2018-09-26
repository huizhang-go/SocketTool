package util

import (
	"time"
	"reflect"
	"log"
	"strconv"
)

var timeUnit = []string{"h","s"}

// 方法
type funcInfoStruce struct {
	Func interface{}
	FuncTime int
	FuncParams []interface{}
	TimeUnit string
}

// 定时器结构
type TimerStruct struct {
	RegisterFuncs map[string]*funcInfoStruce
	WaitTimerFinsh chan struct{} // 测试使用
	TimerNum int8 // 定时器数量充当map的键值
}

// 初始化返回定时器结构
func NewTimer() *TimerStruct {
	return &TimerStruct{
		RegisterFuncs:make(map[string]*funcInfoStruce),
		WaitTimerFinsh:make(chan struct{}),
		TimerNum:0,
	}
}

// 注册定时器
func (t *TimerStruct) RegisterTimer(registerFunc interface{}, timeFormat string, params ... interface{}) bool {
	// 判断时间格式
	if len([]rune(timeFormat)) == 0 {
		log.Print("时间格式没有输入")
		return false
	}
	// 截取出左右两部分，左数值，右单位
	lTimeFormat := timeFormat[0 : len(timeFormat)-1]
	time,error := strconv.Atoi(lTimeFormat)
	if error != nil{
		log.Println("字符串转换成整数失败")
	}
	// 判断是否为数值
	if !(time > 0) {
		log.Print("不合法时间值")
		return false
	}
	rTimeFormat := timeFormat[len(timeFormat)-1 : len(timeFormat)]
	var issetFormat = false;
	for _, v := range timeUnit {
		if v == rTimeFormat {
			issetFormat = true
			break
		}
	}
	if !issetFormat {
		log.Print("使用默认时间格式s")
		rTimeFormat = "s"
	}
	log.Print(time)
	t.TimerNum++
	fInfo := &funcInfoStruce{
		Func:registerFunc,
		FuncTime:time,
		FuncParams:params,
		TimeUnit:rTimeFormat,
	}
	t.RegisterFuncs[string(t.TimerNum)] = fInfo
	return true
}


// 执行定时器
func (t *TimerStruct) ExecTimer()  {
	t.WaitTimerFinsh = make(chan struct{})
	go func() {
		t.rangeTimer()
	}()
}

// 遍历定时器
func (t *TimerStruct) rangeTimer()  {
	// 遍历出注册的定时器方法
	for _, v := range t.RegisterFuncs {
		go t.execFunc(v)
	}
}

// 执行方法
func (t *TimerStruct) execFunc(v *funcInfoStruce) {
	var unit time.Duration
	if v.TimeUnit == "h" {
		unit = time.Hour
	} else {
		unit = time.Second
	}
	timerTool := time.NewTicker(time.Duration(v.FuncTime) * unit)
	for {
		select {
		case <-timerTool.C:
			f := reflect.ValueOf(v.Func)
			in := make([]reflect.Value, len(v.FuncParams))
			// 将方法参数拼装出来
			for k, param := range v.FuncParams {
				in[k] = reflect.ValueOf(param)
			}
			f.Call(in)
		}
	}
}
