/*  completed till  coinsurer agree to the lead insurer quote,
next step is to provide client to set the capacity of different
 coinsurer and then policy and all  and one more thing ONLY BROKER HAS THE FUCNTIONALITY TO PROVIDE RFQ    */

package main

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

type InsuranceManagement struct {
}

//=============================main=====================================================

func main() {
	err := shim.Start(new(InsuranceManagement))
	if err != nil {
		fmt.Println("error starting chaincode :%s", err)
	}
}

//======================================init========================================

func (t *InsuranceManagement) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init called")
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:Init::Wrong number of arguments"))
	}
	return shim.Success(nil)
}

//==============================invoke====================================================

func (t *InsuranceManagement) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "initClient" {

		return t.InitClient(stub, args) //done

	} else if function == "initInsurer" {

		return t.InitInsurer(stub, args) //done

	} else if function == "generateRFQByClient" {

		return t.GenerateRFQByClient(stub, args) //done

	} else if function == "provideQuote" {

		return t.ProvideQuote(stub, args) //done

	} else if function == "initBroker" {

		return t.InitBroker(stub, args) //done

	} else if function == "generateRFQByBroker" {

		return t.GenerateRFQByBroker(stub, args) //done

	} else if function == "initClientByBroker" {

		return t.InitClientByBroker(stub, args) //done

	} else if function == "selectLeadInsurer" {

		return t.SelectLeadInsurer(stub, args) //done
	} else if function == "readAcc" {
		return t.ReadAcc(stub, args)
	} else if function == "readAllRFQ" {
		return t.ReadAllRFQ(stub, args)
	}

	return shim.Error(fmt.Sprintf("chaincode:Invoke::NO such function exists"))

}

func (t *InsuranceManagement) ReadAcc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:readAcc:0 arguments expected"))
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readAcc:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:readAcc:couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	invokerAddress := hex.EncodeToString(invokerhash[:])

	invokerAsBytes, err := stub.GetState(invokerAddress)
	if err != nil || len(invokerAsBytes) == 0 {
		shim.Error(fmt.Sprintf("chaincode:readAcc::account doesnt exists"))

	}
	return shim.Success(invokerAsBytes)

}

///func (t *InsuranceManagement) ReadRFQListForinvokerByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

//}

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
	insurer := Insurer{}

	err = json.Unmarshal(insurerAsBytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt unmarshal insurer "))
	}

	rfqAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::generateRfQ:RFQ doesnt exists"))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote::couldnt unmarshal rfq "))
	}

	quote := Quote{}
	quoteHash := sha256.Sum256([]byte(cert.Subject.CommonName + args[0]))
	quoteAddress := hex.EncodeToString(quoteHash[:])

	quote.QuoteId = quoteAddress
	quote.InsurerName = insurer.InsurerName
	quote.InsurerId = insurerAddress
	quote.Premium = args[1]
	quote.Capacity = args[2]
	quote.RFQId = rfq.RFQId

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

	err = stub.PutState(args[0], finalRFQAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:ProvideQuote:couldnt put state rfq "))
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

//============================================SelectLeadInsurer=====================================================================

func (t *InsuranceManagement) SelectLeadInsurer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//args[0]=RFQId
	//args[1]=QuoteId

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
		if client.RFQArray[i] == args[0] {
			flag = 0
			break
		}
	}
	if flag == 1 {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::invalid RFQ ID"))
	}

	rfqAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt get RFQ from the state "))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt unmarshal rfq "))
	}

	quoteArrayLength := len(rfq.Quotes)
	for i := 0; i < quoteArrayLength; i++ {
		if rfq.Quotes[i] == args[1] {
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

	rfq.LeadInsurer = args[1]

	finalRFQAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt marshal finalrfqasbytes "))
	}

	err = stub.PutState(args[0], finalRFQAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:SelectLeadInsurer::couldnt put RFQ to the state "))
	}

	return shim.Success(nil)

}

//============================================AcceptLeadQuote=====================================================================

func (t *InsuranceManagement) AcceptLeadQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args[0]=rfqId
	//args[1]=quoteId

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

	rfqAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AcceptLeadQuote::RFQ doesnt exists"))
	}

	rfq := RFQ{}

	err = json.Unmarshal(rfqAsBytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("Chaincode:AcceptLeadQuote:couldnt unmarshal rfq "))
	}

	rfq.FinalInsurer = append(rfq.FinalInsurer, insurerAddress)

	finalRfqAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::AcceptLeadQuote:couldnt marshal rfq"))
	}

	err = stub.PutState(args[0], finalRfqAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode::AcceptLeadQuote:couldnt putstate rfq"))
	}

	return shim.Success(nil)

}
