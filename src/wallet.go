package pigcoin

import (
    "crypto"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "errors"
)

type Wallet struct {
    PrivateKey *rsa.PrivateKey
    PublicKey  *rsa.PublicKey
    WalletName string
    UserName   string
    Email      string
}

func GenerateWallet(wname string, uname string, email string) (*Wallet, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, err
    }
    publicKey := privateKey.Public().(*rsa.PublicKey)
    return &Wallet{privateKey, publicKey, wname, uname, email}, nil
}

func (w *Wallet) ToPEM() ([]byte, error) {
    if w == nil {
        return nil, errors.New("wallet is nil")
    }
    var typ string = "RSA PRIVATE KEY"
    var bytes []byte = nil
    var bytesPEM []byte = nil
    var err error
    if w.PrivateKey != nil {
        bytes, err = x509.MarshalPKCS8PrivateKey(w.PrivateKey)
        if err != nil {
            return nil, err
        }
    } else if w.PublicKey != nil {
        bytes, err = x509.MarshalPKIXPublicKey(w.PublicKey)
        typ = "RSA PUBLIC KEY"
        if err != nil {
            return nil, err
        }
    }
    if bytes != nil {
        bytesPEM = pem.EncodeToMemory(&pem.Block{
            Type:  typ,
            Bytes: bytes,
        })
    }
    return bytesPEM, nil
}

var signopts = rsa.PSSOptions{
    SaltLength:  20,
    Hash:        crypto.SHA256,
}

func (w *Wallet) Sign(data []byte) ([]byte, error) {
    hashed := sha256.Sum256(data)
    sig, err := rsa.SignPSS(rand.Reader, w.PrivateKey, crypto.SHA256, hashed[:], &signopts)
    if err != nil {
        return nil, err
    }
    return sig, nil
}

func (w *Wallet) Verify(data []byte, sig []byte) error {
    hashed := sha256.Sum256(data)
    return rsa.VerifyPSS(w.PublicKey, crypto.SHA256, hashed[:], sig, &signopts)
}

// https://blog.logrocket.com/learn-golang-encryption-decryption/
var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(s)
    if err != nil {
        return nil, err
    }
    return data, nil
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text []byte, pwd []byte) ([]byte, error) {
    hashed := sha256.Sum256(pwd)
    block, err := aes.NewCipher(hashed[:16])
    if err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cipherText := make([]byte, len(text))
    cfb.XORKeyStream(cipherText, text)
    return cipherText, nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(cipherText []byte, pwd []byte) ([]byte, error) {
    hashed := sha256.Sum256(pwd)
    block, err := aes.NewCipher(hashed[:16])
    if err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBDecrypter(block, iv)
    plainText := make([]byte, len(cipherText))
    cfb.XORKeyStream(plainText, cipherText)
    return plainText, nil
}
