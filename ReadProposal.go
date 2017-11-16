package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"bytes"
	//"strconv"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//=================================== Read All Proposal ================================================================
func (t *InsuranceManagement) ReadAllProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal::account doesnt exists"))
	
		}
		
		invokerClient := Client{}

		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal::Error Unmarshalling %s",err.Error()))
		}

		proposalArr:= invokerClient.ProposalArray

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false
		proposalobj:=Proposal{}
		for i:=0; i < len(proposalArr) ; i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			proposalAsBytes,err := stub.GetState(proposalArr[i])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readAllProposal::Could not get state of %s",proposalArr[i]))
			}
			buffer.WriteString(string(proposalAsBytes))
			err=json.Unmarshal(proposalAsBytes,&proposalobj)
			RFQAsBytes,_:=stub.GetState(proposalobj.RFQId)
			buffer.Truncate(buffer.Len()-1)
			buffer.WriteString(",")
			buffer.WriteString("\"rfqDetails\": ")
			buffer.WriteString(string(RFQAsBytes))
			buffer.WriteString("}")
			flag = true
		}
		buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}


//=================================== Read Single Proposal ================================================================

func (t *InsuranceManagement) ReadSingleProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::account doesnt exists"))
	
		}
		
		invokerClient := Client{}

		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::Error Unmarshalling %s",err.Error()))
		}

		proposalArr:= invokerClient.ProposalArray

		if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:: No proposals found in acount"))
		}

		flag:=0

		for i:=0;i < len(proposalArr);i++ {
			if proposalArr[i] == args[0] {
				flag = 1
				break
			}
		}
		if flag == 0 { return shim.Error("chaincode:readSingleProposal:: Propo Not found in account")}
		var buffer bytes.Buffer
		//buffer.WriteString("[")
		flag:=false
		proposalobj:=Proposal{}
		
			proposalAsBytes,err := stub.GetState(args[0])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::Could not get state of %s",args[0]))
			}

			buffer.WriteString(string(proposalAsBytes))
			err=json.Unmarshal(proposalAsBytes,&proposalobj)
			RFQAsBytes,_:=stub.GetState(proposalobj.RFQId)
			buffer.Truncate(buffer.Len()-1)
			buffer.WriteString(",")
			buffer.WriteString("\"rfqDetails\": ")
			buffer.WriteString(string(RFQAsBytes))
			buffer.WriteString("}")
			flag = true
		
		//buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}