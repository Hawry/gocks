package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hawry/gomocka/jwt"
)

const (
	defaultPort          = 8080
	defaultGeneratedFile = "example.json"
)

//Settings describe the settings.json file contents
type Settings struct {
	ListenPort    int      `json:"port" yaml:"port"`
	Authorization AuthData `json:"authorization,omitempty" yaml:"authorization,omitempty"`
	Mocks         []Mock   `json:"mocks" yaml:"mocks"`
}

// AuthData describes the possible authorization headers required
type AuthData struct {
	Basic  BasicAuth         `json:"basic_auth" yaml:"basic_auth"`
	Header map[string]string `json:"header" yaml:"header"`
	OpenID OpenIDAuth        `json:"openid" yaml:"openid"`
}

//BasicAuth contains username/password field in settings file
type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

//OpenIDAuth contains an endpoint to JWKS keys
type OpenIDAuth struct {
	JWKSEndpoint string `json:"jwks" yaml:"jwks"`
}

//New returns a new configuration from the given file
func New(f *os.File) (s *Settings, err error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	if strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml") {
		err = yaml.Unmarshal(b, &s)
	} else {
		err = json.Unmarshal(b, &s)
	}
	return
}

//Port returns the port number from the configuration, or a default port number if not set
func (s *Settings) Port() int {
	if s.ListenPort > 0 {
		return s.ListenPort
	}
	return defaultPort
}

//VerifyBasicAuth returns true if username & password matches the settings file
func (s *Settings) VerifyBasicAuth(username, password string) bool {
	if !s.hasBasicAuthEnabled() {
		return false
	}
	return (s.Authorization.Basic.Username == username && s.Authorization.Basic.Password == password)
}

//VerifyHeaderAuth returns true if the specified header has the correct value in the request
func (s *Settings) VerifyHeaderAuth(h http.Header) bool {
	if !s.hasHeaderAuthEnabled() {
		return false
	}
	for k, v := range s.Authorization.Header {
		if h.Get(k) == v {
			return true
		}
	}
	return false
}

//VerifyOpenIDAuth returns true if the given bearer token have been signed with the specified key(s)
func (s *Settings) VerifyOpenIDAuth(h http.Header) bool {
	if !s.hasOpenIDAuthEnabled() {
		return false
	}
	err := jwt.VerifyToken(h.Get("Authorization"), s.Authorization.OpenID.JWKSEndpoint)
	if err != nil {
		log.Printf("warn: %v", err)
		return false
	}
	return true
}

//RequireAuthentication returns true if the settings require authentication to be provided, or false if not
func (s *Settings) RequireAuthentication() bool {
	return s.hasBasicAuthEnabled() || s.hasHeaderAuthEnabled() || s.hasOpenIDAuthEnabled()
}

func (s *Settings) hasBasicAuthEnabled() bool {
	return (len(s.Authorization.Basic.Password) > 0 && len(s.Authorization.Basic.Username) > 0)
}

func (s *Settings) hasHeaderAuthEnabled() bool {
	return (len(s.Authorization.Header) > 0)
}

func (s *Settings) hasOpenIDAuthEnabled() bool {
	return (len(s.Authorization.OpenID.JWKSEndpoint) > 0)
}

//CreateDefault creates a default configuration and returns the struct representation as well, default filename is example.json
func CreateDefault() (*Settings, error) {
	f, err := os.Create(defaultGeneratedFile)
	if err != nil {
		return nil, err
	}
	s := Settings{
		ListenPort: 8080,
		Mocks: []Mock{
			Mock{
				Path:         "/",
				Method:       "GET",
				ResponseBody: "{\"hello\":\"world\"}",
				ResponseCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			Mock{
				Path:         "/users/{userid}",
				Method:       "GET",
				ResponseBody: "{\"user\":\"{userid}\"}",
				ResponseCode: 200,
			},
			Mock{
				Path:         "/health",
				Method:       "GET",
				ResponseBody: "{\"status\": \"OK\"}",
				ResponseCode: 200,
				DisableAuth:  true,
			},
		},
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(&s)
	return &s, err
}
