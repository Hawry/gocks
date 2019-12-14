package settings

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

const (
	settingsWithBasicAuth = `{
		"port": 8080,
		"authorization": {
			"basic_auth": {
				"username":"ausername",
				"password":"apassword"
			}
		}}`
	settingsWithHeaderAuth = `{
		"port": 8080,
		"authorization": {
			"header": {
				"x-api-key": "thisisanauthheader"
			}
		}}`
	settingsWithoutAuth = `{
		"port": 8080
		}`
	settingsWithOpenIDAuth = `{
		 "authorization": {
			"openid": {
				"jwks":"http://localhost/.well-known/openid-configuration"
			}
		 }
		}`
)

func toStruct(js string) *Settings {
	var s Settings
	err := json.Unmarshal([]byte(js), &s)
	if err != nil {
		log.Fatal(err)
	}
	return &s
}

func TestRequireAuthentication_BasicAuth(t *testing.T) {
	sut := toStruct(settingsWithBasicAuth)
	if !sut.RequireAuthentication() {
		t.Fail()
	}
	if actual := sut.hasBasicAuthEnabled(); !actual {
		t.Errorf("expected %t, got %t", true, actual)
	}
}

func TestRequireAuthentication_HeaderAuth(t *testing.T) {
	sut := toStruct(settingsWithHeaderAuth)
	if !sut.RequireAuthentication() {
		t.Fail()
	}
	if actual := sut.hasHeaderAuthEnabled(); !actual {
		t.Errorf("expected %t, got %t", true, actual)
	}
}

func TestRequireAuthentication_NoAuth(t *testing.T) {
	sut := toStruct(settingsWithoutAuth)
	if sut.RequireAuthentication() {
		t.Fail()
	}
}

func TestRequireAuthentication_OpenIDAuth(t *testing.T) {
	sut := toStruct(settingsWithOpenIDAuth)
	if !sut.RequireAuthentication() {
		t.Fail()
	}
	if actual := sut.hasOpenIDAuthEnabled(); !actual {
		t.Errorf("expected %t, got %t", true, actual)
	}
}

func TestVerifyBasicAuth_CorrectData(t *testing.T) {
	sut := toStruct(settingsWithBasicAuth)
	if !sut.VerifyBasicAuth("ausername", "apassword") {
		t.Fail()
	}
}

func TestVerifyBasicAuth_IncorrectData(t *testing.T) {
	sut := toStruct(settingsWithBasicAuth)
	if sut.VerifyBasicAuth("wrong", "apassword") {
		t.Fail()
	}
}

func TestVerifyHeaderAuth_CorrectData(t *testing.T) {
	r := &http.Request{}
	r.Header = http.Header{}
	r.Header.Add("x-api-key", "thisisanauthheader")
	sut := toStruct(settingsWithHeaderAuth)
	if !sut.VerifyHeaderAuth(r.Header) {
		t.Fail()
	}
}

func TestVerifyHeaderAuth_IncorrectData(t *testing.T) {
	r := &http.Request{}
	r.Header = http.Header{}
	r.Header.Add("x-api-key", "thisiswrongheader")
	sut := toStruct(settingsWithHeaderAuth)
	if sut.VerifyHeaderAuth(r.Header) {
		t.Fail()
	}
}

func TestVerifyOpenIDAuth_IncorrectData(t *testing.T) {
	r := &http.Request{}
	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")
	sut := toStruct(settingsWithOpenIDAuth)
	if sut.VerifyOpenIDAuth(r.Header) {
		t.Fail()
	}
}
