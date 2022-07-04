/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var mutex sync.RWMutex

func expectApi(status int, testFunc string) {
	str := "******"
	file := "./chaincode_test_result.txt"
	switch status {
	case 1:
		str += "测试通过：pass"
		mutex.Lock()
		writeFile(file, testFunc+":PASS\n")
		mutex.Unlock()
	case 2:
		str += "测试失败：fail"
		mutex.Lock()
		writeFile(file, testFunc+":FAIL\n")
		mutex.Unlock()
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

func GetNewStub() *shim.MockStub {
	var scc = new(PerishableFood)
	var stub = shim.NewMockStub("ex01", scc)
	return stub
}

// 链码初始化-检查商品初始化
func Test_Init1(t *testing.T) {
	stub := GetNewStub()
	stub.MockInit("1", nil)
	results, _ := stub.GetStateByPartialCompositeKey("commodity", []string{})
	if results.HasNext() {
		kv, _ := results.Next()
		commodity := new(Commodity)
		json.Unmarshal(kv.GetValue(), commodity)
		if commodity.Id == "" {
			expectApi(2, "Init1")
			t.FailNow()
		}
		expectApi(1, "Init1")
	} else {
		expectApi(2, "Init1")
		t.FailNow()
	}
}

// 链码初始化-检查账户初始化
func Test_Init2(t *testing.T) {
	stub := GetNewStub()
	stub.MockInit("1", nil)
	results, _ := stub.GetStateByPartialCompositeKey("account", []string{"1"})
	if results.HasNext() {
		kv, _ := results.Next()
		account := new(Account)
		json.Unmarshal(kv.GetValue(), account)
		if account.Id == "" || account.Balance == 0 {
			expectApi(2, "Init2")
			t.FailNow()
		}
		expectApi(1, "Init2")
	} else {
		expectApi(2, "Init2")
		t.FailNow()
	}
}

// 链码初始化-检查账户初始化
func Test_Init3(t *testing.T) {
	stub := GetNewStub()
	stub.MockInit("1", nil)
	results, _ := stub.GetStateByPartialCompositeKey("account", []string{"2"})
	if results.HasNext() {
		kv, _ := results.Next()
		account := new(Account)
		json.Unmarshal(kv.GetValue(), account)
		if account.Id == "" || account.Balance == 0 {
			expectApi(2, "Init3")
			t.FailNow()
		}
		expectApi(1, "Init3")
	} else {
		expectApi(2, "Init3")
		t.FailNow()
	}
}

// 链码初始化-检查账户初始化
func Test_Init4(t *testing.T) {
	stub := GetNewStub()
	stub.MockInit("1", nil)
	results, _ := stub.GetStateByPartialCompositeKey("account", []string{"3"})
	if results.HasNext() {
		kv, _ := results.Next()
		account := new(Account)
		json.Unmarshal(kv.GetValue(), account)
		if account.Id == "" || account.Balance == 0 {
			expectApi(2, "Init4")
			t.FailNow()
		}
		expectApi(1, "Init4")
	} else {
		expectApi(2, "Init4")
		t.FailNow()
	}
}

// 创建商品-参数非空校验
func Test_createCommodity1(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 1, 7, 1, "createCommodity", "1")
}

// 创建商品-参数个数校验
func Test_createCommodity2(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 8, 8, 1, "createCommodity", "2")
}

// 创建商品-检测是否已存在商品
func Test_createCommodity3(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 9, 9, 2, "createCommodity", "3")
}

// 创建商品-创建商品是否成功
func Test_createCommodity4(t *testing.T) {
	stub := GetNewStub()
	res := testSomeTx(stub, 1, "createCommodity", 9)
	if res.Status != shim.OK {
		expectApi(2, "createCommodity4")
		t.FailNow()
	}
	res = getTr(stub, []string{"commodity", "20211001001"})
	if res.Payload == nil {
		expectApi(2, "createCommodity4")
		t.FailNow()
	}
	expectApi(1, "createCommodity4")
}

// 新建订单-参数非空校验
func Test_createOrder1(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 1, 8, 2, "createOrder", "1")
}

// 新建订单-参数个数校验
func Test_createOrder2(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 9, 9, 2, "createOrder", "2")
}

// 新建订单-检测是否存在此商品
func Test_createOrder3(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 10, 10, 2, "createOrder", "3")
}

// 新建订单-检测是否已存在订单
func Test_createOrder4(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 11, 11, 3, "createOrder", "4")
}

// 新建订单-创建订单是否成功
func Test_createOrder5(t *testing.T) {
	stub := GetNewStub()
	res := testSomeTx(stub, 2, "createOrder", 11)
	if res.Status != shim.OK {
		expectApi(2, "createCommodity5")
		t.FailNow()
	}
	res = getTr(stub, []string{"order", "20211001101"})
	if res.Payload == nil {
		expectApi(2, "createCommodity5")
		t.FailNow()
	}
	expectApi(1, "createCommodity5")
}

// 查询商品列表-查询成功
func Test_queryCommodityList(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 1)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryCommodityList"),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(2, "queryCommodityList")
		t.FailNow()
	} else {
		expectApi(1, "queryCommodityList")
	}
}

// 查询订单列表-参数非空校验
func Test_queryOrderList1(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 2)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryOrderList"),
		[]byte(""),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(1, "queryOrderList1")
	} else {
		expectApi(2, "queryOrderList1")
		t.FailNow()
	}
}

// 查询订单列表-参数个数校验
func Test_queryOrderList2(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 2)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryOrderList"),
		[]byte("1"),
		[]byte("2"),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(1, "queryOrderList2")
	} else {
		expectApi(2, "queryOrderList2")
		t.FailNow()
	}
}

// 查询订单列表-查询成功
func Test_queryOrderList3(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 2)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryOrderList"),
		[]byte("20211001101"),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(2, "queryOrderList3")
		t.FailNow()
	} else {
		expectApi(1, "queryOrderList3")
	}
}

// 查询账户列表-参数非空校验
func Test_queryAccount1(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 3)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryAccount"),
		[]byte(""),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(1, "queryAccount1")
	} else {
		expectApi(2, "queryAccount1")
		t.FailNow()
	}
}

// 查询账户列表-参数个数校验
func Test_queryAccount2(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 3)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryAccount"),
		[]byte("111"),
		[]byte("2222"),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(1, "queryAccount2")
	} else {
		expectApi(2, "queryAccount2")
		t.FailNow()
	}
}

// 查询账户列表-查询成功
func Test_queryAccount3(t *testing.T) {
	stub := GetNewStub()

	putStateTransaction(stub, 3)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("queryAccount"),
		[]byte("1"),
	})
	t.Log(resp.Message)
	if len(resp.Payload) == 0 || bytes.NewBuffer(resp.Payload).String() == "[]" || resp.Status == shim.ERRORTHRESHOLD || resp.Status == shim.ERROR {
		expectApi(2, "queryAccount3")
		t.FailNow()
	} else {
		expectApi(1, "queryAccount3")
	}
}

// 更新订单温度列表-参数非空校验
func Test_updateOrderTemperature1(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 1, 3, 3, "updateOrderTemperature", "1")
}

// 更新订单温度列表-参数个数校验
func Test_updateOrderTemperature2(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 4, 4, 3, "updateOrderTemperature", "2")
}

// 更新订单温度列表-更新成功
func Test_updateOrderTemperature3(t *testing.T) {
	stub := GetNewStub()
	res := testSomeTx(stub, 3, "updateOrderTemperature", 5)
	if res.Status != shim.OK {
		expectApi(2, "updateOrderTemperature3")
		t.FailNow()
	}
	res = getTr(stub, []string{"order", "20211001101"})
	order := new(Order)
	_ = json.Unmarshal(res.Payload, order)
	if len(order.TemperatureVariation) != 0 {
		if order.TemperatureVariation[0].Temperature != 26 {
			expectApi(2, "updateOrderTemperature3")
			t.FailNow()
		} else {
			expectApi(1, "updateOrderTemperature3")
		}
	} else {
		expectApi(2, "updateOrderTemperature3")
		t.FailNow()
	}
}

// 更新订单状态-参数非空校验
func Test_updateOrderStatus1(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 1, 2, 3, "updateOrderStatus", "1")
}

// 更新订单状态-参数个数校验
func Test_updateOrderStatus2(t *testing.T) {
	stub := GetNewStub()
	loopCheck(stub, t, 3, 3, 3, "updateOrderStatus", "2")
}

// 更新订单状态-更新状态成功
func Test_updateOrderStatus3(t *testing.T) {
	stub := GetNewStub()
	res := testSomeTx(stub, 3, "updateOrderStatus", 4)
	if res.Status != shim.OK {
		expectApi(2, "updateOrderStatus3")
		t.FailNow()
	}
	res = getTr(stub, []string{"order", "20211001101"})
	order := new(Order)
	_ = json.Unmarshal(res.Payload, order)
	if order.Status != "" && (order.Status == "Processing" || order.Status == "运送中") {
		expectApi(1, "updateOrderStatus3")
	} else {
		expectApi(2, "updateOrderStatus3")
		t.FailNow()
	}
}

// 更新订单状态-订单完成时买家扣款成功
func Test_updateOrderStatus4(t *testing.T) {
	stub := GetNewStub()
	putStateTransaction(stub, 3)
	putStateTransaction(stub, 4)

	_ = stub.MockInvoke("1", [][]byte{
		[]byte("updateOrderStatus"),
		[]byte("20211001101"),
		[]byte("Done"),
	})
	key, _ := stub.CreateCompositeKey("account", []string{"3"})
	value, _ := stub.GetState(key)
	acc := new(Account)
	_ = json.Unmarshal(value, acc)
	if acc.Balance == 930 {
		expectApi(1, "updateOrderStatus4")
	} else {
		expectApi(2, "updateOrderStatus4")
		t.FailNow()
	}
}

// 更新订单状态-订单完成时供应商收款成功
func Test_updateOrderStatus5(t *testing.T) {
	stub := GetNewStub()
	putStateTransaction(stub, 3)
	putStateTransaction(stub, 4)

	_ = stub.MockInvoke("1", [][]byte{
		[]byte("updateOrderStatus"),
		[]byte("20211001101"),
		[]byte("Done"),
	})
	key, _ := stub.CreateCompositeKey("account", []string{"1"})
	value, _ := stub.GetState(key)
	acc := new(Account)
	_ = json.Unmarshal(value, acc)
	if acc.Balance == 1056 {
		expectApi(1, "updateOrderStatus5")
	} else {
		expectApi(2, "updateOrderStatus5")
		t.FailNow()
	}
}

// 更新订单状态-订单完成时物流商收款成功
func Test_updateOrderStatus6(t *testing.T) {
	stub := GetNewStub()
	putStateTransaction(stub, 3)
	putStateTransaction(stub, 4)

	_ = stub.MockInvoke("1", [][]byte{
		[]byte("updateOrderStatus"),
		[]byte("20211001101"),
		[]byte("Done"),
	})
	key, _ := stub.CreateCompositeKey("account", []string{"2"})
	value, _ := stub.GetState(key)
	acc := new(Account)
	_ = json.Unmarshal(value, acc)
	if acc.Balance == 1014 {
		expectApi(1, "updateOrderStatus6")
	} else {
		expectApi(2, "updateOrderStatus6")
		t.FailNow()
	}
}

func putStateTransaction(stub *shim.MockStub, status int) {
	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	switch status {
	case 1: // 买家购买商品
		commodity := &Commodity{
			Name:            "testBuy",
			Id:              "20211001001",
			Origin:          "中国",
			LowTemperature:  -10,
			HighTemperature: 30,
			Price:           10, //单价
			OwnerId:         "1",
		}
		bytes, _ := json.Marshal(commodity)
		compositeKey, _ := stub.CreateCompositeKey("commodity", []string{commodity.Id})
		_ = stub.PutState(compositeKey, bytes)
	case 2:
		commodity := &Commodity{
			Name:            "testBuy",
			Id:              "20211001001",
			Origin:          "中国",
			LowTemperature:  -10,
			HighTemperature: 30,
			Price:           10, //单价
			OwnerId:         "1",
		}
		bytes, _ := json.Marshal(commodity)
		compositeKey, _ := stub.CreateCompositeKey("commodity", []string{commodity.Id})
		_ = stub.PutState(compositeKey, bytes)
		order := &Order{
			Commodity: commodity,
			Id:        "20211001101",
			Quantity:  10,
			Status:    enumStatus.New,
			BuyerId:   "3",
			SellerId:  "1",
		}
		orderBytes, _ := json.Marshal(order)
		orderCompositeKey, _ := stub.CreateCompositeKey("order", []string{order.Id})
		_ = stub.PutState(orderCompositeKey, orderBytes)
	case 3:
		var accountsName = [3]string{"供货商", "物流商", "买家"}
		for i, val := range accountsName {
			account := &Account{
				Name:    val,
				Id:      strconv.Itoa(i + 1),
				Balance: 1000,
			}
			// 序列化对象
			accountBytes, _ := json.Marshal(account)
			accountCompositeKey, _ := stub.CreateCompositeKey("account", []string{account.Id})
			_ = stub.PutState(accountCompositeKey, accountBytes)
		}
	case 4:
		commodity := &Commodity{
			Name:            "testBuy",
			Id:              "20211001001",
			Origin:          "中国",
			LowTemperature:  -10,
			HighTemperature: 30,
			Price:           10, //单价
			OwnerId:         "1",
		}
		bytes, _ := json.Marshal(commodity)
		compositeKey, _ := stub.CreateCompositeKey("commodity", []string{commodity.Id})
		_ = stub.PutState(compositeKey, bytes)

		temperatures := make([]*Temperature, 0)
		temperature1 := &Temperature{
			Temperature: -11,
			RecordTime:  time.Now(),
		}
		temperature2 := &Temperature{
			Temperature: 31,
			RecordTime:  time.Now(),
		}
		temperatures = append(temperatures, temperature1)
		temperatures = append(temperatures, temperature2)

		order := &Order{
			Commodity:            commodity,
			Id:                   "20211001101",
			OrderTime:            time.Now(),
			DeliverAddress:       "testAddress",
			DeliverTime:          time.Now().Add(time.Hour * 72),
			Quantity:             10,
			Status:               enumStatus.Processing,
			TemperatureVariation: temperatures,
			BuyerId:              "3",
			SellerId:             "1",
		}
		orderBytes, _ := json.Marshal(order)
		orderCompositeKey, _ := stub.CreateCompositeKey("order", []string{order.Id})
		_ = stub.PutState(orderCompositeKey, orderBytes)
	}
}

func GetTxArgs(stub *shim.MockStub, funcName string, number int) [][]byte {
	switch funcName {
	case "createCommodity":
		switch number {
		case 1:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte(""),
				[]byte("20211001001"),
				[]byte("中国"),
				[]byte("-10"),
				[]byte("30"),
				[]byte("10"),
				[]byte("1"),
			}
		case 2:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte(""),
				[]byte("中国"),
				[]byte("-10"),
				[]byte("30"),
				[]byte("10"),
				[]byte("1"),
			}
		case 3:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte("20211001001"),
				[]byte(""),
				[]byte("-10"),
				[]byte("30"),
				[]byte("10"),
				[]byte("1"),
			}
		case 4:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte("20211001001"),
				[]byte("中国"),
				[]byte(""),
				[]byte("30"),
				[]byte("10"),
				[]byte("1"),
			}
		case 5:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte("20211001001"),
				[]byte("中国"),
				[]byte("-10"),
				[]byte(""),
				[]byte("10"),
				[]byte("1"),
			}
		case 6:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte("20211001001"),
				[]byte("中国"),
				[]byte("-10"),
				[]byte("30"),
				[]byte(""),
				[]byte("1"),
			}
		case 7:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte("20211001001"),
				[]byte("中国"),
				[]byte("-10"),
				[]byte("30"),
				[]byte("10"),
				[]byte(""),
			}
		case 8:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
			}
		case 9:
			return [][]byte{
				[]byte("createCommodity"),
				[]byte("testBuy"),
				[]byte("20211001001"),
				[]byte("中国"),
				[]byte("-10"),
				[]byte("30"),
				[]byte("10"),
				[]byte("1"),
			}
		default:
			return [][]byte{}
		}
	case "createOrder":
		switch number {
		case 1:
			return [][]byte{
				[]byte("createOrder"),
				[]byte(""),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 2:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte(""),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 3:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte(""),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 4:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(""),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 5:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte(""),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 6:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte(""),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 7:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte(""),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 8:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(""),
			}
		case 9:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
			}
		case 10:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("error"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 11:
			return [][]byte{
				[]byte("createOrder"),
				[]byte("20211001001"),
				[]byte("20211001101"),
				[]byte("testAddress"),
				[]byte(time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05")),
				[]byte("10"),
				[]byte("3"),
				[]byte("1"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		default:
			return [][]byte{}
		}
	case "updateOrderTemperature":
		switch number {
		case 1:
			return [][]byte{
				[]byte("updateOrderTemperature"),
				[]byte(""),
				[]byte("26"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 2:
			return [][]byte{
				[]byte("updateOrderTemperature"),
				[]byte("20211001101"),
				[]byte(""),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		case 3:
			return [][]byte{
				[]byte("updateOrderTemperature"),
				[]byte("20211001101"),
				[]byte("26"),
				[]byte(""),
			}
		case 4:
			return [][]byte{
				[]byte("updateOrderTemperature"),
				[]byte("20211001101"),
			}
		case 5:
			return [][]byte{
				[]byte("updateOrderTemperature"),
				[]byte("20211001101"),
				[]byte("26"),
				[]byte(time.Now().Format("2006-01-02 15:04:05")),
			}
		default:
			return [][]byte{}
		}
	case "updateOrderStatus":
		switch number {
		case 1:
			return [][]byte{
				[]byte("updateOrderStatus"),
				[]byte(""),
				[]byte("New"),
			}
		case 2:
			return [][]byte{
				[]byte("updateOrderStatus"),
				[]byte("20211001101"),
				[]byte(""),
			}
		case 3:
			return [][]byte{
				[]byte("updateOrderStatus"),
				[]byte("20211001101"),
			}
		case 4:
			return [][]byte{
				[]byte("updateOrderStatus"),
				[]byte("20211001101"),
				[]byte("Processing"),
			}
		case 5:
			return [][]byte{
				[]byte("updateOrderStatus"),
				[]byte("20211001101"),
				[]byte("Done"),
			}
		default:
			return [][]byte{}
		}
	default:
		return [][]byte{}
	}
}

func testSomeTx(stub *shim.MockStub, status int, funcName string, number int) pb.Response {
	if status != 1 {
		putStateTransaction(stub, status-1)
	}
	res := stub.MockInvoke("1", GetTxArgs(stub, funcName, number))
	return res
}

func loopCheck(stub *shim.MockStub, t *testing.T, start, end, status int, funcName string, flag string) {
	var res pb.Response
	for i := start; i <= end; i++ {
		res = testSomeTx(stub, status, funcName, i)
		t.Log(res.Status)
		t.Log(res.Message)
		if res.Status != shim.ERRORTHRESHOLD && res.Status != shim.ERROR {
			expectApi(2, funcName+flag)
			t.FailNow()
		}
	}
	expectApi(1, funcName+flag)
}

func getTr(stub *shim.MockStub, args []string) pb.Response {
	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	compositeKey, _ := stub.CreateCompositeKey(args[0], []string{args[1]})
	trByte, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(trByte)
}
