'use strict'; // Enforce strict mode to catch common coding mistakes

const { Contract } = require('fabric-contract-api'); // Import the Contract class from Hyperledger Fabric API

/**
 * ResidentManagement smart contract
 * Manages residents, visitors, delivery workers, and maintenance workers in a private collection ledger.
 */
class ResidentManagement extends Contract {


    /**
     * Initializes the ledger (executed when the chaincode is instantiated).
     */
    async initLedger(ctx) {
        console.log('ResidentManagement Chaincode Initialized');
        return;
    }

    /**
     * تسجيل دخول زائر بناءً على موافقة الساكن.
     * Registers a visitor entry if a resident approves.
     *
     * @param {Context} ctx - Hyperledger Fabric transaction context
     * @param {String} visitorId - Unique ID of the visitor
     * @param {String} residentId - ID of the resident approving entry
     * @param {String} approval - Approval status ('approved' required)
     */
    async visitorEntry(ctx, visitorId, residentId, approval) {
        if (approval !== 'approved') {
            throw new Error(`Visitor ${visitorId} is not approved by resident ${residentId}`);
        }

        const visitor = {
            visitorId,
            residentId,
            entryStatus: 'Approved',
            timestamp: Date.now()
        };

        // Store the visitor entry in the 'VisitorApprovalCollection' private data collection
        await ctx.stub.putPrivateData('VisitorApprovalCollection', visitorId, Buffer.from(JSON.stringify(visitor)));
        return JSON.stringify(visitor);
    }

    /**
     * تسجيل دخول زائر بناءً على موافقة الساكن والمدير.
     * Registers a visitor entry if both the resident and the manager approve.
     *
     * @param {Context} ctx
     * @param {String} visitorId
     * @param {String} residentId
     * @param {String} residentApproval - Approval from the resident
     * @param {String} managerApproval - Approval from the manager
     */
    async visitorEntryWithManager(ctx, visitorId, residentId, residentApproval, managerApproval) {
        if (residentApproval !== 'approved' || managerApproval !== 'approved') {
            throw new Error(`Visitor ${visitorId} needs both resident and manager approval`);
        }

        const visitor = {
            visitorId,
            residentId,
            entryStatus: 'Approved by Manager',
            timestamp: Date.now()
        };

        // Store visitor entry in private data collection
        await ctx.stub.putPrivateData('VisitorApprovalCollection', visitorId, Buffer.from(JSON.stringify(visitor)));
        return JSON.stringify(visitor);
    }

    /**
     * تسجيل مقيم دائم
     * Registers a new resident.
     *
     * @param {Context} ctx
     * @param {String} residentId - Unique ID of the resident
     * @param {String} name - Name of the resident
     */
    async registerResident(ctx, residentId, name) {
        // Check if resident already exists in the private collection
        const exists = await ctx.stub.getPrivateData('ResidentPrivateCollection', residentId);
        if (exists && exists.length > 0) {
            throw new Error(`Resident ${residentId} is already registered.`);
        }

        const resident = {
            residentId,
            name,
            status: 'Permanent Resident',
            timestamp: Date.now()
        };

        // Store resident in private data collection
        await ctx.stub.putPrivateData('ResidentPrivateCollection', residentId, Buffer.from(JSON.stringify(resident)));
        return JSON.stringify(resident);
    }

    /**
     * دخول فني الصيانة أو عامل الخدمات بموافقة الساكن.
     * Registers a maintenance worker entry if approved by a resident.
     *
     * @param {Context} ctx
     * @param {String} workerId - Unique ID of the worker
     * @param {String} residentId - ID of the resident approving the entry
     * @param {String} approval - Approval status ('approved' required)
     */
    async maintenanceEntry(ctx, workerId, residentId, approval) {
        if (approval !== 'approved') {
            throw new Error(`Maintenance worker ${workerId} is not approved by resident ${residentId}`);
        }

        const worker = {
            workerId,
            residentId,
            role: 'Maintenance/Service',
            entryStatus: 'Approved',
            timestamp: Date.now()
        };

        // Store maintenance worker entry in private data collection
        await ctx.stub.putPrivateData('ServiceWorkerCollection', workerId, Buffer.from(JSON.stringify(worker)));
        return JSON.stringify(worker);
    }

    /**
     * تسجيل موظف التوصيل في القائمة المؤقتة للزوار.
     * Registers a delivery worker with a temporary access period.
     *
     * @param {Context} ctx
     * @param {String} workerId - Unique ID of the worker
     * @param {Number} duration - Access duration in hours
     */
    async registerDeliveryWorker(ctx, workerId, duration) {
        const expiryTime = Date.now() + duration * 60 * 60 * 1000; // Convert hours to milliseconds

        const deliveryWorker = {
            workerId,
            role: 'Delivery Worker',
            entryStatus: 'Temporary Access',
            expiry: expiryTime
        };

        // Store delivery worker entry in private data collection
        await ctx.stub.putPrivateData('DeliveryWorkerCollection', workerId, Buffer.from(JSON.stringify(deliveryWorker)));
        return JSON.stringify(deliveryWorker);
    }

    /**
     * تحديث حالة الدخول
     * Updates the status of an existing entry in the private collection.
     *
     * @param {Context} ctx
     * @param {String} collectionName - The name of the private collection
     * @param {String} id - Unique ID of the entry
     * @param {String} newStatus - The new status to be set
     */
    async updateEntryStatus(ctx, collectionName, id, newStatus) {
        const entryAsBytes = await ctx.stub.getPrivateData(collectionName, id);
        if (!entryAsBytes || entryAsBytes.length === 0) {
            throw new Error(`Entry with ID ${id} not found`);
        }

        let entry = JSON.parse(entryAsBytes.toString());
        entry.entryStatus = newStatus;
        entry.updatedAt = Date.now();

        // Store updated entry in private collection
        await ctx.stub.putPrivateData(collectionName, id, Buffer.from(JSON.stringify(entry)));
        return JSON.stringify(entry);
    }

    /**
     * حذف تسجيل دخول
     * Deletes an entry from the private collection.
     *
     * @param {Context} ctx
     * @param {String} collectionName - The name of the private collection
     * @param {String} id - Unique ID of the entry to be deleted
     */
    async deleteEntry(ctx, collectionName, id) {
        const exists = await ctx.stub.getPrivateData(collectionName, id);
        if (!exists || exists.length === 0) {
            throw new Error(`Entry with ID ${id} does not exist`);
        }

        // Delete entry from private collection
        await ctx.stub.deletePrivateData(collectionName, id);
        return `Entry ${id} deleted successfully.`;
    }

    /**
     * الاستعلام عن بيانات الدخول
     * Retrieves an entry from the private collection by ID.
     *
     * @param {Context} ctx
     * @param {String} collectionName - The name of the private collection
     * @param {String} id - Unique ID of the entry
     */
    async queryEntry(ctx, collectionName, id) {
        const entryAsBytes = await ctx.stub.getPrivateData(collectionName, id);
        if (!entryAsBytes || entryAsBytes.length === 0) {
            throw new Error(`Entry with ID ${id} not found`);
        }
        return entryAsBytes.toString();
    }

    /**
     * استعلام عن جميع الإدخالات
     * Retrieves all entries from a given private collection.
     *
     * @param {Context} ctx
     * @param {String} collectionName - The name of the private collection
     */
    async queryAllEntries(ctx, collectionName) {
        const iterator = await ctx.stub.getPrivateDataByRange(collectionName, '', '');
        const results = [];

        for await (const res of iterator) {
            results.push(JSON.parse(res.value.toString()));
        }

        return JSON.stringify(results);
    }
}

// Export the ResidentManagement contract so it can be used in Hyperledger Fabric
 module.exports = ResidentManagement;
