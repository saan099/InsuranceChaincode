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

	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readAllProposal::Error Unmarshalling %s", err.Error()))
	}

	proposalArr := invokerClient.ProposalArray

	/*if len(proposalArr) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readAllProposal:: No proposals found in acount"))
	}*/

	var buffer bytes.Buffer
	buffer.WriteString("[")
	flag := false
	proposalobj := Proposal{}
	for i := 0; i < len(proposalArr); i++ {
		if flag == true {
			buffer.WriteString(",")
		}
		proposalAsBytes, err := stub.GetState(proposalArr[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllProposal::Could not get state of %s", proposalArr[i]))
		}
		buffer.WriteString(string(proposalAsBytes))
		err = json.Unmarshal(proposalAsBytes, &proposalobj)
		RFQAsBytes, _ := stub.GetState(proposalobj.RFQId)
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString(",")
		buffer.WriteString("\"rfqDetails\": ")
		buffer.WriteString(string(RFQAsBytes))
		buffer.WriteString("}")
		flag = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (t *InsuranceManagement) ReadSingleProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:1 argument expected"))
	}
	proposalNum := args[0]

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

	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::Error Unmarshalling %s", err.Error()))
	}

	proposalArr := invokerClient.ProposalArray

	/*if len(proposalArr) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:: No proposals found in acount"))
	}*/

	flag := 0

	for i := 0; i < len(proposalArr); i++ {
		if proposalArr[i] == args[0] {
			flag = 1
			break
		}
	}
	if flag == 0 {
		return shim.Error("chaincode:readSingleProposal:: Proposal Not found in account")
	}

	type readRFQ struct {
		RFQId              string              `json:"rfqId"`
		ClientId           string              `json:"clientId"`
		InsuredName        string              `json:"insuredName"`
		TypeOfInsurance    string              `json:"typeOfInsurance"`
		RiskAmount         float64             `json:"riskAmount"`
		RiskLocation       string              `json:"riskLocation"`
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

	type readProposal struct {
		ProposalNum        string              `json:"proposalNum"`
		RFQId              string              `json:"rfqId"`
		Status             string              `json:"status"`
		PolicyNum          string              `json:"policyNum"`
		RFQDetails         readRFQ             `json:"RFQDetails"`
		TransactionHistory []TransactionRecord `json:"transactionHistory"`
	}
	proposal := Proposal{}
	proposalAsbytes, err := stub.GetState(proposalNum)
	if err != nil || len(proposalAsbytes) == 0 {
		return shim.Error("chaincode:readSingleProposal:: Proposal Not found in account")
	}
	err = json.Unmarshal(proposalAsbytes, &proposal)
	if err != nil {
		return shim.Error("chaincode:readSingleProposal:: Proposal Not unmarshalled")
	}

	var rfqId = proposal.RFQId

	rRFQ := readRFQ{}
	rfq := RFQ{}
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt get rfq state "))
	}

	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt unmarshall rfq"))
	}
	rRFQ.ClientId = rfq.ClientId
	rRFQ.InsuredName = rfq.InsuredName
	rRFQ.EndDate = rfq.EndDate
	rRFQ.Intermediary = rfq.Intermediary
	rRFQ.ProposalDocHash = rfq.ProposalDocHash
	rRFQ.ProposalNum = rfq.ProposalNum
	rRFQ.RFQId = rfq.RFQId
	rRFQ.RiskAmount = rfq.RiskAmount
	rRFQ.RiskLocation = rfq.RiskLocation
	rRFQ.StartDate = rfq.StartDate
	rRFQ.Status = rfq.Status
	rRFQ.TransactionHistory = rfq.TransactionHistory
	rRFQ.TypeOfInsurance = rfq.TypeOfInsurance

	for j := range rfq.SelectedInsurer {
		insurer := Insurer{}
		insurerAsbytes, err := stub.GetState(rfq.SelectedInsurer[j])
		if err != nil || len(insurerAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt get insurer"))
		}
		err = json.Unmarshal(insurerAsbytes, &insurer)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt unmarshall insurer"))
		}
		readInsurer := Reads{}
		readInsurer.Id = insurer.InsurerId
		readInsurer.Name = insurer.InsurerName
		rRFQ.SelectedInsurer = append(rRFQ.SelectedInsurer, readInsurer)
		if insurer.InsurerId == rfq.LeadInsurer {
			rRFQ.LeadInsurer = readInsurer
		}
	}

	for j := range rfq.Quotes {
		quote := Quote{}
		quoteAsbytes, err := stub.GetState(rfq.Quotes[j])
		if err != nil || len(quoteAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt read quote"))
		}
		err = json.Unmarshal(quoteAsbytes, &quote)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt unmarshall quote"))
		}
		rRFQ.Quotes = append(rRFQ.Quotes, quote)
		if rfq.Quotes[j] == rfq.LeadQuote {
			rRFQ.LeadQuote = quote
		}
	}

	rProposal := readProposal{}
	rProposal.PolicyNum = proposal.PolicyNum
	rProposal.ProposalNum = proposal.ProposalNum
	rProposal.RFQDetails = rRFQ
	rProposal.RFQId = proposal.RFQId
	rProposal.Status = proposal.Status
	rProposal.TransactionHistory = proposal.TransactionHistory

	readProposalAsBytes, err := json.Marshal(rProposal)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt marshal read proposal"))
	}

	return shim.Success(readProposalAsBytes)
}

//=================================== Read Single Proposal ================================================================

// func (t *InsuranceManagement) ReadSingleProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//
// 	if len(args) != 1 {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:1 argument expected"))
// 	}
//
// 	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
// 	id := &mspprotos.SerializedIdentity{}
// 	err = proto.Unmarshal(creator, id)
//
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:couldnt unmarshal creator"))
// 	}
// 	block, _ := pem.Decode(id.GetIdBytes())
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:couldnt parse certificate"))
// 	}
// 	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
// 	invokerAddress := hex.EncodeToString(invokerhash[:])
//
// 	invokerAsBytes, err := stub.GetState(invokerAddress)
// 	if err != nil || len(invokerAsBytes) == 0 {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::account doesnt exists"))
//
// 	}
//
// 	invokerClient := Client{}
//
// 	err = json.Unmarshal(invokerAsBytes, &invokerClient)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::Error Unmarshalling %s", err.Error()))
// 	}
//
// 	proposalArr := invokerClient.ProposalArray
//
// 	/*if len(proposalArr) == 0 {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal:: No proposals found in acount"))
// 	}*/
//
// 	flag := 0
//
// 	for i := 0; i < len(proposalArr); i++ {
// 		if proposalArr[i] == args[0] {
// 			flag = 1
// 			break
// 		}
// 	}
// 	if flag == 0 {
// 		return shim.Error("chaincode:readSingleProposal:: Proposal Not found in account")
// 	}
// 	var buffer bytes.Buffer
// 	//buffer.WriteString("[")
// 	//flag=false
// 	proposalobj := Proposal{}
//
// 	proposalAsBytes, err := stub.GetState(args[0])
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleProposal::Could not get state of %s", args[0]))
// 	}
//
// 	buffer.WriteString(string(proposalAsBytes))
// 	err = json.Unmarshal(proposalAsBytes, &proposalobj)
// 	RFQAsBytes, _ := stub.GetState(proposalobj.RFQId)
// 	buffer.Truncate(buffer.Len() - 1)
// 	buffer.WriteString(",")
// 	buffer.WriteString("\"rfqDetails\": ")
// 	buffer.WriteString(string(RFQAsBytes))
// 	buffer.WriteString("}")
// 	//flag = true
//
// 	//buffer.WriteString("]")
//
// 	return shim.Success(buffer.Bytes())
// }

//=================================== Read Proposal By Range ================================================================

func (t *InsuranceManagement) ReadProposalByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange:2 arguments expected"))
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	invokerAddress := hex.EncodeToString(invokerhash[:])

	invokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || len(invokerAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange::account doesnt exists"))

	}

	invokerClient := Client{}

	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange::Error Unmarshalling %s", err.Error()))
	}

	proposalArr := invokerClient.ProposalArray

	/*if len(proposalArr) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange:: No proposals found in acount"))
	}*/

	//flag:=0

	start, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange::Could not convert %s to int", args[0]))
	}
	end, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readProposalByRange::Could not conver %s to int", args[1]))
	}

	if end > len(proposalArr) {
		//return shim.Error(fmt.Sprintf("chaincode:readProposalByRange:: End limit exceeded"))
		end = len(proposalArr)
	}

	if start > len(proposalArr) {
		start = 0
		end = 0
	}

	//if flag == 0 { return shim.Error("chaincode:readProposalByRange:: Propo Not found in account")}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	flag := false
	proposalobj := Proposal{}

	for i := start; i < end; i++ {
		if flag == true {
			buffer.WriteString(",")
		}
		proposalAsBytes, err := stub.GetState(proposalArr[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readProposalByRange::Could not get state of %s", proposalArr[i]))
		}

		buffer.WriteString(string(proposalAsBytes))
		err = json.Unmarshal(proposalAsBytes, &proposalobj)
		RFQAsBytes, _ := stub.GetState(proposalobj.RFQId)
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString(",")
		buffer.WriteString("\"rfqDetails\": ")
		buffer.WriteString(string(RFQAsBytes))
		buffer.WriteString("}")
		flag = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
