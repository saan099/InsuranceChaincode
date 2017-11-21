package main

import (
	"bytes"
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

	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readAllPolicy::Error Unmarshalling %s", err.Error()))
	}

	policyArr := invokerClient.Policies

	/*if len(proposalArr) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
	}*/

	var buffer bytes.Buffer
	buffer.WriteString("[")
	flag := false
	//proposalobj:=Proposal{}
	for i := len(policyArr) - 1; i >= 0; i-- {
		if flag == true {
			buffer.WriteString(",")
		}
		policyAsBytes, err := stub.GetState(policyArr[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllPolicy::Could not get state of %s", policyArr[i]))
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

	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::Error Unmarshalling %s", err.Error()))
	}

	policyArr := invokerClient.Policies

	/*if len(proposalArr) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
	}*/
	start := len(policyArr) - 1
	end := 0
	lowerLimit, err := strconv.Atoi(args[0])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:lowerlimit not integer"))
		}
		upperLimit, err := strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:upperlimit not integer"))
		}
		if upperLimit < lowerLimit {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:upperlimit is not bigger than lowerlimit"))
		}
		if upperLimit >= len(policyArr) {
			end = 0

		} else {
			end = len(policyArr) - 1 - upperLimit
		}
		if lowerLimit < 0 {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange:lowerlimit is less than 0"))
		}
		start = len(policyArr) - 1 - lowerLimit

	var buffer bytes.Buffer
	buffer.WriteString("[")
	flag := false
	//proposalobj:=Proposal{}
	for i := start; i >= end; i-- {
		if flag == true {
			buffer.WriteString(",")
		}
		policyAsBytes, err := stub.GetState(policyArr[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readPolicyByRange::Could not get state of %s", policyArr[i]))
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

	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Error Unmarshalling "))
	}

	var policyArr []string

	policyArr = invokerClient.Policies

	/*if len(proposalArr) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
	}*/

	//var buffer bytes.Buffer
	//buffer.WriteString("[")
	flag := false
	//proposalobj:=Proposal{}
	for i := 0; i < len(policyArr); i++ {
		if policyArr[i] == args[0] {
			flag = true
			break
		}
	}

	if flag == false {
		return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy:: Policy not found in account"))
	}
	type readRFQ struct {
		RFQId              string              `json:"rfqId"`
		ClientId           string              `json:"clientId"`
		InsuredName        string              `json:"insuredName"`
		TypeOfInsurance    string              `json:"typeOfInsurance"`
		RiskAmount         float64             `json:"riskAmount"`
		RiskLocation       string              `json:"riskLocation"`
		Premium            float64             `json:"premium"`
		StartDate          string              `json:"startDate"`
		EndDate            string              `json:"endDate"`
		Status             string              `json:"status"`
		Quotes             []Quote             `json:"quotes"`
		LeadQuote          Quote               `json:"leadQuote"`
		SelectedInsurer    []Reads             `json:"selectedInsurer"`
		LeadInsurer        Reads               `json:"leadInsurer"`
		ProposalDocHash    string              `json:"proposalDocHash"`
		ProposalNum        string              `json:"proposalNum"`
		Intermediary       string              `json:"intermediary"`
		
		TransactionHistory []TransactionRecord `json:"transactionHistory"`
	}

	type readPolicy struct {
		PolicyNumber       string              `json:"policyNumber"`
		ProposalNum        string              `json:"proposalNum"`
		PolicyDocHash      string              `json:"policyDocHash"`
		Details            readRFQ             `json:"details"`
		Status             string              `json:"status"`
		Claim			   Claim			   `json:"claim"`
		TransactionHistory []TransactionRecord `json:"transactionHistory"`
	}

	//var arr []byte
	//var quoteArr []Quote
	policy := Policy{}
	policyAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Could not get state of %s", args[0]))
	}
	
	err = json.Unmarshal(policyAsBytes, &policy)
	claimAsBytes,err:= stub.GetState(policy.Claim)
	claim:=Claim{}
	err = json.Unmarshal(claimAsBytes,&claim)
	readpolicy := readPolicy{}
	readrfq := readRFQ{}
	rfq := policy.Details
	readpolicy.PolicyDocHash = policy.PolicyDocHash
	readpolicy.PolicyNumber = policy.PolicyNumber
	readpolicy.ProposalNum = policy.ProposalNum
	readpolicy.Status = policy.Status
	readpolicy.TransactionHistory = policy.TransactionHistory
	readpolicy.Claim = claim

	readrfq.ClientId = rfq.ClientId
	readrfq.InsuredName = rfq.InsuredName
	readrfq.EndDate = rfq.EndDate
	readrfq.Intermediary = rfq.Intermediary
	readrfq.ProposalDocHash = rfq.ProposalDocHash
	readrfq.ProposalNum = rfq.ProposalNum
	readrfq.RFQId = rfq.RFQId
	readrfq.RiskAmount = rfq.RiskAmount
	readrfq.RiskLocation = rfq.RiskLocation
	readrfq.Premium = rfq.Premium
	readrfq.StartDate = rfq.StartDate
	readrfq.Status = rfq.Status
	readrfq.TransactionHistory = rfq.TransactionHistory
	readrfq.TypeOfInsurance = rfq.TypeOfInsurance

	var quote Quote
	for i := range rfq.Quotes {
		quoteAsBytes, err := stub.GetState(rfq.Quotes[i])
		if err != nil || len(quoteAsBytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Could not get state of %s", rfq.Quotes[i]))
		}
		err = json.Unmarshal(quoteAsBytes, &quote)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Could not marshal readpolicy "))
		}
		if quote.QuoteId == rfq.LeadQuote {
			readrfq.LeadQuote = quote
		}
		readrfq.Quotes = append(readrfq.Quotes, quote)
	}
	readpolicy.Details = readrfq
	policyAsBytes, err = json.Marshal(readpolicy)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readSinglePolicy::Could not marshal readpolicy "))
	}
	//buffer.WriteString(string(quoteAsBytes))

	//flag = true

	//buffer.WriteString("]")

	return shim.Success(policyAsBytes)
}
