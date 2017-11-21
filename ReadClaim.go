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

//=================================== Read All Claim ================================================================
	func (t *InsuranceManagement) ReadAllClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Claim Id
		//args[1]= Claim report hash
			
		if len(args) != 0 && len(args)!=2{
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim: 0 or 2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim::account doesnt exists"))	
		}

		client:= Client{}

		err = json.Unmarshal(invokerAsBytes,&client)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:couldnt unmarshal invoker"))
		}
		claim:=Claim{}
		var claimsArr []Claim

		start:=0
		end:=len(client.Claims)-1
		
		if len(args) == 2 {
			start = strconv.Atoi(args[0])
			end = strconv.Atoi(args[1])
			if start > end {
				start = 0
				end = 0
			}
			if end > 
		}

		for i:=end ; i >=start ;i-- {
			claimsAsBytes ,err := stub.GetState(client.Claims[i])
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:couldnt getstate of "+client.Claims[i]))
			}
			err = json.Unmarshal(claimsAsBytes,&claim)
			claimsArr = append(claimsArr,claim)
		}
		claimsArrAsBytes,err := json.Marshal(claimsArr)

		return shim.Success(claimsArrAsBytes)
		
	}
