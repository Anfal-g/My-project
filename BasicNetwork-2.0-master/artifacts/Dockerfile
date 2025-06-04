# Use the base Hyperledger Fabric container image
FROM hyperledger/fabric-ccenv:2.4

# Switch to root user in case it's needed to install packages
USER root

# Install curl and bash first, which are necessary for the next steps using apk
RUN apk update && \
    apk add --no-cache curl bash

# Install Node.js and npm (needed for your chaincode)
RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apk add --no-cache nodejs

# Set up the working directory
WORKDIR /usr/local/src

