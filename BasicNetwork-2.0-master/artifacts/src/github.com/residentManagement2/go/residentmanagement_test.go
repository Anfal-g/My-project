package main_test

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"github.com/residentManagement2/go"
	"time"
    "math/rand"

    "github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)


// ==================== Mock Implementations ====================

type MockTransactionContext struct {
	contractapi.TransactionContextInterface
	stub *MockChaincodeStub
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

type MockChaincodeStub struct {
	shim.ChaincodeStubInterface
	state       map[string][]byte
	privateData map[string]map[string][]byte
}

func (m *MockChaincodeStub) GetState(key string) ([]byte, error) {
	return m.state[key], nil
}

func (m *MockChaincodeStub) PutState(key string, value []byte) error {
	if m.state == nil {
		m.state = make(map[string][]byte)
	}
	m.state[key] = value
	return nil
}

func (m *MockChaincodeStub) GetPrivateData(collection, key string) ([]byte, error) {
	if col, ok := m.privateData[collection]; ok {
		return col[key], nil
	}
	return nil, nil
}

func (m *MockChaincodeStub) PutPrivateData(collection, key string, value []byte) error {
	if m.privateData == nil {
		m.privateData = make(map[string]map[string][]byte)
	}
	if _, ok := m.privateData[collection]; !ok {
		m.privateData[collection] = make(map[string][]byte)
	}
	m.privateData[collection][key] = value
	return nil
}

// ==================== Test Cases ====================

func TestInitLedger(t *testing.T) {
	chaincode := &res.ResidentManagement{}
	stub := &MockChaincodeStub{}
	ctx := &MockTransactionContext{stub: stub}

	result, err := chaincode.InitLedger(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Success", result)
}

func TestRegisterResident(t *testing.T) {
	chaincode := &ResidentManagement{}
	stub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	ctx := &MockTransactionContext{stub: stub}

	resident, err := chaincode.RegisterResident(ctx, "res1", "user1", "male", "single", "owner", "apt101")
	assert.NoError(t, err)
	assert.Equal(t, "res1", resident.ResidentID)

	data, err := stub.GetPrivateData("ResidentPrivateCollection", "res1")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var stored Resident
	err = json.Unmarshal(data, &stored)
	assert.NoError(t, err)
	assert.Equal(t, "res1", stored.ResidentID)

	// Duplicate should fail
	_, err = chaincode.RegisterResident(ctx, "res1", "user1", "male", "single", "owner", "apt101")
	assert.Error(t, err)
}


// TestVisitorEntry tests visitor registration
func TestVisitorEntry(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful visitor entry
	visitor, err := chaincode.VisitorEntry(mockCtx, "vis1", "res1", "approved")
	assert.NoError(t, err)
	assert.Equal(t, "Approved", visitor.EntryStatus)

	// Test invalid status
	_, err = chaincode.VisitorEntry(mockCtx, "vis2", "res1", "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid approval status")
}

// TestMaintenanceEntry tests maintenance worker registration
func TestMaintenanceEntry(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful entry
	worker, err := chaincode.MaintenanceEntry(mockCtx, "worker1", "res1", "approved")
	assert.NoError(t, err)
	assert.Equal(t, "Maintenance/Service", worker.Role)

	// Test invalid status
	_, err = chaincode.MaintenanceEntry(mockCtx, "worker2", "res1", "invalid")
	assert.Error(t, err)
}

// TestUpdateVisitorApproval tests visitor status updates
func TestUpdateVisitorApproval(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: map[string]map[string][]byte{
			"VisitorApprovalCollection": {
				"vis1": []byte(`{"visitorId":"vis1","residentId":"res1","entryStatus":"Pending"}`),
			},
		},
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful update
	updated, err := chaincode.UpdateVisitorApproval(mockCtx, "vis1", "approved")
	assert.NoError(t, err)
	assert.Equal(t, "Approved", updated.EntryStatus)

	// Test invalid status
	_, err = chaincode.UpdateVisitorApproval(mockCtx, "vis1", "invalid")
	assert.Error(t, err)

	// Test non-existent visitor
	_, err = chaincode.UpdateVisitorApproval(mockCtx, "vis2", "approved")
	assert.Error(t, err)
}

// TestRegisterApartment tests apartment registration
func TestRegisterApartment(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful registration
	apt, err := chaincode.RegisterApartment(mockCtx, "apt101", "Apartment 101", "2BHK")
	assert.NoError(t, err)
	assert.Equal(t, "apt101", apt.ApartmentID)

	// Test duplicate registration
	_, err = chaincode.RegisterApartment(mockCtx, "apt101", "Apartment 101", "2BHK")
	assert.Error(t, err)
}

// TestSetResidentPolicy tests policy setting
func TestSetResidentPolicy(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful policy setting
	err := chaincode.SetResidentPolicy(mockCtx, "res1", 5, true, "admin1")
	assert.NoError(t, err)

	// Verify the policy was stored
	policyData, err := mockStub.GetPrivateData("ResidentPolicyCollection", "res1")
	assert.NoError(t, err)
	
	var policy ResidentPolicy
	err = json.Unmarshal(policyData, &policy)
	assert.NoError(t, err)
	assert.Equal(t, 5, policy.MaxVisitors)
}

// TestQueryFunctions tests query operations
func TestQueryFunctions(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: map[string]map[string][]byte{
			"ResidentPrivateCollection": {
				"res1": []byte(`{"residentId":"res1","userId":"user1"}`),
			},
		},
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test QueryEntry
	result, err := chaincode.QueryEntry(mockCtx, "ResidentPrivateCollection", "res1")
	assert.NoError(t, err)
	assert.Contains(t, result, "res1")

	// Test QueryEntry not found
	_, err = chaincode.QueryEntry(mockCtx, "ResidentPrivateCollection", "res2")
	assert.Error(t, err)
}

// TestDeleteEntry tests deletion functionality
func TestDeleteEntry(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: map[string]map[string][]byte{
			"ResidentPrivateCollection": {
				"res1": []byte(`{"residentId":"res1"}`),
			},
		},
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful deletion
	result, err := chaincode.DeleteEntry(mockCtx, "ResidentPrivateCollection", "res1")
	assert.NoError(t, err)
	assert.Contains(t, result, "deleted successfully")

	// Test deletion of non-existent entry
	_, err = chaincode.DeleteEntry(mockCtx, "ResidentPrivateCollection", "res2")
	assert.Error(t, err)
}

// TestCheckResidentStatus tests resident status verification
func TestCheckResidentStatus(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: map[string]map[string][]byte{
			"ResidentPrivateCollection": {
				"res1": []byte(`{"residentId":"res1"}`),
			},
		},
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test existing resident
	exists, err := chaincode.CheckResidentStatus(mockCtx, "res1")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test non-existent resident
	exists, err = chaincode.CheckResidentStatus(mockCtx, "res2")
	assert.NoError(t, err)
	assert.False(t, exists)
}

// TestGenerateSecureQR tests QR code generation
func TestGenerateSecureQR(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test QR generation
	qr, err := chaincode.GenerateSecureQR(mockCtx)
	assert.NoError(t, err)
	assert.Contains(t, qr, "SECURE-")
	assert.Len(t, qr, 21) // "SECURE-" + 16 hex chars
}

// TestDeliveryWorkerEntry tests delivery worker registration
func TestDeliveryWorkerEntry(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful entry
	worker, err := chaincode.DeliveryWorkerEntry(mockCtx, "del1", "res1", "approved")
	assert.NoError(t, err)
	assert.Equal(t, "Delivery Worker", worker.Role)

	// Test invalid status
	_, err = chaincode.DeliveryWorkerEntry(mockCtx, "del2", "res1", "invalid")
	assert.Error(t, err)
}

// TestRegisterDeliveryWorker tests temporary delivery worker registration
func TestRegisterDeliveryWorker(t *testing.T) {
	chaincode := &ResidentManagement{}
	mockStub := &MockChaincodeStub{
		privateData: make(map[string]map[string][]byte),
	}
	mockCtx := &MockTransactionContext{chaincodeStub: mockStub}

	// Test successful registration
	worker, err := chaincode.RegisterDeliveryWorker(mockCtx, "del1", 2) // 2 hours
	assert.NoError(t, err)
	assert.Equal(t, "Temporary Access", worker.EntryStatus)
	assert.Greater(t, worker.Expiry, time.Now().UnixNano()/int64(time.Millisecond))
}