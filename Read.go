package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"bytes"
	"strconv"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//===================================Read All RFQ ================================================================
func (t *InsuranceManagement) ReadAllRFQ(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::account doesnt exists"))
	
		}
		invokerClient := Client{}
		//invokerInsurer:= Insurer{}
		//invokerBroker:= Broker{}
		//flag_type:= 0
		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::Unmarshal error:: %s",err.Error()))
		}
		var rfqId []string
		//var rfqArr []RFQ	
		
		rfqId = invokerClient.RFQArray
		if len(rfqId) == 0 {
		
			return shim.Error("chaincode:readAllRFQ::No RFQs to read")
		}

		
		var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false
		for i:=0 ; i < len(rfqId); i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			rfqAsBytes,_ := stub.GetState(rfqId[i])
			buffer.WriteString(string(rfqAsBytes))
			flag = true
			//err = json.Unmarshal(rfqAsBytes,&rfqArr[i])		
		}
		buffer.WriteString("]")
		
		return shim.Success(buffer.Bytes())

	}


	func (t *InsuranceManagement) ReadSingleRFQ(stub shim.ChaincodeStubInterface,args []string) pb.Response{
		//args[0]=RFQId
		
		
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ:1 argument expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ::account doesnt exists"))
	
		}
		invokerClient := Client{}
		//invokerInsurer:= Insurer{}
		//invokerBroker:= Broker{}
		//flag_type:= 0
		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil { 			
			return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::Error Unmarshalling %s",err.Error()))				
		}
		var rfqId []string
		
		rfqId = invokerClient.RFQArray
		if len(rfqId) == 0 {		
			return shim.Error("chaincode:readSingleRFQ::No RFQs to read in this account")
		}

		flag:=0

		for i:=0;i < len(rfqId);i++ {
			if rfqId[i] == args[0] {
				flag = 1
			}
		}

		if flag == 0 { return shim.Error("chaincode:readSingleRFQ:: RFQId Not found in account")}

		rfqAsBytes,err:=stub.GetState(args[0])
		if err != nil {
			return shim.Error("chaincode:readSingleRFQ:: RFQId Not found in WorldState")
		}

		return shim.Success(rfqAsBytes)	
	}


func (t *InsuranceManagement) ReadRFQByRange(stub shim.ChaincodeStubInterface,args []string) pb.Response{
		//args[0]=start
		//args[1]=end
		
		
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("chaincode:readRFQByRange:2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readRFQByRange:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readRFQByRange:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::account doesnt exists"))
	
		}
		invokerClient := Client{}
		//invokerInsurer:= Insurer{}
		//invokerBroker:= Broker{}
		//flag_type:= 0
		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil { 
			return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::Error Unmarshalling %s",err.Error()))			
		}
		var rfqId []string
		//var rfqArr []RFQ
		
		rfqId = invokerClient.RFQArray
		
		if len(rfqId) == 0 {		
			return shim.Error("chaincode:readRFQByRange::No RFQs to read in this account")
		}

		start,err:=strconv.Atoi(args[0])
		if err!=nil {
			return shim.Error("cannot convert start string to int")
		}
		end,err:=strconv.Atoi(args[1])
		if err!=nil {
			return shim.Error("cannot convert start string to int")
		}
		
		if end > len(rfqId) {
			return shim.Error("chaincode:readRFQByRange::End range exceeded")
		}
		var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false

		for i:=start; i < end ; i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			rfqAsBytes,err := stub.GetState(rfqId[i])
			if err!=nil { 
				return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::Could not get state of %s",rfqId[i]))}
			buffer.WriteString(string(rfqAsBytes))
			flag = true
		}
		buffer.WriteString("]")
		
		return shim.Success(buffer.Bytes())
	}

	func (t *InsuranceManagement) ReadClientOfBroker(stub shim.ChaincodeStubInterface,args []string) pb.Response{
		
		
		
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker:0 argument expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker::account doesnt exists"))
	
		}
		//client := Client{}
		invokerBroker:= Broker{}
		//flag_type:= 0
		
		err = json.Unmarshal(invokerAsBytes,&invokerBroker)

		if len(invokerBroker.Clients) == 0 {		
			return shim.Error("chaincode:readClientOfBroker::No Clients to read in this account")
		}
			var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false
		for i:=0 ; i < len(invokerBroker.Clients); i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			clientAsBytes,err := stub.GetState(invokerBroker.Clients[i])
			if err!=nil { 
				return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker::Could not get state of %s",invokerBroker.Clients[i]))}
			buffer.WriteString(string(clientAsBytes))
			flag = true
		}
		buffer.WriteString("]")
		return shim.Success(buffer.Bytes())
	}