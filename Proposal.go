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

func (t *InsuranceManagement) UploadProposalFormByClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::wrong number of arguments"))
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	clientAddress := hex.EncodeToString(invokerhash[:])

	clientAsBytes, err := stub.GetState(clientAddress)
	if err != nil || clientAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::account doesnt exists"))

	}
	client := Client{}

	err = json.Unmarshal(clientAsBytes, &client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::couldnt unmarshal client "))
	}
	rfqId := args[0]
	proposalFormHash := args[1]
	var found bool = false
	for i := range client.RFQArray {
		if client.RFQArray[i] == rfqId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq not found in client's stack"))
	}
	rfq := RFQ{}
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq not read or doesnt exit"))
	}
	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq couldnt get unmarshalled"))
	}
	rfq.Status = RFQ_PROPOSAL_FINALIZED
	rfq.ProposalDocHash = proposalFormHash

	newRfqAsbytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq couldnt get marshalled"))
	}
	err = stub.PutState(rfqId, newRfqAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq didnt put its state"))
	}

	return shim.Success(nil)
}

func (t *InsuranceManagement) UploadProposalFormByBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::wrong number of arguments"))
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	brokerAddress := hex.EncodeToString(invokerhash[:])

	brokerAsBytes, err := stub.GetState(brokerAddress)
	if err != nil || brokerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::account doesnt exists"))

	}
	broker := Client{}

	err = json.Unmarshal(brokerAsBytes, &broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::couldnt unmarshal client "))
	}
	rfqId := args[0]
	proposalFormHash := args[1]
	var found bool = false
	for i := range broker.RFQArray {
		if broker.RFQArray[i] == rfqId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq not found in client's stack"))
	}
	rfq := RFQ{}
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq not read or doesnt exit"))
	}
	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq couldnt get unmarshalled"))
	}
	rfq.Status = RFQ_PROPOSAL_FINALIZED
	rfq.ProposalDocHash = proposalFormHash

	newRfqAsbytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq couldnt get marshalled"))
	}
	err = stub.PutState(rfqId, newRfqAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:UploadProposalFormByClient::rfq didnt put its state"))
	}

	return shim.Success(nil)
}
