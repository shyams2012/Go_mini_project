# Mini Project

# How to run project
Create a database having name "mini_project" in mysql

& run below command 

```bash 
    $ go run .
```


## This is sample project not meant project production use.

1. Provide a REST endpoint that authenticates a user credentials. The endpoint receives a user’s email address and password and compares the passed email/password against what is stored in the database. If the operation succeeds, return a JWT in the response payload. The JWT should be asymmetrically signed with a private key generated by an appropriate algorithm such as RSA or ECDSA.

2. Provide an authorized REST endpoint that returns the user profile matching the bearer of the request. In other words, we need to embed the bearer token in JWT format in the Authorization HTTP header of the request eg. Authorization: Bearer <JWT here>. If the user is authorized, we return the user profile. Keep the user profile simple with no more than 4 user profile attributes eg. email, name, location, etc. Validate the JWT by using a public key.

3. Persist the users and key pairs in a relational database ie. sqlite should suffice. To keep it simple, the server can seed the data in the data store when the API starts up.

4. If you end up using asymmetrical signing, then we need to have a private/public key pair. We can generate the public/private key pair using a script or/and Makefile. Optionally, provide a JWKS endpoint (see section below) that vents out a set of JWKs. Each JWK contains the public key. A common best practice is to rotate the keys. This means that the JWKS endpoint will return more than one public key with expirations that overlap each other. You can generate the public/private key pairs and persist them in the data store, and remove them in the background when keys expire. NOTE: Since this is a monolith API server, a JWKS endpoint really has no practical use other than a coding exercise. In the real world, we may use JWKS in a microservice system or secure webhook setup.
