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

//=================================== Read All Quote ================================================================
func (t *InsuranceManagement) ReadAllQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllQuote:0 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllQuote:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllQuote:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllQuote::account doesnt exists"))
	
		}
		
		invokerInsurer := Insurer{}

		err = json.Unmarshal(invokerAsBytes,&invokerInsurer)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllQuote::Error Unmarshalling %s",err.Error()))
		}

		quoteArr:= invokerInsurer.Quotes

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		var buffer bytes.Buffer
		buffer.WriteString("[")
		flag:=false
		//proposalobj:=Proposal{}
		for i:=0; i < len(quoteArr) ; i++ {
			if flag == true {
				buffer.WriteString(",")
			}
			quoteAsBytes,err := stub.GetState(quoteArr[i])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readAllQuote::Could not get state of %s",quoteArr[i]))
			}
			buffer.WriteString(string(quoteAsBytes))
			
			flag = true
		}
		buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}


//=================================== Read Quote By Range ================================================================
func (t *InsuranceManagement) ReadQuoteByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange:2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange::account doesnt exists"))
	
		}
		
		invokerInsurer := Insurer{}

		err = json.Unmarshal(invokerAsBytes,&invokerInsurer)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange::Error Unmarshalling %s",err.Error()))
		}

		quoteArr:= invokerInsurer.Quotes

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		start,err:= strconv.Atoi(args[0])
		if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange::Could not convert %s to int",args[0]))
			}
		end,err:= strconv.Atoi(args[1])
		if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange::Could not convert %s to int",args[1]))
			}
		if end > len(quoteArr) {
			//return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange::End limit exceeded"))
			end= len(quoteArr)
		}

		if start > len(quoteArr) {
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
			quoteAsBytes,err := stub.GetState(quoteArr[i])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readQuoteByRange::Could not get state of %s",quoteArr[i]))
			}
			buffer.WriteString(string(quoteAsBytes))
			
			flag = true
		}
		buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}


//=================================== Read Single Quote ================================================================
func (t *InsuranceManagement) ReadSingleQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleQuote:1 argument expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleQuote:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleQuote:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleQuote::account doesnt exists"))
	
		}
		
		invokerInsurer := Insurer{}

		err = json.Unmarshal(invokerAsBytes,&invokerInsurer)
		if err!=nil {
			return shim.Error(fmt.Sprintf("chaincode:readSingleQuote::Error Unmarshalling %s",err.Error()))
		}

		var quoteArr []string
		quoteArr = invokerInsurer.Quotes

		/*if len(proposalArr) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
		}*/

		//var buffer bytes.Buffer
		//buffer.WriteString("[")
		flag:=0
		//proposalobj:=Proposal{}
		for i:=0; i < len(quoteArr) ; i++ {
			if quoteArr[i] == args[0] {
				flag = 1
				//break
			}
		}
		var buffer bytes.Buffer
		if flag == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSingleQuote:: Quote not found in account")) 
		}
			quoteObj:=Quote{}
			quoteAsBytes,err := stub.GetState(args[0])
			if err!=nil {
				return shim.Error(fmt.Sprintf("chaincode:readSingleQuote::Could not get state of %s",args[0]))
			}
			buffer.WriteString(string(quoteAsBytes))
			err = json.Unmarshal(quoteAsBytes,&quoteObj)
			rfqAsBytes,err :=stub.GetState(quoteObj.RFQId)
			str:= "\""+quoteObj.RFQId + "\""
				var a string
				//if rfqObj.LeadInsurer == rfqObj.SelectedInsurer[j]{
				a=string(bytes.Replace(buffer.Bytes(),[]byte(str),[]byte(rfqAsBytes),1)) 
				buffer.Reset()
				buffer.WriteString(a)
			//quoteObj.RFQId = string(rfqAsBytes)
			//quoteAsBytes,err = json.Marshal(quoteObj)
			//buffer.WriteString(string(quoteAsBytes))
			
			//flag = true
		
		//buffer.WriteString("]")

		return shim.Success(buffer.Bytes())
}