# Simple gateway to Keycloak

Intermediary between your app and Keycloak

To configure this gateway use file `config.env`

## Endpoints

- POST: / : endpoints list
- POST: /token/get : get token from Keycloak
- POST: /token/verify : verify Keycloak token

## Examples
#### Verify token
`curl -X POST http://localhost:8787/token/verify -H 'Authorization: Bearer eyJhbG...qzm61uA`
```json
{
  "sub": "c636d74c-fa1c-4d79-b3de-35edbf3a3434",
  "email_verified": true,
  "name": "Reiterus Github",
  "preferred_username": "reiterus",
  "given_name": "Reiterus",
  "family_name": "Github",
  "email": "some@email"
}
```

#### Get token
`curl -X POST http://localhost:8787/token/get`
```json
{
  "access_token": "eyJhbG...EIiwg",
  "expires_in": 18000,
  "refresh_expires_in": 36000,
  "refresh_token": "eyJhbG...qz1s",
  "token_type": "Bearer",
  "not-before-policy": 0,
  "session_state": "ead9fa51-a437-4e60-ab4f-87d5adeef655",
  "scope": "profile email"
}
```

#### Get token only
`curl -X POST http://localhost:8787/token/only`
```shell
eyJhbG...EIiwg
```

#### Endpoints list
`curl -X POST http://localhost:8787`
```json
[
  {
    "method": "POST",
    "path": "/",
    "name": "main.main.func1"
  },
  {
    "method": "POST",
    "path": "/token/get",
    "name": "main.tokenGet"
  },
  {
    "method": "POST",
    "path": "/token/verify",
    "name": "main.tokenVerify"
  },
  {
    "method": "POST",
    "path": "/token/only",
    "name": "main.tokenOnly"
  }
]
```