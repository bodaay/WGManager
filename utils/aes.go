package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

//GenerateRandomBytes GenerateRandomBytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// //GenerateRandomBytes128Bit GenerateRandomBytes with key size 256-bit
// func GenerateRandomBytes128Bit() ([]byte, error) {
// 	b := make([]byte, 16)
// 	_, err := rand.Read(b)
// 	// Note that err == nil only if we read len(b) bytes.
// 	if err != nil {
// 		return nil, err
// 	}

// 	return b, nil
// }

//GenerateRandomBytes256Bit GenerateRandomBytes with key size 256-bit
func GenerateRandomAESKey256bit() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

//GenerateRandomAESKey256bitBase64Encoded GenerateRandomAESKeyBase64Encoded
func GenerateRandomAESKey256bitBase64Encoded() string {
	key, _ := GenerateRandomString(32)
	return key
}

// //GenerateRandomAESKey128bitBase64Encoded GenerateRandomAESKeyBase64Encoded
// func GenerateRandomAESKey128bitBase64Encoded() string {
// 	key, _ := GenerateRandomString(16)
// 	return key
// }

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// AESKey256BitToHexEncodedString Convert AES Key bytes to hex encoded string
func AESKey256BitToHexEncodedString(aeskey []byte) (string, error) {
	if len(aeskey) != 32 {
		return "", errors.New("Key size must be 256 bit")
	}
	key := hex.EncodeToString(aeskey)
	return key, nil
}

// AESKey256BitFromHexEncodedString Convert AES Key bytes from hex encoded string
func AESKey256BitFromHexEncodedString(aeskey string) ([]byte, error) {
	if len(aeskey) != 64 { //since its hex string, it will be 64 charcters
		return nil, errors.New("Key size must be 256 bit")
	}
	key, err := hex.DecodeString(aeskey)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Most the below code taken from: https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

// CreateHashStringMD5 create hash string using MD5, this can later be used as AES-256 key
func CreateHashStringMD5(passphrase string) (KeyEncodedHexString string, KeyBytes []byte) {
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	data := hasher.Sum(nil)
	return hex.EncodeToString(data), data
}

// CreateHashStringSHA256 create hash string using sha256, this can later be used as AES-256 key
func CreateHashStringSHA256(passphrase string) (KeyEncodedHexString string, KeyBytes []byte) {
	hasher := sha256.New()
	hasher.Write([]byte(passphrase))
	data := hasher.Sum(nil)
	return hex.EncodeToString(data), data
}

// Encrypt Data with AES256-Bit Key
func Encrypt(data []byte, AESKey256Bit []byte) ([]byte, error) {
	block, _ := aes.NewCipher(AESKey256Bit)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt Data with AES256-Bit Key
func Decrypt(data []byte, AESKey256Bit []byte) ([]byte, error) {
	key := AESKey256Bit
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// EncryptFile Decrypt File with AES256-Bit Key
func EncryptFile(sourcefileName string, encryptedOutputFileName string, AESKey256Bit []byte) error {
	data, err := ioutil.ReadFile(sourcefileName)
	if err != nil {
		return err
	}
	encrypted, err := Encrypt(data, AESKey256Bit)
	if err != nil {
		return err
	}
	f, err := os.Create(encryptedOutputFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(encrypted)
	if err != nil {
		return err
	}
	return nil
}

// DecryptFile Decrypt File with AES256-Bit Key
func DecryptFile(encryptedInputFileName string, decryptedOutFileName string, AESKey256Bit []byte) error {
	data, err := ioutil.ReadFile(encryptedInputFileName)
	if err != nil {
		return err
	}
	f, err := os.Create(decryptedOutFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	decrypted, err := Decrypt(data, AESKey256Bit)
	if err != nil {
		return err
	}
	_, err = f.Write(decrypted)
	if err != nil {
		return err
	}
	return nil
}

/* Full Test With RSA
aesKey, err := utils.GenerateRandomAESKey256bit() // you can ofcourse encrypt any data, not only AES key
if err != nil {
	log.Fatal(err)
}
fmt.Println(aesKey)
private, public, err := utils.GenerateKeyPair2048bits()
if err != nil {
	log.Fatal(err)
}
encrypted, err := utils.EncryptWithPublicKey(aesKey, public)
if err != nil {
	log.Fatal(err)
}
fmt.Println(len(encrypted))
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
//Testing Encrypting File
err = utils.EncryptFile("D:\\DownloadedYoutube\\Brother VC 500W Colour Label Printer.mkv", "D:\\DownloadedYoutube\\Brother VC 500W Colour Label Printer.mkv.encrypted", decrypted)
if err != nil {
	log.Fatal(err)
}

err = utils.DecryptFile("D:\\DownloadedYoutube\\Brother VC 500W Colour Label Printer.mkv.encrypted", "D:\\DownloadedYoutube\\Brother VC 500W Colour Label Printer.decrypted.mkv", decrypted)
if err != nil {
	log.Fatal(err)
}

*/
