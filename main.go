package main

func main() {
	//SA_SERVER_URL := "https://yijianqingxin.datasink.sensorsdata.cn/sa?project=sztayuan&token=4a2ee4d0f20b1b44"
	//
	//var err error
	//sensorConsumer, err := sensorsanalytics.InitDefaultConsumer(SA_SERVER_URL, 30*1000)
	//if err != nil {
	//	panic(err)
	//}
	//sensorV := sensorsanalytics.InitSensorsAnalytics(sensorConsumer, "sztayuan", false)
	//
	//now := time.Now()
	//mid := 1
	//event := "send_be_accost_msg"
	//
	//for i := 0; i < 20; i++ {
	//	properties := map[string]interface{}{
	//		"StarLevel":   1,
	//		"GirlType":    "typ1",
	//		"ChatUpTime":  now.Add(time.Second).Format(`2006-01-02 15:04:05`),
	//		"ManUserId":   strconv.Itoa(mid),
	//		"ManType":     "typ2",
	//		"WomanUserId": strconv.Itoa(22),
	//		"PoolId":      2,
	//	}
	//	err = sensorV.Track("22", event, properties, false)
	//	if err != nil {
	//		panic(err)
	//	}
	//	time.Sleep(time.Millisecond * 50)
	//	fmt.Println(i)
	//}

}
