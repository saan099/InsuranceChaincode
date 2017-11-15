package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//============================================AcceptLeadQuote=====================================================================

func (t *InsuranceManagement) AcceptLeadQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args[0]=rfqId
	//args[1]=quoteId
	RFQId := args[0]

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::AcceptLeadQuote:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AcceptLeadQuote:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(invokerhash[:])

	insurerAsBytes, err := stub.GetState(insurerAddress)
	if err != nil || insurerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:AcceptLeadQuote:account doesnt exists"))

	}
	insurer := Insurer{}

	err = json.Unmarshal(insurerAsBytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt unmarshal insurer "))
	}
	var found bool = false
	for i := range insurer.RFQArray {
		if insurer.RFQArray[i] == RFQId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:rfq doesnt exist"))
	}

	rfqAsBytes, err := stub.GetState(RFQId)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AcceptLeadQuote::RFQ doesnt exists"))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt unmarshal rfq "))
	}

	if rfq.Status != LEAD_ASSIGNED {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:Lead Insurer not resolved or Incorrect state"))
	}

	leadQuote := Quote{}
	leadQuoteAsbytes, err := stub.GetState(rfq.LeadQuote)
	if err != nil || len(leadQuoteAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:Lead Quote Not Found"))
	}
	err = json.Unmarshal(leadQuoteAsbytes, &leadQuote)
	if err != nil {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:Lead Quote wasnt Unmarshalled"))
	}

	for i := range rfq.Quotes {
		if rfq.Quotes[i] != leadQuote.QuoteId {

			quote := Quote{}
			quoteAsBytes, err := stub.GetState(rfq.Quotes[i])
			if err != nil || len(quoteAsBytes) == 0 {
				return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt read quote"))
			}
			err = json.Unmarshal(quoteAsBytes, &quote)
			if err != nil || len(quoteAsBytes) == 0 {
				return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt unmarshal a quote"))
			}
			if quote.InsurerId == insurerAddress {
				quote.Premium = leadQuote.Premium
				quote.Status = QUOTE_ACCEPTED
				newQuoteAsbytes, err := json.Marshal(quote)
				if err != nil {
					return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt marshal a quote"))
				}
				err = stub.PutState(quote.QuoteId, newQuoteAsbytes)
				if err != nil {
					return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt write state for a quote"))
				}
			}
		}
	}

	finalRfqAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::AcceptLeadQuote:couldnt marshal rfq"))
	}

	err = stub.PutState(RFQId, finalRfqAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::AcceptLeadQuote:couldnt putstate rfq"))
	}

	return shim.Success(nil)

}

//============================================RejectLeadQuote=====================================================================

func (t *InsuranceManagement) RejectLeadQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args[0]=rfqId
	//args[1]=quoteId
	RFQId := args[0]

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::RejectLeadQuote:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:RejectLeadQuote:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(invokerhash[:])

	insurerAsBytes, err := stub.GetState(insurerAddress)
	if err != nil || insurerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:RejectLeadQuote:account doesnt exists"))

	}
	insurer := Insurer{}

	err = json.Unmarshal(insurerAsBytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:couldnt unmarshal insurer "))
	}
	var found bool = false
	for i := range insurer.RFQArray {
		if insurer.RFQArray[i] == RFQId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:rfq doesnt exist"))
	}

	rfqAsBytes, err := stub.GetState(RFQId)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:RejectLeadQuote::RFQ doesnt exists"))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:couldnt unmarshal rfq "))
	}

	if rfq.Status != LEAD_ASSIGNED {
		return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:Lead Insurer not resolved or Incorrect state"))
	}

	var updatedQuotes []string

	for i := range rfq.Quotes {
		if rfq.Quotes[i] != rfq.LeadQuote {

			quote := Quote{}
			quoteAsBytes, err := stub.GetState(rfq.Quotes[i])
			if err != nil || len(quoteAsBytes) == 0 {
				return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:couldnt read quote"))
			}
			err = json.Unmarshal(quoteAsBytes, &quote)
			if err != nil || len(quoteAsBytes) == 0 {
				return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:couldnt unmarshal a quote"))
			}
			if quote.InsurerId == insurerAddress {

				quote.Status = QUOTE_REJECTED
				newQuoteAsbytes, err := json.Marshal(quote)
				if err != nil {
					return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:couldnt marshal a quote"))
				}
				err = stub.PutState(quote.QuoteId, newQuoteAsbytes)
				if err != nil {
					return shim.Error(fmt.Sprintf("Chaincode:RejectLeadQuote:couldnt write state for a quote"))
				}
			} else {
				updatedQuotes = append(updatedQuotes, rfq.Quotes[i])
			}
		}
	}

	rfq.Quotes = updatedQuotes

	finalRfqAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::RejectLeadQuote:couldnt marshal rfq"))
	}

	err = stub.PutState(RFQId, finalRfqAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::RejectLeadQuote:couldnt putstate rfq"))
	}

	return shim.Success(nil)

}
