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

func (t *InsuranceManagement) AllotProposalNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::wrong number of arguments"))
	}
	rfqId := args[0]
	proposalNumber := args[1]

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(invokerhash[:])

	insurerAsbytes, err := stub.GetState(insurerAddress)
	if err != nil || insurerAsbytes == nil {
		shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::account doesnt exists"))

	}
	insurer := Insurer{}

	err = json.Unmarshal(insurerAsbytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt unmarshal client "))
	}
	var found bool = false
	for i := range insurer.RFQArray {
		if insurer.RFQArray[i] == rfqId {
			found = true
			break
		}
	}
	if found == false {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::RFQ not found in stack"))
	}

	rfq := RFQ{}
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt read rfq or rfq doesnt exist"))
	}
	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt unmarshal rfq"))
	}
	var clientOrBrokerAddress string = rfq.ClientId

	rfq.Status = RFQ_COMPLETED
	rfq.ProposalNum = proposalNumber
	newRFQAsbytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt marshal RFQ"))
	}
	err = stub.PutState(rfqId, newRFQAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt write state of RFQ"))
	}

	proposal := Proposal{}
	proposal.PolicyNum = proposalNumber
	proposal.RFQId = rfqId
	proposal.Status = PROPOSAL_INITIALIZED
	proposalAsBytes, err := json.Marshal(proposal)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt marshal proposal"))
	}
	err = stub.PutState(proposalNumber, proposalAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt write state of proposal"))
	}

	insurer.ProposalArray = append(insurer.ProposalArray, proposal.ProposalNum)
	newInsurerAsbytes, err := json.Marshal(insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt marshal insurer"))
	}
	err = stub.PutState(insurerAddress, newInsurerAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt put state for insurer"))
	}

	if rfq.Intermediary == INTERMEDIARY_CLIENT {
		client := Client{}
		clientAsbytes, err := stub.GetState(clientOrBrokerAddress)
		if err != nil || len(clientAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt read client account or account doesnt exist"))
		}
		client.ProposalArray = append(client.ProposalArray, proposalNumber)
		newClientAsbytes, err := json.Marshal(client)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt marshal client"))
		}
		err = stub.PutState(clientOrBrokerAddress, newClientAsbytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt put client state"))
		}
	} else {
		broker := Broker{}
		brokerAsbytes, err := stub.GetState(clientOrBrokerAddress)
		if err != nil || len(brokerAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt read broker address or broker doesnt exist"))
		}
		err = json.Unmarshal(brokerAsbytes, &broker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt unmarshal broker"))
		}
		broker.ProposalArray = append(broker.ProposalArray, proposalNumber)
		newBrokerAsbytes, err := json.Marshal(broker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt marshal broker"))
		}
		err = stub.PutState(clientOrBrokerAddress, newBrokerAsbytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt put state for broker"))
		}

	}

	return shim.Success(nil)
}

func (t *InsuranceManagement) MarkPaymentAndGeneratePolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(invokerhash[:])

	insurerAsbytes, err := stub.GetState(insurerAddress)
	if err != nil || insurerAsbytes == nil {
		shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::account doesnt exists"))

	}
	insurer := Insurer{}

	err = json.Unmarshal(insurerAsbytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:AllotProposalNumber::couldnt unmarshal client "))
	}

	return shim.Success(nil)
}
