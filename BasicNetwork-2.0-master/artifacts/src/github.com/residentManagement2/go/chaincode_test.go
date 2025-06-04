package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking TransactionContext and ChaincodeStub
type MockTransactionContext struct {
	mock.Mock
	contractapi.TransactionContextInterface
}

type MockChaincodeStub struct {
	mock.Mock
	contractapi.ChaincodeStubInterface
}

func (m *MockChaincodeStub) GetState(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockChaincodeStub) PutState(key string, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func TestRegisterResident(t *testing.T) {
	// Create mock context and stub
	ctx := new(MockTransactionContext)
	stub := new(MockChaincodeStub)
	ctx.On("GetStub").Return(stub)

	// Create an instance of the chaincode
	chaincode := new(ResidentManagement)

	// Define test data
	residentID := "res1"
	residentName := "John Doe"
	resident := Resident{
		ID:       residentID,
		Name:     residentName,
		Visitors: []Visitor{},
		Workers:  []Worker{},
	}

	// Convert to JSON
	residentBytes, _ := json.Marshal(resident)

	// Set expected behavior for PutState
	stub.On("PutState", residentID, residentBytes).Return(nil)

	// Call function and check for errors
	err := chaincode.RegisterResident(ctx, residentID, residentName)
	assert.NoError(t, err)

	// Verify that PutState was called with expected arguments
	stub.AssertCalled(t, "PutState", residentID, residentBytes)
}
