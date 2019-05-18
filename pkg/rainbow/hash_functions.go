package rainbow

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
)

type HashFunction struct {
	hash.Hash
}

func (hf HashFunction) Apply(plaintext string) []byte {
	io.WriteString(hf, plaintext)
	result := hf.Sum(nil)
	hf.Reset()
	return result
}

type HashFunctionProvider struct {
	newFunc func() hash.Hash
}

func (provider HashFunctionProvider) NewHashFunction() HashFunction {
	return HashFunction{Hash: provider.newFunc()}
}

const (
	MD5    = "MD5"
	SHA1   = "SHA1"
	SHA256 = "SHA256"
	SHA384 = "SHA384"
	SHA512 = "SHA512"
)

var hashFunctionProvidersByName = map[string]HashFunctionProvider{
	MD5:    {newFunc: md5.New},
	SHA1:   {newFunc: sha1.New},
	SHA256: {newFunc: sha256.New},
	SHA384: {newFunc: sha512.New384},
	SHA512: {newFunc: sha512.New},
}

func GetHashFunctionProvider(hashFunctionName string) (HashFunctionProvider, error) {
	hashFunctionProvider, found := hashFunctionProvidersByName[hashFunctionName]
	if !found {
		return HashFunctionProvider{}, fmt.Errorf("invalid hash function provided %s", hashFunctionName)
	}

	return hashFunctionProvider, nil
}
