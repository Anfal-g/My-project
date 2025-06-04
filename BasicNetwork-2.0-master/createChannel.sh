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

# Define channel names
export RESIDENTS_CHANNEL=residentschannel
export ACCESS_CONTROL_CHANNEL=accesscontrolchannel

# Set environment variables for Org1 (Residents) Peer0
setGlobalsForPeer0Org1(){
    export CORE_PEER_LOCALMSPID="ResidentsMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=/mnt/c/Users/MICRO/Desktop/My-Project/BasicNetwork-2.0-master/artifacts/channel/crypto-config/peerOrganizations/residents.example.com/users/Admin@residents.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export FABRIC_CFG_PATH=$ORG1_CFG_PATH
}
# set -x  # Enable debug mode
# setGlobalsForPeer0Org1
# set +x  # Disable debug mode

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

echo "Using ORDERER_CA=$ORDERER_CA"
openssl x509 -in $ORDERER_CA -noout -subject -issuer -dates

# Function to create a channel
createChannel() {
    CHANNEL_NAME=$1
    CHANNEL_TX_FILE=$2
    
    echo "Creating channel: $CHANNEL_NAME"
    peer channel create -o localhost:7050 \
        -c $CHANNEL_NAME \
        -f $CHANNEL_TX_FILE \
        --outputBlock ./artifacts/channel/${CHANNEL_NAME}.block \
        --tls --cafile $ORDERER_CA
    
    if [ $? -ne 0 ]; then
        echo "Failed to create channel $CHANNEL_NAME  !!!"
        exit 1
    fi
    echo "Channel $CHANNEL_NAME created successfully"
}


# Function for a peer to join a channel
joinChannel() {
    CHANNEL_NAME=$1
    PEER_FUNCTION=$2  

    $PEER_FUNCTION
    echo "Peer $CORE_PEER_ADDRESS joining channel: $CHANNEL_NAME"
    
    peer channel join -b ./artifacts/channel/${CHANNEL_NAME}.block
    
    if [ $? -ne 0 ]; then
        echo "Failed to join $CORE_PEER_ADDRESS to channel $CHANNEL_NAME"
        exit 1
    fi
    echo "Peer $CORE_PEER_ADDRESS joined channel $CHANNEL_NAME successfully"
}

# Function to set anchor peers
updateAnchorPeers() {
    CHANNEL_NAME=$1
    PEER_FUNCTION=$2
    ANCHOR_TX_FILE=$3
    
    $PEER_FUNCTION
    echo "Updating anchor peers for $CORE_PEER_LOCALMSPID on channel: $CHANNEL_NAME"
    
    peer channel update -o localhost:7050 \
        -c $CHANNEL_NAME \
        -f $ANCHOR_TX_FILE \
        --tls --cafile $ORDERER_CA
    
    if [ $? -ne 0 ]; then
        echo "Failed to update anchor peers for $CORE_PEER_LOCALMSPID on channel $CHANNEL_NAME"
        exit 1
    fi
    echo "Anchor peers updated for $CORE_PEER_LOCALMSPID on channel $CHANNEL_NAME successfully"
}



# # # # Create both channels

      setGlobalsForPeer0Org2
           createChannel $RESIDENTS_CHANNEL ./artifacts/channel/${RESIDENTS_CHANNEL}.tx
    #    setGlobalsForPeer0Org2
    #    createChannel $ACCESS_CONTROL_CHANNEL ./artifacts/channel/${ACCESS_CONTROL_CHANNEL}.tx
  

     
     


 
# # # # # # # # # Org1 (Residents) joins both 
 


        # joinChannel $ACCESS_CONTROL_CHANNEL setGlobalsForPeer0Org1
        # joinChannel $ACCESS_CONTROL_CHANNEL setGlobalsForPeer1Org1

        joinChannel $RESIDENTS_CHANNEL setGlobalsForPeer1Org1
       joinChannel $RESIDENTS_CHANNEL setGlobalsForPeer0Org1


# # # # # # # # # Org2 (Building Manager) joins both channels
       joinChannel $RESIDENTS_CHANNEL setGlobalsForPeer0Org2
    #   joinChannel $ACCESS_CONTROL_CHANNEL setGlobalsForPeer0Org2


# # # # # # # # # Org3 (Entry System) joins only AccessControlChannel
    #    joinChannel $RESIDENTS_CHANNEL setGlobalsForPeer0Org3
    #   joinChannel $ACCESS_CONTROL_CHANNEL setGlobalsForPeer0Org3

# # # # # Update anchor peers for each organization
      updateAnchorPeers $RESIDENTS_CHANNEL setGlobalsForPeer0Org2 ./artifacts/channel/artifacts/channel/ManagerMSP_ResidentsChannel_anchors.tx



    #  updateAnchorPeers $ACCESS_CONTROL_CHANNEL setGlobalsForPeer0Org2 ./artifacts/channel/artifacts/channel/ManagerMSP_AccessControlChannel_anchors.tx
    #  updateAnchorPeers $ACCESS_CONTROL_CHANNEL setGlobalsForPeer0Org3 ./artifacts/channel/artifacts/channel/EntrySystemMSPanchors.tx
    #  updateAnchorPeers $ACCESS_CONTROL_CHANNEL setGlobalsForPeer0Org1 ./artifacts/channel/artifacts/channel/ResidentsMSP_AccessControlChannel_anchors.tx

    updateAnchorPeers $RESIDENTS_CHANNEL setGlobalsForPeer0Org1 ./artifacts/channel/artifacts/channel/ResidentsMSPanchors.tx
        # updateAnchorPeers $RESIDENTS_CHANNEL setGlobalsForPeer1Org1 ./artifacts/channel/artifacts/channel/ResidentsMSPanchors.tx

    # updateAnchorPeers $RESIDENTS_CHANNEL setGlobalsForPeer0Org3 ./artifacts/channel/artifacts/channel/EntrySystemMSP_ResidentsChannel_anchors.tx

#  peer channel fetch config config_block.pb -o localhost:7050 -c $ACCESS_CONTROL_CHANNEL --tls --cafile $ORDERER_CA
#  peer channel fetch config config_block.pb -o localhost:7050 -c $RESIDENTS_CHANNEL --tls --cafile $ORDERER_CA


#  peer channel update -o localhost:7050 --channelID ACCESS_CONTROL_CHANNEL --file empty.pb --tls --cafile $ORDERER_CA
# peer channel leave -o localhost:7050 -c $ACCESS_CONTROL_CHANNEL --tls --cafile $ORDERER_CA
#  peer channel update -o localhost:7050 --channelID RESIDENTS_CHANNEL --file empty.pb --tls --cafile $ORDERER_CA
# peer channel leave -o localhost:7050 -c $RESIDENTS_CHANNEL --tls --cafile $ORDERER_CA

# #   rm -rf ./artifacts/channel/*.block
# #   rm -rf ./artifacts/channel/*.tx


#   docker exec -it peer0.residents.example.com rm -rf /var/hyperledger/production/ledgersData
#   docker exec -it peer1.residents.example.com rm -rf /var/hyperledger/production/ledgersData
#   docker exec -it peer0.manager.example.com rm -rf /var/hyperledger/production/ledgersData
#   docker exec -it peer0.entrysystem.example.com rm -rf /var/hyperledger/production/ledgersData








