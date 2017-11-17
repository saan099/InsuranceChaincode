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

//=================================== Read All Policy ================================================================
func (t *InsuranceManagement) ReadAllPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllPolicy:0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllPolicy:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllPolicy:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllPolicy::account doesnt exists"))
	
		}
		
		invokerClient := Client{}

		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllPolicy::Error Unmarshalling %s",err.Error()))
		}

		policyArr:= invokerClient.Policies

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false
		//proposalobj:=Proposal{}
		for i:=0; i < len(policyArr) ; i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			policyAsBytes,err := stub.GetState(policyArr[i])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readAllPolicy::Could not get state of %s",policyArr[i]))
			}
			buffer.WriteString(string(policyAsBytes))
			
			flag = true
		}
		buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}


//=================================== Read Policy By Range ================================================================
func (t *InsuranceManagement) ReadPolicyByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::account doesnt exists"))
	
		}
		
		invokerClient := Client{}

		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::Error Unmarshalling %s",err.Error()))
		}

		policyArr:= invokerClient.Policies

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		start,err:= strconv.Atoi(args[0])
		if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::Could not convert %s to int",args[0]))
			}
		end,err:= strconv.Atoi(args[1])
		if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::Could not convert %s to int",args[1]))
			}
		if end > len(policyArr) {
			//return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::End limit exceeded"))
			end=len(policyArr)
		}

		if start > len(policyArr) {
			start =0 
			end =0
		}

		var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false
		//proposalobj:=Proposal{}
		for i:=start; i < end ; i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			policyAsBytes,err := stub.GetState(policyArr[i])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::Could not get state of %s",policyArr[i]))
			}
			buffer.WriteString(string(policyAsBytes))
			
			flag = true
		}
		buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}


//=================================== Read Single Polciy ================================================================
func (t *InsuranceManagement) ReadSinglePolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy:1 argument expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::account doesnt exists"))
	
		}
		
		invokerClient := Client{}

		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Error Unmarshalling "))
		}

		var policyArr []string
		
		policyArr = invokerClient.Policies

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		//var buffer bytes.Buffer
		//buffer.WriteString("[")
		flag:=false
		//proposalobj:=Proposal{}
		for i:=0; i < len(policyArr) ; i++ {
			if policyArr[i] == args[0] {
				flag = true
				break
			}
		}

		if flag == false {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy:: Policy not found in account")) 
		}

			policyAsBytes,err := stub.GetState(args[0])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Could not get state of %s",args[0]))
			}
			//buffer.WriteString(string(quoteAsBytes))
			
			//flag = true
		
		//buffer.WriteString("]")

		return shim.Success(policyAsBytes)
}