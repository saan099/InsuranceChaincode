package main 

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	//"bytes"
	"strconv"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)


//=================================== Generate Claim ================================================================
func (t *InsuranceManagement) GenerateClaimByClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Intimation Date
		//args[1]= Loss Date
		//args[2]= Loss Description
		//args[3]= Policy Number
		//args[4]= Claim amount
		//args[5]= Insured Phone
		//args[6]= Insured Address
		//args[7]= Insured Email
		//args[8]= CLaim Type
		if len(args) != 9 {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim:9 arguments expected"))
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
		
		err = json.Unmarshal(invokerAsBytes,&client)
		for i:=0 ; i < len(client.Policies); i++ {
			if args[3] == client.Policies[i] {
				flag = true 
				break
			}
		}

		if flag == false {
			return shim.Error("chaincode:generateClaim::Policy number doesnt exist in account")
		}

		//err = json.Unmarshal(invokerAsBytes,&client)
		policyAsBytes,err:=stub.GetState(args[3])
		policy:=Policy{}
		err=json.Unmarshal(policyAsBytes,&policy)
		if len(policy.Claim) != 0 {
			return shim.Error("chaincode:generateClaim::Claim already initiated for this policy")
		}
		

		claim:=Claim{}
		// add claim details
		claim.IntimationDate = args[0]
		claim.LossDate = args[1]
		claim.LossDescription = args[2]
		claim.PolicyNumber = args[3]
		claim.ClaimAmount ,err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim::Claim Amount not float"))
		}
		claim.InsuredPhone = args[5]
		claim.InsuredAddress = args[6]
		claim.InsuredEmail = args[7]
		claim.ClaimType = args[8]
		claim.InsuredName= client.ClientName
		claim.ClientId = client.ClientId
		claim.Status = CLAIM_INITIALIZED
		claim.ClaimId = stub.GetTxID()

		//assign claim to policy
		policy.Claim = claim.ClaimId
		

		client.Claims = append(client.Claims,claim.ClaimId)

		
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

		//update Lead insurer 
		insurerAsBytes,err:=stub.GetState(policy.Details.LeadInsurer)
		err = json.Unmarshal(insurerAsBytes,&insurer)
		insurer.Claims = append(insurer.Claims, claim.ClaimId)
		insurerAsBytes,_ = json.Marshal(insurer)
		err = stub.PutState(insurer.InsurerId,insurerAsBytes)

		transactionRecord := TransactionRecord{}
		transactionRecord.TxId = stub.GetTxID()
		timestamp, err := stub.GetTxTimestamp()
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:generateClaim::couldnt get timestamp for transaction"))
		}
		transactionRecord.Timestamp = timestamp.String()
		transactionRecord.Message = "Generated Claim of Id- " + claim.ClaimId + " by " + invokerAddress

		claim.TransactionHistory = append(claim.TransactionHistory, transactionRecord)
		
		//append tx history to tx stack of policy
		policy.TransactionHistory = append(policy.TransactionHistory, transactionRecord)
		policyAsBytes ,err = json.Marshal(policy)
		err = stub.PutState(policy.PolicyNumber,policyAsBytes)

		claimAsBytes,err:= json.Marshal(claim)//update claim
		err = stub.PutState(claim.ClaimId,claimAsBytes)
		
		clientAsBytes,err := json.Marshal(client)// update client
		err = stub.PutState(client.ClientId,clientAsBytes)

		return shim.Success(nil)
	}

//=================================== Generate Claim By Broker================================================================
func (t *InsuranceManagement) GenerateClaimByBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Intimation Date
		//args[1]= Loss Date
		//args[2]= Loss Description
		//args[3]= Policy Number
		//args[4]= Claim amount
		//args[5]= Insured Phone
		//args[6]= Insured Address
		//args[7]= Insured Email
		//args[8]= CLaim Type
		if len(args) != 9 {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker:9 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		

		brokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || brokerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker::account doesnt exists"))

	}
	broker := Broker{}

	err = json.Unmarshal(brokerAsBytes, &broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker:couldnt unmarshal broker "))
	}
		flag:=false
		
		//err = json.Unmarshal(invokerAsBytes,&client)
		for i:=0 ; i < len(broker.Policies); i++ {
			if args[3] == broker.Policies[i] {
				flag = true 
				break
			}
		}

		if flag == false {
			return shim.Error("chaincode:GenerateClaimByBroker::Policy number doesnt exist in account")
		}

		//err = json.Unmarshal(invokerAsBytes,&client)
		policyAsBytes,err:=stub.GetState(args[3])
		policy:=Policy{}
		err=json.Unmarshal(policyAsBytes,&policy)
		if len(policy.Claim) != 0 {
			return shim.Error("chaincode:GenerateClaimByBroker::Claim already initiated for this policy")
		}
		

		claim:=Claim{}
		// add claim details
		claim.IntimationDate = args[0]
		claim.LossDate = args[1]
		claim.LossDescription = args[2]
		claim.PolicyNumber = args[3]
		claim.ClaimAmount ,err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker::Claim Amount not float"))
		}
		claim.InsuredPhone = args[5]
		claim.InsuredAddress = args[6]
		claim.InsuredEmail = args[7]
		claim.ClaimType = args[8]
		claim.InsuredName= policy.Details.InsuredName
		claim.ClientId = invokerAddress
		claim.Status = CLAIM_INITIALIZED
		claim.ClaimId = stub.GetTxID()

		//assign claim to policy
		policy.Claim = claim.ClaimId
		

		broker.Claims = append(broker.Claims,claim.ClaimId)

		
		insurer:=Insurer{}

		//update every insurer with new generated claim
		for i:=0 ; i < len(policy.Details.SelectedInsurer) ; i++ {
			insurerAsBytes,err:= stub.GetState(policy.Details.SelectedInsurer[i])
			if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker:couldnt get Insurers"))
		}
			err = json.Unmarshal(insurerAsBytes,&insurer)
			if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker:couldnt UNmarshal insurer"))
		}
			insurer.Claims = append(insurer.Claims,claim.ClaimId)
			insurerAsBytes,err = json.Marshal(insurer)
				err=stub.PutState(insurer.InsurerId,insurerAsBytes)
		}

		//update Lead insurer 
		insurerAsBytes,err:=stub.GetState(policy.Details.LeadInsurer)
		err = json.Unmarshal(insurerAsBytes,&insurer)
		insurer.Claims = append(insurer.Claims, claim.ClaimId)
		insurerAsBytes,_ = json.Marshal(insurer)
		err = stub.PutState(insurer.InsurerId,insurerAsBytes)

		transactionRecord := TransactionRecord{}
		transactionRecord.TxId = stub.GetTxID()
		timestamp, err := stub.GetTxTimestamp()
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateClaimByBroker::couldnt get timestamp for transaction"))
		}
		transactionRecord.Timestamp = timestamp.String()
		transactionRecord.Message = "Generated Claim of Id- " + claim.ClaimId + " by " + invokerAddress

		claim.TransactionHistory = append(claim.TransactionHistory, transactionRecord)
		
		//append tx history to tx stack of policy
		policy.TransactionHistory = append(policy.TransactionHistory, transactionRecord)
		policyAsBytes ,err = json.Marshal(policy)
		err = stub.PutState(policy.PolicyNumber,policyAsBytes)

		claimAsBytes,err:= json.Marshal(claim)//update claim
		err = stub.PutState(claim.ClaimId,claimAsBytes)
		
		brokerAsBytes,err = json.Marshal(broker)// update broker
		err = stub.PutState(broker.BrokerId,brokerAsBytes)

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
		err=json.Unmarshal(claimAsBytes,&claim)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim:couldnt unmarshal claim"))
		}
		//check if the invoker is Lead Insurer
		policyAsBytes,err:=stub.GetState(claim.PolicyNumber)
		policy:=Policy{}
		err = json.Unmarshal(policyAsBytes,&policy)
		if policy.Details.LeadInsurer != invokerAddress {
			return shim.Error("chaincode:asignSurveyorToClaim::Only Lead insurer is allowed to assign surveyor")
		}
		if claim.Status == CLAIM_SURVEYOR_ASSIGNED{
			return shim.Error("chaincode:asignSurveyorToClaim:: Surveyor already assigned for this claim")
		}
		transactionRecord := TransactionRecord{}
		transactionRecord.TxId = stub.GetTxID()
		timestamp, err := stub.GetTxTimestamp()
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim::couldnt get timestamp for transaction"))
		}
		transactionRecord.Timestamp = timestamp.String()
		transactionRecord.Message = "Surveyor assigned for claim - " + claim.ClaimId + " by " + invokerAddress

		claim.TransactionHistory = append(claim.TransactionHistory, transactionRecord)
		
		//append tx history to tx stack of policy
		policy.TransactionHistory = append(policy.TransactionHistory, transactionRecord)
		
		
		//check if Surveyor exists
		surveyorAsBytes,err:= stub.GetState(args[1])
		if err != nil || len(surveyorAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:asignSurveyorToClaim::Surveyor doesnt exist"))	
		}
		surveyor:= Surveyor{}
		err = json.Unmarshal(surveyorAsBytes,&surveyor)
		surveyor.PendingInspection = append(surveyor.PendingInspection,claim.ClaimId)

		//assign surveyor
		claim.Surveyor = args[1]
		claim.Status = CLAIM_SURVEYOR_ASSIGNED	
		claimAsBytes,err = json.Marshal(claim)	
		policyAsBytes,err = json.Marshal(policy)
		surveyorAsBytes,err = json.Marshal(surveyor)
		err = stub.PutState(policy.PolicyNumber,policyAsBytes)
		err = stub.PutState(claim.ClaimId,claimAsBytes) // update state of claim
		err = stub.PutState(surveyor.SurveyorId,surveyorAsBytes)

		return shim.Success(nil)	
	}

//=================================== Send Claim ================================================================
func (t *InsuranceManagement) SendClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Claim Id
		//args[1]= amount
			
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim:2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim::account doesnt exists"))	
		}
		claim:=Claim{}
		claimAsBytes,err := stub.GetState(args[0])
		if err != nil || len(claimAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim::Claim Id doesnt exist"))	
		}
		err=json.Unmarshal(claimAsBytes,&claim)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim:couldnt unmarshal claim"))
		}
		//check if the invoker is Lead Insurer
		policyAsBytes,err:=stub.GetState(claim.PolicyNumber)
		policy:=Policy{}
		err = json.Unmarshal(policyAsBytes,&policy)
		if policy.Details.LeadInsurer != invokerAddress {
			return shim.Error("chaincode:SendClaim::Only Lead insurer is allowed to send Claim")
		}
		transactionRecord := TransactionRecord{}
		transactionRecord.TxId = stub.GetTxID()
		timestamp, err := stub.GetTxTimestamp()
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim::couldnt get timestamp for transaction"))
		}
		transactionRecord.Timestamp = timestamp.String()
		transactionRecord.Message = "Claim amount - "+args[1] +" sent for claim - " + claim.ClaimId + " by " + invokerAddress

		claim.TransactionHistory = append(claim.TransactionHistory, transactionRecord)
		
		//append tx history to tx stack of policy
		policy.TransactionHistory = append(policy.TransactionHistory, transactionRecord)

		if claim.Status != CLAIM_INSPECTION_COMPLETED {
			return shim.Error("chaincode:SendClaim:Inspection report is not done yet")
		}
		if claim.Status == CLAIM_COMPLETED {
			return shim.Error("chaincode:SendClaim:Claim amount already Sent")
		}
		claim.ApprovedAmount,err = strconv.ParseFloat(args[1],64)
		claim.Status = 	CLAIM_COMPLETED	

		claimAsBytes ,err = json.Marshal(claim)
		policyAsBytes,err = json.Marshal(policy)
		err = stub.PutState(policy.PolicyNumber,policyAsBytes)

		err = stub.PutState(claim.ClaimId,claimAsBytes)

		return shim.Success(nil)
}

//=================================== Upload Claim Report ================================================================
	func (t *InsuranceManagement) UploadClaimReport(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		//args[0]= Claim Id
		//args[1]= Claim report hash
			
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport: 2 arguments expected"))
		}
	
		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
		id := &mspprotos.SerializedIdentity{}
		err = proto.Unmarshal(creator, id)
	
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt unmarshal creator"))
		}
		block, _ := pem.Decode(id.GetIdBytes())
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt parse certificate"))
		}
		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
		invokerAddress := hex.EncodeToString(invokerhash[:])
	
		invokerAsBytes, err := stub.GetState(invokerAddress)
		if err != nil || len(invokerAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport::account doesnt exists"))	
		}

		claimAsBytes,err:= stub.GetState(args[0])
		if err != nil || len(claimAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport::Claim Id doesnt exists"))	
		}
		claim:=Claim{}
		err = json.Unmarshal(claimAsBytes,&claim)

		surveyor := Surveyor{}
		err = json.Unmarshal(invokerAsBytes,&surveyor)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt unmarshal surveyor"))
		}
		
		flag:=false
		var i int
		for i=0; i<  len(surveyor.PendingInspection);i++ {
			if args[0] == surveyor.PendingInspection[i] {
				flag = true
				break
			}
		}
		if flag == false {
			return shim.Error("chaincode:UploadClaimReport:Claim Id not found in account")
		}

		if claim.Status == CLAIM_INSPECTION_COMPLETED {
			return shim.Error("chaincode:UploadClaimReport:Inspection report already uploaded for this claim")
		}

		//add report to completed
		claim.Report = args[1] 		
		claim.Status = CLAIM_INSPECTION_COMPLETED
		transactionRecord := TransactionRecord{}
		transactionRecord.TxId = stub.GetTxID()
		timestamp, err := stub.GetTxTimestamp()
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:SendClaim::couldnt get timestamp for transaction"))
		}
		transactionRecord.Timestamp = timestamp.String()
		transactionRecord.Message = "Claim Survey report uploaded for claim - " + claim.ClaimId + " by " + invokerAddress

		claim.TransactionHistory = append(claim.TransactionHistory, transactionRecord)
		
		policy:=Policy{}
		policyAsBytes,_:=stub.GetState(claim.PolicyNumber)
		err = json.Unmarshal(policyAsBytes,&policy)
		//append tx history to tx stack of policy
		policy.TransactionHistory = append(policy.TransactionHistory, transactionRecord)

		surveyor.CompletedInspection = append(surveyor.CompletedInspection,args[0])
		//remove claim from pending
		surveyor.PendingInspection = append(surveyor.PendingInspection[:i],surveyor.PendingInspection[i+1:]...)
		
		claimAsBytes, err = json.Marshal(claim)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt unmarsahl claim2"))
		}
		err= stub.PutState(claim.ClaimId,claimAsBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt putstate claim"))
		}
		surveyorAsBytes, err := json.Marshal(surveyor)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt unmarshal surveyor2"))
		}
		err = stub.PutState(surveyor.SurveyorId,surveyorAsBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:UploadClaimReport:couldnt putstate surveyor"))
		}
		policyAsBytes,err = json.Marshal(policy)
		err = stub.PutState(policy.PolicyNumber,policyAsBytes)

		return shim.Success(nil)
	}	