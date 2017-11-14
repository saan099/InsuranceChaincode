
package main
/*
import (
	"crypto/sha256"
	//"reflect"
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

//============================================SelectLeadInsurer=====================================================================

func (t *InsuranceManagement) SelectLeadInsurerByClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//args[0]=RFQId
	//args[1]=QuoteId
	RFQId := args[0]
	QuoteId := args[1]
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

	RFQArrayLength := len(client.RFQArray)
	flag := 1
	for i := 0; i < RFQArrayLength; i++ {
		if client.RFQArray[i] == RFQId {
			flag = 0
			break
		}
	}
	if flag == 1 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::invalid RFQ ID"))
	}

	rfqAsBytes, err := stub.GetState(RFQId)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt get RFQ from the state "))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal rfq "))
	}
	flag = 1
	var leadInsurerId string
	quoteArrayLength := len(rfq.Quotes)
	for i := 0; i < quoteArrayLength; i++ {
		if rfq.Quotes[i] == QuoteId {

			flag = 0
			break
		}
	}

	if flag == 1 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::invalid Quote Id"))
	}

	if len(rfq.LeadInsurer) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::lead insurer already selected for this RFQ"))
	}

	quote := Quote{}
	quoteAsBytes, err := stub.GetState(QuoteId)
	if err != nil || len(quoteAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::quote not found"))
	}
	err = json.Unmarshal(quoteAsBytes, &quote)
	if err != nil || len(quoteAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::error unmarshalling quote"))
	}
	rfq.LeadQuote = quote.QuoteId
	rfq.LeadInsurer = quote.InsurerId
	rfq.Status = LEAD_ASSIGNED

	finalRFQAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt marshal finalrfqasbytes "))
	}

	err = stub.PutState(RFQId, finalRFQAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt put RFQ to the state "))
	}
	quote.Status = QUOTE_ACCEPTED
	newQuoteAsbytes, err := json.Marshal(quote)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt marshal quote"))
	}
	err = stub.PutState(quote.QuoteId, newQuoteAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt write state"))
	}

	return shim.Success(nil)

}

//============================================SelectLeadInsurer=====================================================================

func (t *InsuranceManagement) SelectLeadInsurerByBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//args[0]=RFQId
	//args[1]=QuoteId
	RFQId := args[0]
	QuoteId := args[1]
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
	if err != nil || brokerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::account doesnt exists"))

	}
	broker := Broker{}

	err = json.Unmarshal(brokerAsBytes, &broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal client "))
	}

	RFQArrayLength := len(broker.RFQArray)
	flag := 1
	for i := 0; i < RFQArrayLength; i++ {
		if broker.RFQArray[i] == RFQId {
			flag = 0
			break
		}
	}
	if flag == 1 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::invalid RFQ ID"))
	}

	rfqAsBytes, err := stub.GetState(RFQId)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt get RFQ from the state "))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal rfq "))
	}
	flag = 1
	var leadInsurerId string
	quoteArrayLength := len(rfq.Quotes)
	for i := 0; i < quoteArrayLength; i++ {
		if rfq.Quotes[i] == QuoteId {

			flag = 0
			break
		}
	}

	if flag == 1 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::invalid Quote Id"))
	}

	if len(rfq.LeadInsurer) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::lead insurer already selected for this RFQ"))
	}

	quote := Quote{}
	quoteAsBytes, err := stub.GetState(QuoteId)
	if err != nil || len(quoteAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::quote not found"))
	}
	err = json.Unmarshal(quoteAsBytes, &quote)
	if err != nil || len(quoteAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::error unmarshalling quote"))
	}
	rfq.LeadQuote = quote.QuoteId
	rfq.LeadInsurer = quote.InsurerId
	rfq.Status = LEAD_ASSIGNED

	finalRFQAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt marshal finalrfqasbytes "))
	}

	err = stub.PutState(RFQId, finalRFQAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt put RFQ to the state "))
	}

	return shim.Success(nil)

}
*/