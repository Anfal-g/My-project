package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	
)

type ResidentsContract struct{}



type BuildingConfig struct {
    ResidentsPerApartment int `json:"residentsPerApartment"`
}
type VisitorInfo struct {
    VisitorId      string `json:"visitorId"`
    FullName       string `json:"fullName"`
    Phone          string `json:"phone"`
    VisitTimeFrom  string `json:"visitTimeFrom"`
    VisitTimeTo    string `json:"visitTimeTo"`
    Relationship   string `json:"relationship"`
    QRCodeData     string `json:"qrCodeData"`
    Status         string `json:"status"` // "Active" or "Blocked"
	CurrentBlock  string `json:"currentBlock"` // BlockId if blocked
    CreatedAt      int64  `json:"createdAt"`
}
type Resident struct {
    DocType       string `json:"docType"`     // For CouchDB queries
    ResidentID    string `json:"residentId"`  // Matching API
    UserID        string `json:"userId"`      // Matching API
    Name          string `json:"name"`
    Email         string `json:"email"`
    Phone         string `json:"phone"`
    Gender        string `json:"gender"`
    Apartment     string `json:"apartment"`
    MaritalStatus string `json:"maritalStatus"`
    ResidentType  string `json:"residentType"`
    QRCodeData    string `json:"qrCodeData"`  // Matching API
    QRCodeImage   string `json:"qrCodeImage"` // Matching API
    CreatedAt     int64  `json:"createdAt"`
    UpdatedAt     int64  `json:"updatedAt"`
	Visitors     []VisitorInfo `json:"visitors"`
    IsBlocked bool `json:"isBlocked"`
}

type User struct {
	UserID           string    `json:"userId"`
	Email            string   `json:"email"`
	Password          string   `json:"password"`  
	Role              string   `json:"role"`       
	ProfileImage      string   `json:"image"`     
	
}

type AccessRequest struct {
	UserID  string   `json:"userId"` //resident_name
	VisitorID   string   `json:"visitorId"`
	Status      string   `json:"status"`
	ApprovedBy  []string `json:"approvedBy"`
	Timestamp   string   `json:"timestamp"`
	QRCode      string   `json:"qrCode"`  // 🔥 إضافة رمز QR
	Note        string    `json:"note"` 
	QRCodeTimestamp  int64    `json:"qrCodeTimestamp"` // ✅ Add this
}
// to block resident BlockVisitor function
type BlockRecord struct {
    DocType      string `json:"docType"`
    ResidentID   string `json:"residentId"`
    Reason       string `json:"reason"`
    BlockedBy    string `json:"blockedBy"`
    FromDateTime string `json:"fromDateTime"`
    ToDateTime   string `json:"toDateTime"`
    CreatedAt    int64  `json:"createdAt"`
}
// to block visitor BlockVisitor function
type BlockInfo struct {
    BlockId       string `json:"blockId"`
    VisitorId     string `json:"visitorId"`
    Reason        string `json:"reason"`
    BlockedBy     string `json:"blockedBy"`
    FromDateTime  int64  `json:"fromDateTime"`
    ToDateTime    int64  `json:"toDateTime"`
    CreatedAt     int64  `json:"createdAt"`
}
type VisitRequest struct {
	RequestID      string `json:"requestId"`
	CreatedBy      string `json:"createdBy"`
	TargetResident string `json:"targetResident"`
	VisitorName    string `json:"visitorName"`
	VisitorPhone   string `json:"visitorPhone"`
	Type           string `json:"type"`
	VisitPurpose   string `json:"visitPurpose"`
	CustomReason   string `json:"customReason"`
	VisitTimeFrom  string `json:"visitTimeFrom"`
	VisitTimeTo    string `json:"visitTimeTo"`
	VisitDate      string `json:"visitDate"`
	Status         string `json:"status"` // "Pending", "Approved", "Rejected"
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"` 
	StatusChangedBy string `json:"statusChangedBy,omitempty"`
}
type EntryLog struct {
    DocType     string `json:"docType"`    // For CouchDB queries
    LogID       string `json:"logId"`      // Unique ID
    RequestID  string `json:"requestId"` // Matches MongoDB
    Type        string `json:"type"`       // "enter" or "leave"
    Timestamp   int64  `json:"timestamp"`  // Unix timestamp
    APIEndpoint string `json:"apiEndpoint"`// For tracking
}


func (rc *ResidentsContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (rc *ResidentsContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "RegisterResident":
		return rc.RegisterResident(stub, args)  //me
	case "RegisterUser":
		return rc.RegisterUser(stub,args)   
	case "UpdateResident":
		return rc.UpdateResident(stub,args)  //me
    case "BlockResident":
		return rc.BlockResident(stub,args) //me
    case "UnblockResident":
		return rc.UnblockResident(stub,args) //me
	case "AddVisitor":
		return rc.AddVisitor(stub, args)  //me
	case "RequestApproval":
		return rc.RequestApproval(stub, args)
	case "ApproveRequest":
		return rc.ApproveRequest(stub, args)
	case "RejectRequest":
		return rc.RejectRequest(stub, args)  
	case "GetRequestStatus":
		return rc.GetRequestStatus(stub, args)
	case "GetAllResidents":  
		return rc.GetAllResidents(stub, args) //me
	case "GetResident":  
		return rc.GetResident(stub, args) //me
	case "GetVisitor":   
		return rc.GetVisitor(stub, args)  //me
	case "GetVisitors":  
		return rc.GetVisitors(stub, args)  //me
	case "BlockVisitor":   
		return rc.BlockVisitor(stub, args)  //me
	case "UnblockVisitor":   
		return rc.UnblockVisitor(stub, args) //me
	case "UpdateVisitor":   
		return rc.UpdateVisitor(stub, args)  //me
	case "AddVisitRequest":   
		return rc.AddVisitRequest(stub, args)   //me
	case "GetVisitRequest":
    return rc.GetVisitRequest(stub, args)  //me
	case "UpdateVisitRequestStatus":   
		return rc.UpdateVisitRequestStatus(stub, args)  //me
	case "CheckVisitorAuthorization":  
		return rc.CheckVisitorAuthorization(stub, args)
	case "RequestServiceAccess":
		return rc.RequestServiceAccess(stub, args)
	case "RequestDeliveryAccess":
		return rc.RequestDeliveryAccess(stub, args)
	case "ApproveServiceAccess":
		return rc.ApproveServiceAccess(stub, args)
	case "ApproveDeliveryAccess":
		return rc.ApproveDeliveryAccess(stub, args)
	case "GrantEmergencyAccess":
		return rc.GrantEmergencyAccess(stub, args)
	case "RejectRequestForWorker":
		return rc.RejectRequestForWorker(stub, args)
	case "RejectRequestForDelivery":
		return rc.RejectRequestForDelivery(stub, args)
	case "GetWorkerRequestStatus":
		return rc.GetWorkerRequestStatus(stub, args)
	case "GetDeliveryRequestStatus":
		return rc.GetDeliveryRequestStatus(stub,args)
	case "GetEmergencyRequest":
		return rc.GetEmergencyRequest(stub,args)	
	case "SaveLogToChain":
		return rc.SaveLogToChain(stub,args)  // me
	case "GetLastLogByResident":
		return rc.GetLastLogByResident(stub,args)  // me but not used
	case "SimulateQRCodeExpiry":
		// التأكد من أن المعاملات تحتوي على visitorID
		if len(args) < 1 {
			return shim.Error("❌ Visitor ID is required for SimulateQRCodeExpiry.")
		}
		visitorID := args[0]  // استخراج visitorID من args
		return rc.SimulateQRCodeExpiry(stub, visitorID)

	case "SimulateWorkerQRCodeExpiry":
		// التحقق من وجود معرف العامل
		if len(args) < 1 {
			return shim.Error("❌ Worker ID is required for SimulateQRCodeExpiryWorker.")
		}
		workerID := args[0]
		return rc.SimulateWorkerQRCodeExpiry(stub, workerID)

	case "SimulateDeliveryQRCodeExpiry":
		// التحقق من وجود معرف موظف التوصيل
		if len(args) < 1 {
			return shim.Error("❌ Delivery ID is required for SimulateDeliveryQRCodeExpiry.")
		}
		deliveryID := args[0]
		return rc.SimulateDeliveryQRCodeExpiry(stub, deliveryID)

	default:
		return shim.Error("❌ Invalid function name.")
	}
}



// Checks if a visitor is already in the list (VisitorInfo slice)
func containsVisitorInfo(visitors []VisitorInfo, visitorId string) bool {
	for _, v := range visitors {
		if v.VisitorId == visitorId {
			return true
		}
	}
	return false
}

// Checks if a string slice contains a value (e.g., ApprovedBy)
func containsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func generatePermanentQRCode(userId string) string {
	return fmt.Sprintf("QR-RESIDENT-%s", userId)
}

func (rc *ResidentsContract) RegisterResident(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    // Require all 8 fields that the API sends
    if len(args) < 8 {
        return shim.Error("❌ Required: ResidentID, Name, Email, Phone, Gender, MaritalStatus, ResidentType, Apartment")
    }

    // Parse all arguments
    residentId := args[0]
    name := args[1]
    email := args[2]
    phone := args[3]
    gender := args[4]
    maritalStatus := args[5]
    residentType := args[6]
    apartment := args[7]

    // Check if resident already exists (now using residentId as key)
    existingResident, err := stub.GetState("RESIDENT_" + residentId)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to check existing resident: %s", err.Error()))
    }
    if existingResident != nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s already exists", residentId))
    }

    // Generate QR code (using residentId as you do in API)
    qrCode := residentId // Or generate differently if needed
    timestamp := time.Now().Unix()

    // Create resident object matching API structure
    resident := Resident{
        DocType:       "resident", // Good practice for CouchDB queries
        ResidentID:    residentId,
        UserID:        residentId, // Or use separate user ID if needed
        Name:          name,
        Email:         email,
        Phone:         phone,
        Gender:        gender,
        Apartment:     apartment,
        MaritalStatus: maritalStatus,
        ResidentType:  residentType,
        QRCodeData:    qrCode,
        QRCodeImage:   residentId + ".png", // Matching API format
        CreatedAt:     timestamp,
        UpdatedAt:     timestamp,
    }

    // Marshal and store resident
    residentBytes, err := json.Marshal(resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal resident: %s", err.Error()))
    }

    // Store using ResidentID as key
    err = stub.PutState("RESIDENT_"+residentId, residentBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to store resident: %s", err.Error()))
    }

    // Return the full resident object as JSON
    return shim.Success(residentBytes)
}
func (rc *ResidentsContract) UpdateResident(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    // Validate arguments - now matches all fields from RegisterResident
    if len(args) < 8 {
        return shim.Error("❌ Required: ResidentID, Name, Email, Phone, Gender, MaritalStatus, ResidentType, Apartment")
    }

    residentId := args[0]
    name := args[1]
    email := args[2]
    phone := args[3]
    gender := args[4]
    maritalStatus := args[5]
    residentType := args[6]
    apartment := args[7]

    // Check if resident exists - using same key format as RegisterResident
    residentKey := "RESIDENT_" + residentId
    existingResidentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if existingResidentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s does not exist", residentId))
    }

    // Unmarshal existing resident
    var existingResident Resident
    err = json.Unmarshal(existingResidentBytes, &existingResident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Check if apartment is changing
    if existingResident.Apartment != apartment {
        // Get building capacity (same as RegisterResident)
        buildingBytes, err := stub.GetState("BUILDING_CONFIG")
        if err != nil {
            return shim.Error(fmt.Sprintf("❌ Failed to read building config: %s", err.Error()))
        }
        if buildingBytes == nil {
            return shim.Error("❌ Building configuration not found")
        }

        var building BuildingConfig
        err = json.Unmarshal(buildingBytes, &building)
        if err != nil {
            return shim.Error(fmt.Sprintf("❌ Failed to unmarshal building config: %s", err.Error()))
        }

        // Count residents in new apartment
        residentsCount, err := rc.getResidentsCountInApartment(stub, apartment)
        if err != nil {
            return shim.Error(fmt.Sprintf("❌ Failed to count residents: %s", err.Error()))
        }

        if residentsCount >= building.ResidentsPerApartment {
            return shim.Error("❌ Maximum number of residents for this apartment reached")
        }
    }

    // Update ALL resident fields (not just some) to match RegisterResident
    existingResident.Name = name
    existingResident.Email = email
    existingResident.Phone = phone
    existingResident.Gender = gender
    existingResident.MaritalStatus = maritalStatus
    existingResident.ResidentType = residentType
    existingResident.Apartment = apartment
    existingResident.UpdatedAt = time.Now().Unix()

    // Marshal and save updated resident
    updatedResidentBytes, err := json.Marshal(existingResident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal updated resident: %s", err.Error()))
    }

    err = stub.PutState(residentKey, updatedResidentBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update resident: %s", err.Error()))
    }

    // Return the full updated resident object as JSON (same as RegisterResident)
    return shim.Success(updatedResidentBytes)
}

// BlockResident adds a block record to the resident
func (rc *ResidentsContract) BlockResident(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    // Validate arguments
    if len(args) < 6 {
        return shim.Error("❌ Required: ResidentID, Reason, BlockedBy, FromDate, FromTime, ToDate, ToTime")
    }

    residentId := args[0]
    reason := args[1]
    blockedBy := args[2]
    fromDate := args[3]
    fromTime := args[4]
    toDate := args[5]
    toTime := args[6]

    // Check if resident exists
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s does not exist", residentId))
    }

    // Check if already blocked
    blockKey := "BLOCK_" + residentId
    existingBlock, _ := stub.GetState(blockKey)
    if existingBlock != nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s is already blocked", residentId))
    }

    // Create block record
    block := BlockRecord{
        DocType:      "block",
        ResidentID:   residentId,
        Reason:       reason,
        BlockedBy:    blockedBy,
        FromDateTime: fmt.Sprintf("%sT%s", fromDate, fromTime),
        ToDateTime:   fmt.Sprintf("%sT%s", toDate, toTime),
        CreatedAt:    time.Now().Unix(),
    }

    blockBytes, err := json.Marshal(block)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal block record: %s", err.Error()))
    }

    // Save block record
    err = stub.PutState(blockKey, blockBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to block resident: %s", err.Error()))
    }

    // Update resident status
    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    resident.IsBlocked = true
    resident.UpdatedAt = time.Now().Unix()

    updatedResidentBytes, _ := json.Marshal(resident)
    stub.PutState(residentKey, updatedResidentBytes)

    return shim.Success(blockBytes)
}
func (rc *ResidentsContract) GetResident(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    // Validate input
    if len(args) != 1 {
        return shim.Error("❌ Required: ResidentID")
    }

    residentId := args[0]
    
    // Retrieve resident from ledger
    residentBytes, err := stub.GetState("RESIDENT_" + residentId)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s does not exist", residentId))
    }

    // Return the resident as-is (already JSON formatted)
    return shim.Success(residentBytes)
}
// UnblockResident removes a block from a resident
func (rc *ResidentsContract) UnblockResident(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 1 {
        return shim.Error("❌ Required: ResidentID")
    }

    residentId := args[0]

    // Check if resident exists
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s does not exist", residentId))
    }

    // Check if blocked
    blockKey := "BLOCK_" + residentId
    blockBytes, err := stub.GetState(blockKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read block record: %s", err.Error()))
    }
    if blockBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s is not blocked", residentId))
    }

    // Remove block
    err = stub.DelState(blockKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unblock resident: %s", err.Error()))
    }

    // Update resident status
    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    resident.IsBlocked = false
    resident.UpdatedAt = time.Now().Unix()

    updatedResidentBytes, _ := json.Marshal(resident)
    stub.PutState(residentKey, updatedResidentBytes)

    return shim.Success(nil)
}



// Helper function to count residents in an apartment
func (rc *ResidentsContract) getResidentsCountInApartment(stub shim.ChaincodeStubInterface, apartment string) (int, error) {
    // This is a simplified implementation - you might need a more efficient way
    // in production, like using a composite key index
    
    queryString := fmt.Sprintf(`{
        "selector": {
            "docType": "resident",
            "Apartment": "%s"
        }
    }`, apartment)

    resultsIterator, err := stub.GetQueryResult(queryString)
    if err != nil {
        return 0, err
    }
    defer resultsIterator.Close()

    count := 0
    for resultsIterator.HasNext() {
        _, err := resultsIterator.Next()
        if err != nil {
            return 0, err
        }
        count++
    }

    return count, nil
}
func (rc *ResidentsContract) RegisterUser(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 5 {
		return shim.Error("❌ Required: User Email, Password, Role, Image")
	}
    userId :=args[0]
	email := args[1]
	password := args[2]
	role := args[3]
	image := args[4]

	// Check if the user already exists
	userKey := "USER_" +userId
	existingUser, _ := stub.GetState(userKey)
	if existingUser != nil {
		return shim.Error(fmt.Sprintf("❌ User %s already exists.", userId))
	}

	// Create the new user object
	user := User{
		UserID:   userId,
		Email:    email,
		Password: password,
		Role:     role,
		ProfileImage:   image,
	}

	// Store it on the ledger
	userBytes, _ := json.Marshal(user)
	err := stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error("❌ Failed to store user.")
	}

	return shim.Success([]byte(fmt.Sprintf("✅ User %s registered successfully.", userId)))
}


// ✅ توليد رمز QR دائم للزائر مع الطابع الزمني
func generateVisitorQRCode(visitorId string) string {
	// Get the current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05") // Format: YYYY-MM-DD HH:MM:SS

	// Generate a permanent QR code for the visitor, including the timestamp
	return fmt.Sprintf("QR-VISITOR-%s-%s", visitorId, currentTime)
}


func (rc *ResidentsContract) AddVisitor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 7 {
        return shim.Error("❌ Required: ResidentID, VisitorID, FullName, Phone, VisitTimeFrom, VisitTimeTo, Relationship")
    }

    residentId := args[0]
    visitorId := args[1]
    fullName := args[2]
    phone := args[3]
    visitTimeFrom := args[4]
    visitTimeTo := args[5]
    relationship := args[6]

    // Get resident from ledger
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s not found", residentId))
    }

    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Check if visitor already exists
    for _, v := range resident.Visitors {
        if v.VisitorId == visitorId {
            return shim.Error(fmt.Sprintf("❌ Visitor %s already exists for resident %s", visitorId, residentId))
        }
    }

    // Create new visitor (without storing image)
    newVisitor := VisitorInfo{
        VisitorId:      visitorId,
        FullName:       fullName,
        Phone:          phone,
        VisitTimeFrom:  visitTimeFrom,
        VisitTimeTo:    visitTimeTo,
        Relationship:   relationship,
        QRCodeData:     visitorId, // Using visitorId as QR data
        Status:         "Active",
        CreatedAt:      time.Now().Unix(),
    }

    // Add visitor to resident's list
    resident.Visitors = append(resident.Visitors, newVisitor)
    resident.UpdatedAt = time.Now().Unix()

    // Save updated resident
    updatedResidentBytes, err := json.Marshal(resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal resident: %s", err.Error()))
    }

    err = stub.PutState(residentKey, updatedResidentBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update resident: %s", err.Error()))
    }

    // Return success with visitor details
    response := map[string]interface{}{
        "success": true,
        "visitor": newVisitor,
    }
    responseBytes, _ := json.Marshal(response)
    return shim.Success(responseBytes)
}

func (rc *ResidentsContract) GetVisitors(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 1 {
        return shim.Error("❌ ResidentID required")
    }

    residentId := args[0]
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s not found", residentId))
    }

    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Return visitors list
    response := map[string]interface{}{
        "success":  true,
        "visitors": resident.Visitors,
    }
    responseBytes, _ := json.Marshal(response)
    return shim.Success(responseBytes)
}
func (rc *ResidentsContract) GetVisitor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 {
        return shim.Error("❌ ResidentID and VisitorID required")
    }

    residentId := args[0]
    visitorId := args[1]

    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s not found", residentId))
    }

    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Find the specific visitor
    for _, visitor := range resident.Visitors {
        if visitor.VisitorId == visitorId {
            response := map[string]interface{}{
                "success": true,
                "visitor": visitor,
            }
            responseBytes, _ := json.Marshal(response)
            return shim.Success(responseBytes)
        }
    }

    return shim.Error(fmt.Sprintf("❌ Visitor %s not found for resident %s", visitorId, residentId))
}
func (rc *ResidentsContract) BlockVisitor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 6 {
        return shim.Error("❌ Required: VisitorID, ResidentID, Reason, FromDate, FromTime, ToDate, ToTime, BlockedBy")
    }

    visitorId := args[0]
    residentId := args[1]
    reason := args[2]
    fromDate := args[3]
    fromTime := args[4]
    toDate := args[5]
    toTime := args[6]
    blockedBy := args[7]

    // Parse datetime strings to timestamps
    fromDateTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", fromDate, fromTime))
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Invalid FromDateTime format: %s", err.Error()))
    }

    toDateTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", toDate, toTime))
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Invalid ToDateTime format: %s", err.Error()))
    }

    // Get resident from ledger
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s not found", residentId))
    }

    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Find visitor in resident's list
    visitorFound := false
    for i, visitor := range resident.Visitors {
        if visitor.VisitorId == visitorId {
            // Check if already blocked
            if visitor.Status == "Blocked" {
                return shim.Error(fmt.Sprintf("❌ Visitor %s is already blocked", visitorId))
            }

            // Create block record
            blockId := fmt.Sprintf("BLOCK_%s_%d", visitorId, time.Now().Unix())
            block := BlockInfo{
                BlockId:      blockId,
                VisitorId:    visitorId,
                Reason:       reason,
                BlockedBy:    blockedBy,
                FromDateTime: fromDateTime.Unix(),
                ToDateTime:   toDateTime.Unix(),
                CreatedAt:    time.Now().Unix(),
            }

            // Update visitor status
            resident.Visitors[i].Status = "Blocked"
            resident.Visitors[i].CurrentBlock = blockId

            // Save block record
            blockBytes, err := json.Marshal(block)
            if err != nil {
                return shim.Error(fmt.Sprintf("❌ Failed to marshal block: %s", err.Error()))
            }
            err = stub.PutState(blockId, blockBytes)
            if err != nil {
                return shim.Error(fmt.Sprintf("❌ Failed to save block: %s", err.Error()))
            }

            visitorFound = true
            break
        }
    }

    if !visitorFound {
        return shim.Error(fmt.Sprintf("❌ Visitor %s not found for resident %s", visitorId, residentId))
    }

    // Save updated resident
    updatedResidentBytes, err := json.Marshal(resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal resident: %s", err.Error()))
    }

    err = stub.PutState(residentKey, updatedResidentBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update resident: %s", err.Error()))
    }

    return shim.Success([]byte(fmt.Sprintf("✅ Visitor %s blocked successfully", visitorId)))
}

func (rc *ResidentsContract) UnblockVisitor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 {
        return shim.Error("❌ Required: VisitorID, ResidentID")
    }

    visitorId := args[0]
    residentId := args[1]

    // Get resident from ledger
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s not found", residentId))
    }

    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Find visitor in resident's list
    visitorFound := false
    for i, visitor := range resident.Visitors {
        if visitor.VisitorId == visitorId {
            if visitor.Status != "Blocked" {
                return shim.Error(fmt.Sprintf("❌ Visitor %s is not blocked", visitorId))
            }

            // Delete block record
            err = stub.DelState(visitor.CurrentBlock)
            if err != nil {
                return shim.Error(fmt.Sprintf("❌ Failed to delete block record: %s", err.Error()))
            }

            // Update visitor status
            resident.Visitors[i].Status = "Active"
            resident.Visitors[i].CurrentBlock = ""

            visitorFound = true
            break
        }
    }

    if !visitorFound {
        return shim.Error(fmt.Sprintf("❌ Visitor %s not found for resident %s", visitorId, residentId))
    }

    // Save updated resident
    updatedResidentBytes, err := json.Marshal(resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal resident: %s", err.Error()))
    }

    err = stub.PutState(residentKey, updatedResidentBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update resident: %s", err.Error()))
    }

    return shim.Success([]byte(fmt.Sprintf("✅ Visitor %s unblocked successfully", visitorId)))
}


func (rc *ResidentsContract) UpdateVisitor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    // Validate arguments
    if len(args) < 4 {
        return shim.Error("❌ Required: ResidentID, VisitorID, Phone, VisitTimeFrom, VisitTimeTo")
    }

    residentId := args[0]
    visitorId := args[1]
    phone := args[2]
    visitTimeFrom := args[3]
    visitTimeTo := args[4]

    // Get resident from ledger
    residentKey := "RESIDENT_" + residentId
    residentBytes, err := stub.GetState(residentKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read resident: %s", err.Error()))
    }
    if residentBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Resident %s not found", residentId))
    }

    // Unmarshal resident data
    var resident Resident
    err = json.Unmarshal(residentBytes, &resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
    }

    // Find and update visitor
    visitorFound := false
    for i, visitor := range resident.Visitors {
        if visitor.VisitorId == visitorId {
            // Update visitor fields
            resident.Visitors[i].Phone = phone
            resident.Visitors[i].VisitTimeFrom = visitTimeFrom
            resident.Visitors[i].VisitTimeTo = visitTimeTo
            resident.UpdatedAt = time.Now().Unix()
            visitorFound = true
            break
        }
    }

    if !visitorFound {
        return shim.Error(fmt.Sprintf("❌ Visitor %s not found for resident %s", visitorId, residentId))
    }

    // Save updated resident
    updatedResidentBytes, err := json.Marshal(resident)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal resident: %s", err.Error()))
    }

    err = stub.PutState(residentKey, updatedResidentBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update resident: %s", err.Error()))
    }

    return shim.Success([]byte(fmt.Sprintf("✅ Visitor %s updated successfully", visitorId)))
}

func (rc *ResidentsContract) AddVisitRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 11 {
        return shim.Error("❌ Required: RequestID, CreatedBy, TargetResident, VisitorName, VisitorPhone, Type, VisitPurpose, CustomReason, VisitTimeFrom, VisitTimeTo, VisitDate")
    }

    // Use the provided requestId instead of generating a new one
    requestId := args[0]
    createdBy := args[1]
    targetResident := args[2]
    visitorName := args[3]
    visitorPhone := args[4]
    visitType := args[5]
    visitPurpose := args[6]
    customReason := args[7]
    visitTimeFrom := args[8]
    visitTimeTo := args[9]
    visitDate := args[10]

    // Create new visit request
    newRequest := VisitRequest{
        RequestID:     requestId, // Use the provided ID
        CreatedBy:     createdBy,
        TargetResident: targetResident,
        VisitorName:   visitorName,
        VisitorPhone:  visitorPhone,
        Type:          visitType,
        VisitPurpose:  visitPurpose,
        CustomReason:  customReason,
        VisitTimeFrom: visitTimeFrom,
        VisitTimeTo:   visitTimeTo,
        VisitDate:     visitDate,
        Status:        "Pending",
        CreatedAt:     time.Now().Unix(),
    }

    // Marshal and store the request
    requestBytes, err := json.Marshal(newRequest)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal visit request: %s", err.Error()))
    }

    requestKey := fmt.Sprintf("VISITREQUEST_%s", requestId)
    err = stub.PutState(requestKey, requestBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to store visit request: %s", err.Error()))
    }

    return shim.Success(requestBytes)
}
// like handleStatusChange in request.js
func (rc *ResidentsContract) UpdateVisitRequestStatus(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 3 {
        return shim.Error("❌ Required arguments: RequestID, Status, RequestedBy")
    }

    requestID := args[0]
    status := args[1]
    requestedBy := args[2]

    // Validate status
    allowedStatuses := map[string]bool{
        "accepted": true,
        "rejected": true,
    }

    if !allowedStatuses[status] {
        return shim.Error("❌ Invalid status value.")
    }

    requestKey := fmt.Sprintf("VISITREQUEST_%s", requestID)

    // Get the existing visit request
    requestBytes, err := stub.GetState(requestKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to get visit request: %s", err.Error()))
    }
    if requestBytes == nil {
        return shim.Error("❌ Visit request not found.")
    }

    var request VisitRequest
    err = json.Unmarshal(requestBytes, &request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to parse visit request: %s", err.Error()))
    }

    // Update status and audit fields
    request.Status = status
    request.UpdatedAt = time.Now().Unix()
    request.StatusChangedBy = requestedBy

    updatedBytes, err := json.Marshal(request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal updated visit request: %s", err.Error()))
    }

    // Save the updated request
    err = stub.PutState(requestKey, updatedBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to store updated request: %s", err.Error()))
    }

    return shim.Success(updatedBytes)
}

// ✅ 3️⃣ إرسال طلب دخول زائر للموافقة عليه.
func (rc *ResidentsContract) RequestApproval(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Resident ID and Visitor ID are required.")
	}
	userId, visitorId := args[0], args[1]

	// التحقق مما إذا كان الزائر موجودًا مسبقًا في قائمة الزوار
	residentBytes, _ := stub.GetState("RESIDENT_" + userId)
	//residentBytes, _ := stub.GetPrivateData("ResidentPrivateCollection", "RESIDENT_"+residentId)
	if residentBytes != nil {
		var resident Resident
		json.Unmarshal(residentBytes, &resident)

		// Check if visitor is already in the list
		if containsVisitorInfo(resident.Visitors, visitorId) {
			return shim.Success([]byte(fmt.Sprintf("✅ Visitor %s is already approved for %s.", visitorId, userId)))
		}
	}
	

	requestKey := "REQ_" + visitorId
	requestBytes, _ := stub.GetState(requestKey)
	if requestBytes != nil {
		return shim.Error(fmt.Sprintf("❌ Access request for %s already exists.", visitorId))
	}

	timestamp := time.Now().Format(time.RFC3339)
	request := AccessRequest{UserID: userId, VisitorID: visitorId, Status: "Pending", ApprovedBy: []string{}, Timestamp: timestamp}
	requestBytes, _ = json.Marshal(request)
	stub.PutState(requestKey, requestBytes)
	//stub.PutPrivateData("VisitorApprovalCollection", requestKey, requestBytes)
	return shim.Success([]byte(fmt.Sprintf("✅ Request for visitor %s submitted.", visitorId)))
}
//send request for worker
// ✅ العامل يرسل طلب دخول ليتم مراجعته من المدير
func (rc *ResidentsContract) RequestServiceAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Worker ID and Requester ID are required.")
	}
	workerId := args[0]
	requesterId := args[1]

	requestKey := "REQ_SERVICE_" + workerId
	existingBytes, _ := stub.GetState(requestKey)
	if existingBytes != nil {
		return shim.Error(fmt.Sprintf("❌ Access request for worker %s already exists.", workerId))
	}

	timestamp := time.Now().Format(time.RFC3339)
	request := AccessRequest{
		UserID:     requesterId,
		VisitorID:  workerId,
		Status:     "Pending",
		ApprovedBy: []string{},
		Timestamp:  timestamp,
		Note:       "Awaiting approval for service access.",
	}

	requestBytes, _ := json.Marshal(request)
	stub.PutState(requestKey, requestBytes)
	return shim.Success([]byte(fmt.Sprintf("✅ Request for worker %s submitted.", workerId)))
}
//delivery send request to resident
func (rc *ResidentsContract) RequestDeliveryAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Delivery ID and Resident ID are required.")
	}
	deliveryId := args[0]
	residentId := args[1]
	requestKey := "REQ_DELIVERY_" + deliveryId

	// التحقق من وجود طلب مسبق
	existingBytes, _ := stub.GetState(requestKey)
	if existingBytes != nil {
		return shim.Error(fmt.Sprintf("❌ Request for delivery ID %s already exists.", deliveryId))
	}

	timestamp := time.Now().Format(time.RFC3339)
	request := AccessRequest{
		VisitorID: deliveryId,
		UserID:    residentId,
		Status:    "Pending",
		ApprovedBy: []string{},
		Timestamp: timestamp,
		Note:      "Delivery access request pending resident approval.",
	}

	requestBytes, _ := json.Marshal(request)
	stub.PutState(requestKey, requestBytes)

	return shim.Success([]byte(fmt.Sprintf("✅ Delivery access request from %s submitted and awaiting approval.", deliveryId)))
}


func (rc *ResidentsContract) ApproveRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 {
        return shim.Error("❌ Visitor ID and Approver ID are required.")
    }
    visitorId, approverId := args[0], args[1]
    requestKey := "REQ_" + visitorId
    requestBytes, err := stub.GetState(requestKey)
    if err != nil || requestBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Request for %s not found.", visitorId))
    }

    var request AccessRequest
    json.Unmarshal(requestBytes, &request)

    // ✅ If the status is "Approved", check the validity of the QR code
    if request.Status == "Approved" {
        if isQRCodeExpired(request.QRCodeTimestamp) {
            // ❌ QR code expired → Change status to "Pending", remove QR code and add reason
            request.Status = "Pending"
            request.QRCode = ""
            request.QRCodeTimestamp = 0
            request.Note = "QR code is invalid because the time has expired. Re-approval from the resident and manager is required."

            requestBytes, _ = json.Marshal(request)
            stub.PutState(requestKey, requestBytes)

            return shim.Success([]byte(fmt.Sprintf("⚠️ QR code for %s expired. Status set to Pending. Awaiting new approvals.", visitorId)))
        }
        return shim.Error(fmt.Sprintf("❌ Request for %s is already approved and still valid.", visitorId))
    }

    // ✅ Add the approver to the list
    if !contains(request.ApprovedBy, approverId) {
        request.ApprovedBy = append(request.ApprovedBy, approverId)
    }

    // ✅ Check if both the resident and manager have approved
    if contains(request.ApprovedBy, request.UserID) && contains(request.ApprovedBy, "Manager") {
        request.Status = "Approved"
        request.QRCode = generateQRCode(visitorId)
        request.QRCodeTimestamp = time.Now().Unix()
        request.Note = "Request approved and QR code is valid for 10 hours."
    }

    // ✅ Update the request in the ledger
    requestBytes, _ = json.Marshal(request)
    stub.PutState(requestKey, requestBytes)

    return shim.Success([]byte(fmt.Sprintf("✅ Approval from %s recorded for %s.", approverId, visitorId)))
}


// Function to check if the QR code is expired (older than 10 hours)
func isQRCodeExpired(timestamp int64) bool {
	expiryDuration := int64(10 * 60 * 60) // 10 hours in seconds(available for 10 hours)
	currentTime := time.Now().Unix()
	return currentTime > (timestamp + expiryDuration)
}

// Function to generate a unique QR code with timestamp
func generateQRCode(visitorId string) string {
	timestamp := time.Now().Unix()
	qrCode := fmt.Sprintf("QR-%s-%d", visitorId, timestamp)
	return qrCode
}

// دالة لتعديل وقت صلاحية رمز QR (لأغراض اختبارية)
func (s *ResidentsContract) SimulateQRCodeExpiry(stub shim.ChaincodeStubInterface, visitorID string) sc.Response {
    // ✅ استخدام نفس المفتاح المستخدم في ApproveRequest
    requestKey := "REQ_" + visitorID

    // الحصول على بيانات الطلب
    requestAsBytes, err := stub.GetState(requestKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to get request for %s: %v", visitorID, err))
    }

    if requestAsBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Request for %s does not exist", visitorID))
    }

    // تحويل البيانات إلى هيكل الطلب
    var request AccessRequest
    err = json.Unmarshal(requestAsBytes, &request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal request: %v", err))
    }

    // ⏳ محاكاة انتهاء صلاحية رمز QR بتعيين وقت قديم جدًا (منذ 11 ساعة)
    expiredTimestamp := time.Now().Add(-11 * time.Hour).Unix()
    request.QRCodeTimestamp = expiredTimestamp

    // تحديث السجل مع الوقت الجديد
    requestAsBytes, err = json.Marshal(request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal updated request: %v", err))
    }

    err = stub.PutState(requestKey, requestAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update request: %v", err))
    }

    return shim.Success([]byte(fmt.Sprintf("✅ QR code expiry simulated for %s", visitorID)))
}

//approve request for worker
// ✅ المدير يوافق على طلب الدخول ويتم توليد QR صالح لساعتين
func (rc *ResidentsContract) ApproveServiceAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Worker ID and Approver ID are required.")
	}
	workerId := args[0]
	approverId := args[1]
	requestKey := "REQ_SERVICE_" + workerId

	requestBytes, err := stub.GetState(requestKey)
	if err != nil || requestBytes == nil {
		return shim.Error(fmt.Sprintf("❌ Request for worker %s not found.", workerId))
	}

	var request AccessRequest
	json.Unmarshal(requestBytes, &request)

	if request.Status == "Approved" {
		if isQRCodeExpiredCustomfor2(request.QRCodeTimestamp, 2*60*60) {
			request.Status = "Pending"
			request.QRCode = ""
			request.QRCodeTimestamp = 0
			request.Note = "QR expired. Awaiting re-approval."
		} else {
			return shim.Success([]byte("✅ QR still valid for worker."))
		}
	}

	if !contains(request.ApprovedBy, approverId) {
		request.ApprovedBy = append(request.ApprovedBy, approverId)
	}

	request.Status = "Approved"
	request.QRCode = generateQRCodeWorker(workerId)
	request.QRCodeTimestamp = time.Now().Unix()
	request.Note = "Service access granted. QR valid for 2 hours."

	requestBytes, _ = json.Marshal(request)
	stub.PutState(requestKey, requestBytes)

	return shim.Success([]byte(fmt.Sprintf("✅ Approval recorded. QR generated for worker %s.", workerId)))
}

func generateQRCodeWorker(workerId string) string {
	timestamp := time.Now().Unix()
	qrCode := fmt.Sprintf("QR-%s-%d", workerId, timestamp)
	return qrCode
}

//just for test
// ✅ دالة لمحاكاة انتهاء صلاحية رمز QR لفني صيانة (Worker)
func (s *ResidentsContract) SimulateWorkerQRCodeExpiry(stub shim.ChaincodeStubInterface, workerID string) sc.Response {
    requestKey := "REQ_SERVICE_" + workerID

    // جلب الطلب من السجل
    requestAsBytes, err := stub.GetState(requestKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to get request for %s: %v", workerID, err))
    }

    if requestAsBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Request for worker %s does not exist", workerID))
    }

    // تحويل JSON إلى هيكل AccessRequest
    var request AccessRequest
    err = json.Unmarshal(requestAsBytes, &request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal request: %v", err))
    }

    // محاكاة انتهاء الصلاحية بتعيين وقت سابق (مثلاً 3 ساعات ماضية)
    expiredTimestamp := time.Now().Add(-3 * time.Hour).Unix()
    request.QRCodeTimestamp = expiredTimestamp

    // إعادة حفظ التحديث في السجل
    updatedBytes, err := json.Marshal(request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal updated request: %v", err))
    }

    err = stub.PutState(requestKey, updatedBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update request in ledger: %v", err))
    }

    return shim.Success([]byte(fmt.Sprintf("✅ QR code expiry simulated for worker %s", workerID)))
}

//resident approve request and generate qrcode limited with 30min

func (rc *ResidentsContract) ApproveDeliveryAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Delivery ID and Resident ID are required.")
	}
	deliveryId := args[0]
	residentId := args[1]
	requestKey := "REQ_DELIVERY_" + deliveryId

	requestBytes, err := stub.GetState(requestKey)
	if err != nil || requestBytes == nil {
		return shim.Error(fmt.Sprintf("❌ Request for delivery %s not found.", deliveryId))
	}

	var request AccessRequest
	json.Unmarshal(requestBytes, &request)

	if request.Status == "Approved" {
		if isQRCodeExpiredCustomfor2(request.QRCodeTimestamp, 30*60) {
			request.Status = "Pending"
			request.QRCode = ""
			request.QRCodeTimestamp = 0
			request.Note = "⚠️ QR code expired. Request set to Pending again."
			requestBytes, _ := json.Marshal(request)
			stub.PutState(requestKey, requestBytes)
			return shim.Success([]byte("⚠️ QR expired. Please approve again."))
		}
		return shim.Success([]byte("✅ QR code is still valid for delivery access."))
	}

	 // ✅ Add the approver to the list
	 if !contains(request.ApprovedBy,residentId) {
        request.ApprovedBy = append(request.ApprovedBy,residentId)
    }

    

	request.Status = "Approved"
	request.QRCode = generateQRCodeDelivery(deliveryId)
	request.QRCodeTimestamp = time.Now().Unix()
	//request.ApprovedBy = append(request.ApprovedBy, residentId)
	request.Note = "Delivery access approved. QR code valid for 30 minutes."

	requestBytes, _ = json.Marshal(request)
	stub.PutState(requestKey, requestBytes)

	return shim.Success([]byte(fmt.Sprintf("✅ Delivery access for %s approved. QR generated.", deliveryId)))
}
// Helper to check QR expiration with custom duration (in seconds)
func isQRCodeExpiredCustomfor2(timestamp int64, validDuration int64) bool {
	currentTime := time.Now().Unix()
	return currentTime > (timestamp + validDuration)
}
func generateQRCodeDelivery(deliveryId string) string {
	timestamp := time.Now().Unix()
	qrCode := fmt.Sprintf("QR-%s-%d",deliveryId, timestamp)
	return qrCode
}
//just for test
// ✅ دالة لمحاكاة انتهاء صلاحية رمز QR لموظف توصيل (Delivery)
func (s *ResidentsContract) SimulateDeliveryQRCodeExpiry(stub shim.ChaincodeStubInterface, deliveryID string) sc.Response {
    requestKey := "REQ_DELIVERY_" + deliveryID

    // جلب الطلب من السجل
    requestAsBytes, err := stub.GetState(requestKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to get request for %s: %v", deliveryID, err))
    }

    if requestAsBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Request for delivery %s does not exist", deliveryID))
    }

    // تحويل JSON إلى هيكل AccessRequest
    var request AccessRequest
    err = json.Unmarshal(requestAsBytes, &request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to unmarshal request: %v", err))
    }

    // محاكاة انتهاء الصلاحية بتعيين وقت سابق (مثلاً 1 ساعة ماضية)
    expiredTimestamp := time.Now().Add(-1 * time.Hour).Unix()
    request.QRCodeTimestamp = expiredTimestamp

    // إعادة حفظ التحديث في السجل
    updatedBytes, err := json.Marshal(request)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal updated request: %v", err))
    }

    err = stub.PutState(requestKey, updatedBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to update request in ledger: %v", err))
    }

    return shim.Success([]byte(fmt.Sprintf("✅ QR code expiry simulated for delivery %s", deliveryID)))
}

// for emergency creation authorization fawriya
// ✅ Grant Emergency Access (By Building Manager or Emergency System)
func (rc *ResidentsContract) GrantEmergencyAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 {
        return shim.Error("❌ Resident ID and Visitor ID are required.")
    }
    visitorId, residentId := args[0], args[1]
	
    // Create a unique request key for the emergency access
    requestKey := "EMERGENCY_" + visitorId

    // Check if the emergency access request already exists
    existingBytes, _ := stub.GetState(requestKey)
    if existingBytes != nil {
        return shim.Error("❌ Emergency access for this visitor has already been granted.")
    }

    // Create an emergency access request
    timestamp := time.Now().Format(time.RFC3339)

    request := AccessRequest{
        UserID:          residentId,
        VisitorID:       visitorId,
        Status:          "Approved",
        ApprovedBy:      []string{"Manager"},
		Timestamp:        timestamp,
        QRCode:          "", // No QR code for emergency access
        QRCodeTimestamp: 0,
        Note:            "Emergency access granted immediately without prior approval.",
    }

    // Store the emergency access request in the ledger
    bytes, _ := json.Marshal(request)
    stub.PutState(requestKey, bytes)

    return shim.Success([]byte(fmt.Sprintf("✅ Emergency access granted for visitor %s.", visitorId)))
}


// ✅ 5️⃣ رفض الطلب
func (rc *ResidentsContract) RejectRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Visitor ID and Rejector ID are required.")
	}
	visitorId := args[0]
	requestKey := "REQ_" + visitorId
	requestBytes, err := stub.GetState(requestKey)
	//requestBytes, err := stub.GetPrivateData("VisitorApprovalCollection", requestKey)
	if err != nil || requestBytes == nil {
		return shim.Error(fmt.Sprintf("❌ Request for %s not found.", visitorId))
	}
	var request AccessRequest
	json.Unmarshal(requestBytes, &request)

	if request.Status == "Approved" {
		return shim.Error("❌ Cannot reject an already approved request.")
	}

	request.Status = "Rejected"
	requestBytes, _ = json.Marshal(request)
	stub.PutState(requestKey, requestBytes)
	//stub.PutPrivateData("VisitorApprovalCollection", requestKey, requestBytes)
	return shim.Success([]byte(fmt.Sprintf("❌ Request for %s has been rejected.", visitorId)))
}

// ✅ Reject Worker Access Request (Rejected by Manager)
func (rc *ResidentsContract) RejectRequestForWorker(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 {
        return shim.Error("❌ Worker ID and Rejector ID (Manager) are required.")
    }
    workerId := args[0]
    rejectorId := args[1]
    requestKey := "REQ_SERVICE_" + workerId
    requestBytes, err := stub.GetState(requestKey)
    
    if err != nil || requestBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Request for worker %s not found.", workerId))
    }

    var request AccessRequest
    json.Unmarshal(requestBytes, &request)

    if request.Status == "Approved" {
        return shim.Error("❌ Cannot reject an already approved request for worker.")
    }

    // If already rejected, prevent re-rejection
    if request.Status == "Rejected" {
        return shim.Error(fmt.Sprintf("❌ Request for worker %s is already rejected.", workerId))
    }

    // Reject the worker's access request
    request.Status = "Rejected"
    //request.RejectedBy = append(request.RejectedBy, rejectorId)
    request.Note = fmt.Sprintf("Worker access rejected by manager %s.", rejectorId)

    // Update the request in the ledger
    requestBytes, _ = json.Marshal(request)
    stub.PutState(requestKey, requestBytes)

    return shim.Success([]byte(fmt.Sprintf("❌ Worker access request for %s has been rejected by manager %s.", workerId, rejectorId)))
}


// ✅ Reject Delivery Access Request (Rejected by Resident)
func (rc *ResidentsContract) RejectRequestForDelivery(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 {
        return shim.Error("❌ Delivery ID and Rejector ID (Resident) are required.")
    }
    deliveryId := args[0]
    rejectorId := args[1]
    requestKey := "REQ_DELIVERY_" + deliveryId
    requestBytes, err := stub.GetState(requestKey)
    
    if err != nil || requestBytes == nil {
        return shim.Error(fmt.Sprintf("❌ Request for delivery %s not found.", deliveryId))
    }

    var request AccessRequest
    json.Unmarshal(requestBytes, &request)

    if request.Status == "Approved" {
        return shim.Error("❌ Cannot reject an already approved request for delivery.")
    }

    // If already rejected, prevent re-rejection
    if request.Status == "Rejected" {
        return shim.Error(fmt.Sprintf("❌ Request for delivery %s is already rejected.", deliveryId))
    }

    // Reject the delivery access request
    request.Status = "Rejected"
   // request.RejectedBy = append(request.RejectedBy, rejectorId)
    request.Note = fmt.Sprintf("Delivery access rejected by resident %s.", rejectorId)

    // Update the request in the ledger
    requestBytes, _ = json.Marshal(request)
    stub.PutState(requestKey, requestBytes)

    return shim.Success([]byte(fmt.Sprintf("❌ Delivery access request for %s has been rejected by resident %s.", deliveryId, rejectorId)))
}


//GetRequestStatus → استعلام عن حالة طلب معين.
/*func (rc *ResidentsContract) GetRequestStatus(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("Visitor ID is required.")
	}
	visitorId := args[0]
	requestBytes, _ := stub.GetState("REQ_" + visitorId)
	if requestBytes == nil {
		return shim.Error(fmt.Sprintf("No request found for visitor %s.", visitorId))
	}
	return shim.Success(requestBytes)
}*/
func (rc *ResidentsContract) GetRequestStatus(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("❌ Visitor ID is required.")
	}
	visitorId := args[0]
	requestBytes, _ := stub.GetState("REQ_" + visitorId)
	//requestBytes, _ := stub.GetPrivateData("VisitorApprovalCollection", "REQ_" + visitorId)
	if requestBytes == nil {
		return shim.Error(fmt.Sprintf("❌ No request found for visitor %s.", visitorId))
	}

	var request AccessRequest
	json.Unmarshal(requestBytes, &request)

	if request.Status == "Approved" {
		response := fmt.Sprintf(`{"status": "Approved", "qrCode": "%s"}`, request.QRCode)
		return shim.Success([]byte(response))
	}

	return shim.Success(requestBytes)
}

// ✅ Get Worker Request Status
func (rc *ResidentsContract) GetWorkerRequestStatus(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 1 {
        return shim.Error("❌ Worker ID is required.")
    }
    workerId := args[0]
    requestKey := "REQ_SERVICE_" + workerId

    requestBytes, err := stub.GetState(requestKey)
    if err != nil || requestBytes == nil {
        return shim.Error(fmt.Sprintf("❌ No request found for worker %s.", workerId))
    }

    var request AccessRequest
    json.Unmarshal(requestBytes, &request)

    if request.Status == "Approved" {
        // If approved, return QR code and status
        response := fmt.Sprintf(`{"status": "Approved", "qrCode": "%s"}`, request.QRCode)
        return shim.Success([]byte(response))
    }

    // If not approved (pending or rejected), return the full request data
    return shim.Success(requestBytes)
}
// ✅ Get Delivery Request Status
func (rc *ResidentsContract) GetDeliveryRequestStatus(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 1 {
        return shim.Error("❌ Delivery ID is required.")
    }
    deliveryId := args[0]
    requestKey := "REQ_DELIVERY_" + deliveryId

    requestBytes, err := stub.GetState(requestKey)
    if err != nil || requestBytes == nil {
        return shim.Error(fmt.Sprintf("❌ No request found for delivery %s.", deliveryId))
    }

    var request AccessRequest
    json.Unmarshal(requestBytes, &request)

    if request.Status == "Approved" {
        // If approved, return QR code and status
        response := fmt.Sprintf(`{"status": "Approved", "qrCode": "%s"}`, request.QRCode)
        return shim.Success([]byte(response))
    }

    // If not approved (pending or rejected), return the full request data
    return shim.Success(requestBytes)
}
func (rc *ResidentsContract) GetEmergencyRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("❌ Visitor ID is required.")
	}
	emergencyId := args[0]
	requestKey := "EMERGENCY_" + emergencyId

	bytes, err := stub.GetState(requestKey)
	if err != nil || bytes == nil {
		return shim.Error("❌ No emergency access found for this visitor.")
	}

	return shim.Success(bytes)
}


// ✅ Fonction pour récupérer tous les résidents
func (rc *ResidentsContract) GetAllResidents(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    // CouchDB query for all documents of type "resident"
    query := `{
        "selector": {
            "DocType": "resident"
        }
    }`

    // Execute query
    resultsIterator, err := stub.GetQueryResult(query)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Query failed: %s", err.Error()))
    }
    defer resultsIterator.Close()

    // Build array of residents
    var residents []Resident
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(fmt.Sprintf("❌ Failed to parse resident: %s", err.Error()))
        }

        var resident Resident
        err = json.Unmarshal(queryResponse.Value, &resident)
        if err != nil {
            return shim.Error(fmt.Sprintf("❌ Failed to unmarshal resident: %s", err.Error()))
        }
        residents = append(residents, resident)
    }

    // Marshal results to JSON
    residentsJSON, err := json.Marshal(residents)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal residents: %s", err.Error()))
    }

    return shim.Success(residentsJSON)
}

// ✅ التحقق مما إذا كان الزائر مسموحًا له بالدخول مباشرةً
func (rc *ResidentsContract) CheckVisitorAuthorization(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Resident ID and Visitor ID are required.")
	}

	userId, visitorId := args[0], args[1]

	// استرجاع بيانات المقيم
	residentBytes, err := stub.GetState("RESIDENT_" + userId)
	//residentBytes, err := stub.GetPrivateData("ResidentPrivateCollection", "RESIDENT_" + residentId)
	if err != nil || residentBytes == nil {
		return shim.Error(fmt.Sprintf("❌ Resident %s not found.", userId))
	}

	// تحويل البيانات من JSON إلى كائن Go
	var resident Resident
	json.Unmarshal(residentBytes, &resident)

	// التحقق مما إذا كان الزائر موجودًا في القائمة
	if containsVisitorInfo(resident.Visitors, visitorId) {
		return shim.Success([]byte(fmt.Sprintf("✅ Visitor %s is authorized to enter.", visitorId)))
	}

	return shim.Error(fmt.Sprintf("❌ Visitor %s is NOT authorized to enter. Request approval is required.", visitorId))
}
// SaveLogToChain invokes a transaction to save log
func (rc *ResidentsContract) SaveLogToChain(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) != 3 {
        return shim.Error("❌ Required: RequestID, ActionType, Timestamp")
    }

    logId := fmt.Sprintf("LOG_%s_%d", args[0], time.Now().UnixNano())

    log := struct {
        LogID      string `json:"logId"`
        RequestID string `json:"requestId"`
        Type       string `json:"type"`
        Timestamp  int64  `json:"timestamp"`
    }{
        LogID:      logId,
        RequestID: args[0],
        Type:       args[1],
    }

    // Use timestamp from args instead of current time
    parsedTime, err := strconv.ParseInt(args[2], 10, 64)
    if err != nil {
        return shim.Error("❌ Invalid timestamp format")
    }
    log.Timestamp = parsedTime

    logBytes, err := json.Marshal(log)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to marshal log: %s", err.Error()))
    }

    err = stub.PutState(logId, logBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to save log to chain: %s", err.Error()))
    }

    return shim.Success([]byte(logId))
}

func (rc *ResidentsContract) GetLastLogByResident(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) != 1 {
        return shim.Error("❌ Required: ResidentID")
    }

    residentID := args[0]

    // Build CouchDB rich query (assuming CouchDB state database)
    queryString := fmt.Sprintf(`{
        "selector": {
            "residentId": "%s"
        },
        "sort": [{"timestamp": "desc"}],
        "limit": 1
    }`, residentID)

    resultsIterator, err := stub.GetQueryResult(queryString)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to query logs: %s", err.Error()))
    }
    defer resultsIterator.Close()

    if !resultsIterator.HasNext() {
        return shim.Error("❌ No logs found for the given resident ID")
    }

    result, err := resultsIterator.Next()
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read query result: %s", err.Error()))
    }

    // Parse the log entry
    var logEntry struct {
        Type string `json:"type"`
    }
    err = json.Unmarshal(result.Value, &logEntry)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to parse log entry: %s", err.Error()))
    }

    return shim.Success([]byte(logEntry.Type))
}
func (rc *ResidentsContract) GetVisitRequest(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 1 {
        return shim.Error("❌ Required: RequestID")
    }

    requestId := args[0]
    requestKey := fmt.Sprintf("VISITREQUEST_%s", requestId)

    requestBytes, err := stub.GetState(requestKey)
    if err != nil {
        return shim.Error(fmt.Sprintf("❌ Failed to read visit request: %s", err.Error()))
    }

    if requestBytes == nil {
        return shim.Error("❌ Visit request not found")
    }

    return shim.Success(requestBytes)
}


func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	err := shim.Start(new(ResidentsContract))
	if err != nil {
		fmt.Printf("❌ Error starting ResidentsContract: %s", err)
	}
}
