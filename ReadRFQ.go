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

func (t *InsuranceManagement) ReadAllRFQ(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 && len(args) != 2 {
		return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:0 arguments expected"))
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
		return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::account doesnt exists"))

	}
	invokerClient := Client{}
	//invokerInsurer:= Insurer{}
	//invokerBroker:= Broker{}
	//flag_type:= 0
	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::Unmarshal error:: %s", err.Error()))
	}
	var rfqIds []string
	//var rfqArr []RFQ

	rfqIds = invokerClient.RFQArray
	var rRFQList []readRFQ
	start := len(rfqIds) - 1
	end := 0
	if len(args) == 2 {
		lowerLimit, err := strconv.Atoi(args[0])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:lowerlimit not integer"))
		}
		upperLimit, err := strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:upperlimit not integer"))
		}
		if upperLimit < lowerLimit {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:upperlimit is not bigger than lowerlimit"))
		}
		if upperLimit > len(rfqIds) {
			end = 0

		} else {
			end = len(rfqIds) - 1 - upperLimit
		}
		if lowerLimit < 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:lowerlimit is less than 0"))
		}
		start = len(rfqIds) - 1 - lowerLimit
	}

	for i := start; i >= end; i-- {
		rRFQ := readRFQ{}
		rfq := RFQ{}
		rfqAsbytes, err := stub.GetState(rfqIds[i])
		if err != nil || len(rfqAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt get rfq state "))
		}

		err = json.Unmarshal(rfqAsbytes, &rfq)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt unmarshall rfq"))
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
				return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt get insurer"))
			}
			err = json.Unmarshal(insurerAsbytes, &insurer)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt unmarshall insurer"))
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
				return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt read quote"))
			}
			err = json.Unmarshal(quoteAsbytes, &quote)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt unmarshall quote"))
			}
			rRFQ.Quotes = append(rRFQ.Quotes, quote)
			if rfq.Quotes[j] == rfq.LeadQuote {
				rRFQ.LeadQuote = quote
			}
		}
		rRFQList = append(rRFQList, rRFQ)
	}
	rRFQAsbytes, err := json.Marshal(rRFQList)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::couldnt marshal rRFQ list"))
	}

	return shim.Success(rRFQAsbytes)
}

func (t *InsuranceManagement) ReadSingleRFQ(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args[0]=RFQId

	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ:1 argument expected"))
	}
	var rfqId = args[0]

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
	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	invokerAddress := hex.EncodeToString(invokerhash[:])

	invokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || len(invokerAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::account doesnt exists"))

	}
	invokerClient := Client{}
	//invokerInsurer:= Insurer{}
	//invokerBroker:= Broker{}
	//flag_type:= 0
	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::Unmarshal error:: %s", err.Error()))
	}

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
	rRFQ.Premium = rfq.Premium
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

	rRFQAsbytes, err := json.Marshal(rRFQ)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ReadSingleRFQ::couldnt marshal rRFQ list"))
	}

	return shim.Success(rRFQAsbytes)

}

//===================================Read All RFQ ================================================================
// func (t *InsuranceManagement) ReadAllRFQ(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//
// 		if len(args) != 0 {
// 			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:0 arguments expected"))
// 		}
//
// 		creator, err := stub.GetCreator() // it'll give the certificate of the invoker
// 		id := &mspprotos.SerializedIdentity{}
// 		err = proto.Unmarshal(creator, id)
//
// 		if err != nil {
// 			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:couldnt unmarshal creator"))
// 		}
// 		block, _ := pem.Decode(id.GetIdBytes())
// 		cert, err := x509.ParseCertificate(block.Bytes)
// 		if err != nil {
// 			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ:couldnt parse certificate"))
// 		}
// 		invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
// 		invokerAddress := hex.EncodeToString(invokerhash[:])
//
// 		invokerAsBytes, err := stub.GetState(invokerAddress)
// 		if err != nil || len(invokerAsBytes) == 0 {
// 			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::account doesnt exists"))
//
// 		}
// 		invokerClient := Client{}
// 		//invokerInsurer:= Insurer{}
// 		//invokerBroker:= Broker{}
// 		//flag_type:= 0
// 		err = json.Unmarshal(invokerAsBytes,&invokerClient)
// 		if err!=nil {
// 			return shim.Error(fmt.Sprintf("chaincode:readAllRFQ::Unmarshal error:: %s",err.Error()))
// 		}
// 		var rfqId []string
// 		//var rfqArr []RFQ
//
// 		rfqId = invokerClient.RFQArray
// 		if len(rfqId) == 0 {
//
// 			return shim.Error("chaincode:readAllRFQ::No RFQs to read")
// 		}
//
//
// 		var bigbuffer bytes.Buffer
// 		bigbuffer.WriteString("[")
// 		flag:=false
// 		rfqObj:= RFQ{}
// 		readsObj:= Reads{}
// 		insurerObj:=Insurer{}
// 		for i:=0 ; i < len(rfqId); i++ {
// 			var buffer bytes.Buffer
// 			if flag == true {
// 				bigbuffer.WriteString(",")
// 			}
// 			rfqAsBytes,_ := stub.GetState(rfqId[i])
// 			err= json.Unmarshal(rfqAsBytes,&rfqObj)
// 			buffer.WriteString(string(rfqAsBytes))
//
// 			for j:=0 ; j < len(rfqObj.SelectedInsurer); j++ {
// 				insurerAsBytes,_:= stub.GetState(rfqObj.SelectedInsurer[j])
// 				err = json.Unmarshal(insurerAsBytes,&insurerObj)
// 				readsObj.Id = insurerObj.InsurerId
// 				readsObj.Name = insurerObj.InsurerName
// 				readsAsBytes,_ := json.Marshal(readsObj)
// 				str:= "\""+rfqObj.SelectedInsurer[j] + "\""
// 				var a string
// 				if rfqObj.LeadInsurer == rfqObj.SelectedInsurer[j]{
// 				a=string(bytes.Replace(buffer.Bytes(),[]byte(str),[]byte(readsAsBytes),2))
// 				}else{
// 					a=string(bytes.Replace(buffer.Bytes(),[]byte(str),[]byte(readsAsBytes),1))
// 				}
// 				buffer.Reset()
// 				buffer.WriteString(a)
// 			}
// 			flag = true
// 			bigbuffer.WriteString(buffer.String())
// 		}
// 		bigbuffer.WriteString("]")
//
// 		return shim.Success(bigbuffer.Bytes())
//
// 	}

//=================================== ReadSingleRFQ =================================================================
// func (t *InsuranceManagement) ReadSingleRFQ(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	//args[0]=RFQId
//
// 	if len(args) != 1 {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ:1 argument expected"))
// 	}
//
// 	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
// 	id := &mspprotos.SerializedIdentity{}
// 	err = proto.Unmarshal(creator, id)
//
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ:couldnt unmarshal creator"))
// 	}
// 	block, _ := pem.Decode(id.GetIdBytes())
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ:couldnt parse certificate"))
// 	}
// 	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
// 	invokerAddress := hex.EncodeToString(invokerhash[:])
//
// 	invokerAsBytes, err := stub.GetState(invokerAddress)
// 	if err != nil || len(invokerAsBytes) == 0 {
// 		return shim.Error(fmt.Sprintf("chaincode:readSingleRFQ::account doesnt exists"))
//
// 	}
// 	invokerClient := Client{}
// 	//invokerInsurer:= Insurer{}
// 	//invokerBroker:= Broker{}
// 	//flag_type:= 0
// 	err = json.Unmarshal(invokerAsBytes, &invokerClient)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::Error Unmarshalling %s", err.Error()))
// 	}
// 	var rfqId []string
//
// 	rfqId = invokerClient.RFQArray
// 	if len(rfqId) == 0 {
// 		return shim.Error("chaincode:readSingleRFQ::No RFQs to read in this account")
// 	}
//
// 	flag := 0
//
// 	for i := 0; i < len(rfqId); i++ {
// 		if rfqId[i] == args[0] {
// 			flag = 1
// 		}
// 	}
//
// 	if flag == 0 {
// 		return shim.Error("chaincode:readSingleRFQ:: RFQId Not found in account")
// 	}
// 	var buffer bytes.Buffer
// 	//buffer.WriteString("[")
// 	rfqAsBytes, err := stub.GetState(args[0])
// 	if err != nil {
// 		return shim.Error("chaincode:readSingleRFQ:: RFQId Not found in WorldState")
// 	}
// 	rfqObj := RFQ{}
// 	readsObj := Reads{}
// 	insurerObj := Insurer{}
// 	buffer.WriteString(string(rfqAsBytes))
// 	err = json.Unmarshal(rfqAsBytes, &rfqObj)
// 	for j := 0; j < len(rfqObj.SelectedInsurer); j++ {
// 		insurerAsBytes, _ := stub.GetState(rfqObj.SelectedInsurer[j])
// 		err = json.Unmarshal(insurerAsBytes, &insurerObj)
// 		readsObj.Id = insurerObj.InsurerId
// 		readsObj.Name = insurerObj.InsurerName
// 		readsAsBytes, _ := json.Marshal(readsObj)
// 		str := "\"" + rfqObj.SelectedInsurer[j] + "\""
// 		var a string
// 		if rfqObj.LeadInsurer == rfqObj.SelectedInsurer[j] {
// 			a = string(bytes.Replace(buffer.Bytes(), []byte(str), []byte(readsAsBytes), 2))
// 		} else {
// 			a = string(bytes.Replace(buffer.Bytes(), []byte(str), []byte(readsAsBytes), 1))
// 		}
// 		buffer.Reset()
// 		buffer.WriteString(a)
// 	}
// 	for j := 0; j < len(rfqObj.Quotes); j++ {
// 		quotesAsBytes, _ := stub.GetState(rfqObj.Quotes[j])
// 		//err = json.Unmarshal(insurerAsBytes,&insurerObj)
//
// 		//readsAsBytes,_ := json.Marshal(readsObj)
// 		str := "\"" + rfqObj.Quotes[j] + "\""
// 		var a string
// 		//if rfqObj.LeadInsurer == rfqObj.Quotes[j]{
// 		a = string(bytes.Replace(buffer.Bytes(), []byte(str), []byte(quotesAsBytes), 1))
// 		//}else{
// 		//	a=string(bytes.Replace(buffer.Bytes(),[]byte(str),[]byte(quotesAsBytes),1))
// 		//}
// 		buffer.Reset()
// 		buffer.WriteString(a)
// 	}
// 	//flag = true
// 	//bigbuffer.WriteString(buffer.String())
//
// 	return shim.Success(buffer.Bytes())
// }

//================================= ReadRFQByRange ===================================================================
func (t *InsuranceManagement) ReadRFQByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args[0]=start
	//args[1]=end

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("chaincode:readRFQByRange:2 arguments expected"))
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readRFQByRange:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readRFQByRange:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	invokerAddress := hex.EncodeToString(invokerhash[:])

	invokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || len(invokerAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::account doesnt exists"))

	}
	invokerClient := Client{}
	//invokerInsurer:= Insurer{}
	//invokerBroker:= Broker{}
	//flag_type:= 0
	err = json.Unmarshal(invokerAsBytes, &invokerClient)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readRFQByRange::Error Unmarshalling %s", err.Error()))
	}
	var rfqId []string
	//var rfqArr []RFQ

	rfqId = invokerClient.RFQArray

	if len(rfqId) == 0 {
		return shim.Error("chaincode:readRFQByRange::No RFQs to read in this account")
	}

	start, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("cannot convert start string to int")
	}
	end, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("cannot convert start string to int")
	}

	if end >= len(rfqId) {
		//return shim.Error("chaincode:readRFQByRange::End range exceeded")
		end = len(rfqId) - 1
	}

	if start > len(rfqId) {
		start = 0
		end = 0
	}
	var bigbuffer bytes.Buffer
	bigbuffer.WriteString("[")
	flag := false
	rfqObj := RFQ{}
	readsObj := Reads{}
	insurerObj := Insurer{}

	for i := end; i >= start; i-- {
		var buffer bytes.Buffer
		if flag == true {
			bigbuffer.WriteString(",")
		}
		rfqAsBytes, _ := stub.GetState(rfqId[i])
		err = json.Unmarshal(rfqAsBytes, &rfqObj)
		buffer.WriteString(string(rfqAsBytes))

		for j := 0; j < len(rfqObj.SelectedInsurer); j++ {
			insurerAsBytes, _ := stub.GetState(rfqObj.SelectedInsurer[j])
			err = json.Unmarshal(insurerAsBytes, &insurerObj)
			readsObj.Id = insurerObj.InsurerId
			readsObj.Name = insurerObj.InsurerName
			readsAsBytes, _ := json.Marshal(readsObj)
			str := "\"" + rfqObj.SelectedInsurer[j] + "\""
			var a string
			if rfqObj.LeadInsurer == rfqObj.SelectedInsurer[j] {
				a = string(bytes.Replace(buffer.Bytes(), []byte(str), []byte(readsAsBytes), 2))
			} else {
				a = string(bytes.Replace(buffer.Bytes(), []byte(str), []byte(readsAsBytes), 1))
			}
			buffer.Reset()
			buffer.WriteString(a)
		}
		flag = true
		bigbuffer.WriteString(buffer.String())
	}
	bigbuffer.WriteString("]")

	return shim.Success(bigbuffer.Bytes())
}

//================================= ReadClientOfBroker ==================================================================
func (t *InsuranceManagement) ReadClientOfBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker:0 argument expected"))
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	invokerAddress := hex.EncodeToString(invokerhash[:])

	invokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || len(invokerAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker::account doesnt exists"))

	}
	//client := Client{}
	invokerBroker := Broker{}
	//flag_type:= 0

	err = json.Unmarshal(invokerAsBytes, &invokerBroker)

	if len(invokerBroker.Clients) == 0 {
		return shim.Error("chaincode:readClientOfBroker::No Clients to read in this account")
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	flag := false
	for i := 0; i < len(invokerBroker.Clients); i++ {
		if flag == true {
			buffer.WriteString(",")
		}
		clientAsBytes, err := stub.GetState(invokerBroker.Clients[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:readClientOfBroker::Could not get state of %s", invokerBroker.Clients[i]))
		}
		buffer.WriteString(string(clientAsBytes))
		flag = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}
