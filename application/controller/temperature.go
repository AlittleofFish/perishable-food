package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	bc "gdzce.cn/perishable-food/application/blockchain"
	"gdzce.cn/perishable-food/application/fbeecloud"
	"gdzce.cn/perishable-food/application/util"
	"github.com/gin-gonic/gin"
)

var Fbee fbeecloud.Fbee // fbee传感器，应用启动时初始化，获取温度时使用

type orderTemperatureRequest struct {
	OrderId     string `form:"order_id" json:"order_id" binding:"required"`
	Temperature string `form:"temperature" json:"temperature" binding:"required"`
}

// 更新订单温度
func UpdateOrderTemperature(ctx *gin.Context) {
	// 解析请求体
	req := new(orderTemperatureRequest)
	if err := ctx.ShouldBind(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 打印请求体
	fmt.Println("请求参数：")
	marshal, err := json.Marshal(req)
	fmt.Println(string(marshal))

	// 格式化时间
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 调用链码
	resp, err := bc.ChannelExecute("updateOrderTemperature", [][]byte{
		[]byte(req.OrderId),
		[]byte(req.Temperature),
		[]byte(timestamp),
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 将结果返回
	ctx.JSON(http.StatusOK, resp)
}

// 调用fbeeCloud获取传感器温度并上链
type FbeeTemperatureUpdater struct {
	OrderId  string    // 订单id
	Interval int       // 获取温度时间间隔（单位毫秒）
	stopFlag chan bool // 需要停止时传入true
}

func (f *FbeeTemperatureUpdater) Start() {
	// f.stopFlag这个通道接收true时，退出定时
	f.stopFlag = util.SetInterval(func() {
		// 若未初始化则初始化
		if Fbee == (fbeecloud.Fbee{}) {
			Fbee, _ = fbeecloud.InitFbeeCloud()
		}

		// 获取温度信息
		deviceInfo, err := Fbee.TemperatureGetDevice()
		if err != nil {
			fmt.Printf("get temperature error: %s\n", err.Error())
		} else {
			// 格式化时间
			timestamp := time.Now().Format("2006-01-02 15:04:05")

			// 调用链码更新温度信息到相应订单
			resp, err := bc.ChannelExecute("updateOrderTemperature", [][]byte{
				[]byte(f.OrderId),
				[]byte(fmt.Sprintf("%f", deviceInfo.LastTemperature)),
				[]byte(timestamp),
			})
			if err != nil {
				fmt.Printf("update temperature error! orderId: %s, error: %s\n", f.OrderId, err.Error())
			} else {
				respStr, _ := json.Marshal(resp)
				fmt.Printf("successfully update order's temperature! orderId: %s, temperature: %s, time: %s\nblockchain response:%s\n",
					f.OrderId, fmt.Sprintf("%f", deviceInfo.LastTemperature), timestamp, respStr)
			}
		}

	}, f.Interval, true)
}

// 停止定时器
func (f *FbeeTemperatureUpdater) Stop() {
	f.stopFlag <- true
}
