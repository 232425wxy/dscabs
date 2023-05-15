package chaincode

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

func (s *Gatekeeper) FuncXXX1(ctx contractapi.TransactionContextInterface, userID string, sig string) (bool, error) {
	return s.Access(ctx, userID, s.ContractName, "FuncXXX1", sig)
}

func (s *Gatekeeper) FuncXXX2(ctx contractapi.TransactionContextInterface, userID string, sig string) (bool, error) {
	return s.Access(ctx, userID, s.ContractName, "FuncXXX2", sig)
}

func (s *Gatekeeper) FuncXXX3(ctx contractapi.TransactionContextInterface, userID string, sig string) (bool, error) {
	return s.Access(ctx, userID, s.ContractName, "FuncXXX3", sig)
}

func (s *Gatekeeper) FuncXXX4(ctx contractapi.TransactionContextInterface, userID string, sig string) (bool, error) {
	return s.Access(ctx, userID, s.ContractName, "FuncXXX4", sig)
}
