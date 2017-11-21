package main

import (
	//"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	//"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//===================================Client Read Meta Data ================================================================
	func (t *InsuranceManagement) ReadMetaDataClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
					
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataClient: 0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataClient:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataClient:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataClient::account doesnt exists"))	
		}

		type ReadData struct {
			Id			      string 			`json:"id"`
			RegisteredName	  string 			`json:"registeredName"`
			UserName		  string			`json:"userName"`
		}

		

		invoker:=Client{}
		err = json.Unmarshal(invokerAsBytes,&invoker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataClient:couldnt unmarshal client"))
		}
		read:=ReadData{}
		read.Id = invoker.ClientId
		read.RegisteredName = invoker.ClientName
		read.UserName = invoker.UserName

		readAsBytes,err := json.Marshal(read)

		return shim.Success(readAsBytes)

	}

	//===================================Broker Read Meta Data ================================================================
	func (t *InsuranceManagement) ReadMetaDataBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
					
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataBroker: 0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataBroker:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataBroker:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataBroker::account doesnt exists"))	
		}

		type ReadData struct {
			Id			      string 			`json:"id"`
			RegisteredName	  string 			`json:"registeredName"`
			UserName		  string			`json:"userName"`
		}

		

		invoker:=Broker{}
		err = json.Unmarshal(invokerAsBytes,&invoker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataBroker:couldnt unmarshal broker"))
		}
		read:=ReadData{}
		read.Id = invoker.BrokerId
		read.RegisteredName = invoker.BrokerName
		read.UserName = invoker.UserName

		readAsBytes,err := json.Marshal(read)

		return shim.Success(readAsBytes)

	}
	//===================================Insurer Read Meta Data ================================================================
	func (t *InsuranceManagement) ReadMetaDataInsurer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
					
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataInsurer: 0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataInsurer:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataInsurer:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataInsurer::account doesnt exists"))	
		}

		type ReadData struct {
			Id			      string 			`json:"id"`
			RegisteredName	  string 			`json:"registeredName"`
			UserName		  string			`json:"userName"`
		}

		

		invoker:=Insurer{}
		err = json.Unmarshal(invokerAsBytes,&invoker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataInsurer:couldnt unmarshal insurer"))
		}
		read:=ReadData{}
		read.Id = invoker.InsurerId
		read.RegisteredName = invoker.InsurerName
		read.UserName = invoker.UserName

		readAsBytes,err := json.Marshal(read)

		return shim.Success(readAsBytes)

	}
	//===================================Surveyor Read Meta Data ================================================================
	func (t *InsuranceManagement) ReadMetaDataSurveyor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
					
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataSurveyor: 0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataSurveyor:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataSurveyor:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataSurveyor::account doesnt exists"))	
		}

		type ReadData struct {
			Id			      string 			`json:"id"`
			RegisteredName	  string 			`json:"registeredName"`
			UserName		  string			`json:"userName"`
		}

		

		invoker:=Surveyor{}
		err = json.Unmarshal(invokerAsBytes,&invoker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadMetaDataSurveyor:couldnt unmarshal surveyor"))
		}
		read:=ReadData{}
		read.Id = invoker.SurveyorId
		read.RegisteredName = invoker.SurveyorName
		read.UserName = invoker.UserName

		readAsBytes,err := json.Marshal(read)

		return shim.Success(readAsBytes)

	}