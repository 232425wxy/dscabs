package chaincode

import (
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/232425wxy/dscabs/algorithm"
	"github.com/232425wxy/dscabs/compoments"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Gatekeeper struct {
	contractapi.Contract
	ContractName string
}

func (s *Gatekeeper) InitLedger(ctx contractapi.TransactionContextInterface, sl string, contractName string) error {
	s.ContractName = contractName
	securityLevel, _ := strconv.Atoi(sl)

	params := algorithm.Setup(securityLevel)

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(DSCABSMSK, paramsJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *Gatekeeper) ExtractAK(ctx contractapi.TransactionContextInterface, userID string, attributes string) (string, error) {
	if userID == "" {
		return "", errors.New("user id must be different from \"\"")
	}

	if attributes == "" {
		return "", errors.New("attributes must be different from \"\"")
	}

	params := &algorithm.SystemParams{Curve: new(elliptic.CurveParams)}

	paramsJSON, err := ctx.GetStub().GetState(DSCABSMSK)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(paramsJSON, &params)
	if err != nil {
		return "", err
	}

	var ak *algorithm.AttributeKey

	if strings.Contains(attributes, ",") {
		attributesSlice := strings.Split(attributes, ",")
		ak = compoments.AddUserAttributes(params, userID, attributesSlice)
	} else {
		ak = compoments.AddUserAttributes(params, userID, []string{attributes})
	}

	akJSON, err := json.Marshal(ak)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(AKTag(userID), akJSON)
	if err != nil {
		return "", err
	}

	return ak.SecretKey.String(), nil
}

func (s *Gatekeeper) GenPK(ctx contractapi.TransactionContextInterface, contractName string, functionName string, policy string) error {
	if contractName == "" {
		return errors.New("contract name must be different from \"\"")
	}

	if functionName == "" {
		return errors.New("function name must be different from \"\"")
	}

	if policy == "" {
		return errors.New("policy must be different from \"\"")
	}

	params := &algorithm.SystemParams{Curve: new(elliptic.CurveParams)}

	paramsJSON, err := ctx.GetStub().GetState(DSCABSMSK)
	if err != nil {
		return err
	}

	err = json.Unmarshal(paramsJSON, &params)
	if err != nil {
		return err
	}

	compoments.AddSmartContractFunctionPolicy(params, contractName, functionName, policy)

	return nil
}

func (s *Gatekeeper) Access(ctx contractapi.TransactionContextInterface, userID string, contractName string, functionName string, sig string) (bool, error) {
	params := &algorithm.SystemParams{Curve: new(elliptic.CurveParams)}

	paramsJSON, err := ctx.GetStub().GetState(DSCABSMSK)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(paramsJSON, &params)
	if err != nil {
		return false, err
	}

	return compoments.GateKeeper(params, userID, contractName, functionName, sig)
}
