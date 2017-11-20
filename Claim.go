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


//=================================== Generate Claim ================================================================
func (t *InsuranceManagement) GenerateClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Intimation Date
		//args[1]= Loss Date
		//args[2]= Loss Description
		//args[3]= Policy Number
		//args[4]= Claim amount
	
		if len(args) != 5 {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim:5 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim::account doesnt exists"))	
		}

		client:= Client{}
		flag:=false
		
		for i:=0 ; i < len(client.Policies); i++ {
			if args[3] == client.Policies[i] {
				flag = true 
				break
			}
		}

		if flag == false {
			return shim.Error("chaincode:generateClaim::Policy number doesnt exist in account")
		}

		err = json.Unmarshal(invokerAsBytes,&client)

		claim:=Claim{}
		// add claim details
		claim.IntimationDate = args[0]
		claim.LossDate = args[1]
		claim.LossDescription = args[2]
		claim.PolicyNumber = args[3]
		claim.ClaimAmount = args[4]
		claim.InsuredName= client.ClientName
		claim.ClientId = client.ClientId
		claim.Status = CLAIM_INITIALIZED
		claim.ClaimId = stub.GetTxID()

		client.Claims = append(client.Claims,claim.ClaimId)

		policyAsBytes,err:=stub.GetState(claim.PolicyNumber)
		policy:=Policy{}
		err=json.Unmarshal(policyAsBytes,&policy)
		insurer:=Insurer{}

		//update every insurer with new generated claim
		for i:=0 ; i < len(policy.Details.SelectedInsurer) ; i++ {
			insurerAsBytes,err:= stub.GetState(policy.Details.SelectedInsurer[i])
			if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim:couldnt get Insurers"))
		}
			err = json.Unmarshal(insurerAsBytes,&insurer)
			if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim:couldnt UNmarshal insurer"))
		}
			insurer.Claims = append(insurer.Claims,claim.ClaimId)
			insurerAsBytes,err = json.Marshal(insurer)
				err=stub.PutState(insurer.InsurerId,insurerAsBytes)
		}
		insurerAsBytes,err:=stub.GetState(policy.Details.LeadInsurer)
		err = json.Unmarshal(insurerAsBytes,&insurer)
		insurer.Claims = append(insurer.Claims, claim.ClaimId)

		claimAsBytes,err:= json.Marshal(claim)
		err = stub.PutState(claim.ClaimId,claimAsBytes)
		
		clientAsBytes,err := json.Marshal(client)
		err = stub.PutState(client.ClientId,clientAsBytes)

		return shim.Success(nil)
	}


	//=================================== Assign Surveyor To Claim ================================================================
func (t *InsuranceManagement) AssignSurveyorToClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Claim Id
		//args[1]= Surveyor Id
			
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim:2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim::account doesnt exists"))	
		}
		claim:=Claim{}
		claimAsBytes,err := stub.GetState(args[0])
		if err != nil || len(claimAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim::Claim Id doesnt exist"))	
		}
		err=json.Unmarshal(claimAsBytes,claim)

		//check if the invoker is Lead Insurer
		policyAsBytes,err:=stub.GetState(claim.PolicyNumber)
		policy:=Policy{}
		err = json.Unmarshal(policyAsBytes,&policy)
		if policy.Details.LeadInsurer != invokerAddress {
			return shim.Error("chaincode:asignSurveyorToClaim::Only Lead insurer is allowed to assign surveyor")
		}

		//assign surveyor
		claim.Surveyor = args[1]
		claim.Status = CLAIM_SURVEYOR_ASSIGNED	
		claimAsBytes,err = json.Marshal(claim)	

		err = stub.PutState(claim.ClaimId,claimAsBytes) // update state of claim

		return shim.Success(nil)	
	}