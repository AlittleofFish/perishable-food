#!/bin/bash

# 遇到错误直接退出程序
set -e


echo "八、链码安装"
echo "node1.organization1 安装链码"
docker exec cli peer chaincode install -n mycc -v 1.0.0 -l golang -p github.com/chaincode/chaincode_example02/go

echo "node1.organization3 安装链码"
docker exec -e CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051 -e CORE_PEER_LOCALMSPID=Organization3MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp cli  peer chaincode install -n mycc -v 1.0.0 -l golang -p github.com/chaincode/chaincode_example02/go

echo "node1.organization2 安装链码"
docker exec -e CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051 -e CORE_PEER_LOCALMSPID=Organization2MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp cli  peer chaincode install -n mycc -v 1.0.0 -l golang -p github.com/chaincode/chaincode_example02/go


echo "九、实例化链码"
docker exec cli peer chaincode instantiate -o orderer.gdzce.cn:7050 -C mychannel -n mycc -l golang -v 1.0.0 -c '{"Args":["init","a","10000","b","20000"]}' -P 'AND("Organization3MSP.member","Organization1MSP.member","Organization2MSP.member")'
sleep 10
#请注意，安装链码是文件的复制，其实不等于我们电脑的安装，实例化才是真正的安装

# 进行链码查询，查询刚刚执行的智能合约是否将数据写入区块链
echo "十、通过调用链码使organization3与organization2上的链码实例化"
docker exec -e CORE_PEER_ADDRESS=node1.organization3.gdzce.cn:11051 -e CORE_PEER_LOCALMSPID=Organization3MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization3.gdzce.cn/users/Admin@organization3.gdzce.cn/msp cli peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
docker exec -e CORE_PEER_ADDRESS=node1.organization2.gdzce.cn:9051 -e CORE_PEER_LOCALMSPID=Organization2MSP -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/organization2.gdzce.cn/users/Admin@organization2.gdzce.cn/msp cli peer chaincode query -C mychannel -n mycc -c '{"Args":["query","b"]}'
