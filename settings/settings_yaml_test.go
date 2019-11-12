package settings

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	settingsWithPort = `
port: 8081
authorization:
  basic_auth:
    username: user
    password: pass
`
	settingsWithoutPort = `
authorization:
  header:
    x-api-key: hello`
)

func TestReturnsSetPort(t *testing.T) {
	file := toFile(settingsWithPort)
	defer os.Remove(file.Name())

	sut, err := New(file)
	assert.Nil(t, err)
	assert.Equal(t, 8081, sut.Port())
}

func TestReturnsDefaultPort(t *testing.T) {
	file := toFile(settingsWithoutPort)
	defer os.Remove(file.Name())

	sut, err := New(file)
	assert.Nil(t, err)
	assert.Equal(t, 8080, sut.Port())
}

func toFile(body string) *os.File {
	file, err := ioutil.TempFile("", "gomocka-*.yaml")
	if err != nil {
		log.Fatal(err)
	}
	file.Write([]byte(body))
	filename := file.Name()
	file.Close()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return f
}
