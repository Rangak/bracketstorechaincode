package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInvoke(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInvoke("1", "UPSERT", args)
	if err != nil {
		fmt.Println("UPSERT", args, "failed", err)
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes, err := stub.MockQuery("QUERY", []string{name})
	if err != nil {
		fmt.Println("Query", name, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func TestBracketStore(t *testing.T) {
	bst := new(bracketStore)
	stub := shim.NewMockStub("bst", bst)
	checkInvoke(t, stub, []string{`{"uuid":"1234","title":"test"}`})
	checkQuery(t, stub, "1234", `{"uuid":"1234","title":"test"}`)
}
