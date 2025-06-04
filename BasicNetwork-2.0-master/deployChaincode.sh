export CORE_PEER_TLS_ENABLED=true
export GOCACHE=/tmp/go-cache
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
    pushd ./artifacts/src/github.com/residentManagement2/go
    GO111MODULE=on go mod vendor
    popd
    echo Finished vendoring Go dependencies
}

#  presetup

export CC_RUNTIME_LANGUAGE="golang"
export VERSION="1"
export CC_SRC_PATH="./artifacts/src/github.com/residentManagement2/go" 
export CC_NAME="residentManagement"
export CHANNEL_NAME="residentschannel"
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
    export GOCACHE=/tmp/go-cache  # Set a writable GOCACHE once

    # Make sure GOCACHE directory exists and is writable
    mkdir -p $GOCACHE
    chmod -R 777 $GOCACHE

    for peer in Peer0Org1 Peer0Org2; do
        setGlobalsFor${peer}
        peer lifecycle chaincode install ${CC_NAME}.tar.gz
        echo "==== Chaincode installed on ${peer} ===="
    done
}



#   installChaincode


# Checks if the chaincode is installed on peer0.org1
# Extracts and prints the Package ID
# Prepares the Package ID for the next steps (approval & commit)


queryInstalled() {
    setGlobalsForPeer0Org1
    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}_${VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo "PackageID is ${PACKAGE_ID}"
    echo "===================== Query installed successful on peer0.org1 ====================="
}


approveForAllOrgs() {
    for peer in Peer0Org1 Peer0Org2; do
        setGlobalsFor${peer}
        peer lifecycle chaincode approveformyorg -o localhost:7050 \
            --ordererTLSHostnameOverride orderer.example.com --tls \
            --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
            --version ${VERSION} \
            --init-required --package-id ${PACKAGE_ID} \
            --sequence ${sequence} \
            --signature-policy "OR('ResidentsMSP.peer','ManagerMSP.peer')"
        echo "===================== Chaincode approved from ${peer} ====================="
    done
}

# approveForAllOrgs

getBlock() {
    setGlobalsForPeer0Org1
    peer channel getinfo -c residentschannel -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls \
        --cafile $ORDERER_CA
}

# getBlock

checkCommitReadyness() {
    setGlobalsForPeer0Org1
    peer lifecycle chaincode checkcommitreadiness \
        --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --version ${VERSION} --sequence ${sequence} \
        --signature-policy "OR('ResidentsMSP.peer','ManagerMSP.peer')" \
        --output json --init-required
    echo "===================== Checking commit readiness from org 1 ====================="
}

# checkCommitReadyness

commitChaincodeDefination() {
    setGlobalsForPeer0Org1
    peer lifecycle chaincode commit -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com --tls \
        --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
        --peerAddresses localhost:8151 --tlsRootCertFiles $PEER0_ORG2_CA \
        --version ${VERSION} --sequence ${sequence} --init-required \
        --signature-policy "OR('ResidentsMSP.peer','ManagerMSP.peer')"
    echo "==== Chaincode committed successfully ===="
}

# commitChaincodeDefination

queryCommitted() {
    setGlobalsForPeer0Org1
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME}
}

# queryCommitted

chaincodeInvokeInit() {
    setGlobalsForPeer0Org1
    peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
        -C residentschannel -n residentManagement \
        --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA \
        --peerAddresses localhost:8151 --tlsRootCertFiles $PEER0_ORG2_CA \
        --isInit -c '{"Args":[]}'
}

# chaincodeInvokeInit




# Run this function if you add any new dependency in chaincode
# presetup

   packageChaincode
   installChaincode
    queryInstalled
    # queryInstalled
echo "Waiting for Raft leader election..."
sleep 15   # <-- Add this lineS
    approveForAllOrgs
    checkCommitReadyness
    commitChaincodeDefination
    queryCommitted
     chaincodeInvokeInit
# sleep 5
# chaincodeInvoke
# sleep 3
# chaincodeQuery
