const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')

const helper = require('./helper')

const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name, transientData) => {
    try {
        logger.debug(util.format('\n============ invoke transaction on channel %s ============\n', channelName));

        // load the network configuration
        // const ccpPath =path.resolve(__dirname, '..', 'config', 'connection-org1.json');
        // const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
        const ccp = await helper.getCCP(org_name) //JSON.parse(ccpJSON);

        // Create a new file system based wallet for managing identities.
        const walletPath = await helper.getWalletPath(org_name) //path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        

        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true },
            eventHandlerOptions: {
                commitTimeout: 100,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX
            }
            // transaction: {
            //     strategy: createTransactionEventhandler()
            // }
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(channelName);

        const contract = network.getContract(chaincodeName);

        let result
        let message;
        if (fcn === "RegisterResident") {
            // args: [userId, gender, apartment, maritalStatus, residentType, ...visitorIds]
            result = await contract.submitTransaction(
                fcn, 
                args[0], // userId
                args[1], // gender
                args[2], // apartment
                args[3], // maritalStatus
                args[4], // residentType
                ...args.slice(5) // visitorIds (optional)
            );
            message = `✅ Resident ${args[0]} registered successfully.`;

        } else if (fcn === "RegisterUser") {
            // args: [userId, email, password, role, image]
            result = await contract.submitTransaction(
                fcn,
                args[0], // userId
                args[1], // email
                args[2], // password
                args[3], // role
                args[4]  // image
            );
            message = `✅ User ${args[0]} registered successfully.`;

        } else if (fcn === "AddVisitor") {
            // args: [userId, visitorId]
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            message = `✅ Visitor ${args[1]} added to resident ${args[0]}'s list.`;

        } else if (fcn === "RequestApproval") {
            // args: [userId, visitorId]
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            message = `✅ Visitor access request for ${args[1]} submitted.`;

        } else if (fcn === "RequestServiceAccess") {
            // args: [workerId, requesterId]
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            message = `✅ Service access request for worker ${args[0]} submitted.`;

        } else if (fcn === "RequestDeliveryAccess") {
            // args: [deliveryId, residentId]
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            message = `✅ Delivery access request from ${args[0]} submitted.`;

        } else if (fcn === "ApproveRequest") {
            // args: [visitorId, approverId]
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            message = `✅ Approval from ${args[1]} recorded for visitor ${args[0]}.`;



        } else {
            throw new Error(`❌ Invalid function name: ${fcn}. Supported functions: 
                RegisterResident, RegisterUser, AddVisitor, RequestApproval, 
                RequestServiceAccess, RequestDeliveryAccess, ApproveRequest, UpdateResident.`);
        }

        await gateway.disconnect();

        result = JSON.parse(result.toString());

        let response = {
            message: message,
            result
        }

        return response;


    } catch (error) {

        console.log(`Getting error: ${error}`)
        return error.message

    }
}

exports.invokeTransaction = invokeTransaction;