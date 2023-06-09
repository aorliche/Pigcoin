package main

import (
    "encoding/json"
    //"fmt"
    "os"
    "net/http"
    //"strconv"
    //"time"
    pdb "github.com/aorliche/pigcoin/db"
)

type Wallet struct {
    Id int
    Name string
    User string
    Balance int64
}

var uiFname = os.Getenv("BASEDIR") + "/ui.html"
var jsFname = os.Getenv("BASEDIR") + "/ui.js"

func HTML(w http.ResponseWriter, req *http.Request) {
    http.ServeFile(w, req, uiFname)
}

func JS(w http.ResponseWriter, req *http.Request) {
    http.ServeFile(w, req, jsFname)
}

func Wallets(w http.ResponseWriter, req *http.Request) {
    wallets, err := pdb.ListWallets(20)
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
    _, err := os.ReadFile(uiFname)
    if err != nil {
        panic(err)
    }
    http.HandleFunc("/ui", Headers(HTML))
    http.HandleFunc("/ui.js", Headers(JS))
    http.HandleFunc("/api/wallets", Headers(Wallets))
    http.ListenAndServe("0.0.0.0:8888", nil)
}
