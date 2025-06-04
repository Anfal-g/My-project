'use strict';

const { Contract } = require('fabric-contract-api');
const { Shim } = require('fabric-shim');  // Import the shim
const ResidentManagement = require('./residentManagement');

// Instead of just exporting the contract,
// call Shim.start() with your contract instance so that the chaincode remains running.
Shim.start(new ResidentManagement());
