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

func (t *InsuranceManagement) MarkPaymentAndGeneratePolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::Wrong number of arguments"))
	}

	proposalNumber := args[0]
	policyNumber := args[1]
	policyDocHash := args[2]

	policybytes,err:=stub.GetState(policyNumber)
	if len(policybytes) != 0 {
		return shim.Error("chaincode:MarkPaymentAndGeneratePolicy::Policy already exists")
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	insurerAddress := hex.EncodeToString(invokerhash[:])

	insurerAsbytes, err := stub.GetState(insurerAddress)
	if err != nil || insurerAsbytes == nil {
		shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::account doesnt exists"))

	}
	insurer := Insurer{}

	err = json.Unmarshal(insurerAsbytes, &insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt unmarshal client "))
	}
	var found bool = false
	for i := range insurer.ProposalArray {
		if insurer.ProposalArray[i] == proposalNumber {
			found = true
			break
		}
	}

	if found == false {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt find proposal in insurer stack"))
	}

	proposal := Proposal{}
	proposalAsbytes, err := stub.GetState(proposalNumber)
	if err != nil || len(proposalAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt read proposal or proposal doesnt exist"))
	}
	err = json.Unmarshal(proposalAsbytes, &proposal)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::proposal couldnt unmarshal"))
	}
	if proposal.Status != PROPOSAL_INITIALIZED {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::proposal already resolved or rejected"))
	}
	proposal.PolicyNum = policyNumber
	proposal.Status = PROPOSAL_PAYMENT_MARKED
	newProposalAsbytes, err := json.Marshal(proposal)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::proposal couldnt marshal"))
	}
	err = stub.PutState(proposalNumber, newProposalAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::proposal couldnt put state"))
	}
	rfq := RFQ{}
	rfqId := proposal.RFQId
	rfqAsbytes, err := stub.GetState(rfqId)
	if err != nil || len(rfqAsbytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt read rfq or rfq doesnt exist"))
	}
	err = json.Unmarshal(rfqAsbytes, &rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt unmarshal rfq"))
	}
	transactionRecord := TransactionRecord{}
	transactionRecord.TxId = stub.GetTxID()
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt get timestamp for transaction"))
	}
	transactionRecord.Timestamp = timestamp.String()
	transactionRecord.Message = "Lead made policy of address- " + policyNumber

	clientOrBrokerAddress := rfq.ClientId

	policy := Policy{}
	policy.PolicyNumber = policyNumber
	policy.ProposalNum = proposalNumber
	policy.Details = rfq
	policy.PolicyDocHash = policyDocHash
	policy.Status = POLICY_INITIALIZED
	policy.TransactionHistory = rfq.TransactionHistory
	policy.TransactionHistory = append(policy.TransactionHistory, transactionRecord)

	policyAsbytes, err := json.Marshal(policy)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt marshal policy"))
	}

	err = stub.PutState(policyNumber, policyAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt put state of policy"))
	}

	insurer.Policies = append(insurer.Policies, policyNumber)
	newInsurerasbytes, err := json.Marshal(insurer)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt marshal insurer"))
	}
	err = stub.PutState(insurerAddress, newInsurerasbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt put state insurer"))
	}

	if rfq.Intermediary == INTERMEDIARY_CLIENT {
		client := Client{}
		clientAsbytes, err := stub.GetState(clientOrBrokerAddress)
		if err != nil || len(clientAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt read client account or account doesnt exist"))
		}
		err = json.Unmarshal(clientAsbytes, &client)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt unmarshal client"))
		}
		client.Policies = append(client.Policies, policyNumber)
		newClientAsbytes, err := json.Marshal(client)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt marshal client"))
		}
		err = stub.PutState(clientOrBrokerAddress, newClientAsbytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt put client state"))
		}
	} else {
		broker := Broker{}
		brokerAsbytes, err := stub.GetState(clientOrBrokerAddress)
		if err != nil || len(brokerAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt read broker address or broker doesnt exist"))
		}
		err = json.Unmarshal(brokerAsbytes, &broker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt unmarshal broker"))
		}
		broker.Policies = append(broker.Policies, policyNumber)
		newBrokerAsbytes, err := json.Marshal(broker)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt marshal broker"))
		}
		err = stub.PutState(clientOrBrokerAddress, newBrokerAsbytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:MarkPaymentAndGeneratePolicy::couldnt put state for broker"))
		}

	}

	return shim.Success(nil)
}
