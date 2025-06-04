FROM hyperledger/fabric-nodeenv:2.4
WORKDIR /usr/local/src
COPY package.json .  
RUN npm install  
COPY . .  
CMD ["node", "index.js"]
