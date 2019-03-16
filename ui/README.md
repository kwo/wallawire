# Wallawire UI

### Build

    yarn install
    yarn build

## Authentication Mechanism
The client posts the username and password as JSON to the backend authentication endpoint.
If successful, the service responds with a 200 OK and a cookie "auth" containing a JWT.
The client can check if it is authenticated by checking the for presence of the "auth" cookie
and by parsing the the JWT on the client, checking if the expiration date is past the current time.
Other possible return codes are
 - 400 Bad Request (no content-type application/json, missing payload or payload missing fields or fields empty)
 - 401 Unauthorized (bad username and/or password)
 - 403 Forbidden (correct username and password but incorrect role or permissions).
