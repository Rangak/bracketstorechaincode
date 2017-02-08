/*
Copyright (c) 2016 Skuchain,Inc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// This chaincode implements the ledger operations for the proofchaincode

// ProofChainCode example simple Chaincode implementation
type bracketStore struct {
}

type bracketSummary struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

func (t *bracketStore) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

//ProofChainCode.Invoke runs a transaction against the current state
func (t *bracketStore) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//Proofs Chaincode should have one transaction argument. This is body of serialized protobuf
	if len(args) == 0 {
		fmt.Println("Zero arguments found")
		return nil, errors.New("Zero arguments found")
	}
	jsonBracket := args[0]
	brack := bracketSummary{}

	err := json.Unmarshal([]byte(jsonBracket), brack)
	if err != nil {
		fmt.Printf("Error getting summary JSON from %s", jsonBracket)
		return nil, errors.New("Error getting summary JSON")
	}

	switch function {

	case "UPSERT":
		fmt.Println("Inserting at UUID: %s", brack.UUID)
		stub.PutState(brack.UUID, []byte(jsonBracket))
		err = stub.SetEvent("evtsender", []byte(brack.UUID))
		if err != nil {
			return nil, err
		}
		return nil, nil

	default:
		fmt.Println("Invalid function type")
		return nil, errors.New("Invalid function type")
	}
}

// Query callback representing the query of a chaincode
func (t *bracketStore) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Printf("function: %s", function)
	switch function {
	case "QUERY":
		if len(args) != 1 {
			return nil, fmt.Errorf("No argument specified")
		}
		UUID := args[0]
		fmt.Println("Querying bracket json for UUID: %s", args[0])
		bracketBytes, err := stub.GetState(UUID)
		if err != nil {
			return nil, err
		}
		fmt.Println("Returning : %s", bracketBytes)
		return bracketBytes, nil

	default:
		return nil, errors.New("Unsupported operation")
	}
}

func main() {
	err := shim.Start(new(bracketStore))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err)
	}
}
