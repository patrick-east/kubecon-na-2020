package builtins

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

func Register() {
	rego.RegisterBuiltin1(myBuiltinDecl, myBuiltinImpl)
}


var myBuiltinDecl = &rego.Function{
	Name: "custom.fetch_jwks",
	Decl: types.NewFunction(
		types.Args(types.S),    // Single string argument
		types.S),               // Returns a string
}

// Use a custom cache key type to avoid collisions with other builtins caching data!!
type myCacheKeyType string

// myBuiltinImpl will attempt to fetch a jwks using well known locations for
// the key using a base url. On an error it will *not* report an error, but
// will be undefined.
func myBuiltinImpl(bctx rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
	var baseURL string
	if err := ast.As(a.Value, &baseURL); err != nil {
		return nil, err
	}

	// Check if it is already cached, assume they never become invalid.
	var cacheKey = myCacheKeyType(baseURL)
	cached, ok := bctx.Cache.Get(cacheKey)
	if ok {
		return ast.NewTerm(cached.(ast.Value)), nil
	}

	// Guess the JWKS URL
	jwksURL := baseURL + "./well-known/jwks.json"

	// See if there is an openid config, use the `jwks_url` if available
	resp, err := http.Get(baseURL + "/.well-known/openid-configuration")
	if err == nil {
		bs, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			// Parse out only the jwks_uri, if it exists.
			var config struct{
				URI string `json:"jwks_uri"`
			}
			err = json.Unmarshal(bs, &config)
			if err == nil {
				jwksURL = config.URI
			}
		}
	}

	resp, err = http.Get(jwksURL)
	if err == nil {
		bs, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			// Return the JWKS string
			return ast.StringTerm(string(bs)), nil
		}
	}

	// undefined
	return nil, nil
}