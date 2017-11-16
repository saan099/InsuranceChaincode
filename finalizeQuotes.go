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

func (t *InsuranceManagement) FinalizeQuotesByClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	rfqId := args[0]
	numOfQuotes, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Expected integer for quote number"))
	}
	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	clientAddress := hex.EncodeToString(invokerhash[:])

	clientAsBytes, err := stub.GetState(clientAddress)
	if err != nil || clientAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::account doesnt exists"))

	}
	client := Client{}

	err = json.Unmarshal(clientAsBytes, &client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal client "))
	}

	var found bool = false
	for i := range client.RFQArray {
		if client.RFQArray[i] == rfqId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:You dont have any such RFQ in stack"))
	}

	rfq := RFQ{}
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:couldnt read rfq or rfq doesnt exist"))
	}
	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:couldnt unmarshal rfq"))
	}
	initalQuoteList := rfq.Quotes
	var finalizedQuotes []string
	var insurerList []string
	var totalCapacity float64 = 0
	var count int = 0
	for i := 2; i < numOfQuotes*2+2; i++ {
		quoteId := args[i]
		i++
		capacity, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Expected float capacity"))
		}
		for j := range rfq.Quotes {
			if rfq.Quotes[j] == quoteId {
				quote := Quote{}
				quoteAsbytes, err := stub.GetState(quoteId)
				if err != nil || len(quoteAsbytes) == 0 {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Quote was not found"))
				}
				err = json.Unmarshal(quoteAsbytes, &quote)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Quote was not able to unmarshall"))
				}
				if capacity > quote.Capacity {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Exceeded maximum capacity set by supplier"))
				}
				count++
				if quote.Status != QUOTE_ACCEPTED {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Quote already rejected by insurer"))
				}
				quote.Capacity = capacity
				quote.Status = QUOTES_FINALIZED
				newQuoteAsbytes, err := json.Marshal(quote)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Quote couldnt marshal"))
				}
				err = stub.PutState(quoteId, newQuoteAsbytes)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A quote couldnt PutState"))
				}
				finalizedQuotes = append(finalizedQuotes, quoteId)
				insurerList = append(insurerList, quote.InsurerId)
				totalCapacity += capacity
				break
			}
		}
	}
	if count != numOfQuotes {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Some quote didnt get accepted as it is not in the stack. Maybe it was rejected"))
	}
	if totalCapacity > 100 {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Total capacity exceeding 100%"))
	}

	rfq.Quotes = finalizedQuotes
	rfq.SelectedInsurer = insurerList
	rfq.Status = RFQ_QUOTES_FINALIZED

	newRfqAsbytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:RFQ couldnt marshal"))
	}
	err = stub.PutState(rfqId, newRfqAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:couldnt putstate rfq"))
	}

	for i := range initalQuoteList {
		var quoteFound bool = false
		for j := range finalizedQuotes {
			if initalQuoteList[i] == finalizedQuotes[j] {
				quoteFound = true
				break
			}
		}
		if quoteFound == false {
			quote := Quote{}
			quoteAsbytes, err := stub.GetState(initalQuoteList[i])
			if err != nil || len(quoteAsbytes) == 0 {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not found"))
			}
			err = json.Unmarshal(quoteAsbytes, &quote)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not able to unmarshall"))
			}
			quote.Status = QUOTE_REJECTED
			newQuoteAsbytes, err := json.Marshal(quote)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not Marshalled"))
			}
			err = stub.PutState(initalQuoteList[i], newQuoteAsbytes)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not Put in state"))
			}
		}
	}

	return shim.Success(nil)
}

func (t *InsuranceManagement) FinalizeQuotesByBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	rfqId := args[0]
	numOfQuotes, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Expected integer for quote number"))
	}
	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	brokerAddress := hex.EncodeToString(invokerhash[:])

	brokerAsBytes, err := stub.GetState(brokerAddress)
	if err != nil || len(brokerAsBytes) == 0 {
		shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::account doesnt exists"))

	}
	broker := Broker{}

	err = json.Unmarshal(brokerAsBytes, &broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal client "))
	}

	if len(args) == 0 {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Wrong number of quotes"))
	}

	var found bool = false
	for i := range broker.RFQArray {
		if broker.RFQArray[i] == rfqId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:You dont have any such RFQ in stack"))
	}

	rfq := RFQ{}
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:couldnt read rfq or rfq doesnt exist"))
	}
	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:couldnt unmarshal rfq"))
	}
	initalQuoteList := rfq.Quotes
	var finalizedQuotes []string
	var insurerList []string
	var totalCapacity float64 = 0
	for i := 2; i < numOfQuotes*2+2; i++ {
		quoteId := args[i]
		i++
		capacity, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Expected float capacity"))
		}
		for j := range rfq.Quotes {
			if rfq.Quotes[j] == quoteId {
				quote := Quote{}
				quoteAsbytes, err := stub.GetState(quoteId)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Quote was not found"))
				}
				err = json.Unmarshal(quoteAsbytes, &quote)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Quote was not able to unmarshall"))
				}
				if capacity > quote.Capacity {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Exceeded maximum capacity set by supplier"))
				}
				quote.Capacity = capacity
				quote.Status = QUOTES_FINALIZED
				newQuoteAsbytes, err := json.Marshal(quote)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Quote couldnt marshal"))
				}
				err = stub.PutState(quoteId, newQuoteAsbytes)
				if err != nil {
					return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A quote couldnt PutState"))
				}
				finalizedQuotes = append(finalizedQuotes, quoteId)
				insurerList = append(insurerList, quote.InsurerId)
				totalCapacity += capacity
				break
			}
		}
	}
	if totalCapacity > 100 {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:Total capacity exceeding 100%"))
	}

	rfq.Quotes = finalizedQuotes
	rfq.SelectedInsurer = insurerList
	rfq.Status = RFQ_QUOTES_FINALIZED

	newRfqAsbytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:RFQ couldnt marshal"))
	}
	err = stub.PutState(rfqId, newRfqAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:couldnt putstate rfq"))
	}
	for i := range initalQuoteList {
		var quoteFound bool = false
		for j := range finalizedQuotes {
			if initalQuoteList[i] == finalizedQuotes[j] {
				quoteFound = true
				break
			}
		}
		if quoteFound == false {
			quote := Quote{}
			quoteAsbytes, err := stub.GetState(initalQuoteList[i])
			if err != nil || len(quoteAsbytes) == 0 {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not found"))
			}
			err = json.Unmarshal(quoteAsbytes, &quote)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not able to unmarshall"))
			}
			quote.Status = QUOTE_REJECTED
			newQuoteAsbytes, err := json.Marshal(quote)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not Marshalled"))
			}
			err = stub.PutState(initalQuoteList[i], newQuoteAsbytes)
			if err != nil {
				return shim.Error(fmt.Sprintf("chaincode::FinalizeQuote:A Rejected Quote was not Put in state"))
			}
		}
	}

	return shim.Success(nil)
}
