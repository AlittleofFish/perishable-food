#!/bin/bash

# 遇到错误直接退出程序
set -e

# 一、环境清理
echo "一、环境清理"
mkdir -p channel-artifacts
mkdir -p crypto-config
rm -fr channel-artifacts/*
rm -fr crypto-config/*
echo "清理完毕"

# 二、生成证书和起始区块信息
echo "二、生成证书和起始区块信息"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile ThreeOrgsOrdererGenesis -channelID syschannel -outputBlock ./channel-artifacts/genesis.block

# 二、生成通道(这个动作会创建一个创世交易，也是该通道的创世交易)
# channel === 通道
echo "三、生成通道的TX文件(这个动作会创建一个创世交易，也是该通道的创世交易)"
configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannel

# 三.一 生成锚节点配置更新文件
echo "四、 生成锚节点配置更新文件"
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Organization1MSPanchors.tx -channelID mychannel -asOrg Organization1MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Organization2MSPanchors.tx -channelID mychannel -asOrg Organization2MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Organization3MSPanchors.tx -channelID mychannel -asOrg Organization3MSP

# 四、启动区块链网络
echo "区块链 ： 启动"
docker-compose -f docker-compose-cli.yaml up -d        # 按照docker-compose.yaml的配置启动区块链网络并在后台运行
echo "正在等待节点的启动完成，等待10秒"
sleep 10                    # 启动整个区块链网络需要一点时间，所以此处等待15s，让区块链网络完全启动

# 五、在区块链上按照刚刚生成的TX文件去创建通道
# 该操作和上面操作不一样的是，这个操作会写入区块链
echo "五、在区块链上按照刚刚生成的TX文件去创建通道"

docker exec cli peer channel create -o orderer.gdzce.cn:7050 -c mychannel -f ./channel-artifacts/mychannel.tx


# 六、让节点去加入到通道
# 所有节点都要加入通道中
echo "六、让节点去加入到通道"

docker exec cli peer channel join -b mychannel.block
# 修改环境变量链接到其他节点
docker exec -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051" cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node2.organization3.gdzce.cn:12051" cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Organization1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization1.gdzce.cn/users/Admin@organization1.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node2.organization1.gdzce.cn:8051" cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Organization2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051" cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Organization2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node2.organization2.gdzce.cn:10051" cli peer channel join -b mychannel.block

# 六.一 更新锚节点通道
echo "七、更新锚节点到通道"
docker exec -e "CORE_PEER_LOCALMSPID=Organization3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051" cli peer channel update -o orderer.gdzce.cn:7050 -c mychannel -f ./channel-artifacts/Organization3MSPanchors.tx
docker exec -e "CORE_PEER_LOCALMSPID=Organization1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization1.gdzce.cn/users/Admin@organization1.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization1.gdzce.cn:7051" cli peer channel update -o orderer.gdzce.cn:7050 -c mychannel -f ./channel-artifacts/Organization1MSPanchors.tx
docker exec -e "CORE_PEER_LOCALMSPID=Organization2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp" -e "CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051" cli peer channel update -o orderer.gdzce.cn:7050 -c mychannel -f ./channel-artifacts/Organization2MSPanchors.tx

# 七、链码安装
# -n 是链码的名字，可以自己随便设置
# -v 就是版本号，就是composer的bna版本
# -p 是目录，目录是基于cli这个docker里面的$GOPATH相对的
# 背书节点需要安装链码
echo "八、链码安装"
echo "node1.organization1 安装链码"
docker exec cli peer chaincode install -n mychaincode -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

echo "node1.organization3 安装链码"
docker exec -e CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051 -e CORE_PEER_LOCALMSPID=Organization3MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp cli  peer chaincode install -n mychaincode -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

echo "node1.organization2 安装链码"
docker exec -e CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051 -e CORE_PEER_LOCALMSPID=Organization2MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp cli  peer chaincode install -n mychaincode -v 1.0.0 -l golang -p github.com/chaincode/perishable-food

#八、实例化链码
#-n 对应前文安装链码的名字 其实就是composer network start bna名字
#-v 为版本号，相当于composer network start bna名字@版本号
#-C 是通道，在fabric的世界，一个通道就是一条不同的链，composer并没有很多提现这点，composer提现channel也就在于多组织时候的数据隔离和沟通使用
#-c 为传参，传入init参数
echo "九、实例化链码"
docker exec cli peer chaincode instantiate -o orderer.gdzce.cn:7050 -C mychannel -n mychaincode -l golang -v 1.0.0 -c '{"Args":["init"]}' -P 'AND("Organization3MSP.member","Organization1MSP.member","Organization2MSP.member")'
sleep 10
#请注意，安装链码是文件的复制，其实不等于我们电脑的安装，实例化才是真正的安装

# 进行链码查询，查询刚刚执行的智能合约是否将数据写入区块链
echo "十、通过调用链码使organization3与organization2上的链码实例化"
docker exec -e CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051 -e CORE_PEER_LOCALMSPID=Organization3MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp cli peer chaincode query -C mychannel -n mychaincode -c '{"Args":["queryAccount","1"]}'
docker exec -e CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051 -e CORE_PEER_LOCALMSPID=Organization2MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp cli peer chaincode query -C mychannel -n mychaincode -c '{"Args":["queryAccount","2"]}'

echo "区块链网络部署完成！"
