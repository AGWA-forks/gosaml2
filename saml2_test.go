package saml2

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/russellhaering/gosaml2/types"
)

var cert tls.Certificate
var pk crypto.PrivateKey

func init() {
	var err error
	pfx := "./testdata/test"
	cert, err = tls.LoadX509KeyPair(fmt.Sprintf("%s.crt", pfx), fmt.Sprintf("%s.key", pfx))
	if err != nil {
		log.Fatal(err)
	}
	pk = cert.PrivateKey
}

func TestDecode(t *testing.T) {
	f, err := ioutil.ReadFile("./testdata/saml.post")
	if err != nil {
		t.Fatalf("could not open test file: %v\n", err)
	}
	decoded := make([]byte, len(f))

	base64.StdEncoding.Decode(decoded, f)
	response := &types.Response{}

	err = xml.Unmarshal(decoded, response)
	if err != nil {
		t.Fatalf("error decoding test saml: %v", err)
	}

	ea := response.EncryptedAssertions[0]

	k, err := ea.EncryptedKey.DecryptSymmetricKey(&cert)
	if err != nil {
		t.Fatalf("could not get symmetric key: %v\n", err)
	}

	if k == nil {
		t.Fatalf("no symmetric key")
	}

	bs, err := ea.Decrypt(&cert)
	if err != nil {
		t.Fatalf("error decrypting saml data: %v\n", err)
	}

	f2, err := ioutil.ReadFile("./testdata/saml.xml")
	if err != nil {
		t.Fatalf("could not read expected output")
	}

	if !bytes.Equal(f2, bs) {
		t.Errorf("decrypted assertion did not match expectation")
	}
}
