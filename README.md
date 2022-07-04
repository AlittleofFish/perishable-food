# 生鲜食品 Fabric 区块链网络

## 环境依赖版本

 * Ubuntu 20.04 LTS 64位
 * Hyperledger Fabric 1.4
 * Golang 1.17 以上
 * Docker 20.10 以上

## 目录结构

 * `application` 外部应用程序（调用fabric-sdk与区块链交互）
 * `chaincode` 链码
 * `deploy` 区块链网络配置相关文件及启动关闭脚本
 * `vendor` 依赖

## 启动区块链网络

本项目应放在 `$GOPATH/src/gdzce.cn/` 目录下

```bash
mkdir -p $GOPATH/src/gdzce.cn
cd $GOPATH/src/gdzce.cn
```

`git clone` 之后 `cd $GOPATH/src/gdzce.cn/perishable-food` 执行以下命令，将 `.sh` 脚本转为可执行文件

```bash
cd $GOPATH/src/gdzce.cn/perishable-food
chmod +x deploy/*.sh
```

需要联网拉取相关依赖

```bash
cd $GOPATH/src/gdzce.cn/perishable-food
## 拉取依赖
go mod tidy
```

deploy 目录下文件相关操作

```bash
cd deploy
# 启动
./start.sh
# 停止并删除
./stop.sh

# 其他
docker-compose down     # 停止
docker-compose restart  # 重新启动（不执行stop.sh脚本的情况下）
# 推荐使用.sh脚本
```

## 编译应用程序(app)

```bash
cd $GOPATH/src/gdzce.cn/perishable-food/application
go build
```

成功编译后会出现一个可执行文件 `application`

启动 app server

```bash
./application
```

页面访问

localhost:8080/web