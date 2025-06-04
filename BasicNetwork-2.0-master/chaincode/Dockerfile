# Set base image to Node.js 16.20
ARG NODE_IMAGE_TAG=16.20-alpine
FROM node:${NODE_IMAGE_TAG}

# Install build dependencies
RUN apk add --no-cache make python3 g++

# Create necessary directories
RUN mkdir -p /artifacts/src/github.com/residentManagement/javascript/input && mkdir -p /artifacts/src/github.com/residentManagement/javascript/output && mkdir -p /usr/local/src

# Add custom build/start scripts
ADD  start.sh /artifacts/src/github.com/residentManagement/javascript/

# Install any necessary chaincode dependencies (example: if you are using fabric-shim)
WORKDIR /artifacts/src/github.com/residentManagement/javascript
RUN npm install fabric-shim

# Set default command to start node
CMD ["node"]
