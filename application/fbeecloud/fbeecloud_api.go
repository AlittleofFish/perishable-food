package fbeecloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 获取温度的请求体
type Fbee struct {
	Url    string `json:"Url"`
	Port   string `json:"Port"`
	DevKey string `json:"DevKey"`
}

// 通过文件读取设备账号信息，通过配置文件的方式，修改时就不需要重复build应用
func InitFbeeCloud() (Fbee, error) {
	var fbee = Fbee{}
	filebyte, err := ioutil.ReadFile("./fbeecloud.json")
	fmt.Println(string(filebyte))
	if err != nil {
		return Fbee{}, err
	}
	err = json.Unmarshal(filebyte, &fbee)
	if err != nil {
		return Fbee{}, err
	}
	// 打印设备数据
	fmt.Printf("%+v", fbee)
	return fbee, nil
}

// 获取的温度数据
type FbeeTemperature struct {
	DevKey          string  `json:"DevKey"`
	DevTempValue    string  `json:"Temperature"`
	LastTemperature float32 `json:"LastTemperature"`
}

// 获取温度数据
func (fbee Fbee) TemperatureGetDevice() (FbeeTemperature, error) {
	resp, err := http.Get("http://" + fbee.Url + ":" + fbee.Port) //服务器地址 http://120.77.210.77:8082
	if err != nil {
		return FbeeTemperature{"", "", 100}, nil //服务器离线时，返回温度100
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FbeeTemperature{}, err
	}

	// 解析数据
	// 解析成数组
	var arr []interface{}
	err = json.Unmarshal(content, &arr)
	if err != nil {
		return FbeeTemperature{}, err
	}
	fmt.Println("数据个数：", len(arr))

	// 解析每个json数据
	var temp FbeeTemperature
	for i := 0; i < len(arr); i++ {
		someone, err := json.Marshal(arr[i])
		if err != nil {
			return FbeeTemperature{}, err
		}
		_ = json.Unmarshal(someone, &temp)
		z, _ := strconv.ParseFloat(temp.DevTempValue, 32)
		temp.LastTemperature = float32(z)
		// 设备正确时，返回温度
		if temp.DevKey == fbee.DevKey {
			// 打印数据
			fmt.Printf("温度数据：%+v\n", temp)
			return temp, nil
		}
	}
	return FbeeTemperature{}, errors.New("找不到设备，请检查DevKey(or 节点ID)")
}
