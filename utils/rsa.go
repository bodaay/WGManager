package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return &privkey.PublicKey, privkey, nil
}

// GenerateKeyPair2048bits generates a new key pair with key size 2048
func GenerateKeyPair2048bits() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return &privkey.PublicKey, privkey, nil
}

// GenerateKeyPair4096bits generates a new key pair with key size 4096
func GenerateKeyPair4096bits() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	return &privkey.PublicKey, privkey, nil
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// PublicKeyToHexEncodedString Convert public key to Hex Encoded String
func PublicKeyToHexEncodedString(pub *rsa.PublicKey) (string, error) {
	pbytes, err := PublicKeyToBytes(pub)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pbytes), nil
}

// PrivateKeyToHexEncodedString Convert private key to Hex Encoded String
func PrivateKeyToHexEncodedString(priv *rsa.PrivateKey) string {
	pbytes := PrivateKeyToBytes(priv)
	return hex.EncodeToString(pbytes)
}

// PublicKeyFromHexEncodedString Convert public key from Hex Encoded String
func PublicKeyFromHexEncodedString(phex string) (*rsa.PublicKey, error) {
	pbytes, err := hex.DecodeString(phex)
	if err != nil {
		return nil, err
	}
	pub, err := BytesToPublicKey(pbytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

// PrivateKeyFromHexEncodedString Convert private key from Hex Encoded String
func PrivateKeyFromHexEncodedString(phex string) (*rsa.PrivateKey, error) {
	pbytes, err := hex.DecodeString(phex)
	if err != nil {
		return nil, err
	}
	priv, err := BytesToPrivateKey(pbytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}
	return key, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// SignEncryptedDataWithPrivateKey Sign the Message With Private key
func SignEncryptedDataWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, string, error) {
	rng := rand.Reader

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := sha256.Sum256(ciphertext)

	signature, err := rsa.SignPKCS1v15(rng, priv, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, "", err
	}
	return signature, hex.EncodeToString(signature), err
}

// VerifyEncodedSignatureWithPublicKey Sign the Message With Public key, if you pass encoded signa
func VerifyEncodedSignatureWithPublicKey(ciphertext []byte, signature string, pub *rsa.PublicKey) error {
	signatureData, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}
	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := sha256.Sum256(ciphertext)

	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signatureData)
	if err != nil {
		return err
	}
	return nil
}

// VerifySignatureWithPublicKey Sign the Message With Public key, if you pass encoded signed
func VerifySignatureWithPublicKey(ciphertext []byte, signature []byte, pub *rsa.PublicKey) error {

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := sha256.Sum256(ciphertext)

	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return err
	}
	return nil
}

/*
Full Example for main driver showing how I encrypt aes key, sign the message, verify the encrypted data with signature:

	aesKey := utils.GenerateRandomAESKey256bitBase64Encoded() // you can of course encrypt any data, not only AES key
	fmt.Println(aesKey)
	public,private, err := utils.GenerateKeyPair2048bits()
	if err != nil {
		log.Fatal(err)
	}
	encrypted, err := utils.EncryptWithPublicKey([]byte(aesKey), public)
	if err != nil {
		log.Fatal(err)
	}
	decrypted, err := utils.DecryptWithPrivateKey(encrypted, private)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(decrypted))
	singByte, signEncodedString, err := utils.SignEncryptedDataWithPrivateKey(encrypted, private)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(singByte))
	fmt.Println(signEncodedString)
	err = utils.VerifyEncodedSignatureWithPublicKey(encrypted, signEncodedString, public)
	if err != nil {
		log.Fatal(err)
	}

*/
