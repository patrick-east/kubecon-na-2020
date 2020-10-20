package authz

default allow = false

allow {
    jwt := input.jwt
    jwks := custom.fetch_jwks("https://dev-vpda1l4e.us.auth0.com")

    # Verify signature only
    io.jwt.verify_rs256(jwt, jwks)

    # Check specific parts of the payload
    [_, payload, _] := io.jwt.decode(jwt)
    payload.iss == "https://dev-vpda1l4e.us.auth0.com/"
    payload.aud == "https://opa4fun.com/myapp"
}
