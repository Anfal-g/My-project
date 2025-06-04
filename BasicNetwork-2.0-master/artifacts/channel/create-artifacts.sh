# # Delete existing artifacts
#  rm -rf crypto-config/*
# rm -rf artifacts/channel/*.block artifacts/channel/*.tx
#    rm -f $RESIDENTS_CHANNEL.tx $ACCESS_CONTROL_CHANNEL.tx
#   # rm -f  $ACCESS_CONTROL_CHANNEL.tx

#    rm -f ./genesis.block
# # rm -f ./*.tx

# #Generate Crypto artifactes for organizations
#    rm -rf crypto-config/*
    # cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/


# # System channel name
    SYS_CHANNEL="sys-channel"

# # Define channels
   RESIDENTS_CHANNEL="residentschannel"
  ACCESS_CONTROL_CHANNEL="accesscontrolchannel"

 echo "Creating channels: $RESIDENTS_CHANNEL and $ACCESS_CONTROL_CHANNEL"

# # Generate System Genesis block
    # configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL -outputBlock ./genesis.block

# # Generate channel configuration block for ResidentsChannel
   configtxgen -profile ResidentsChannel -configPath . -outputCreateChannelTx ./$RESIDENTS_CHANNEL.tx -channelID $RESIDENTS_CHANNEL

# # Generate channel configuration block for AccessControlChannel
    #  configtxgen -profile AccessControlChannel -configPath . -outputCreateChannelTx ./$ACCESS_CONTROL_CHANNEL.tx -channelID $ACCESS_CONTROL_CHANNEL

# # Generate anchor peer updates for each organization in each channel



  configtxgen -profile ResidentsChannel \
       -configPath . -outputAnchorPeersUpdate ./ResidentsMSPanchors.tx \
      -channelID $RESIDENTS_CHANNEL \
     -asOrg ResidentsMSP


 configtxgen -profile ResidentsChannel \
     -outputAnchorPeersUpdate ./ManagerMSP_ResidentsChannel_anchors.tx \
     -channelID $RESIDENTS_CHANNEL \
     -asOrg ManagerMSP





#  configtxgen -profile AccessControlChannel \
#      -configPath . -outputAnchorPeersUpdate ./ResidentsMSP_AccessControlChannel_anchors.tx \
#      -channelID $ACCESS_CONTROL_CHANNEL \
#      -asOrg ResidentsMSP


#  configtxgen -profile AccessControlChannel \
#     -configPath . -outputAnchorPeersUpdate ./ManagerMSP_AccessControlChannel_anchors.tx \
#      -channelID $ACCESS_CONTROL_CHANNEL \
#      -asOrg ManagerMSP



#  configtxgen -profile AccessControlChannel \
#      -configPath . -outputAnchorPeersUpdate ./EntrySystemMSPanchors.tx \
#      -channelID $ACCESS_CONTROL_CHANNEL \
#      -asOrg EntrySystemMSP

 configtxgen -profile ResidentsChannel \
     -outputAnchorPeersUpdate ./EntrySystemMSP_ResidentsChannel_anchors.tx \
     -channelID $RESIDENTS_CHANNEL \
     -asOrg EntrySystemMSP