package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"gdzce.cn/perishable-food/application/blockchain"
	"gdzce.cn/perishable-food/application/repository"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// 根据特定请求uri，发起get请求返回响应
func get(uri string, router *gin.Engine) ([]byte, int) {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	// 提取响应
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(result.Body)
	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body, result.StatusCode
}

// 根据特定请求uri和参数param，以表单形式传递参数，发起post请求返回响应
func postForm(uri string, param []byte, router *gin.Engine) ([]byte, int) {
	// 构造post请求
	req := httptest.NewRequest("POST", uri, strings.NewReader(bytes.NewBuffer(param).String()))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	// 提取响应
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(result.Body)
	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body, result.StatusCode
}

func expectApi(status int, testFunc string) {
	str := "******"
	file := "./api_test_result.txt"
	switch status {
	case 1:
		str += "测试通过：pass"
		writeFile(file, testFunc+":PASS\n")
	case 2:
		str += "测试失败：fail"
		writeFile(file, testFunc+":FAIL\n")
	}
	fmt.Println(str)
}

func writeFile(filePath string, data string) {
	//os.Remove(filePath)
	//os.Create(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("文件打开失败", err)
		file, _ = os.Create(filePath)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(data)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

var routers *gin.Engine

// testing 初始化内容
func init() {
	routers = setupRouter()
	blockchain.Init()
	// 加载存储[{orderId，txid}]的json文件
	repository.TransactionRecordList.FilePath = TransactionRecordFileName
	_ = repository.TransactionRecordList.LoadTransactionRecords()
}

// Test_SDK SDK能否访问区块链网络
func Test_SDK1(t *testing.T) {

	blockchain.ChaincodeName = "mycc"
	defer func() {
		blockchain.ChaincodeName = "mychaincode"
	}()

	resp, errString := blockchain.ChannelExecute1("invoke", [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("10"),
	})
	if resp.ChaincodeStatus == 200 && errString == nil {
		time.Sleep(2 * time.Second)
		resp, _ := blockchain.ChannelQuery1("query", [][]byte{
			[]byte("a"),
		})
		t.Logf("resp.Payload: %+v\n", bytes.NewBuffer(resp.Payload).String())
		t.Logf("resp.ChaincodeStatus: %+v\n", resp.ChaincodeStatus)
		if len(resp.Payload) != 0 {
			expectApi(1, "Test_SDK1")
		} else {
			expectApi(2, "Test_SDK1")
			t.Errorf("err: %+v\n", len(resp.Payload))
			t.FailNow()
		}
	} else {
		expectApi(2, "Test_SDK1")
		t.Errorf("err: %+v\n", len(resp.Payload))
		t.FailNow()
	}
}

// Test_SDK SDK是否符合全组织背书策略
func Test_SDK2(t *testing.T) {
	routers.GET("/testGet", func(c *gin.Context) {
		key := c.Query("key")
		blockchain.ChaincodeName = "mycc"
		defer func() {
			blockchain.ChaincodeName = "mychaincode"
		}()
		resp, _ := blockchain.ChannelQuery1("query", [][]byte{
			[]byte(key),
		})
		// 将结果返回
		c.JSON(int(resp.ChaincodeStatus), resp)
	})

	value, statusCode := get("/testGet?key=a", routers)
	var resp channel.Response
	_ = json.Unmarshal(value, &resp)
	if statusCode == 200 && len(resp.Responses) == 3 {
		expectApi(1, "Test_SDK2")
	} else {
		expectApi(2, "Test_SDK2")
		t.Errorf("请求的组织数量: %+v\n，期望：3", len(resp.Responses))
		t.FailNow()
	}
}

// 能调用链码创建商品
func Test_createCommodity(t *testing.T) {
	data := commodityRequest2{
		Name:            "testCommodity1",
		Id:              "20211001001",
		LowTemperature:  "5",
		HighTemperature: "20",
		Price:           10,
		OwnerId:         "1",
		Origin:          "testOrigin",
	}
	arrByte, _ := json.Marshal(data)
	_, status := postForm("/createCommodity", arrByte, routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_createCommodity")
	} else {
		expectApi(2, "Test_createCommodity")
		t.FailNow()
	}
}

// 能调用链码查询所有商品
func Test_commodityList(t *testing.T) {
	_, status := get("/commodityList", routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_commodityList")
	} else {
		expectApi(2, "Test_commodityList")
		t.FailNow()
	}
}

// 能调用链码创建订单
func Test_createOrder(t *testing.T) {
	data := orderRequest2{
		CommodityId:    "testcreateOrder",
		Id:             "20211001001",
		DeliverAddress: "testAddress",
		DeliverTime:    20211001001,
		Quantity:       10,
		BuyerId:        "3",
		SellerId:       "1",
	}
	arrByte, _ := json.Marshal(data)
	_, status := postForm("/createOrder", arrByte, routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_createOrder")
	} else {
		expectApi(2, "Test_createOrder")
		t.FailNow()
	}
}

// 能根据订单id查询订单
func Test_orderList1(t *testing.T) {
	_, status := get("/orderList?orderId=1570885832799", routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_orderList1")
	} else {
		expectApi(2, "Test_orderList1")
		t.FailNow()
	}
}

// 能查询所有订单
func Test_orderList2(t *testing.T) {
	_, status := get("/orderList", routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_orderList2")
	} else {
		expectApi(2, "Test_orderList2")
		t.FailNow()
	}
}

// 能查询账户信息
func Test_accountList(t *testing.T) {
	data := accountListRequestBody2{
		AccountId: "1",
	}
	arrByte, _ := json.Marshal(data)
	_, status := postForm("/accountList", arrByte, routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_accountList")
	} else {
		expectApi(2, "Test_accountList")
		t.FailNow()
	}
}

func Test_updateOrderStatus(t *testing.T) {
	data := updateOrderStatusRequest2{
		OrderId: "1",
		Status:  "Processing",
	}
	arrByte, _ := json.Marshal(data)
	_, status := postForm("/updateOrderStatus", arrByte, routers)
	t.Log(status)
	if status == 200 {
		expectApi(1, "Test_updateOrderStatus")
	} else {
		expectApi(2, "Test_updateOrderStatus")
		t.FailNow()
	}
}

type commodityRequest2 struct {
	Name            string  `json:"name" form:"name" binding:"required"`                       // 商品名
	Id              string  `json:"id" form:"id" binding:"required"`                           // id
	Origin          string  `json:"origin" form:"origin" binding:"required"`                   // 产地
	LowTemperature  string  `json:"lowTemperature" form:"lowTemperature" binding:"required"`   // 最低温
	HighTemperature string  `json:"highTemperature" form:"highTemperature" binding:"required"` // 最高温
	Price           float64 `json:"price" form:"price" binding:"required"`                     // 单价
	OwnerId         string  `json:"owner" form:"owner" binding:"required"`                     // 所有者
}

type orderRequest2 struct {
	CommodityId    string `json:"commodity_id" binding:"required"`   // 商品id
	Id             string `json:"id" binding:"required"`             // 订单id
	DeliverAddress string `json:"deliverAddress" binding:"required"` // 配送地址
	DeliverTime    int64  `json:"deliverTime" binding:"required"`    // 预计送达时间（时间戳）
	Quantity       int64  `json:"quantity" binding:"required"`       // 数量
	BuyerId        string `json:"buyer" binding:"required"`          // 买家
	SellerId       string `json:"seller" binding:"required"`         // 卖家
}

type accountListRequestBody2 struct {
	AccountId string `form:"account_id" json:"account_id" binding:"required"`
}

type updateOrderStatusRequest2 struct {
	OrderId string `form:"order_id" json:"order_id" binding:"required"`
	Status  string `form:"status" json:"status" binding:"required"`
}
