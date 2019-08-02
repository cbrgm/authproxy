# Example app using authproxy

This is an example of using the authproxy with a fake provider and calling the API through a command line client.

## Getting started

### Generate tls certificates

As authproxy works with tls ecnryption only, fake certificates  must be generated. Please execute the following bashscript to generate a CA, server and client certificate.

```bash
export SAN=DNS.1:localhost

# generate ca, client and server certs + keys
./gencerts.sh
```

You will find the files `ca.crt`, `ca.key`, `server.crt`, `server.key`, `client.crt`, `client.key` in this directory as an output from the script

### Start the api

We can now use the certificates to start the api first. 

```
./api/api --tls-key server.key --tls-cert server.crt --tls-ca-cert ca.crt --log-level debug

```

On another terminal you can query the api using the cli client tool.

***login unsuccessful: username: foo, password: fail***
```bash 
./client/cli --tls-key client.key --tls-cert client.crt --tls-ca-cert ca.crt login foo fail
> failed to run api: unauthorized: invalid authentication credentials%  
```

***login successful: username: foo, password: bar***
```bash 
./client/cli --tls-key client.key --tls-cert client.crt --tls-ca-cert ca.crt login foo bar
> Received token for user: AbCdEf123456
```

***authenticate the token***
```bash 
./client/cli --tls-key client.key --tls-cert client.crt --tls-ca-cert ca.crt authenticate AbCdEf123456
client successfully authenticated, token is valid
```

