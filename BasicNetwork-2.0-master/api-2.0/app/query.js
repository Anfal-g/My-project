const { Gateway, Wallets, } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')


const helper = require('./helper')
const query = async (channelName, chaincodeName, args, fcn, username, org_name) => {
    try {
        const ccp = await helper.getCCP(org_name);
        const walletPath = await helper.getWalletPath(org_name);
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`Identity for the user ${username} does not exist in the wallet. Registering...`);
            await helper.getRegisteredUser(username, org_name, true);
            identity = await wallet.get(username);
            if (!identity) {
                throw new Error(`❌ Failed to register user ${username}`);
            }
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true }
        });

        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        const resultBuffer = await contract.evaluateTransaction(fcn, ...args);
        const resultString = resultBuffer.toString();
        console.log(`✅ Transaction evaluated. Result: ${resultString}`);

        const parsedResult = JSON.parse(resultString);
        return parsedResult;

    } catch (error) {
        console.error(`❌ Failed to evaluate transaction: ${error}`);
        return { error: error.message }; // return error in expected format
    }
};


exports.query = query