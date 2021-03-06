package main

import (
	//"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//=================================== Read All Claim ================================================================
	func (t *InsuranceManagement) ReadAllClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		
			
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
		//claim:=Claim{}
		//var claimsArr []Claim

		start:=len(client.Claims)-1
		end:=0
		
		if len(args) == 2 {
		lowerLimit, err := strconv.Atoi(args[0])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:lowerlimit not integer"))
		}
		upperLimit, err := strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:upperlimit not integer"))
		}
		if upperLimit < lowerLimit {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:upperlimit is not bigger than lowerlimit"))
		}
		if upperLimit >= len(client.Claims) {
			end = 0

		} else {
			end = len(client.Claims) - 1 - upperLimit
		}
		if lowerLimit < 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:lowerlimit is less than 0"))
		}
		start = len(client.Claims) - 1 - lowerLimit
	}
type ReadClaim struct {
	ClaimId         string  `json:"claimId"`
	ClaimType       string  `json:"claimType"`
	ClientId        string  `json:"clientId"`
	IntimationDate  string  `json:"intimationDate"`
	LossDate        string  `json:"lossDate"`
	PolicyNumber    string  `json:"policyNumber"`
	InsuredName     string  `json:"insuredName"`
	InsuredPhone    string  `json:"insuredPhone"`
	InsuredAddress  string  `json:"insuredAddress"`
	InsuredEmail    string  `json:"insuredEmail"`     
	LossDescription string  `json:"lossDescription"`
	ClaimAmount     float64 `json:"claimAmount"`
	ApprovedAmount  float64 `json:"approvedAmount"`
	Status          string  `json:"status"`
	Surveyor        string  `json:"surveyor"`
	Report          string  `json:"report"`
	Policy			Policy 	`json:"policy"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}
readclaim:=ReadClaim{}
var readclaimsArr []ReadClaim
policy:= Policy{}
		for i:=start ; i >=end ;i-- {
			claimsAsBytes ,err := stub.GetState(client.Claims[i])
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim:couldnt getstate of "+client.Claims[i]))
			}
			err = json.Unmarshal(claimsAsBytes,&readclaim)
			
			policyAsBytes,err := stub.GetState(readclaim.PolicyNumber)
			err = json.Unmarshal(policyAsBytes,&policy)
			readclaim.Policy = policy

			readclaimsArr = append(readclaimsArr,readclaim)
			
		}
		claimsArrAsBytes,err := json.Marshal(readclaimsArr)

		return shim.Success(claimsArrAsBytes)	
	}

//=================================== Read Single Claim ================================================================
	func (t *InsuranceManagement) ReadSingleClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
			
		if len(args) != 1{
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaim: 1 argument expected"))
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
		//claim:=Claim{}
		flag:=false

		for i:=range client.Claims {
			if args[0] == client.Claims[i]{
				flag=true
				break
			}
		}
		if flag == false {
			return shim.Error("chaincode:ReadAllClaim:claim id not found")
		}
		
		claimsAsBytes,err:= stub.GetState(args[0])

		type ReadClaim struct {
	ClaimId         string  `json:"claimId"`
	ClaimType       string  `json:"claimType"`
	ClientId        string  `json:"clientId"`
	IntimationDate  string  `json:"intimationDate"`
	LossDate        string  `json:"lossDate"`
	PolicyNumber    string  `json:"policyNumber"`
	InsuredName     string  `json:"insuredName"`
	InsuredPhone    string  `json:"insuredPhone"`
	InsuredAddress  string  `json:"insuredAddress"`
	InsuredEmail    string  `json:"insuredEmail"`     
	LossDescription string  `json:"lossDescription"`
	ClaimAmount     float64 `json:"claimAmount"`
	ApprovedAmount  float64 `json:"approvedAmount"`
	Status          string  `json:"status"`
	Surveyor        string  `json:"surveyor"`
	Report          string  `json:"report"`
	Policy			Policy 	`json:"policy"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}

	readclaim:= ReadClaim{}
	err = json.Unmarshal(claimsAsBytes,&readclaim)
	claim:=Claim{}
	err = json.Unmarshal(claimsAsBytes,&claim)

	policyAsBytes,err :=  stub.GetState(claim.PolicyNumber)
	policy:=Policy{}
	err = json.Unmarshal(policyAsBytes,&policy)

	readclaim.Policy=policy
	readAsBytes,err := json.Marshal(readclaim)

		return shim.Success(readAsBytes)

	}


//=================================== Read All Claim Surveyor ================================================================
	func (t *InsuranceManagement) ReadAllClaimSurveyor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	

		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor: 1 argument expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor::account doesnt exists"))	
		}

		surveyor:= Surveyor{}

		err = json.Unmarshal(invokerAsBytes,&surveyor)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor:couldnt unmarshal invoker"))
		}
		//claim:=Claim{}
		//var claimsArr []Claim

		// start:=len(surveyor.PendingInspection)-1
		// end:=0
		
		
type ReadClaim struct {
	ClaimId         string  `json:"claimId"`
	ClaimType       string  `json:"claimType"`
	ClientId        string  `json:"clientId"`
	IntimationDate  string  `json:"intimationDate"`
	LossDate        string  `json:"lossDate"`
	PolicyNumber    string  `json:"policyNumber"`
	InsuredName     string  `json:"insuredName"`
	InsuredPhone    string  `json:"insuredPhone"`
	InsuredAddress  string  `json:"insuredAddress"`
	InsuredEmail    string  `json:"insuredEmail"`     
	LossDescription string  `json:"lossDescription"`
	ClaimAmount     float64 `json:"claimAmount"`
	ApprovedAmount  float64 `json:"approvedAmount"`
	Status          string  `json:"status"`
	Surveyor        string  `json:"surveyor"`
	Report          string  `json:"report"`
	Policy			Policy 	`json:"policy"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}
readclaim:=ReadClaim{}
var readclaimsArr []ReadClaim
policy:= Policy{}
		if args[0] == READ_CLAIM_PENDING {
			start:=len(surveyor.PendingInspection)-1
		end:=0                                                                                                                 
			for i:=start ; i >=end ;i-- {
				claimsAsBytes ,err := stub.GetState(surveyor.PendingInspection[i])
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor:couldnt getstate of "+surveyor.PendingInspection[i]))
				}
				err = json.Unmarshal(claimsAsBytes,&readclaim)

				policyAsBytes,err := stub.GetState(readclaim.PolicyNumber)
				err = json.Unmarshal(policyAsBytes,&policy)
				readclaim.Policy = policy

				readclaimsArr = append(readclaimsArr,readclaim)

			}	
		}else if args[0] == READ_CLAIM_COMPLETED {
			start:=len(surveyor.CompletedInspection)-1
		end:=0
			for i:=start ; i >=end ;i-- {
				claimsAsBytes ,err := stub.GetState(surveyor.CompletedInspection[i])
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode:ReadAllClaimSurveyor:couldnt getstate of "+surveyor.CompletedInspection[i]))
				}
				err = json.Unmarshal(claimsAsBytes,&readclaim)	

				policyAsBytes,err := stub.GetState(readclaim.PolicyNumber)
				err = json.Unmarshal(policyAsBytes,&policy)
				readclaim.Policy = policy

				readclaimsArr = append(readclaimsArr,readclaim)

			}
		}else {
			return shim.Error("chaincode:ReadAllClaimSurveyor:Keyword not valid")
		}
		
		claimsArrAsBytes,err := json.Marshal(readclaimsArr)

		return shim.Success(claimsArrAsBytes)	
	}