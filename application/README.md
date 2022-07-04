# 生鲜食品供应链平台应用程序



## 依赖安装

执行 `govendor sync`

*可能需要挂代理*



## 编译

本目录下执行`go build`



## 运行

直接运行编译出的文件即可，默认名称为`application`



## 目录结构说明

 * `blockchain` 封装了fabric的sdk

 * `controller` http服务相关业务逻辑

 * `fbeecloud` 温度传感器相关

 * `lib/type.go` 共用类型定义

 * `public` 前端静态文件

 * `repository/transactionRecord.go` 处理订单id与交易id对应关系的类

 * `util` 工具

