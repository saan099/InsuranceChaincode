package main

import (
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

//========================================================InitInsurer============

func (t *InsuranceManagement) InitInsurer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("init Insurer called")
	fmt.Println("=========================================")
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::couldn't get creator"))
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::error unmarshalling"))
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("chaincode:InitInsurer::couldn pasre ParseCertificate")
	}

	insurerHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(insurerHash[:])

	checkInsurerAsBytes, err := stub.GetState(insurerAddress)
	if err != nil || len(checkInsurerAsBytes) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::Insurer already exist"))
	}

	insurer := Insurer{}
	insurer.InsurerId = insurerAddress
	insurer.InsurerName = cert.Subject.CommonName

	insurerAsBytes, err := json.Marshal(insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::couldn't Unmarsh creator"))
	}

	err = stub.PutState(insurerAddress, insurerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitInsurer::couldn't write state "))
	}
	return shim.Success(nil)

}

//============================provideQuote==============================================

func (t *InsuranceManagement) ProvideQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//args[0]=RFQID
	//args[1]=Premium
	//args[2]=Capacity

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(invokerhash[:])

	insurerAsBytes, err := stub.GetState(insurerAddress)
	if err != nil || insurerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:ProvideQuote::account doesnt exists"))

	}
	rfqId := args[0]

	insurer := Insurer{}

	err = json.Unmarshal(insurerAsBytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt unmarshal insurer "))
	}
	var found bool = false
	for i := range insurer.RFQArray {
		if insurer.RFQArray[i] == rfqId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:You dont have any such RFQ."))
	}

	rfqAsBytes, err := stub.GetState(rfqId)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::generateRfQ:RFQ doesnt exists"))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::couldnt unmarshal rfq "))
	}
	// if rfq.Status!=RFQ_INITIATED {
	// 	return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::rfq not available anymore for quotes"))
	// }

	quote := Quote{}
	quoteHash := sha256.Sum256([]byte(cert.Subject.CommonName + rfqId))
	quoteAddress := hex.EncodeToString(quoteHash[:])

	quote.QuoteId = quoteAddress
	quote.InsurerName = insurer.InsurerName
	quote.InsurerId = insurerAddress
	quote.Premium, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::expected float premium "))
	}
	quote.Capacity, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::expected float capacity"))
	}
	quote.RFQId = rfq.RFQId
	quote.Status = QUOTE_INITIALIZED

	quoteAsBytes, err := json.Marshal(quote)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt marshal quote "))
	}

	err = stub.PutState(quoteAddress, quoteAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt put state quote "))
	}

	rfq.Quotes = append(rfq.Quotes, quoteAddress)

	finalRFQAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::couldnt marshal RFQ "))
	}

	err = stub.PutState(rfqId, finalRFQAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt put state rfq "))
	}
	var exists bool = false
	for i := range insurer.Quotes {
		if insurer.Quotes[i] == quoteAddress {
			exists = true
			break
		}
	}
	if exists == true {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:You've already provided quote for this RFQ "))
	}

	insurer.Quotes = append(insurer.Quotes, quote.QuoteId)
	finalInsurerAsBytes, err := json.Marshal(insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt marshal RFQ "))
	}
	err = stub.PutState(insurerAddress, finalInsurerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt put state client "))
	}

	return shim.Success(nil)

}
