package rainbow

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
)

const (
	MD5    = "MD5"
	SHA1   = "SHA1"
	SHA256 = "SHA256"
	SHA384 = "SHA384"
	SHA512 = "SHA512"
)

type hashFunction struct {
	hash.Hash
}

type hashFunctionProvider struct {
	newFunc func() hash.Hash
}

func (provider hashFunctionProvider) newHashFunction() hashFunction {
	return hashFunction{Hash: provider.newFunc()}
}

func (hf hashFunction) apply(plaintext string) []byte {
	io.WriteString(hf, plaintext)
	result := hf.Sum(nil)
	hf.Reset()
	return result
}

var hashFunctionProvidersByName = map[string]hashFunctionProvider{
	MD5:    {newFunc: md5.New},
	SHA1:   {newFunc: sha1.New},
	SHA256: {newFunc: sha256.New},
	SHA384: {newFunc: sha512.New384},
	SHA512: {newFunc: sha512.New},
}
