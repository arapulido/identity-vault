[![Build Status][travis-image]][travis-url]
# Serial Vault

A Go web service that digitally signs device assertion details.

## Install
Go get it:

  ```bash
  $ go get github.com/ubuntu-core/identity-vault
  ```

Configure it:
- Install PostgreSQL and create a database.
- Set up the config file, using ```settings.yaml``` as a guide.
- Create the database tables:
  ```bash
  $ cd identity-vault
  $ go run tools/createdb.go
  ```

Run it:
  ```bash
  $ cd identity-vault
  $ go run server.go -config=/path/to/settings.yaml -mode=signing
  ```

The application has an admin service that can be run by using mode=admin.

## Try with docker
  ```bash
  $ git clone https://github.com/ubuntu-core/identity-vault
  $ cd identity-vault/
  $ docker-compose up
  # remove containers after try
  $ docker-compose rm
  ```

## Development Environment

### Install Go
Follow the instructions to [install Go](https://golang.org/doc/install).

### Install the React development environment
#### Pre-requisites
- Install the build packages
```bash
sudo apt-get install build-essential libssl-dev
# For TPM2.0
sudo apt-get install tpm2-tools
```

- Install NVM
Install the [Node Version Manager](https://github.com/creationix/nvm) that will allow a specific
version of Node.js to be installed. Follow the installation instructions.

- Install the latest stable Node.js and npm
The latest stable (LTS) version of Node can be found on the [Node website](nodejs.org).
```bash
# Overview of available commands
nvm help

# Install the latest stable version
nvm install v4.4.3

# Select the version to use
nvm ls
nvm use v4.4.3

# Install gulp globally
npm install -g gulp
```

- Install the nodejs dependencies
```bash
cd identity-vault
npm install
```

### Working with React

#### Build the project bundle
```bash
# Select the version to use
nvm ls
nvm use v4.4.3
gulp
```

#### Run the tests
```bash
npm test
```


## API Methods

### /1.0/version (GET)
> Return the version of the serial vault service.

#### Output message
```json
{
  "version":"0.1.0",
}
```
- version: the version of the serial vault service (string)


### /1.0/nonce (POST)
> Returns a nonce that is needed for the signing request.

#### Output message
```json
{
  "nonce": "abc123456",
  "success": true,
  "message": ""
}
```
- success: whether the request was successful (bool)
- message: error message from the request (string)
- nonce: unique string that is needed for signing requests (string)

The nonce can only be used once and must be used before it expires (typically 600 seconds).

### /1.0/sign (POST)
> Clear-sign the device identity details.

Takes the details from the device, formats the data and clear-signs it.

#### Input message
The message must be the serial-request assertion format and is best generated using the snapd libraries.
```
type: serial-request
brand-id: System
model: Router 3400
device-key:
    WkUDQbqFCKZBPvKbwR...
request-id: abc123456
body-length: 10
sign-key-sha3-384: UytTqTvREVhx...
revision: 12
serial: A1228ML

HW-DETAILS

AcLBUgQAAQoABgUCV7R2C...
```
- brand-id: the Account ID of the manufacturer (string)
- model: the name of the device (string)
- device-key: the encoded type and public key of the device (string)
- request-id: the nonce returned from the /1.0/nonce method (string)
- signature: the signed data
- serial: serial number of the device (string)
- revision: the revision of the device (integer)

Though the 'serial' and 'revision' headers are not strictly a part of a serial request, they are needed 
to return a signed serial assertion.

#### Output message
The method returns a signed serial assertion using the key from the vault.


[travis-image]: https://travis-ci.org/ubuntu-core/identity-vault.svg?branch=master
[travis-url]: https://travis-ci.org/ubuntu-core/identity-vault
