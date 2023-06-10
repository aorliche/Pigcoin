package main

import (
    "encoding/json"
    "fmt"
    "os"
    "net/http"
    //"net/http/httputil"
    "strconv"
    //"time"
    pdb "github.com/aorliche/pigcoin/db"
)

type Wallet struct {
    Id int
    Name string
    User string
    Balance int64
}

type Transaction struct {
    Id int
    From int
    To int
    Amount int64
    When int64
}

/*var uiFname = os.Getenv("BASEDIR") + "/ui.html"
var jsFname = os.Getenv("BASEDIR") + "/ui.js"*/

func JsonError(s string) string {
    jsn, _ := json.Marshal(map[string]string{"Error": s})
    return string(jsn);
}

func Static(w http.ResponseWriter, req *http.Request, file string) {
    http.ServeFile(w, req, os.Getenv("BASEDIR") + "/" + file)
}

func Wallets(w http.ResponseWriter, req *http.Request) {
    wallets, err := pdb.ListWallets(200)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    wallets2 := make([]Wallet, len(wallets))
    for i, wallet := range wallets {
        balance, err := pdb.GetBalance(wallet.Id)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        wallets2[i] = Wallet{wallet.Id, wallet.Name, wallet.User, balance}
    }
    jsn, _ := json.Marshal(wallets2)
    w.Write(jsn)
}

func Transactions(w http.ResponseWriter, req *http.Request) {
    transactions, err := pdb.ListTransactions(200)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    transactions2 := make([]Transaction, len(transactions))
    for i, transaction := range transactions {
        transactions2[i] = Transaction{transaction.Id, transaction.From, transaction.To, transaction.Amount, transaction.When} 
    }
    jsn, _ := json.Marshal(transactions2)
    w.Write(jsn)
}

func CreateWallet(w http.ResponseWriter, req *http.Request) {
    err := req.ParseMultipartForm(0)
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    if req.FormValue("create-password") != req.FormValue("create-password2") {
        http.Error(w, JsonError("Passwords do not match"), 500)
        return
    }
    wid, err := pdb.CreateWallet(
        req.FormValue("name"), 
        req.FormValue("user"), 
        req.FormValue("email"), 
        req.FormValue("create-password"))
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    w.Write([]byte(strconv.Itoa(int(wid))))
}

func Transfer(w http.ResponseWriter, req *http.Request) {
    err := req.ParseMultipartForm(0)
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    from, err := strconv.Atoi(req.FormValue("from"))
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    to, err := strconv.Atoi(req.FormValue("to"))
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    amount, err := strconv.ParseInt(req.FormValue("amount"), 10, 64)
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    tid, err := pdb.Transfer(from, to, amount, req.FormValue("password"))
    if err != nil {
        http.Error(w, JsonError(err.Error()), 500)
        return
    }
    w.Write([]byte(strconv.Itoa(int(tid))))
}

type HFunc func (http.ResponseWriter, *http.Request)

func Headers(fn HFunc) HFunc {
    return func (w http.ResponseWriter, req *http.Request) {
        //fmt.Println(req.Method)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        fn(w, req)
    }
}

func main() {
    dir, err := os.Open(os.Getenv("BASEDIR"))
    if err != nil {
        panic(err)
    }
    files, err := dir.Readdir(0)
    if err != nil {
        panic(err)
    }
    for _, v := range files {
        fmt.Println(v.Name(), v.IsDir())
        if v.IsDir() {
            continue
        }
        file := "/" + v.Name()
        http.HandleFunc(file, Headers(func (w http.ResponseWriter, req *http.Request) {Static(w, req, file)}))
    }
    http.HandleFunc("/api/wallets", Headers(Wallets))
    http.HandleFunc("/api/transactions", Headers(Transactions))
    http.HandleFunc("/api/create-wallet", Headers(CreateWallet))
    http.HandleFunc("/api/transfer", Headers(Transfer))
    http.ListenAndServe("0.0.0.0:8888", nil)
}
