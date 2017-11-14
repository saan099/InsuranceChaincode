package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//=========================================================InitClient=================================================

func (t *InsuranceManagement) InitClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("init client called")
	fmt.Println("=========================================")
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClient::couldn't get creator"))
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClient::error unmarshalling"))
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("chaincode:InitClient::couldn pasre ParseCertificate")
	}

	invokerHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	clientAddress := hex.EncodeToString(invokerHash[:])

	checkClientAsBytes, err := stub.GetState(clientAddress)
	if err != nil || len(checkClientAsBytes) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:InitClient::client already exist"))
	}

	client := Client{}
	client.ClientId = clientAddress
	client.ClientName = cert.Subject.CommonName

	clientAsBytes, err := json.Marshal(client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClient::couldn't Unmarsh creator"))
	}

	err = stub.PutState(clientAddress, clientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClient::couldn't write state "))
	}
	return shim.Success(nil)

}
