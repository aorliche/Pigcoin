package main

import (
    "database/sql"
    "errors"
    "flag"
    "fmt"
    "os"
    "time"
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
    
var cfg = mysql.Config{
    User: os.Getenv("DBUSER"),      //"anton",
    Passwd: os.Getenv("DBPASS"),    //"AtlanticCityPass",
    Net: "tcp",
    Addr: "localhost:3306",
    DBName: "pigcoin",
    AllowNativePasswords: true,
}

type Wallet struct {
    Id int
    Name string
    User string
    Email string
    PrivateKey []byte
    PublicKey []byte
}

type Transaction struct {
    Id int
    From int
    To int
    Amount int64
    Sig []byte
    When int64
}

func (w *Wallet) String() string {
    return fmt.Sprintf("Id: %d, Name: %s, User: %s, Email: %s", w.Id, w.Name, w.User, w.Email)
}

func (t *Transaction) String() string {
    tm := time.Unix(t.When, 0)
    tmStr := tm.Format(time.RFC3339)
    return fmt.Sprintf("Id: %d, From: %d, To: %d, Amount: %d, When: %s", t.Id, t.From, t.To, t.Amount, tmStr)
}

func AddMoney(w *Wallet, amount int64) (int64, error) {
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    defer db.Close()
    res, err := db.Exec("INSERT INTO transactions (`from`, `to`, amount) VALUES (?, ?, ?)", 0, w.Id, amount)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }
    return id, nil
}

func FindWalletsByField(name string, value string, id int) ([]*Wallet, error) {
    // Sanitize input
    if name != "id" && name != "name" && name != "user" && name != "email" {
        return nil, errors.New("unknown field")
    }
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return nil, err
    }
    defer db.Close()
    // Note sanitization above
    var rows *sql.Rows
    if name == "id" {
        rows, err = db.Query(fmt.Sprintf("SELECT * FROM wallets WHERE %s = ?", name), id)
    } else {
        rows, err = db.Query(fmt.Sprintf("SELECT * FROM wallets WHERE %s = ?", name), value)
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    wallets := make([]*Wallet, 0)
    for rows.Next() {
        wallet := new(Wallet)
        err := rows.Scan(&wallet.Id, &wallet.Name, &wallet.User, &wallet.Email, &wallet.PrivateKey, &wallet.PublicKey)
        if err != nil {
            return nil, err
        }
        wallets = append(wallets, wallet)
    }
    return wallets, nil
}

func DeleteWallet(id int) (int64, error) {
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    defer db.Close()
    res, err := db.Exec("DELETE FROM wallets WHERE id = ?", id)
    if err != nil {
        return 0, err
    }
    nrows, err := res.RowsAffected()
    if err != nil {
        return 0, err
    }
    return nrows, nil
}

func GetBalance(id int) (int64, error) {
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    defer db.Close()
    var in, out sql.NullInt64
    row := db.QueryRow("SELECT sum(amount) FROM transactions WHERE `to` = ?", id)
    err = row.Scan(&in)
    if err != nil {
        return 0, err
    }
    row = db.QueryRow("SELECT sum(amount) FROM transactions WHERE `from` = ?", id)
    err = row.Scan(&out)
    if err != nil {
        return 0, err
    }
    in2 := int64(0)
    out2 := int64(0)
    if in.Valid {
        in2 = in.Int64
    }
    if out.Valid {
        out2 = out.Int64
    }
    return in2 - out2, nil
}

func ListWallets(num int) ([]*Wallet, error) {
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return nil, err
    }
    defer db.Close()
    rows, err := db.Query("SELECT * FROM wallets LIMIT ?", num)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    wallets := make([]*Wallet, 0)
    for rows.Next() {
        wallet := new(Wallet)
        err := rows.Scan(&wallet.Id, &wallet.Name, &wallet.User, &wallet.Email, &wallet.PrivateKey, &wallet.PublicKey)
        if err != nil {
            return nil, err
        }
        wallets = append(wallets, wallet)
    }
    return wallets, nil
}

func ListTransactions(num int) ([]*Transaction, error) {
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return nil, err
    }
    defer db.Close()
    rows, err := db.Query("SELECT * FROM transactions LIMIT ?", num)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    transactions := make([]*Transaction, 0)
    for rows.Next() {
        transaction := new(Transaction)
        var tBytes []byte
        err := rows.Scan(&transaction.Id, &transaction.From, &transaction.To, &transaction.Amount, &transaction.Sig, &tBytes)
        if err != nil {
            return nil, err
        }
        tm, err := time.Parse("2006-01-02 15:04:05", string(tBytes))
        if err != nil {
            return nil, err
        }
        transaction.When = tm.Unix()
        transactions = append(transactions, transaction)
    }
    return transactions, nil
}

func CreateWallet(name string, user string, email string) (int64, error) {
    wallet, err := pigcoin.GenerateWallet(name, user, email)
    if err != nil {
        return 0, err
    }
    privKey, err := wallet.PrivateKeyPEM()
    if err != nil {
        return 0, err
    }
    pubKey, err := wallet.PublicKeyPEM()
    if err != nil {
        return 0, err
    }
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    defer db.Close()
    res, err := db.Exec("INSERT INTO wallets (name, user, email, private, public) VALUES (?, ?, ?, ?, ?)", 
        wallet.Name, wallet.User, wallet.Email, privKey, pubKey)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }
    return id, nil
    /*fmt.Println(string(privKey))
    fmt.Println(string(pubKey))
    fmt.Println(wallet.Name, wallet.User, wallet.Email)*/
}

func SubtractMoney(w *Wallet, amount int64) (int64, error) {
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    defer db.Close()
    res, err := db.Exec("INSERT INTO transactions (`to`, `from`, amount) VALUES (?, ?, ?)", 0, w.Id, amount)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }
    return id, nil
}

func main() {
    action := flag.String("action", "", 
        "REQUIRED: add-money, create-wallet, delete-wallet, find-wallet, get-balance, list-transactions, list-wallets, send-money, subtract-money")
    name := flag.String("name", "", "wallet name")
    user := flag.String("user", "", "user name")
    email := flag.String("email", "", "user email")
    num := flag.Int("num", 20, "number (wallets or transactions)")
    id := flag.Int("id", 0, "wallet id")
    amount := flag.Int64("amount", 0, "amount")
    /*from := flag.Int64("from", 0, "from wallet id")
    to := flag.Int64("to", 0, "to wallet id")*/
    flag.Parse()
    switch *action {
        case "add-money": {
            if *id == 0 {
                fmt.Println("id is required")
                return
            }
            if *amount == 0 {
                fmt.Println("amount is required")
                return
            }
            wallets, err := FindWalletsByField("id", "", *id)
            if len(wallets) == 0 {
                fmt.Println("wallet not found")
                return
            }
            tid, err := AddMoney(wallets[0], *amount)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Printf("%d: Added %d to wallet %d\n", tid, *amount, *id)
        }
        case "create-wallet": {
            names := [3]string{"name", "user", "email"}
            for i,s := range [3]string{*name, *user, *email} {
                err := Validate(s)
                if err != nil {
                    fmt.Println(names[i], ":", err)
                    return
                }
            }
            wid, err := CreateWallet(*name, *user, *email)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(wid)
        }
        case "delete-wallet": {
            nrows, err := DeleteWallet(*id)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(nrows)
        }
        case "find-wallet": {
            var wallets []*Wallet
            var err error
            if *id != 0 {
                wallets, err = FindWalletsByField("id", "", *id)
            } else if *name != "" {
                wallets, err = FindWalletsByField("name", *name, 0)
            } else if *user != "" {
                wallets, err = FindWalletsByField("user", *user, 0)
            } else {
                fmt.Println("id, name, or user required")
                return
            }
            if err != nil {
                fmt.Println(err)
                return
            }
            for _, w := range wallets {
                fmt.Println(w.String())
            }
        }
        case "get-balance": {
            if *id == 0 {
                fmt.Println("id is required")
                return
            }
            balance, err := GetBalance(*id)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(balance)
        }
        case "list-wallets": {
            wallets, err := ListWallets(*num)
            if err != nil {
                fmt.Println(err)
                return
            }
            for _, w := range wallets {
                fmt.Println(w.String())
            }
        }
        case "list-transactions": {
            transactions, err := ListTransactions(*num)
            if err != nil {
                fmt.Println(err)
                return
            }
            for _, t := range transactions {
                fmt.Println(t.String())
            }
        }
        case "subtract-money": {
            if *id == 0 {
                fmt.Println("id is required")
                return
            }
            if *amount == 0 {
                fmt.Println("amount is required")
                return
            }
            wallets, err := FindWalletsByField("id", "", *id)
            if len(wallets) == 0 {
                fmt.Println("wallet not found")
                return
            }
            tid, err := SubtractMoney(wallets[0], *amount)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Printf("%d: Subtracted %d from wallet %d\n", tid, *amount, *id)
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
