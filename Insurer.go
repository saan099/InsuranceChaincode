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

//========================================================InitInsurer============

func (t *InsuranceManagement) InitInsurer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("init Insurer called")
	fmt.Println("=========================================")
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::couldn't get creator"))
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::error unmarshalling"))
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("chaincode:InitInsurer::couldn pasre ParseCertificate")
	}

	insurerHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(insurerHash[:])

	checkInsurerAsBytes, err := stub.GetState(insurerAddress)
	if err != nil || len(checkInsurerAsBytes) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::Insurer already exist"))
	}

	insurer := Insurer{}
	insurer.InsurerId = insurerAddress
	insurer.InsurerName = cert.Subject.CommonName

	insurerAsBytes, err := json.Marshal(insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::couldn't Unmarsh creator"))
	}

	err = stub.PutState(insurerAddress, insurerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::couldn't write state "))
	}
	return shim.Success(nil)

}
