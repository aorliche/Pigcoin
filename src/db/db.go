package db

import (
    "database/sql"
    "errors"
    "fmt"
    "os"
    "time"
    "github.com/go-sql-driver/mysql"
    "github.com/aorliche/pigcoin"
)

func Validate(field string, s string) error {
    if s == "" {
        return errors.New(field + ": empty string")
    } else if len(s) > 20 {
        return errors.New(field + ": string too long")
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

func CreateWallet(name string, user string, email string, password string) (int64, error) {
    names := [4]string{"name", "user", "email", "password"}
    for i,s := range [4]string{name, user, email, password} {
        err := Validate(names[i], s)
        if err != nil {
            return 0, err
        }
    }
    // Find existing wallet
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    row := db.QueryRow("SELECT name FROM wallets WHERE name = ?", name)
    err = row.Scan(&name) 
    if err == nil {
        return 0, errors.New("wallet already exists")
    }
    defer db.Close()
    wallet, err := pigcoin.GenerateWallet(name, user, email)
    if err != nil {
        return 0, err
    }
    privKey, err := wallet.PrivateKeyPEM()
    if err != nil {
        return 0, err
    }
    privKeyEnc, err := pigcoin.Encrypt(privKey, []byte(password))
    if err != nil {
        return 0, err
    }
    pubKey, err := wallet.PublicKeyPEM()
    if err != nil {
        return 0, err
    }
    res, err := db.Exec("INSERT INTO wallets (name, user, email, private, public) VALUES (?, ?, ?, ?, ?)", 
        wallet.Name, wallet.User, wallet.Email, privKeyEnc, pubKey)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }
    return id, nil
}

func Transfer(from int, to int, amount int64, password string) (int64, error) {
    wallets, err := FindWalletsByField("id", "", from)
    if len(wallets) == 0 {
        return 0, errors.New("from wallet not found")
    }
    towallets, err := FindWalletsByField("id", "", to)
    if len(towallets) == 0 {
        return 0, errors.New("to wallet not found")
    }
    balance, err := GetBalance(from)
    if err != nil {
        return 0, err
    }
    if balance < amount {
        return 0, errors.New("insufficient funds")
    }
    dec, err := pigcoin.Decrypt(wallets[0].PrivateKey, []byte(password))
    if err != nil {
        return 0, err
    }
    if dec[0] != '-' && dec[1] != '-' {
        return 0, errors.New("invalid password")
    }
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return 0, err
    }
    defer db.Close()
    res, err := db.Exec("INSERT INTO transactions (`from`, `to`, amount) VALUES (?, ?, ?)", from, to, amount)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }
    return id, nil
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
