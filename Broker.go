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

//=========================================================InitBroker=================================================

func (t *InsuranceManagement) InitBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("init broker called")
	fmt.Println("=======================================")
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitBroker::couldn't get creator"))
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitBroker::error unmarshalling"))
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("chaincode:InitBroker::couldn pasre ParseCertificate")
	}

	invokerHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	brokerAddress := hex.EncodeToString(invokerHash[:])

	checkBrokerAsBytes, err := stub.GetState(brokerAddress)
	if err != nil || len(checkBrokerAsBytes) != 0 {
		return shim.Error(fmt.Sprintf("chaincode:InitBroker::broker already exist"))
	}

	broker := Broker{}
	broker.BrokerId = brokerAddress
	broker.BrokerName = cert.Subject.CommonName

	brokerAsBytes, err := json.Marshal(broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitBroker::couldn't Unmarsh creator"))
	}

	err = stub.PutState(brokerAddress, brokerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitBroker::couldn't write state "))
	}
	return shim.Success(nil)

}

//=========================================================InitClientByBroker=================================================

func (t *InsuranceManagement) InitClientByBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("init InitClientByBroker called")
	fmt.Println("=======================================")
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::couldn't get creator"))
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::error unmarshalling"))
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error("chaincode:InitClientByBroker::couldn pasre ParseCertificate")
	}

	brokerHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	brokerAddress := hex.EncodeToString(brokerHash[:])

	checkBrokerAsBytes, err := stub.GetState(brokerAddress)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::broker didnt found "))
	}

	ClientName := args[0]

	broker := Broker{}

	err = json.Unmarshal(checkBrokerAsBytes, &broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::couldn't Unmarshal creator"))
	}

	clientHash := sha256.Sum256([]byte(ClientName + cert.Subject.CommonName + cert.Issuer.CommonName))
	clientAddress := hex.EncodeToString(clientHash[:])

	client := Client{}
	client.ClientId = clientAddress
	client.ClientName = ClientName

	checkClientExistsBytes, err := stub.GetState(clientAddress)
	if err != nil || len(checkClientExistsBytes) > 0 {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::This client already exists on system"))
	}

	clientAsBytes, err := json.Marshal(client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::couldn't Marshal client "))
	}

	err = stub.PutState(clientAddress, clientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::couldn't write client state "))
	}

	broker.Clients = append(broker.Clients, clientAddress)

	finalBrokerAsBytes, err := json.Marshal(broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::couldn't Marshal broker "))
	}

	err = stub.PutState(brokerAddress, finalBrokerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:InitClientByBroker::couldn't write state "))
	}
	return shim.Success(nil)

}
