package main 

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	//"bytes"
	//"strconv"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)


//=========================================================Init Surveyor =================================================

func (t *InsuranceManagement) InitSurveyor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Println("init client called")
	//fmt.Println("=========================================")
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitSurveyor::couldn't get creator"))
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitSurveyor::error unmarshalling"))
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("chaincode:InitSurveyor::couldn pasre ParseCertificate")
	}

	invokerHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	invokerAddress := hex.EncodeToString(invokerHash[:])

	invokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || len(invokerAsBytes) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:InitSurveyor::Surveyor already exist"))
	}

	surveyor := Surveyor{}
	surveyor.SurveyorId = invokerAddress
	surveyor.SurveyorName = cert.Subject.CommonName

	invokerAsBytes, err = json.Marshal(surveyor)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitSurveyor::couldn't Unmarsh creator"))
	}

	err = stub.PutState(invokerAddress, invokerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitSurveyor::couldn't write state "))
	}
	return shim.Success(nil)

}