'use strict';
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const bodyParser = require('body-parser');
const http = require('http')
const util = require('util');
const express = require('express')
const app = express();
const { expressjwt: expressJWT } = require('express-jwt');
const jwt = require('jsonwebtoken');
const bearerToken = require('express-bearer-token');
const cors = require('cors');
const constants = require('./config/constants.json')

const host = process.env.HOST || constants.host;
const port = process.env.PORT || constants.port;



// Add this line 👇 before your routes
app.use(express.json());
app.use(express.json({ limit: '10mb' }));  // Increased limit for larger payloads
app.use(express.urlencoded({ extended: true }));

const helper = require('./app/helper')
const invoke = require('./app/invoke')
const qscc = require('./app/qscc')
const query = require('./app/query')
app.options('*', cors());
app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: false
}));
// set secret variable
app.set('secret', 'thisismysecret');
app.use(expressJWT({
    secret: 'thisismysecret',
    algorithms: ['HS256']
}).unless({
    path: ['/users','/users/login', '/register']
}));
app.use(bearerToken());

logger.level = 'debug';

app.use((req, res, next) => {
    logger.debug('New req for %s', req.originalUrl);

    if (req.originalUrl.indexOf('/users') >= 0 || req.originalUrl.indexOf('/users/login') >= 0 || req.originalUrl.indexOf('/register') >= 0) {
        return next();
    }

    const bearerHeader = req.headers['authorization'];
    if (typeof bearerHeader === 'undefined') {
        return res.status(403).send({ success: false, message: 'No token provided.' });
    }

    const token = bearerHeader.split(' ')[1]; // Extract token from "Bearer <token>"
    
    jwt.verify(token, app.get('secret'), (err, decoded) => {
        if (err) {
            console.log(`JWT Error: ${err}`);
            return res.status(403).send({
                success: false,
                message: 'Failed to authenticate token.'
            });
        } else {
            req.username = decoded.sub;
            req.orgname = decoded.orgName;
            console.log("Decoded token:", decoded);
            logger.debug(util.format('Decoded from JWT token: username - %s, orgname - %s', decoded.sub, decoded.orgName));
            return next();
        }
    });
});
// Increase JSON parser limits
app.use(express.json({ 
    limit: '10mb',
    verify: (req, res, buf) => {
        req.rawBody = buf;
    }
}));
app.use(express.urlencoded({ extended: true, limit: '10mb' }));

// Add error handling middleware
app.use((err, req, res, next) => {
    if (err instanceof SyntaxError && err.status === 400 && 'body' in err) {
        return res.status(400).json({ 
            success: false, 
            message: 'Invalid JSON payload' 
        });
    }
    next();
});
var server = http.createServer(app).listen(port, function () { console.log(`Server started on ${port}`) });
logger.info('****************** SERVER STARTED ************************');
logger.info('***************  http://%s:%s  ******************', host, port);
server.timeout = 240000;

function getErrorMessage(field) {
    var response = {
        success: false,
        message: field + ' field is missing or Invalid in the request'
    };
    return response;
}

app.post('/users', async function (req, res) {
      console.log("Incoming body:", req.body); // ✅ Add this
  const { username, orgName, role = 'resident' } = req.body;

  // Validate input
  if (!username || !orgName) {
    return res.status(400).json({ success: false, message: 'Missing username or orgName' });
  }

  const validRoles = ['resident', 'admin', 'visitor'];
  if (!validRoles.includes(role)) {
    return res.status(400).json({ success: false, message: 'Invalid role' });
  }

  try {
    // 1. Register user in Fabric CA
    const caResponse = await helper.getRegisteredUser(username, orgName, role, true);
    if (typeof caResponse === 'string') {
      throw new Error(caResponse);
    }

    let userData = { username, orgName, role };
    let tokenPayload = { sub: username, orgName, role };

    // 2. Register resident-specific data if role is 'resident'
    if (role === 'resident') {
      const residentArgs = [
        username,                      // ResidentID
        req.body.name || '',          // Name
        req.body.email || '',         // Email
        req.body.phone || '',         // Phone
        req.body.gender || '',        // Gender
        req.body.maritalStatus || '', // MaritalStatus
        req.body.residentType || '',  // ResidentType
        req.body.apartment || ''      // Apartment
      ];

      const invokeResponse = await invoke.invokeTransaction(
        'residentschannel',
        'residentManagement',
        'RegisterResident',
        residentArgs,
        username,
        orgName
      ).catch(err => {
        console.error('Invoke Error Details:', err);
        throw new Error(`Chaincode invocation failed: ${err.message}`);
      });

      let residentData;
      try {
        residentData = JSON.parse(invokeResponse.result.toString());
      } catch {
        residentData = {
          UserID: username,
          Name: req.body.name,
          Email: req.body.email,
          Phone: req.body.phone,
          Gender: req.body.gender,
          Apartment: req.body.apartment,
          MaritalStatus: req.body.maritalStatus,
          ResidentType: req.body.residentType,
          ResidentID: username
        };
      }

      tokenPayload = {
        ...tokenPayload,
        name: residentData.Name,
        email: residentData.Email,
        phone: residentData.Phone,
        gender: residentData.Gender,
        apartment: residentData.Apartment,
        maritalStatus: residentData.MaritalStatus,
        residentType: residentData.ResidentType,
        qrCode: residentData.ResidentID || `QR-RESIDENT-${username}`
      };

      userData = { ...userData, ...residentData };
    }

    // 3. Generate JWT
    const token = jwt.sign(
      tokenPayload,
      app.get('secret'),
      { expiresIn: constants.jwt_expiretime }
    );

    res.json({
      success: true,
      token,
      user: userData
    });

  } catch (error) {
    res.status(500).json({
      success: false,
      message: error.message
    });
  }
});

// Register and enroll user
app.post('/register', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName
    }, app.get('secret'));

    console.log(token)

    let response = await helper.registerAndGerSecret(username, orgName);

    logger.debug('-- returned from registering the username %s for organization %s', username, orgName);
    if (response && typeof response !== 'string') {
        logger.debug('Successfully registered the username %s for organization %s', username, orgName);
        response.token = token;
        res.json(response);
    } else {
        logger.debug('Failed to register the username %s for organization %s with::%s', username, orgName, response);
        res.json({ success: false, message: response });
    }

});

// Login and get jwt
app.post('/users/login', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName
    }, app.get('secret'));

    let isUserRegistered = await helper.isUserRegistered(username, orgName);

    if (isUserRegistered) {
        res.json({ success: true, message: { token: token } });

    } else {
        res.json({ success: false, message: `User with username ${username} is not registered with ${orgName}, Please register first.` });
    }
});


// Invoke transaction on chaincode on target peers
// Invoke transaction on chaincode on target peers
app.post('/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var peers = req.body.peers;
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var fcn = req.body.fcn;
        var args = req.body.args;
        var transient = req.body.transient;
        console.log(`Transient data is ;${transient}`)
        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        var org_name = req.headers['orgname'];
        var username = req.headers['username'];

let message = await invoke.invokeTransaction(channelName, chaincodeName, fcn, args, username, org_name, transient);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});


app.get('/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;
        let peer = req.query.peer;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let message = await query.query(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/qscc/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;
        // let peer = req.query.peer;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let response_payload = await qscc.qscc(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        // const response_payload = {
        //     result: message,
        //     error: null,
        //     errorData: null
        // }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});
