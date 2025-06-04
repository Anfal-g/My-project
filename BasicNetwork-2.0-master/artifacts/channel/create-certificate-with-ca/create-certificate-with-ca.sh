createcertificatesForresidents() {
  echo
  echo "Enroll the CA admin"
  echo

  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/
  export FABRIC_CA_CLIENT_HOME=${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/

  # ✅ Ensure the target directory exists before copying
  mkdir -p ${PWD}/fabric-ca/residents

  # ✅ Step 1: Copy the TLS CA cert from the container
  echo "Copying TLS cert from container..."
  docker cp ca.residents.example.com:/etc/hyperledger/fabric-ca-server-tls/tlsca.residents.example.com-cert.pem ${PWD}/fabric-ca/residents/tls-cert.pem

  # ✅ Step 2: Enroll the CA admin
fabric-ca-client enroll -u https://admin:adminpw@ca.residents.example.com:7054 \
  --caname ca.residents.example.com \
  --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem \
  --tls.client.certfile ${PWD}/fabric-ca/residents/tls-cert.pem \
  --tls.client.keyfile ${PWD}/fabric-ca/residents/tls-key.pem \
  --csr.hosts ca.residents.example.com


    

  # ✅ Step 3: Write NodeOUs
  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-residents-example-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-residents-example-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-residents-example-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-residents-example-com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/msp/config.yaml

  echo
  echo "Register peer0"
  echo
  fabric-ca-client register --caname ca.residents.example.com --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  echo
  echo "Register peer1"
  echo
  fabric-ca-client register --caname ca.residents.example.com --id.name peer1 --id.secret peer1pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  echo
  echo "Register user"
  echo
  fabric-ca-client register --caname ca.residents.example.com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  echo
  echo "Register the org admin"
  echo
  fabric-ca-client register --caname ca.residents.example.com --id.name residentsadmin --id.secret residentsadminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/peers

  # -----------------------------------------------------------------------------------
  #  Peer 0
  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com

  echo
  echo "## Generate the peer0 msp"
  echo
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca.residents.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/msp --csr.hosts peer0.residents.example.com --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/msp/config.yaml

  echo
  echo "## Generate the peer0-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca.residents.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls --enrollment.profile tls --csr.hosts peer0.residents.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/msp/tlscacerts/ca.crt

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/tlsca
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/tlsca/tlsca.residents.example.com-cert.pem

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/ca
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer0.residents.example.com/msp/cacerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/ca/ca.residents.example.com-cert.pem

  # ------------------------------------------------------------------------------------------------

  # Peer1

  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com

  echo
  echo "## Generate the peer1 msp"
  echo
  fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 --caname ca.residents.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/msp --csr.hosts peer1.residents.example.com --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/msp/config.yaml

  echo
  echo "## Generate the peer1-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 --caname ca.residents.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls --enrollment.profile tls --csr.hosts peer1.residents.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/peers/peer1.residents.example.com/tls/server.key

  # --------------------------------------------------------------------------------------------------

  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/users
  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/users/User1@residents.example.com

  echo
  echo "## Generate the user msp"
  echo
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca.residents.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/users/User1@residents.example.com/msp --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  mkdir -p crypto-config-ca/peerOrganizations/residents.example.com/users/Admin@residents.example.com

  echo
  echo "## Generate the org admin msp"
  echo
  fabric-ca-client enroll -u https://residentsadmin:residentsadminpw@localhost:7054 --caname ca.residents.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/users/Admin@residents.example.com/msp --tls.certfiles ${PWD}/fabric-ca/residents/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/residents.example.com/users/Admin@residents.example.com/msp/config.yaml

}

# createcertificatesForresidents

createCertificateFormanager() {
  echo
  echo "Enroll the CA admin"
  echo
  mkdir -p /crypto-config-ca/peerOrganizations/manager.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/

   
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca.manager.example.com --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem

   # ✅ Ensure the target directory exists before copying
  mkdir -p ${PWD}/fabric-ca/manager

  # ✅ Step 1: Copy the TLS CA cert from the container
  echo "Copying TLS cert from container..."
  docker cp ca.manager.example.com:/etc/hyperledger/fabric-ca-server-tls/tlsca.manager.example.com-cert.pem ${PWD}/fabric-ca/manager/tls-cert.pem

  # ✅ Step 2: Enroll the CA admin
fabric-ca-client enroll -u https://admin:adminpw@ca.manager.example.com:8054 \
  --caname ca.manager.example.com \
  --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem \
  --tls.client.certfile ${PWD}/fabric-ca/manager/tls-cert.pem \
  --tls.client.keyfile ${PWD}/fabric-ca/manager/tls-key.pem \
  --csr.hosts ca.manager.example.com


  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-manager-example-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-manager-example-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-manager-example-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-manager-example-com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/msp/config.yaml

  echo
  echo "Register peer0"
  echo
   
  fabric-ca-client register --caname ca.manager.example.com --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  echo
  echo "Register peer1"
  echo
   
  fabric-ca-client register --caname ca.manager.example.com --id.name peer1 --id.secret peer1pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  echo
  echo "Register user"
  echo
   
  fabric-ca-client register --caname ca.manager.example.com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  echo
  echo "Register the org admin"
  echo
   
  fabric-ca-client register --caname ca.manager.example.com --id.name manageradmin --id.secret manageradminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  mkdir -p crypto-config-ca/peerOrganizations/manager.example.com/peers
  mkdir -p crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com

  # --------------------------------------------------------------
  # Peer 0
  echo
  echo "## Generate the peer0 msp"
  echo
   
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca.manager.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/msp --csr.hosts peer0.manager.example.com --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/msp/config.yaml

  echo
  echo "## Generate the peer0-tls certificates"
  echo
   
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca.manager.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls --enrollment.profile tls --csr.hosts peer0.manager.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/msp/tlscacerts/ca.crt

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/tlsca
  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/tlsca/tlsca.manager.example.com-cert.pem

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/ca
  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/peers/peer0.manager.example.com/msp/cacerts/* ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/ca/ca.manager.example.com-cert.pem

  

  mkdir -p crypto-config-ca/peerOrganizations/manager.example.com/users
  mkdir -p crypto-config-ca/peerOrganizations/manager.example.com/users/User1@manager.example.com

  echo
  echo "## Generate the user msp"
  echo
   
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca.manager.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/users/User1@manager.example.com/msp --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  mkdir -p crypto-config-ca/peerOrganizations/manager.example.com/users/Admin@manager.example.com

  echo
  echo "## Generate the org admin msp"
  echo
   
  fabric-ca-client enroll -u https://manageradmin:manageradminpw@localhost:8054 --caname ca.manager.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/users/Admin@manager.example.com/msp --tls.certfiles ${PWD}/fabric-ca/manager/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/manager.example.com/users/Admin@manager.example.com/msp/config.yaml

}
createCertificateForentrysystem() {
  echo
  echo "Enroll the CA admin"
  echo
  mkdir -p /crypto-config-ca/peerOrganizations/entrysystem.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/

   
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca.entrysystem.example.com --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   
     # ✅ Ensure the target directory exists before copying
  mkdir -p ${PWD}/fabric-ca/entrysystem

  # ✅ Step 1: Copy the TLS CA cert from the container
  echo "Copying TLS cert from container..."
  docker cp ca.entrysystem.example.com:/etc/hyperledger/fabric-ca-server-tls/tlsca.entrysystem.example.com-cert.pem ${PWD}/fabric-ca/entrysystem/tls-cert.pem

  # ✅ Step 2: Enroll the CA admin
fabric-ca-client enroll -u https://admin:adminpw@ca.entrysystem.example.com:9054 \
  --caname ca.entrysystem.example.com \
  --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem \
  --tls.client.certfile ${PWD}/fabric-ca/entrysystem/tls-cert.pem \
  --tls.client.keyfile ${PWD}/fabric-ca/entrysystem/tls-key.pem \
  --csr.hosts ca.entrysystem.example.com


  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-entrysystem-example-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-entrysystem-example-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-entrysystem-example-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-entrysystem-example-com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/msp/config.yaml

  echo
  echo "Register peer0"
  echo
   
  fabric-ca-client register --caname ca.entrysystem.example.com --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  echo
  echo "Register peer1"
  echo
   
  fabric-ca-client register --caname ca.entrysystem.example.com --id.name peer1 --id.secret peer1pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  echo
  echo "Register user"
  echo
   
  fabric-ca-client register --caname ca.entrysystem.example.com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  echo
  echo "Register the org admin"
  echo
   
  fabric-ca-client register --caname ca.entrysystem.example.com --id.name entrysystemadmin --id.secret entrysystemadminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  mkdir -p crypto-config-ca/peerOrganizations/entrysystem.example.com/peers
  mkdir -p crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com

  # --------------------------------------------------------------
  # Peer 0
  echo
  echo "## Generate the peer0 msp"
  echo
   
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca.entrysystem.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/msp --csr.hosts peer0.entrysystem.example.com --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/msp/config.yaml

  echo
  echo "## Generate the peer0-tls certificates"
  echo
   
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca.entrysystem.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls --enrollment.profile tls --csr.hosts peer0.entrysystem.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/msp/tlscacerts/ca.crt

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/tlsca
  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/tlsca/tlsca.entrysystem.example.com-cert.pem

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/ca
  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/peers/peer0.entrysystem.example.com/msp/cacerts/* ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/ca/ca.entrysystem.example.com-cert.pem

  

  mkdir -p crypto-config-ca/peerOrganizations/entrysystem.example.com/users
  mkdir -p crypto-config-ca/peerOrganizations/entrysystem.example.com/users/User1@entrysystem.example.com

  echo
  echo "## Generate the user msp"
  echo
   
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca.entrysystem.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/users/User1@entrysystem.example.com/msp --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  mkdir -p crypto-config-ca/peerOrganizations/entrysystem.example.com/users/Admin@entrysystem.example.com

  echo
  echo "## Generate the org admin msp"
  echo
   
  fabric-ca-client enroll -u https://entrysystemadmin:entrysystemadminpw@localhost:8054 --caname ca.entrysystem.example.com -M ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/users/Admin@entrysystem.example.com/msp --tls.certfiles ${PWD}/fabric-ca/entrysystem/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/entrysystem.example.com/users/Admin@entrysystem.example.com/msp/config.yaml

}
# createCertificateFormanager

createCretificateForOrderer() {
  echo
  echo "Enroll the CA admin"
  echo
  mkdir -p crypto-config-ca/ordererOrganizations/example.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/crypto-config-ca/ordererOrganizations/example.com

   
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8443 --caname ca-orderer --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8443-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8443-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8443-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8443-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/config.yaml

  echo
  echo "Register orderer"
  echo
   
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  echo
  echo "Register orderer2"
  echo
   
  fabric-ca-client register --caname ca-orderer --id.name orderer2 --id.secret ordererpw --id.type orderer --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  echo
  echo "Register orderer3"
  echo
   
  fabric-ca-client register --caname ca-orderer --id.name orderer3 --id.secret ordererpw --id.type orderer --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  echo
  echo "Register the orderer admin"
  echo
   
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  mkdir -p crypto-config-ca/ordererOrganizations/example.com/orderers
  # mkdir -p crypto-config-ca/ordererOrganizations/example.com/orderers/example.com

  # ---------------------------------------------------------------------------
  #  Orderer

  mkdir -p crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com

  echo
  echo "## Generate the orderer msp"
  echo
   
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/msp --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/msp/config.yaml

  echo
  echo "## Generate the orderer-tls certificates"
  echo
   
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls --enrollment.profile tls --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/keystore/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  mkdir ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  # -----------------------------------------------------------------------
  #  Orderer 2

  mkdir -p crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com

  echo
  echo "## Generate the orderer msp"
  echo
   
  fabric-ca-client enroll -u https://orderer2:ordererpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/msp --csr.hosts orderer2.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/config.yaml

  echo
  echo "## Generate the orderer-tls certificates"
  echo
   
  fabric-ca-client enroll -u https://orderer2:ordererpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls --enrollment.profile tls --csr.hosts orderer2.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/keystore/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  # mkdir ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/tlscacerts
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  # ---------------------------------------------------------------------------
  #  Orderer 3
  mkdir -p crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com

  echo
  echo "## Generate the orderer msp"
  echo
   
  fabric-ca-client enroll -u https://orderer3:ordererpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/msp --csr.hosts orderer3.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/msp/config.yaml

  echo
  echo "## Generate the orderer-tls certificates"
  echo
   
  fabric-ca-client enroll -u https://orderer3:ordererpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls --enrollment.profile tls --csr.hosts orderer3.example.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/signcerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/keystore/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  # mkdir ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/tlscacerts
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  # ---------------------------------------------------------------------------

  mkdir -p crypto-config-ca/ordererOrganizations/example.com/users
  mkdir -p crypto-config-ca/ordererOrganizations/example.com/users/Admin@example.com

  echo
  echo "## Generate the admin msp"
  echo
   
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:8443 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/example.com/users/Admin@example.com/msp --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem
   

  cp ${PWD}/crypto-config-ca/ordererOrganizations/example.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/example.com/users/Admin@example.com/msp/config.yaml

}

#  createCretificateForOrderer

# sudo rm -rf crypto-config-ca/*
#  sudo rm -rf fabric-ca/*
#  createcertificatesForresidents
#
 createCertificateFormanager
 createCertificateForentrysystem
# createCretificateForOrderer

