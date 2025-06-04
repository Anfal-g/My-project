export CORE_PEER_TLS_ENABLED=true

# Define paths for Orderer and Peer certificates
export ORDERER_CA=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/ca.crt
export PEER1_ORG1_CA=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/ca.crt
export PEER0_ORG2_CA=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/ca.crt
export PEER0_ORG3_CA=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/ca.crt

# Define paths for core.yaml per organization
export ORG1_CFG_PATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/config/org1/
export ORG1_PEER1_CFG_PATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/config/org1_peer1/
export ORG2_CFG_PATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/config/org2/
export ORG3_CFG_PATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/config/org3/


export PRIVATE_DATA_CONFIG=${PWD}/artifacts/private-data/collections_config.json #for AccessControle chaincode

# Define channel names
export RESIDENTS_CHANNEL=residentschannel
export ACCESS_CONTROL_CHANNEL=accesscontrolchannel

setGlobalsForOrderer() {
    export CORE_PEER_LOCALMSPID="OrdererMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp

}

# Set environment variables for Org1 (Residents) Peer0
setGlobalsForPeer0Org1(){
    export CORE_PEER_LOCALMSPID="ResidentsMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/residents.example.com/users/Admin@residents.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export FABRIC_CFG_PATH=$ORG1_CFG_PATH
}

# Set environment variables for Org1 (Residents) Peer1
setGlobalsForPeer1Org1(){
    export CORE_PEER_LOCALMSPID="ResidentsMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/residents.example.com/users/Admin@residents.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8051
    export FABRIC_CFG_PATH=$ORG1_PEER1_CFG_PATH
}
# Set environment variables for Org2 (Building Manager) Peer0
setGlobalsForPeer0Org2(){
    export CORE_PEER_LOCALMSPID="ManagerMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/manager.example.com/users/Admin@manager.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8151
    export FABRIC_CFG_PATH=$ORG2_CFG_PATH
}
# Set environment variables for Org3 (Entry System) Peer0
setGlobalsForPeer0Org3(){
    export CORE_PEER_LOCALMSPID="EntrySystemMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
    export CORE_PEER_MSPCONFIGPATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/entrysystem.example.com/users/Admin@entrysystem.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
    export FABRIC_CFG_PATH=$ORG3_CFG_PATH
}

presetup() {
    echo Vendoring Go dependencies ...
    pushd ./artifacts/src/github.com/accessControl/go
    GO111MODULE=on go mod vendor
    popd
    echo Finished vendoring Go dependencies
}

#  presetup

export CC_RUNTIME_LANGUAGE="golang"
export VERSION="1"
export CC_SRC_PATH="./artifacts/src/github.com/accessControl/go" 
export CC_NAME="accessControl"
export CHANNEL_NAME="accesscontrolchannel"
export sequence="1"
# Package the chaincode for deployment
packageChaincode() {
    rm -rf ${CC_NAME}.tar.gz
    setGlobalsForPeer0Org1
    peer lifecycle chaincode package ${CC_NAME}.tar.gz \
        --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} \
        --label ${CC_NAME}_${VERSION}
    echo "===================== Chaincode is packaged on peer0.org1 ===================== "
}
#  packageChaincode

installChaincode() {
    export GOCACHE=/tmp/.cache
   for peer in Peer0Org1 Peer0Org2 Peer0Org3 ; do
        setGlobalsFor${peer}
        GOCACHE=/tmp/.cache peer lifecycle chaincode install ${CC_NAME}.tar.gz
        echo "==== Chaincode installed on ${peer} ===="
    done
}
# installChaincode


# Checks if the chaincode is installed on peer0.org1
# Extracts and prints the Package ID
# Prepares the Package ID for the next steps (approval & commit)


queryInstalled() {
    setGlobalsForPeer0Org1
    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}_${VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo PackageID is ${PACKAGE_ID}
    echo "===================== Query installed successful on peer0.org1 ===================== "
}

# queryInstalled


# --collections-config ./artifacts/private-data/collections_config.json \
#         --signature-policy "OR('Org1MSP.member','Org2MSP.member')" \
# --collections-config $PRIVATE_DATA_CONFIG \

# This function approves the chaincode definition for Org on the channel
# Sends approval to the orderer (-o localhost:7050).
# Uses TLS and orderer certificates (--tls --cafile $ORDERER_CA).
# Approves the specific chaincode (--name ${CC_NAME}) with a version and sequence number.
# Includes private data collections if used (--collections-config $PRIVATE_DATA_CONFIG).
# Links the approval to the installed package (--package-id ${PACKAGE_ID}).

approveForAllOrgs() {
    for peer in Peer0Org1 Peer0Org2 Peer0Org3; do  # Peer1Org1  
        setGlobalsFor${peer}
   peer lifecycle chaincode approveformyorg -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com --tls \
    --collections-config $PRIVATE_DATA_CONFIG \
    --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${VERSION} \
    --init-required --package-id ${PACKAGE_ID} \
    --sequence ${sequence} \
    --signature-policy "OR('ManagerMSP.peer', 'EntrySystemMSP.peer')"
        echo "===================== Chaincode approved from ${peer} ===================== "
    done
}
# approveForAllOrgs

getBlock() {
    setGlobalsForPeer0Org1
    # peer channel fetch 10 -c mychannel -o localhost:7050 \
    #     --ordererTLSHostnameOverride orderer.example.com --tls \
    #     --cafile $ORDERER_CA

    peer channel getinfo  -c accesscontrolchannel -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls \
        --cafile $ORDERER_CA
}

# getBlock



# --signature-policy "OR ('Org1MSP.member')"
# --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA
# --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles $PEER0_ORG2_CA
#--channel-config-policy Channel/Application/Admins
# --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer')"


# The function checkCommitReadyness() checks whether all required organizations have approved the chaincode before committing it
checkCommitReadyness() {

    setGlobalsForPeer0Org1
        peer lifecycle chaincode checkcommitreadiness \
            --collections-config $PRIVATE_DATA_CONFIG \
            --channelID $CHANNEL_NAME --name ${CC_NAME} \
            --version ${VERSION} --sequence ${sequence} \
            --signature-policy "OR('ManagerMSP.peer', 'EntrySystemMSP.peer')" \
            --output json --init-required
    echo "===================== Checking commit readiness from org 1 ===================== "
}

# checkCommitReadyness


# This function commits the chaincode definition to the Fabric network after it has been approved by all required organizations.
commitChaincodeDefination() {
      setGlobalsForPeer0Org1
    peer lifecycle chaincode commit -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --collections-config $PRIVATE_DATA_CONFIG \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
        --peerAddresses localhost:8151 --tlsRootCertFiles $PEER0_ORG2_CA \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG3_CA \
        --version ${VERSION} --sequence ${sequence} --init-required \
        --signature-policy "OR('ManagerMSP.peer', 'EntrySystemMSP.peer')"
    echo "==== Chaincode committed successfully ===="

}

# commitChaincodeDefination
# function checks whether the chaincode definition has been successfully committed to the specified channel. It does this by running
queryCommitted() {
    setGlobalsForPeer0Org1
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME}
    }

# queryCommitted

# The chaincodeInvokeInit() function is used to initialize the chaincode after deployment.
#  It sends an invoke transaction with the --isInit flag to indicate that the chaincode should 
#  execute its initialization logic
chaincodeInvokeInit() {
    setGlobalsForPeer0Org1
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
    -C accesscontrolchannel -n accessControl \
    --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
    --peerAddresses localhost:8151 --tlsRootCertFiles $PEER0_ORG2_CA \
    --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG3_CA \
        --isInit -c '{"Args":[]}'
}

# chaincodeInvokeInit

chaincodeInvoke() {


    # setGlobalsForPeer0Org1
    # peer chaincode invoke -o localhost:7050 \
    #     --ordererTLSHostnameOverride orderer.example.com \
    #     --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
    # -C accesscontrolchannel -n accessControl \
    # --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
    # --peerAddresses localhost:8151 --tlsRootCertFiles $PEER0_ORG2_CA \
    # --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG3_CA \
    # --isInit -c '{"Args":[]}'

   setGlobalsForPeer0Org1
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA \
        -C accesscontrolchannel \
        -n accessControl \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
        --peerAddresses localhost:8151 --tlsRootCertFiles $PEER0_ORG2_CA \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG3_CA \
        -c '{"Args":["checkAccessResidents","user28","QR-RESIDENT-user28"]}'
}

#  chaincodeInvoke

chaincodeQuery() {
    setGlobalsForPeer0Org2

    # Query all cars
    # peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["queryAllCars"]}'

    # Query Car by Id
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "queryCar","Args":["CAR0"]}'
    #'{"Args":["GetSampleData","Key1"]}'

    # Query Private Car by Id
    # peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "readPrivateCar","Args":["1111"]}'
    # peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "readCarPrivateDetails","Args":["1111"]}'
}

# chaincodeQuery

# Run this function if you add any new dependency in chaincode
# presetup

 packageChaincode
 installChaincode
 queryInstalled
approveForAllOrgs
 checkCommitReadyness

 commitChaincodeDefination
 queryCommitted
 chaincodeInvokeInit
# sleep 5
#  chaincodeInvoke
# sleep 3
# chaincodeQuery
