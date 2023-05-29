# DSCABS

## Introduction

如果你想在联盟链中开发智能合约，并且想对智能合约施加动态的、细粒度的访问控制，那么可以选择 `DSCABS`。

`DSCABS` 实现了访问控制逻辑与智能合约的业务逻辑相互解耦的目的，它可以将智能合约的访问策略和用户的属性分别转化成策略密钥和属性密钥，用户基于属性密钥可以对特定的消息进行签名，如果用户的属性集满足智能合约的访问策略，那么用户生成的签名即可通过策略密钥的验证，因此，`DSCABS` 可将复杂多变的权限判决过程（访问控制逻辑）用签名的验证过程代替。这样一来，即便合约的访问策略或者用户的属性发生了变化，访问控制逻辑也无需跟着改变，因此，突破了更改访问策略就得重新部署智能合约的限制，实现了对智能合约的动态访问控制目的。

## Usage

### 开发智能合约

首先，我们利用 `Go` 语言编写一个智能合约，然后利用 `DSCABS` 对我们开发的智能合约进行访问控制。具体步骤如下所示：

1. 新建一个名为 `contracts` 的文件夹 :file_folder:，然后进入 `contracts` 文件夹中，再新建两个文件夹 :file_folder:：`dog` 和 `cat`。

2. 我们进入 `dog` 文件夹中，新建 `dog.go` 文件，并在其中写入以下内容：
```go
package dog

import (
	"encoding/json"
	"errors"

	"github.com/232425wxy/dscabs/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DogContract struct {
	contractapi.Contract
	DSCABS *chaincode.DSCABS
}

type Dog struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Kind string `json:"kind"`
}

func (s *DogContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	dogs := []Dog{
		{Name: "dog1", Age: 1, Kind: "dog"},
		{Name: "dog2", Age: 2, Kind: "dog"},
		{Name: "dog3", Age: 3, Kind: "dog"},
		{Name: "dog4", Age: 4, Kind: "dog"},
	}

	dogsJSON, err := json.Marshal(dogs)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("dogs", dogsJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *DogContract) GetDog(ctx contractapi.TransactionContextInterface, 
	userID string, sig string, name string) (*Dog, error) {
	access, err := s.DSCABS.Access(ctx, userID, "DogContract", "GetDog", sig)
	if !access || err != nil {
		return nil, errors.New("forbidden access")
	}

	dogsJSON, err := ctx.GetStub().GetState("dogs")
	if err != nil {
		return nil, err
	}
	dogs := []Dog{}

	err = json.Unmarshal(dogsJSON, &dogs)
	if err != nil {
		return nil, err
	}

	for _, dog := range dogs {
		if name == dog.Name {
			return &dog, nil
		}
	}

	return nil, nil
}

```

3. 接着进入 `cat` 文件夹中，新建 `cat.go` 文件，然后在其中写入以下内容：
```go
package cat

import (
	"encoding/json"
	"errors"

	"github.com/232425wxy/dscabs/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CatContract struct {
	contractapi.Contract
	DSCABS *chaincode.DSCABS
}

type Cat struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Kind string `json:"kind"`
}

func (s *CatContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cats := []Cat{
		{Name: "cat1", Age: 1, Kind: "cat"},
		{Name: "cat2", Age: 2, Kind: "cat"},
		{Name: "cat3", Age: 3, Kind: "cat"},
		{Name: "cat4", Age: 4, Kind: "cat"},
	}

	catsJSON, err := json.Marshal(cats)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("cats", catsJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *CatContract) GetCat(ctx contractapi.TransactionContextInterface, userID string, sig string, name string) (*Cat, error) {
	access, err := s.DSCABS.Access(ctx, userID, "CatContract", "GetCat", sig)
	if !access || err != nil {
		return nil, errors.New("forbidden access")
	}

	catsJSON, err := ctx.GetStub().GetState("cats")
	if err != nil {
		return nil, err
	}
	cats := []Cat{}

	err = json.Unmarshal(catsJSON, &cats)
	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		if name == cat.Name {
			return &cat, nil
		}
	}

	return nil, nil
}
```

4. 返回到 `contracts` 目录下，新建 `main.go` 文件，并在其中写入以下内容：
```go
package main

import (
	"log"

	"github.com/232425wxy/contracts/cat"
	"github.com/232425wxy/contracts/dog"
	"github.com/232425wxy/dscabs/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	dscabs := &chaincode.DSCABS{}
	cat := &cat.CatContract{DSCABS: dscabs}
	dog := &dog.DogContract{DSCABS: dscabs}
	catContract, err := contractapi.NewChaincode(cat, dog, dscabs)
	if err != nil {
		log.Panicf("Error creating Contract chaincode: %v", err)
	}

	if err := catContract.Start(); err != nil {
		log.Panicf("Error starting Contract chaincode: %v", err)
	}
}
```

5. 打开终端，并切换到 `contracts` 目录下，执行 `go mod init` 命令，然后在新生成的 `go.mod` 文件中添加以下内容：
```go
go 1.16

require (
	github.com/232425wxy/dscabs v1.0.1
	github.com/hyperledger/fabric-contract-api-go v1.1.0
)
```

6. 在终端中先后执行 `go mod tidy` 和 `go mod vendor` 命令，分别从 `GitHub` 中拉取依赖并将依赖打包进 `vendor` 文件夹中。

### 部署智能合约

首先根据 [Hyperledger Fabric 文档](https://hyperledger-fabric.readthedocs.io/en/release-2.3/) 的指示，搭建一个版本为 `2.3` 的联盟链。然后将 `contracts` 文件中的智能合约部署到区块链中。

**部署智能合约的命令**

```sh
./network.sh deployCC -ccn contracts -ccp ../contracts -ccl go
```

注意：`network.sh` 文件是 `fabric-samples/test-network` 内的文件，从上面的命令也可以看出来，`contracts` 文件夹被放在了 `fabric-sample` 目录下。

**初始化DSCABS的命令**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"DSCABS:InitLedger","Args":["256"]}'
```

**初始化Cat和Dog两个合约的命令**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"CatContract:InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"DogContract:InitLedger","Args":[]}'
```

**为用户tom注册属性的命令**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"DSCABS:ExtractAK","Args":["tom","a,b,c"]}'
```

该命令执行完后会返回一段字符串，该段字符串即是 `tom` 的属性私钥，可以用来签名。

**为Cat合约设置访问策略的命令**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"DSCABS:GenPK","Args":["CatContract","GetCat","{a,b,c,d,[4,3]}"]}'
```

**为Dog合约设置访问策略的命令**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"DSCABS:GenPK","Args":["DogContract","GetDog","{a,b,c,d,[4,3]}"]}'
```

**用户tom试图调用Cat合约的GetCat方法**

首先，用户tom需要构建消息 `m`，因为是首次调用合约，因此 `m` 的内容是 `"0:tom:CatContract:GetCat"`，然后 `tom` 利用自己的属性私钥对消息 `m` 进行签名，并将签名放入调用合约时构造的参数列表中，命令如下所示：

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n contracts --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"CatContract:GetCat","Args":["tom","signature","cat2"]}'
```