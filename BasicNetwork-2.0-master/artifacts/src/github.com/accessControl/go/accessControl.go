package main
/*demain nchalleh n3awed min lawel psq lyoum filil dart just modification fi lcode */
import (
	"encoding/json"
	"fmt"
	"time"
    

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

type AccessControlContract struct{}

type AccessResponse struct {
	Access string `json:"access"`
	Reason string `json:"reason"`
}

type AccessRequest struct {
	VisitorID  string   `json:"visitorId"`
	Status     string   `json:"status"`
	ApprovedBy []string `json:"approvedBy"`
	ExpiryTime string   `json:"expiryTime,omitempty"`
	QRCode     string   `json:"qrCode,omitempty"`  // Add QRCode field
	Decision   string   `json:"decision,omitempty"` // ✅ Add this line
    QRCodeTimestamp  int64   `json:"qrCodeTimestamp"` // 👈 أضف هذا السطر
}

func (c *AccessControlContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
/*demain nchalleh n3awed deploy psq badalt fil code*/
func (c *AccessControlContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "invokeResidentsContract" :
		// This is where you need to call your ResidentsContract logic
		return c.invokeResidentsContract(stub, "GetAllResidents", args)
	case "checkAccessResidents":
		// Step 1: Check if it's a resident
		return c.checkAccessResidents(stub, args)
	case "checkAccessVisitorsList":
		// Step 2: Check if the visitor is listed under a resident
		return c.checkAccessVisitorsList(stub, args)
	case "checkAccessRequestVisitor":
		// Step 3: Check if the visitor's request is approved and QR code exists
		return c.checkAccessRequestVisitor(stub, args)
	case "checkAccessRequestWorker":
		return c.checkAccessRequestWorker(stub, args)
	case "checkAccessRequestDelivery":
		return c.checkAccessRequestDelivery(stub, args)
	case "checkEmergencyAccess":
		return c.checkEmergencyAccess(stub, args)
	//case "ForceExpiryTimeService":
	//	return c.ForceExpiryTimeService(stub, args) // ✅ هذا هو السطر الجديد
	//case "ForceExpiryTimeDelivery":
		//return c.ForceExpiryTimeDelivery(stub, args) // ✅ هذا هو السطر الجديد
	//case "GetAllDecisions":
	    //return c.GetAllDecisions(stub)
	default:
		return shim.Error("Invalid function name.")
	}
}

func (c *AccessControlContract) invokeResidentsContract(stub shim.ChaincodeStubInterface, function string, args []string) sc.Response {
    // تحديد القناة التي تحتوي على عقد "ResidentsContract"
    residentsChannel := "residentschannel" // يمكنك تغييرها إلى القناة الصحيحة

    queryArgs := [][]byte{[]byte(function)}
    for _, arg := range args {
        queryArgs = append(queryArgs, []byte(arg))
    }

    // استدعاء عقد "ResidentsContract" عبر القناة المحددة
    response := stub.InvokeChaincode("residentManagement", queryArgs, residentsChannel)
    if response.Status != shim.OK {
        return shim.Error(fmt.Sprintf("Error invoking residentsContract: %s", response.Message))
    }

    return response
}

func (c *AccessControlContract) invokeCheckAuthorizationContract(stub shim.ChaincodeStubInterface, function string, residentId string, visitorId string) sc.Response {
    // تحديد القناة التي تحتوي على عقد "ResidentsContract"
    residentsChannel := "residentschannel" // يمكنك تغييرها إلى القناة الصحيحة

    // تحضير الوسائط المطلوبة لاستدعاء وظيفة "CheckVisitorAuthorization"
    queryArgs := [][]byte{[]byte(function), []byte(residentId), []byte(visitorId)}

    // استدعاء عقد "ResidentsContract" عبر القناة المحددة
    response := stub.InvokeChaincode("residentManagement", queryArgs, residentsChannel)
    if response.Status != shim.OK {
        return shim.Error(fmt.Sprintf("Error invoking residentsContract: %s", response.Message))
    }

    return response
}



/*func (c *AccessControlContract) CheckAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Validate input parameters
	if len(args) < 2 {
		return shim.Error("Visitor ID and NFC UID are required.")
	}
	id := args[0]       // Can be the Resident ID or Visitor ID
	uid := args[1]      // UID for NFC card

	// ================================
	// ✅ First step: Check if it's a resident
	// ================================
	residentBytes, err := stub.GetState("RESIDENT_" + id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get resident data: %s", err.Error()))
	}
	if residentBytes != nil {
		var resident struct {
			NFCUID string `json:"nfcUid"`
		}
		json.Unmarshal(residentBytes, &resident)

		if resident.NFCUID == uid {
			// ✅ UID is correct — grant access
			response := AccessResponse{"Granted", "Resident authenticated via NFC. Access granted."}
			respBytes, _ := json.Marshal(response)
			
			// Store the decision in the ledger
			decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
			stub.PutState(decisionKey, respBytes)
			
			return shim.Success(respBytes)
		} else {
			// ❌ Incorrect UID — deny access
			response := AccessResponse{"Denied", "Resident UID does not match. Access denied."}
			respBytes, _ := json.Marshal(response)
			
			// Store the decision in the ledger
			decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
			stub.PutState(decisionKey, respBytes)
			
			return shim.Success(respBytes)
		}
	}

	// ================================
	// ✅ Second step: Check if the visitor is approved by a resident
	// ================================
	// Querying all residents and their approved visitors
	iterator, _ := stub.GetStateByRange("RESIDENT_", "RESIDENT_~")
	for iterator.HasNext() {
		queryResponse, _ := iterator.Next()
		var resident struct {
			Visitors []string `json:"visitors"`
		}
		json.Unmarshal(queryResponse.Value, &resident)

		for _, visitor := range resident.Visitors {
			if visitor == id {
				// ✅ Visitor approved by one of the residents
				response := AccessResponse{"Granted", "Visitor approved by resident"}
				respBytes, _ := json.Marshal(response)
				
				// Store the decision in the ledger
				decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
				stub.PutState(decisionKey, respBytes)
				
				return shim.Success(respBytes)
			}
		}
	}

	// ================================
	// ✅ Third step: Check previous access requests for the visitor
	// ================================
	// Querying the access request status
	requestResponse := c.invokeResidentsContract(stub, "GetRequestStatus", []string{id})
	if requestResponse.Status != shim.OK {
		return requestResponse
	}

	var request AccessRequest
	json.Unmarshal(requestResponse.Payload, &request)

	switch request.Status {
	case "Pending":
		// ⏳ Request under review
		response := AccessResponse{"Pending", "Your request is under review. Please wait for approval."}
		respBytes, _ := json.Marshal(response)
		
		// Store the decision in the ledger
		decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
		stub.PutState(decisionKey, respBytes)
		
		return shim.Success(respBytes)

	case "Approved":
		// ✅ Request approved — check for QR code
		qrCodeKey := "QR_" + id
		qrCodeBytes, _ := stub.GetState(qrCodeKey)
		if qrCodeBytes != nil {
			// QR code verified
			response := AccessResponse{"Granted", "Visitor approved and QR code verified. Access granted."}
			respBytes, _ := json.Marshal(response)
			
			// Store the decision in the ledger
			decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
			stub.PutState(decisionKey, respBytes)
			
			return shim.Success(respBytes)
		} else {
			// ❌ No QR code available
			response := AccessResponse{"Denied", "QR code is invalid or not available."}
			respBytes, _ := json.Marshal(response)
			
			// Store the decision in the ledger
			decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
			stub.PutState(decisionKey, respBytes)
			
			return shim.Success(respBytes)
		}

	case "Rejected":
		// ❌ Request rejected
		response := AccessResponse{"Denied", "Your request was rejected. Please submit a new request."}
		respBytes, _ := json.Marshal(response)
		
		// Store the decision in the ledger
		decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
		stub.PutState(decisionKey, respBytes)
		
		return shim.Success(respBytes)
	}

	// ================================
	// ❌ No matching permission found
	// ================================
	response := AccessResponse{"Denied", "No valid permission found."}
	respBytes, _ := json.Marshal(response)
	
	// Store the decision in the ledger
	decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
	stub.PutState(decisionKey, respBytes)
	
	return shim.Success(respBytes)
}*/
/*func (c *AccessControlContract) CheckAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Validate input parameters
	if len(args) < 2 {
		return shim.Error("Visitor ID and NFC UID are required.")
	}
	id := args[0]       // Can be the Resident ID or Visitor ID
	uid := args[1]      // UID for NFC card


	var residentId string
	if len(args) >= 3 {
		residentId = args[2] // Provided resident ID for visitor authorization check
	}

	// ================================
	// ✅ First step: Check if it's a resident
	// ================================
	// Call GetResident function from the ResidentsContract
	residentResponse := c.invokeResidentsContract(stub, "GetResident", []string{id})
	if residentResponse.Status != shim.OK {
		return shim.Error(fmt.Sprintf("Failed to get resident data: %s", residentResponse.Message))
	}

	var resident struct {
		NFCUID string `json:"uid"`
	}
	if err := json.Unmarshal(residentResponse.Payload, &resident); err != nil {
		return shim.Error("❌ Failed to parse resident data.")
	}

	// Check if UID matches
	if resident.NFCUID == uid {
		// ✅ UID is correct — grant access
		response := AccessResponse{"Granted", "Resident authenticated via NFC. Access granted."}
		respBytes, _ := json.Marshal(response)
		
		// Store the decision in the ledger
		decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
		stub.PutState(decisionKey, respBytes)
		
		return shim.Success(respBytes)
	} else {
		// ❌ Incorrect UID — deny access
		response := AccessResponse{"Denied", "Resident UID does not match. Access denied."}
		respBytes, _ := json.Marshal(response)
		
		// Store the decision in the ledger
		decisionKey := "DECISION_" + id + "_" + time.Now().Format(time.RFC3339)
		stub.PutState(decisionKey, respBytes)
		
		return shim.Success(respBytes)
	}

   
   // ================================
	// ✅ Second step: Check if visitor is listed under a specific resident
	// ================================
	if residentId != "" {
		authResponse := c.invokeCheckAuthorizationContract(stub, "CheckVisitorAuthorization", residentId, id)
		if authResponse.Status == shim.OK {
			accessResponse := AccessResponse{"Granted", "Visitor is authorized by resident " + residentId + "."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		}else {
			accessResponse := AccessResponse{"Denied", "Visitor is not authorized by resident " + residentId + "."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_" + id + "_" + time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		}
	}
	//================================
    // ✅ Step 3: Check if request is approved and QR exists
    requestResponse := c.invokeResidentsContract(stub, "GetRequestStatus", []string{id})
    if requestResponse.Status != shim.OK {
        return requestResponse
    }

    var request AccessRequest
    if err := json.Unmarshal(requestResponse.Payload, &request); err != nil {
        return shim.Error("❌ Failed to parse request status.")
    }

    switch request.Status {
    case "Pending":
        accessResponse := AccessResponse{"Pending", "Your request is under review."}
        respBytes, _ := json.Marshal(accessResponse)
        stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
        return shim.Success(respBytes)

    case "Approved":
        if request.QRCode != "" {
            accessResponse := AccessResponse{"Granted", "Visitor approved and QR code verified."}
            respBytes, _ := json.Marshal(accessResponse)
            stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
            return shim.Success(respBytes)
        } else {
            accessResponse := AccessResponse{"Denied", "QR code is invalid or not available."}
            respBytes, _ := json.Marshal(accessResponse)
            stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
            return shim.Success(respBytes)
        }

    case "Rejected":
        accessResponse := AccessResponse{"Denied", "Your request was rejected."}
        respBytes, _ := json.Marshal(accessResponse)
        stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
        return shim.Success(respBytes)
    }

    accessResponse := AccessResponse{"Denied", "No valid permission found."}
    respBytes, _ := json.Marshal(accessResponse)
    stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
    return shim.Success(respBytes)
}*/
func (c *AccessControlContract) checkAccessResidents(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Visitor ID and QR Code are required.")
	}

	visitorId := args[0]  // Visitor ID or Resident ID
	qrCode := args[1]     // QR Code provided for authentication

	// Call GetResident function from the ResidentsContract
	residentResponse := c.invokeResidentsContract(stub, "GetResident", []string{visitorId})
	if residentResponse.Status != shim.OK {
		return shim.Error(fmt.Sprintf("Failed to get resident data: %s", residentResponse.Message))
	}

	var Resident struct {
		ResidentID string `json:"residentId"`  // This is the QR Code stored during registration
	}
	if err := json.Unmarshal(residentResponse.Payload, &Resident); err != nil {
		return shim.Error("❌ Failed to parse resident data.")
	}

	// Check if QR Code matches the ResidentID
	if Resident.ResidentID == qrCode {
		response := AccessResponse{"Granted", "Resident authenticated via QR code. Access granted."}
		respBytes, _ := json.Marshal(response)
		decisionKey := "DECISION_" + visitorId + "_" + time.Now().Format(time.RFC3339)
		stub.PutState(decisionKey, respBytes)
		return shim.Success(respBytes)
	} else {
		response := AccessResponse{"Denied", "QR Code does not match. Access denied."}
		respBytes, _ := json.Marshal(response)
		decisionKey := "DECISION_" + visitorId + "_" + time.Now().Format(time.RFC3339)
		stub.PutState(decisionKey, respBytes)
		return shim.Success(respBytes)
	}
}

func (c *AccessControlContract) checkAccessVisitorsList(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 3 {
		return shim.Error("❌ Required: Visitor ID, QR Code, Resident ID.")
	}

	visitorId := args[0]
	qrCode := args[1]//qrcode permanenets of visitors//
	userId := args[2]//name of resident//

	// Step 1: Check if the visitor is authorized by the resident
	authResponse := c.invokeCheckAuthorizationContract(stub, "CheckVisitorAuthorization", userId, visitorId)
	if authResponse.Status != shim.OK {
		accessResponse := AccessResponse{"Denied", "🚫 Visitor is not authorized by resident " + userId + "."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+visitorId+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}

	// Step 2: Fetch resident record from ResidentsContract
	residentResponse := c.invokeResidentsContract(stub, "GetResident", []string{userId})
	if residentResponse.Status != shim.OK {
		return shim.Error("❌ Failed to retrieve resident data.")
	}

	var Resident struct {
		Visitors []struct {
			VisitorId string `json:"visitorId"`
			QRCode    string `json:"qrCode"`
		} `json:"visitors"`
	}

	if err := json.Unmarshal(residentResponse.Payload, &Resident); err != nil {
		return shim.Error("❌ Failed to parse resident data.")
	}

	// Step 3: Check if visitor exists in resident's visitor list and QRCode matches
	for _, v := range Resident.Visitors {
		if v.VisitorId == visitorId && v.QRCode == qrCode {
			accessResponse := AccessResponse{"Granted", "✅ Visitor is authorized and QR code matches. Access granted."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+visitorId+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		}
	}

	// If no match found
	accessResponse := AccessResponse{"Denied", "🚫 Visitor QR code does not match or is not listed under resident " + userId + "."}
	respBytes, _ := json.Marshal(accessResponse)
	stub.PutState("DECISION_"+visitorId+"_"+time.Now().Format(time.RFC3339), respBytes)
	return shim.Success(respBytes)
}

/*for delete espace just  -c "{\"function\": \"checkAccessVisitorsList\", \"Args\":[\"$visitor_id\", \"\", \"$res_id\"]}" here*/
/*func (c *AccessControlContract) checkAccessVisitorsList(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 2 { // Adjusting to expect 2 arguments instead of 3
        return shim.Error("Resident ID and Visitor ID are required.")
    }

    id := args[0]         // Visitor ID
    residentId := args[1] // Resident ID

    // Check if visitor is listed under a specific resident
    authResponse := c.invokeCheckAuthorizationContract(stub, "CheckVisitorAuthorization", residentId, id)
    if authResponse.Status == shim.OK {
        accessResponse := AccessResponse{"Granted", "Visitor is authorized by resident " + residentId + "."}
        respBytes, _ := json.Marshal(accessResponse)
        stub.PutState("DECISION_"+id+"_"+time.Now().Format(time.RFC3339), respBytes)
        return shim.Success(respBytes)
    } else {
        accessResponse := AccessResponse{"Denied", "Visitor is not authorized by resident " + residentId + "."}
        respBytes, _ := json.Marshal(accessResponse)
        stub.PutState("DECISION_" + id + "_" + time.Now().Format(time.RFC3339), respBytes)
        return shim.Success(respBytes)
    }
}*/
/* i add it now but i try tomorrow nchalleh*/
func (c *AccessControlContract) checkAccessRequestVisitor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		// Ensure both Visitor ID and QR code are passed as arguments
		return shim.Error("Both Visitor ID and QR code are required.")
	}

	visitorID := args[0] // Visitor ID
	enteredQRCode := args[1] // Entered QR code by the visitor

	// Get the request status for the visitor (from ResidentsContract)
	requestResponse := c.invokeResidentsContract(stub, "GetRequestStatus", []string{visitorID})
	if requestResponse.Status != shim.OK {
		// Error response from the ResidentsContract
		return requestResponse
	}

	// Parse the response payload to extract the QR code and request status
	var request struct {
		VisitorID string `json:"visitorId"`
		Status    string `json:"status"`
		QRCode    string `json:"qrCode"`
	}
	if err := json.Unmarshal(requestResponse.Payload, &request); err != nil {
		return shim.Error("❌ Failed to parse request status: " + err.Error())
	}
//this informations exist in peers of this smart contract//
	// Log for debugging purposes
	fmt.Println("Visitor ID:", visitorID)
	fmt.Println("Status:", request.Status)
	fmt.Println("Stored QR Code:", request.QRCode)
	fmt.Println("Entered QR Code:", enteredQRCode)

	// Process the request based on its status
	switch request.Status {
	case "Pending":
		accessResponse := AccessResponse{"Pending", "Your request is under review."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)

	case "Approved":
		/*if request.QRCode != "" {*/
			// Compare the QR code from the request with the entered QR code (provided by the visitor)
			if request.QRCode == enteredQRCode {
				// If QR codes match
				accessResponse := AccessResponse{"Granted", "Visitor approved and QR code verified."}
				respBytes, _ := json.Marshal(accessResponse)
				stub.PutState("DECISION_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
				return shim.Success(respBytes)
			} else {
				// If QR codes do not match
				accessResponse := AccessResponse{"Denied", "QR code does not match the one stored in the ledger."}
				respBytes, _ := json.Marshal(accessResponse)
				stub.PutState("DECISION_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
				return shim.Success(respBytes)
			}
		/*} else {
			// If the QR code is not available or invalid
			accessResponse := AccessResponse{"Denied", "QR code is invalid or not available."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		}*/

	case "Rejected":
		// If the request was rejected
		accessResponse := AccessResponse{"Denied", "Your request was rejected."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}

	// Default response if no valid permission found
	accessResponse := AccessResponse{"Denied", "No valid permission found."}
	respBytes, _ := json.Marshal(accessResponse)
	stub.PutState("DECISION_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
	return shim.Success(respBytes)
}

// ✅ Check access for service worker by verifying QR Code
func (c *AccessControlContract) checkAccessRequestWorker(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("❌ Both Worker ID and QR code are required.")
	}

	workerID := args[0]         // Worker ID
	enteredQRCode := args[1]    // QR code provided by the worker

	// Invoke ResidentsContract to get the worker request status
	requestResponse := c.invokeResidentsContract(stub, "GetWorkerRequestStatus", []string{workerID})
	if requestResponse.Status != shim.OK {
		return shim.Error("❌ Failed to get worker request status: " + requestResponse.Message)
	}

	// Parse the response payload to extract QR code and status
	var request struct {

		Status string `json:"status"`
		QRCode string `json:"qrCode"`
	}

	if err := json.Unmarshal(requestResponse.Payload, &request); err != nil {
		return shim.Error("❌ Failed to parse request status: " + err.Error())
	}

	// Debug log
	fmt.Println("Worker ID:", workerID)
	fmt.Println("Status:", request.Status)
	fmt.Println("Stored QR Code:", request.QRCode)
	fmt.Println("Entered QR Code:", enteredQRCode)

	// Process based on status
	switch request.Status {
	case "Pending":
		accessResponse := AccessResponse{"Pending", "Your request is under manager review."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+workerID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)

	case "Approved":
		if request.QRCode == enteredQRCode {
			accessResponse := AccessResponse{"Granted", "Worker approved and QR code verified."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+workerID+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		} else {
			accessResponse := AccessResponse{"Denied", "QR code does not match."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+workerID+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		}

	case "Rejected":
		accessResponse := AccessResponse{"Denied", "Your request was rejected by the manager."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+workerID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}

	// Default case
	accessResponse := AccessResponse{"Denied", "No valid access request found."}
	respBytes, _ := json.Marshal(accessResponse)
	stub.PutState("DECISION_"+workerID+"_"+time.Now().Format(time.RFC3339), respBytes)
	return shim.Success(respBytes)
}

func (c *AccessControlContract) checkAccessRequestDelivery(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		// Ensure both Delivery ID and QR code are passed as arguments
		return shim.Error("❌ Both Delivery ID and QR code are required.")
	}

	deliveryID := args[0]       // Delivery ID
	enteredQRCode := args[1]    // QR code entered by the delivery worker

	// Get the delivery request status from ResidentsContract
	requestResponse := c.invokeResidentsContract(stub, "GetDeliveryRequestStatus", []string{deliveryID})
	if requestResponse.Status != shim.OK {
		// Error response from the ResidentsContract
		return requestResponse
	}

	// Parse the response payload to extract the QR code and request status
	var request struct {
		Status string `json:"status"`
		QRCode string `json:"qrCode"`
	}

	if err := json.Unmarshal(requestResponse.Payload, &request); err != nil {
		return shim.Error("❌ Failed to parse delivery request status: " + err.Error())
	}

	// Log for debugging purposes
	fmt.Println("Delivery ID:", deliveryID)
	fmt.Println("Status:", request.Status)
	fmt.Println("Stored QR Code:", request.QRCode)
	fmt.Println("Entered QR Code:", enteredQRCode)

	// Process the request based on its status
	switch request.Status {
	case "Pending":
		accessResponse := AccessResponse{"Pending", "Your delivery access request is under review."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+deliveryID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)

	case "Approved":
		if request.QRCode == enteredQRCode {
			accessResponse := AccessResponse{"Granted", "Delivery access granted. QR code verified."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+deliveryID+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		} else {
			accessResponse := AccessResponse{"Denied", "QR code does not match the one stored in the ledger."}
			respBytes, _ := json.Marshal(accessResponse)
			stub.PutState("DECISION_"+deliveryID+"_"+time.Now().Format(time.RFC3339), respBytes)
			return shim.Success(respBytes)
		}

	case "Rejected":
		accessResponse := AccessResponse{"Denied", "Your delivery access request was rejected."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+deliveryID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}

	// Default response if status is not recognized
	accessResponse := AccessResponse{"Denied", "No valid delivery access permission found."}
	respBytes, _ := json.Marshal(accessResponse)
	stub.PutState("DECISION_"+deliveryID+"_"+time.Now().Format(time.RFC3339), respBytes)
	return shim.Success(respBytes)
}
 //grant
func (c *AccessControlContract) checkEmergencyAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("❌ Visitor ID is required.")
	}

	visitorID := args[0]

	// Build the emergency request key
	//requestKey := "EMERGENCY_" + visitorID

	// Call the first contract to fetch the emergency access record
	response := c.invokeResidentsContract(stub, "GetEmergencyRequest", []string{visitorID})

	if response.Status != shim.OK {
		return shim.Error("❌ Failed to fetch emergency access info.")
	}

	if len(response.Payload) == 0 {
		return shim.Error("❌ No emergency access record found for this visitor.")
	}

	// Parse the emergency access record
	var request AccessRequest
	err := json.Unmarshal(response.Payload, &request)
	if err != nil {
		return shim.Error("❌ Failed to parse emergency access record: " + err.Error())
	}

	// If approved, allow entry
	if request.Status == "Approved" {
		accessResponse := AccessResponse{"Granted", "✅ Emergency access verified."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_EMERGENCY_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}

	// Otherwise, deny
	accessResponse := AccessResponse{"Denied", "❌ Visitor does not have emergency access."}
	respBytes, _ := json.Marshal(accessResponse)
	stub.PutState("DECISION_EMERGENCY_"+visitorID+"_"+time.Now().Format(time.RFC3339), respBytes)
	return shim.Success(respBytes)
}


/*func (c *AccessControlContract) GrantEmergencyAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Visitor ID and Manager ID are required.")
	}
	visitorId, managerId := args[0], args[1]
	requestKey := "REQ_EMERGENCY_" + visitorId

	// Create the access request with the decision field added
	request := AccessRequest{
		VisitorID:  visitorId,
		Status:     "Approved",
		ApprovedBy: []string{managerId},
		Decision:   fmt.Sprintf("Approved by %s for emergency access", managerId), // Storing the decision in the ledger
	}

	// Convert the access request to JSON
	requestBytes, _ := json.Marshal(request)

	// Store the request in the ledger
	stub.PutState(requestKey, requestBytes)

	// Return success with the decision in the response message
	return shim.Success([]byte(fmt.Sprintf("Emergency access granted for %s. Decision: %s", visitorId, request.Decision)))
}

//demain nchalleh nzid partit ta3 QRLimited with Time//
/*func (c *AccessControlContract) GrantServiceAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) < 3 {
        return shim.Error("Worker ID, Grantor ID, and Expiry Minutes are required.")
    }
    
    workerId, grantorId := args[0], args[1]
    expiryMinutes, _ := time.ParseDuration(args[2] + "m")
    expiryTime := time.Now().Add(expiryMinutes).Format(time.RFC3339)
    requestKey := "REQ_SERVICE_" + workerId

    // Check if there is an existing request for the worker
    existingRequestBytes, err := stub.GetState(requestKey)
    if err != nil {
        return shim.Error("Failed to get existing request")
    }

    // If an existing request exists, handle expiry check
    if existingRequestBytes != nil {
        var existingRequest AccessRequest
        err := json.Unmarshal(existingRequestBytes, &existingRequest)
        if err != nil {
            return shim.Error("Failed to parse existing request")
        }

        // Parse the expiry time from the existing request
        expiryTimeParsed, err := time.Parse(time.RFC3339, existingRequest.ExpiryTime)
        if err != nil {
            return shim.Error("Invalid expiry time format in existing request")
        }

        // If the current time is after the expiry time, deny the access
        if time.Now().After(expiryTimeParsed) {
            existingRequest.Status = "Denied"
            existingRequest.Decision = fmt.Sprintf("Denied: expiry time passed at %s", existingRequest.ExpiryTime)

            // Update the request in the ledger
            updatedRequestBytes, _ := json.Marshal(existingRequest)
            err = stub.PutState(requestKey, updatedRequestBytes)
            if err != nil {
                return shim.Error("Failed to update access request")
            }

            return shim.Error(fmt.Sprintf("Access denied: expiry time passed at %s", existingRequest.ExpiryTime))
        }
    }

    // Create a new request if no request exists, or if the previous request was denied.
    request := AccessRequest{
        VisitorID:  workerId,
        ApprovedBy: []string{grantorId},
        Status:     "Approved",
        ExpiryTime: expiryTime,
        Decision:   fmt.Sprintf("Approved by %s until %s", grantorId, expiryTime),
    }

    // Save the new request in the ledger
    requestBytes, _ := json.Marshal(request)
    err = stub.PutState(requestKey, requestBytes)
    if err != nil {
        return shim.Error("Failed to store access request")
    }

    return shim.Success([]byte(fmt.Sprintf("Temporary access granted for service worker %s until %s", workerId, expiryTime)))
}*/
// Function to generate a unique QR code with timestamp
/*func generateQRCodeWorker(workerId string) string {
	timestamp := time.Now().Unix()
	qrCode := fmt.Sprintf("QR-%s-%d", workerId, timestamp)
	return qrCode
}


func (c *AccessControlContract) GrantServiceAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 3 {
		return shim.Error("Worker ID, Grantor ID, and Expiry Minutes are required.")
	}

	workerId := args[0]
	grantorId := args[1]
	expiryMinutes, err := time.ParseDuration(args[2] + "m")
	if err != nil {
		return shim.Error("Invalid expiry time format.")
	}

	requestKey := "REQ_SERVICE_" + workerId
	currentTime := time.Now()

	// Check for existing request
	existingBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}

	if existingBytes != nil {
		var existing AccessRequest
		if err := json.Unmarshal(existingBytes, &existing); err != nil {
			return shim.Error("Failed to parse existing request")
		}

		// Check if expired
		expiryTime, err := time.Parse(time.RFC3339, existing.ExpiryTime)
		if err != nil {
			return shim.Error("Invalid expiry time format in existing request")
		}

		if currentTime.After(expiryTime) {
			// Mark old request as denied
			existing.Status = "Rejected"
			existing.Decision = fmt.Sprintf("Rejected: expired at %s", existing.ExpiryTime)
			oldBytes, _ := json.Marshal(existing)
			stub.PutState(requestKey+"_EXPIRED", oldBytes) // Save old with a different key
			stub.DelState(requestKey)                      // Delete the old one
		} else {
			// Still valid — no need to regenerate
			return shim.Success([]byte(fmt.Sprintf("Access already granted for %s until %s", workerId, existing.ExpiryTime)))
		}
	}

	// Create new QR code and timestamps using generateQRCode
	newExpiry := currentTime.Add(expiryMinutes).Format(time.RFC3339)
	qrCode := generateQRCodeWorker(workerId)
	qrTimestamp := currentTime.Format(time.RFC3339)

	request := AccessRequest{
		VisitorID:        workerId,
		ApprovedBy:       []string{grantorId},
		Status:           "Approved",
		ExpiryTime:       newExpiry,
		Decision:         fmt.Sprintf("Approved by %s until %s", grantorId, newExpiry),
		QRCode:           qrCode,
		QRCodeTimestamp:  qrTimestamp,
	}

	requestBytes, _ := json.Marshal(request)
	if err := stub.PutState(requestKey, requestBytes); err != nil {
		return shim.Error("Failed to store access request")
	}

	return shim.Success([]byte(fmt.Sprintf("Temporary QR access granted to %s until %s", workerId, newExpiry)))
}


//just for test//
func (c *AccessControlContract) ForceExpiryTimeService(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Worker ID and Forced Expiry Minutes are required.")
	}

	workerId := args[0]
	forcedExpiryMinutes, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid expiry time format.")
	}

	requestKey := "REQ_SERVICE_" + workerId
	currentTime := time.Now()

	// جلب الطلب الحالي
	requestBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}

	if requestBytes == nil {
		return shim.Error("No access request found for the worker ID")
	}

	var request AccessRequest
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		return shim.Error("Failed to parse existing request")
	}

	// تعديل وقت الانتهاء ليكون في الماضي (للمحاكاة فقط)
	request.ExpiryTime = currentTime.Add(-time.Duration(forcedExpiryMinutes) * time.Minute).Format(time.RFC3339)

	// إعادة حفظ الطلب بعد تعديل وقت الانتهاء فقط
	updatedBytes, _ := json.Marshal(request)
	err = stub.PutState(requestKey, updatedBytes)
	if err != nil {
		return shim.Error("Failed to update expiry time")
	}

	return shim.Success([]byte(fmt.Sprintf("Expiry time forcibly set for worker %s to %s", workerId, request.ExpiryTime)))
}

func (c *AccessControlContract) checkAccessOfServiceGrant(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Worker ID and QR Code are required.")
	}

	workerId := args[0]
	inputQRCode := args[1]

	// بناء المفتاح للتحقق من الطلب
	requestKey := "REQ_SERVICE_" + workerId

	// جلب الطلب من السجل
	existingBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}

	if existingBytes == nil {
		return shim.Error("No access request found for the worker ID")
	}

	var existingRequest AccessRequest
	if err := json.Unmarshal(existingBytes, &existingRequest); err != nil {
		return shim.Error("Failed to parse existing request")
	}

	// التحقق مما إذا كان QR Code المدخل يتطابق مع الذي في السجل
	if inputQRCode == existingRequest.QRCode {
		// If QR codes match
		accessResponse := AccessResponse{"Granted", "Worker approved and QR code verified."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+workerId+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	} else {
		// If QR codes do not match
		accessResponse := AccessResponse{"Denied", "QR code does not match the one stored in the ledger."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+workerId+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}
}

/*func (c *AccessControlContract) GrantDeliveryAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Resident ID and Delivery ID are required.")
	}
	residentId, deliveryId := args[0], args[1]
	expiryTime := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
	requestKey := "REQ_DELIVERY_" + deliveryId
	request := AccessRequest{VisitorID: deliveryId, ApprovedBy: []string{residentId}, Status: "Approved", ExpiryTime: expiryTime}
	requestBytes, _ := json.Marshal(request)
	stub.PutState(requestKey, requestBytes)
	return shim.Success([]byte(fmt.Sprintf("Temporary access granted for delivery %s until %s", deliveryId, expiryTime)))
}*/
/*func (c *AccessControlContract) GrantDeliveryAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Resident ID and Delivery ID are required.")
	}

	residentId, deliveryId := args[0], args[1]
	requestKey := "REQ_DELIVERY_" + deliveryId

	// التحقق إذا كان هناك طلب سابق
	existingRequestBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}
	if existingRequestBytes != nil {
		var existingRequest AccessRequest
		err := json.Unmarshal(existingRequestBytes, &existingRequest)
		if err != nil {
			return shim.Error("Failed to parse existing request")
		}

		// تحويل وقت الانتهاء إلى time.Time
		expiryTimeParsed, err := time.Parse(time.RFC3339, existingRequest.ExpiryTime)
		if err != nil {
			return shim.Error("Invalid expiry time format in existing request")
		}

		// مقارنة الوقت الحالي مع expiryTime
		if time.Now().After(expiryTimeParsed) {
			// تغيير الحالة إلى "مرفوض" إذا تجاوز الوقت
			existingRequest.Status = "Denied"
			existingRequest.Decision = fmt.Sprintf("Denied: expiry time passed at %s", existingRequest.ExpiryTime)
			
			// تحديث السجل في الـ ledger
			updatedRequestBytes, _ := json.Marshal(existingRequest)
			err = stub.PutState(requestKey, updatedRequestBytes)
			if err != nil {
				return shim.Error("Failed to update access request")
			}

			return shim.Error(fmt.Sprintf("Access denied: expiry time passed at %s", existingRequest.ExpiryTime))
		}
	}

	// إنشاء إذن جديد صالح لمدة 30 دقيقة من الآن
	expiryTime := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
	request := AccessRequest{
		VisitorID:  deliveryId,
		ApprovedBy: []string{residentId},
		Status:     "Approved",
		ExpiryTime: expiryTime,
		Decision:   fmt.Sprintf("Approved by %s until %s", residentId, expiryTime),
	}
	requestBytes, _ := json.Marshal(request)
	err = stub.PutState(requestKey, requestBytes)
	if err != nil {
		return shim.Error("Failed to store access request")
	}

	return shim.Success([]byte(fmt.Sprintf("Temporary access granted for delivery %s until %s", deliveryId, expiryTime)))
}
// Function to generate a unique QR code with timestamp
func generateQRCode(deliveryId string) string {
	timestamp := time.Now().Unix()
	qrCode := fmt.Sprintf("QR-%s-%d",deliveryId, timestamp)
	return qrCode
}
func (c *AccessControlContract) GrantDeliveryAccess(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Resident ID and Delivery ID are required.")
	}

	residentId, deliveryId := args[0], args[1]
	requestKey := "REQ_DELIVERY_" + deliveryId
	currentTime := time.Now()

	// Check for existing request
	existingRequestBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}

	if existingRequestBytes != nil {
		var existingRequest AccessRequest
		if err := json.Unmarshal(existingRequestBytes, &existingRequest); err != nil {
			return shim.Error("Failed to parse existing request")
		}

		// Check if expired
		expiryTime, err := time.Parse(time.RFC3339, existingRequest.ExpiryTime)
		if err != nil {
			return shim.Error("Invalid expiry time format in existing request")
		}

		if currentTime.After(expiryTime) {
			// Mark old request as denied if expired
			existingRequest.Status = "Rejected"
			existingRequest.Decision = fmt.Sprintf("Rejected: expired at %s", existingRequest.ExpiryTime)
			oldRequestBytes, _ := json.Marshal(existingRequest)

			// Save old request with a different key
			err = stub.PutState(requestKey+"_EXPIRED", oldRequestBytes)
			if err != nil {
				return shim.Error("Failed to save expired request")
			}

			// Delete the old request
			err = stub.DelState(requestKey)
			if err != nil {
				return shim.Error("Failed to delete expired request")
			}

			return shim.Error(fmt.Sprintf("Access denied: expired at %s", existingRequest.ExpiryTime))
		} else {
			// If still valid, return the current access request
			return shim.Success([]byte(fmt.Sprintf("Access already granted for delivery %s until %s", deliveryId, existingRequest.ExpiryTime)))
		}
	}

	// If no existing request, create a new one with 30 minutes expiry
	expiryTime := currentTime.Add(30 * time.Minute).Format(time.RFC3339)
	newQRCode := generateQRCode(deliveryId)
	request := AccessRequest{
		VisitorID:       deliveryId,
		ApprovedBy:      []string{residentId},
		Status:          "Approved",
		ExpiryTime:      expiryTime,
		Decision:        fmt.Sprintf("Approved by %s until %s", residentId, expiryTime),
		QRCode:          newQRCode,
		QRCodeTimestamp: currentTime.Format(time.RFC3339),
	}

	// Save the new request to the ledger
	requestBytes, _ := json.Marshal(request)
	err = stub.PutState(requestKey, requestBytes)
	if err != nil {
		return shim.Error("Failed to store delivery access request")
	}

	return shim.Success([]byte(fmt.Sprintf("Temporary access granted for delivery %s until %s", deliveryId, expiryTime)))
}

//just for test//
func (c *AccessControlContract) ForceExpiryTimeDelivery(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Delivery ID and Forced Expiry Minutes are required.")
	}

	deliveryId := args[0]
	forcedExpiryMinutes, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid expiry time format.")
	}

	requestKey := "REQ_DELIVERY_" + deliveryId
	currentTime := time.Now()

	// جلب الطلب الحالي
	existingRequestBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}

	if existingRequestBytes == nil {
		return shim.Error("No access request found for the delivery ID")
	}

	var existingRequest AccessRequest
	if err := json.Unmarshal(existingRequestBytes, &existingRequest); err != nil {
		return shim.Error("Failed to parse existing request")
	}

	// تعديل وقت الانتهاء ليصبح منتهي (في الماضي)
	existingRequest.ExpiryTime = currentTime.Add(-time.Duration(forcedExpiryMinutes) * time.Minute).Format(time.RFC3339)

	// حفظ التحديث في نفس المفتاح
	updatedBytes, _ := json.Marshal(existingRequest)
	err = stub.PutState(requestKey, updatedBytes)
	if err != nil {
		return shim.Error("Failed to update expiry time for delivery request")
	}

	return shim.Success([]byte(fmt.Sprintf("Expiry time forcibly set for delivery %s to %s", deliveryId, existingRequest.ExpiryTime)))
}
func (c *AccessControlContract) checkAccessOfDeliveryGrant(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Worker ID and QR Code are required.")
	}

	deliveryId := args[0]
	inputQRCode := args[1]

	// بناء المفتاح للتحقق من الطلب
	requestKey := "REQ_SERVICE_" +deliveryId

	// جلب الطلب من السجل
	existingBytes, err := stub.GetState(requestKey)
	if err != nil {
		return shim.Error("Failed to get existing request")
	}

	if existingBytes == nil {
		return shim.Error("No access request found for the worker ID")
	}

	var existingRequest AccessRequest
	if err := json.Unmarshal(existingBytes, &existingRequest); err != nil {
		return shim.Error("Failed to parse existing request")
	}

	// التحقق مما إذا كان QR Code المدخل يتطابق مع الذي في السجل
	if inputQRCode == existingRequest.QRCode {
        // If QR codes match
		accessResponse := AccessResponse{"Granted", "Delivery approved and QR code verified."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+deliveryId+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	} else {
		// If QR codes do not match
		accessResponse := AccessResponse{"Denied", "QR code does not match the one stored in the ledger."}
		respBytes, _ := json.Marshal(accessResponse)
		stub.PutState("DECISION_"+deliveryId+"_"+time.Now().Format(time.RFC3339), respBytes)
		return shim.Success(respBytes)
	}
}

/*func (c *AccessControlContract) GetAllDecisions(stub shim.ChaincodeStubInterface) sc.Response {
	// Query to retrieve all resident records
	iterator, err := stub.GetStateByRange("RESIDENT_", "RESIDENT_~")
	if err != nil {
		return shim.Error(fmt.Sprintf("Error retrieving residents: %s", err.Error()))
	}
	defer iterator.Close()

	var allDecisions []map[string]interface{}

	// Iterate through residents and their approved visitors
	for iterator.HasNext() {
		queryResponse, _ := iterator.Next()
		var resident struct {
			Name     string   `json:"name"`
			Visitors []string `json:"visitors"`
		}
		json.Unmarshal(queryResponse.Value, &resident)

		// Check for each visitor's access status
		for _, visitor := range resident.Visitors {
			// Check the visitor's request status
			requestResponse := c.invokeResidentsContract(stub, "GetRequestStatus", []string{visitor})
			if requestResponse.Status != shim.OK {
				return requestResponse
			}

			var request AccessRequest
			json.Unmarshal(requestResponse.Payload, &request)

			var accessResponse AccessResponse
			switch request.Status {
			case "Pending":
				accessResponse = AccessResponse{"Pending", "Your request is under review. Please wait for approval."}
			case "Approved":
				// Check for QR code availability
				qrCodeKey := "QR_" + visitor
				qrCodeBytes, _ := stub.GetState(qrCodeKey)
				if qrCodeBytes != nil {
					accessResponse = AccessResponse{"Granted", "Visitor approved and QR code verified. Access granted."}
				} else {
					accessResponse = AccessResponse{"Denied", "QR code is invalid or not available."}
				}
			case "Rejected":
				accessResponse = AccessResponse{"Denied", "Your request was rejected. Please submit a new request."}
			}

			// Add the decision with resident name and access status
			decision := map[string]interface{}{
				"ResidentName": resident.Name,
				"VisitorID":    visitor,
				"AccessStatus": accessResponse.Access,
				"Reason":       accessResponse.Reason,
			}
			allDecisions = append(allDecisions, decision)
		}
	}

	// Convert all decisions to JSON
	decisionsBytes, _ := json.Marshal(allDecisions)
	return shim.Success(decisionsBytes)
}*/


func main() {
	err := shim.Start(new(AccessControlContract))
	if err != nil {
		fmt.Printf("Error starting AccessControlContract: %s", err)
	}
}
