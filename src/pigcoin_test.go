package pigcoin 

import (
    "fmt"
    "testing"
)

/*func TestGenerateWallet(t *testing.T) {
    wallet, err := GenerateWallet("name", "user", "email")
    if err != nil {
        t.Fatal(err)
    }
    if wallet.WalletName != "name" {
        t.Errorf("Expected name, got %s", wallet.WalletName)
    }
}

func TestToPEM(t *testing.T) {
    wallet, err := GenerateWallet("name", "user", "email")
    if err != nil {
        t.Fatal(err)
    }
    pem, err := wallet.ToPEM()
    if err != nil {
        t.Fatal(err)
    }
    pemStr := string(pem)
    fmt.Println(pemStr)
}

func TestToPEMPublic(t *testing.T) {
    wallet, err := GenerateWallet("name", "user", "email")
    if err != nil {
        t.Fatal(err)
    }
    wallet.PrivateKey = nil
    pem, err := wallet.ToPEM()
    if err != nil {
        t.Fatal(err)
    }
    pemStr := string(pem)
    fmt.Println(pemStr)
}

func TestSign(t *testing.T) {
    wallet, err := GenerateWallet("name", "user", "email")
    if err != nil {
        t.Fatal(err)
    }
    data := []byte("hello world")
    sig, err := wallet.Sign(data)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(sig)
    err = wallet.Verify(data, sig)
    if err != nil {
        t.Fatal(err)
    }
}*/

func TestEncrypt(t *testing.T) {
    out, err := Encrypt([]byte("hello world"), []byte("mypass"))
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(Encode(out))
    cipher, err := Decode(Encode(out))
    orig, err := Decrypt(cipher, []byte("mypass"))
    if err != nil {
        t.Fatal(err)
    }
    if string(orig) != "hello world" {
        t.Errorf("Expected hello world, got %s", string(orig))
    }
}
