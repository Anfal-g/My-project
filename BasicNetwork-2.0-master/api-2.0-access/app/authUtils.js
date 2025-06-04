// auth/jwtUtils.js

const jwt = require('jsonwebtoken');
const SECRET_KEY = "thisismysecret"; // Store this securely
const generateResidentToken = (residentData) => {
  const payload = {
    sub: residentData.UserID,       // Unique identifier (username)
    role: 'resident',               // Default role for residents
    gender: residentData.Gender,
    apartment: residentData.Apartment,
    maritalStatus: residentData.MaritalStatus,
    residentType: residentData.ResidentType,
    qrCode: residentData.ResidentID, // From generatePermanentQRCode()
    iat: Math.floor(Date.now() / 1000),
    exp: Math.floor(Date.now() / 1000) + (60 * 60 * 2) // 2 hours expiry
  };

  return jwt.sign(payload,SECRET_KEY);
};

module.exports = { generateResidentToken };