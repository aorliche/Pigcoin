package main

import (
    "database/sql"
    "errors"
    "flag"
    "fmt"
    pigcoin "github.com/aorliche/pigcoin"
    "github.com/go-sql-driver/mysql"
)

func Validate(s string) error {
    if s == "" {
        return errors.New("empty string")
    } else if len(s) > 20 {
        return errors.New("string too long")
    }
    return nil
}

func CreateWallet(name string, user string, email string) {
    wallet, err := pigcoin.GenerateWallet(name, user, email)
    if err != nil {
        fmt.Println(err)
        return
    }
    privKey, err := wallet.PrivateKeyPEM()
    if err != nil {
        fmt.Println(err)
        return
    }
    pubKey, err := wallet.PublicKeyPEM()
    if err != nil {
        fmt.Println(err)
        return
    }
    cfg := mysql.Config{
        User: "anton",
        Passwd: "AtlanticCityPass",
        Net: "tcp",
        Addr: "localhost:3306",
        DBName: "pigcoin",
        AllowNativePasswords: true,
    }
    db, err := sql.Open("mysql", cfg.FormatDSN())
    defer db.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
    res, err := db.Exec("INSERT INTO wallets (name, user, email, private, public) VALUES (?, ?, ?, ?, ?)", 
        wallet.Name, wallet.User, wallet.Email, string(privKey), string(pubKey))
    if err != nil {
        fmt.Println(err)
        return
    }
    id, err := res.LastInsertId()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(id)
    /*fmt.Println(string(privKey))
    fmt.Println(string(pubKey))
    fmt.Println(wallet.Name, wallet.User, wallet.Email)*/
}

func main() {
    action := flag.String("action", "", "create,list,delete")
    name := flag.String("name", "", "wallet name")
    user := flag.String("user", "", "user name")
    email := flag.String("email", "", "user email")
    flag.Parse()
    switch *action {
        case "create": {
            names := [3]string{"name", "user", "email"}
            for i,s := range [3]string{*name, *user, *email} {
                err := Validate(s)
                if err != nil {
                    fmt.Println(names[i], ":", err)
                    return
                }
            }
            CreateWallet(*name, *user, *email)
        }
        default: {
            flag.Usage()
        }
    }
    /*enc, err := pigcoin.Encrypt([]byte("hello"), []byte("pass"))
    if err != nil {
        panic(err)
    }
    dec, err := pigcoin.Decrypt(enc, []byte("pass"))
    if err != nil {
        panic(err)
    }
    fmt.Println(string(dec))*/
}
