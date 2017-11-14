package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)


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
			shim.Error(fmt.Sprintf("chaincode:readAllRFQ::account doesnt exists"))
	
		}
		invokerClient := Client{}
		invokerInsurer:= Insurer{}
		invokerBroker:= Broker{}
		flag_type:= 0
		err = json.Unmarshal(invokerAsBytes,&invokerClient)
		if err!=nil { flag_type =1}
		err = json.Unmarshal(invokerAsBytes,&invokerInsurer)
		if err!=nil { flag_type=2 }
		err = json.Unmarshal(invokerAsBytes,&invokerBroker)
		if err!=nil { return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::account is not of any type"))}
		var rfqId []string
		//var rfqArr []RFQ
		
		if flag_type ==0 {rfqId = invokerClient.RFQArray}
		if flag_type ==1 {rfqId = invokerInsurer.RFQArray}
		if flag_type ==2 {rfqId = invokerBroker.RFQArray}

		
		if len(rfqId) == 0 {
		
			 shim.Error("chaincode:readAllRFQ::No RFQs to read")
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