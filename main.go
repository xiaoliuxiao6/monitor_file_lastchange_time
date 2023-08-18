package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// 读取配置文件
func init() {
	viper.SetConfigName("monitor_file_lastchange_time") // 配置文件的文件名(不带扩展名)
	viper.SetConfigType("json")                         // (文件的扩展名)如果 SetConfigName 中没有扩展名，则需要
	viper.AddConfigPath("/usr/local/etc")               // 配置文件路径(可以多次出现以搜索多个路径)

	err := viper.ReadInConfig() // 载入配置文件
	if err != nil {
		log.Panicf("配置文件读取失败: %v \n", err)
	} else {
		log.Printf("配置文件读取成功: %v", viper.ConfigFileUsed())
	}

	viper.WatchConfig() //监视和重新读取配置文件
}

func main() {

	SendWeixin("我是监控程序，我要启动了！！!")
	for {
		if time.Now().Minute() == 0 {
			if time.Now().Hour() == 8 || time.Now().Hour() == 20 {
				SendWeixin("我是监控程序，我正在运行！！!")
			}
		}
		checkLastChangeTime()
		log.Print("休眠一会")
		time.Sleep(600 * time.Second)
	}
}

func checkLastChangeTime() {
	// 所有路径
	paths := viper.GetStringSlice("paths")
	var noUpdate []string

	for _, path := range paths {
		log.Printf("路径: %v, 最后更新时间: %v, %v", path, getLastChangeTime(path), getLastChangeTime(path)/1000000000)
		lastChangeTime := getLastChangeTime(path) / 1000000000 // 文件夹内指定文件的最后更新时间
		if lastChangeTime > 300 {
			noUpdate = append(noUpdate, path)
		}
	}

	if len(noUpdate) > 0 {
		var message = fmt.Sprintf("一共有 %v 个目录超过 5分钟没有更新:\n", len(noUpdate))
		for index, path := range noUpdate {
			message = message + fmt.Sprintf("%v: %v\n", index+1, path)
		}

		fmt.Println("-------------------------------- 1")
		fmt.Println("告警内容：", message)
		fmt.Println("-------------------------------- 2")
		SendWeixin(message)
	}
}

func SendWeixin(message string) {
	// WeixinWebhook := viper.GetString("Channel.Weixin.Webhook")
	// 1. https://github.com/feiyu563/PrometheusAlert/blob/master/doc/readme/conf-wechat.md
	// 2. 登录企业微信网页 - 我的企业 - 微信插件 - 扫描二维码即可手机微信来接收信息

	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    	"msgtype": "text",
    	"text": {
        	"content": "%v"
    	}
	}`, message))

	client := &http.Client{}
	// req, err := http.NewRequest(method, "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=0da619f2-afbf-4d2f-a13c-db23abba986e", payload)
	req, err := http.NewRequest(method, "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=ff8d99d2-5d81-4d20-815f-b55d3ffa0d59", payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// 获取指定目录下最后一个文件的更新时间 - 返回纳秒
func getLastChangeTime(basePath string) int64 {

	// 分钟
	var lastChangeTime int64
	err := filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
		}

		if filepath.Ext(info.Name()) == ".bin" {
			if lastChangeTime == 0 || lastChangeTime > int64(time.Since(info.ModTime())) {
				lastChangeTime = int64(time.Since(info.ModTime()))
			}

			// log.Debugf("文件时间: %v, %v, %v 纳秒, %v 秒\n", info.Name(), time.Since(info.ModTime()), int64(time.Since(info.ModTime())), int64(time.Since(info.ModTime()))/1000000000)

		}
		return nil
	})

	// fmt.Println("最后修改时间; ", lastChangeTime)

	if err != nil {
		log.Fatalf("从 %v 路径获取文件信息失败: %v", basePath, err)
	}

	return lastChangeTime
}
