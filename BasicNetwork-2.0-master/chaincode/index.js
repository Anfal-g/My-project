'use strict';

const ResidentManagement = require('./residentManagement');   // Import the ResidentManagement contract
const { Contract } = require('fabric-contract-api');


module.exports.contracts = [ResidentManagement]; // Export the contract for Fabric to recognize

