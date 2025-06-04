'use strict';

console.log('Chaincode server starting...');
process.on('uncaughtException', err => console.error('UNCAUGHT EXCEPTION:', err));
process.on('unhandledRejection', err => console.error('UNHANDLED REJECTION:', err));

const ResidentManagement = require('./residentManagement');

// Fabric 2.x will handle the server initialization
module.exports = ResidentManagement;
