package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//=======================================================GenerateRFQByClient================================
func (t *InsuranceManagement) GenerateRFQByClient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//args[0]=RFQID generated on the client side
	//args[1]=ClientId
	//args[2]=InsurerName
	//args[3]=TypeOFinsurance
	//args[4]=RiskAmount
	//args[5]=number of insurer
	//args[6].args[7]......insurer addresses
	rfqId := stub.GetTxID()
	InsurerClient := args[0]
	TypeOFinsurance := args[1]
	RiskAmount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ::Risk Amount not float"))
	}
	NumberOfInsurer, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ::number of insurer is not int"))
	}
	if NumberOfInsurer < 1 {
		return shim.Error("chaincode:GenerateRFQ::provide atleast one insurer")
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ::couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	clientAddress := hex.EncodeToString(invokerhash[:])

	clientAsBytes, err := stub.GetState(clientAddress)
	if err != nil || len(clientAsBytes) == 0 {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ::account doesnt exists"))

	}
	client := Client{}

	err = json.Unmarshal(clientAsBytes, &client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ:: couldnt unmarshal client"))
	}
	tym := time.Now()
	tym.Format("Mon Jan _2 15:04:05 2006")
	rfq := RFQ{}
	rfq.ClientId = clientAddress
	rfq.RFQId = rfqId
	rfq.RiskAmount = RiskAmount
	rfq.TypeOfInsurance = TypeOFinsurance
	rfq.InsuredName = InsurerClient
	rfq.Status = "RFQ fired on " + tym.String()

	//var insurerArray []string

	for i := 4; i < NumberOfInsurer+4; i++ {
		rfq.SelectedInsurer = append(rfq.SelectedInsurer, args[i])
		/*insurerAsBytes, err := stub.GetState(args[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("Chaincode:generateRFQ:can't get %dth insurer provided", i-3))
		}
		insurer := Insurer{}
		err = json.Unmarshal(insurerAsBytes, &insurer)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateRFQ::couldnt unmarshal insurer "))
		}
		insurer.RFQArray = append(insurer.RFQArray, rfqId)
		finalInsurerAsBytes, err := json.Marshal(insurer)
		if err != nil {
			return shim.Error(fmt.Sprintf("Chaincode:generateRFQ:can't marshal the finalInsurerAsBytes "))
		}
		err = stub.PutState(args[i], finalInsurerAsBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Chaincode:generateRFQ:couldnt putstate the finalInsurerAsBytes "))
		} */

	}

	rfqAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:generateRfQ couldnt marshal rfq"))
	}

	client.RFQArray = append(client.RFQArray, rfqId)

	finalClientAsBytes, err := json.Marshal(client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:generateRfQ couldnt marshal rfq"))
	}

	err = stub.PutState(rfqId, rfqAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:generateRfQ couldnt putstate rfq"))
	}

	err = stub.PutState(clientAddress, finalClientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:generateRfQ couldnt putstate client"))
	}

	return shim.Success(nil)

}

//=======================================================GenerateRFQByBroker================================

func (t *InsuranceManagement) GenerateRFQByBroker(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//args[0]=RFQID generated on the client side
	//args[1]=ClientId
	//args[2]=InsurerName
	//args[3]=TypeOFinsurance
	//args[4]=RiskAmount
	//args[5]=number of insurer
	//args[6].args[7]......insurer addresses
	rfqId := stub.GetTxID()
	clientId := args[0]
	InsurerClient := args[1]
	TypeOfInsurance := args[2]
	RiskAmount, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:risk amount is not float"))
	}
	NumberOfInsurer, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:number of insurer is not int"))
	}
	if NumberOfInsurer < 1 {
		return shim.Error("chaincode:GenerateRFQByBroker:provide atleast one insurer")
	}

	creator, err := stub.GetCreator() // it'll give the certificate of the invoker
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:couldnt unmarshal creator"))
	}
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker::couldnt parse certificate"))
	}
	invokerhash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	brokerAddress := hex.EncodeToString(invokerhash[:])

	brokerAsBytes, err := stub.GetState(brokerAddress)
	if err != nil || brokerAsBytes == nil {
		shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker::account doesnt exists"))

	}
	broker := Broker{}

	err = json.Unmarshal(brokerAsBytes, &broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:couldnt unmarshal broker "))
	}

	flag := 1
	brokerClientArray := broker.Clients
	lengthOfBrokerClientArray := len(brokerClientArray)
	for i := 0; i < lengthOfBrokerClientArray; i++ {
		if brokerClientArray[i] == clientId {
			flag = 0
			break
		}
	}

	if flag == 1 {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:couldnt find the client in your stack "))
	}

	clientAsBytes, err := stub.GetState(clientId)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker::couldnt get client provided by broker"))
	}

	client := Client{}

	err = json.Unmarshal(clientAsBytes, &client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker::couldnt unmarshal client "))
	}

	tym := time.Now()
	tym.Format("Mon Jan _2 15:04:05 2006")
	rfq := RFQ{}
	rfq.ClientId = brokerAddress
	rfq.RFQId = rfqId
	rfq.RiskAmount = RiskAmount
	rfq.TypeOfInsurance = TypeOfInsurance
	rfq.InsuredName = InsurerClient
	rfq.Status = "RFQ fired on " + tym.String()

	//var insurerArray []string

	for i := 5; i < NumberOfInsurer+5; i++ {
		rfq.SelectedInsurer = append(rfq.SelectedInsurer, args[i])
		insurer := Insurer{}
		insurerAsbytes, err := stub.GetState(args[i])
		if err != nil || len(insurerAsbytes) == 0 {
			return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt get on of the insurers"))
		}
		err = json.Unmarshal(insurerAsbytes, &insurer)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt unmarshal one of the insurers"))
		}
		insurer.RFQArray = append(insurer.RFQArray, rfq.RFQId)
		insurerAsNewbytes, err := json.Marshal(insurer)
		if err != nil {
			return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt marshal new bytes for insurers"))
		}
		stub.PutState(args[i], insurerAsNewbytes)
	}
	broker.RFQArray = append(broker.RFQArray, rfq.RFQId)

	rfqAsBytes, err := json.Marshal(rfq)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt marshal rfq"))
	}

	client.RFQArray = append(client.RFQArray, rfqId)

	finalClientAsBytes, err := json.Marshal(client)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt marshal rfq"))
	}

	err = stub.PutState(rfqId, rfqAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:couldnt putstate rfq"))
	}

	err = stub.PutState(clientId, finalClientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt putstate client"))
	}
	brokerAsbytes, err := json.Marshal(broker)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker: couldnt marshal broker"))
	}
	err = stub.PutState(brokerAddress, brokerAsbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("chaincode:GenerateRFQByBroker:couldnt putstate broker"))
	}

	return shim.Success(nil)

}
