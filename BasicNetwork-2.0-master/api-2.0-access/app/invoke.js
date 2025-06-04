const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')

const helper = require('./helper')

const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name) => {
    try {
        // 1. Get wallet path (will throw error if invalid)
        const walletPath = await helper.getWalletPath(org_name);
        console.log(`Using wallet at: ${walletPath}`);

        // 2. Initialize wallet
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        
        // 3. Check user exists
        const identity = await wallet.get(username);
        if (!identity) {
            throw new Error(`User ${username} not found in wallet at ${walletPath}`);
        }

        // 4. Get network connection
        const gateway = new Gateway();
        await gateway.connect(await helper.getCCP(org_name), {
            wallet,
            identity: username,
            discovery: { enabled: true, asLocalhost: true }
        });

        // 5. Submit transaction
        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);
        const result = await contract.submitTransaction(fcn, ...args);
        
        // 6. Parse response
        try {
            return JSON.parse(result.toString());
        } catch {
            return { status: "success", data: result.toString() };
        }

    } catch (error) {
        console.error(`Transaction failed: ${error.message}`);
        throw error;
    }
};
exports.invokeTransaction = invokeTransaction;