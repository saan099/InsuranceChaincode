/*  completed till  coinsurer agree to the lead insurer quote,
next step is to provide client to set the capacity of different
 coinsurer and then policy and all  and one more thing ONLY BROKER HAS THE FUCNTIONALITY TO PROVIDE RFQ    */

package main

import (
	"crypto/sha256"
	//"reflect"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const LEAD_ASSIGNED string = "Lead Assigned"
const QUOTE_INITIALIZED string = "Quote Initialized"
const QUOTE_ACCEPTED string = "Quote Accepted"
const QUOTE_REJECTED string = "Quote Rejected"

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

	} else if function == "selectLeadInsurerByClient" {
		return t.SelectLeadInsurerByClient(stub, args) //done
	} else if function == "selectLeadInsurerByBroker" {
		return t.SelectLeadInsurerByBroker(stub, args) //done
	} else if function == "readAcc" {
		return t.ReadAcc(stub, args)
	} else if function == "readAllRFQ" {
		return t.ReadAllRFQ(stub, args)
	} else if function == "acceptLeadQuote" {
		return t.AcceptLeadQuote(stub, args)
	} else if function == "rejectLeadQuote" {
		return t.RejectLeadQuote(stub, args)
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
