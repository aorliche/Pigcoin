package main

import (
    "flag"
    "fmt"
    pdb "github.com/aorliche/pigcoin/db"
)

func main() {
    action := flag.String("action", "", 
        "REQUIRED: add-money, create-wallet, delete-wallet, find-wallet, " + 
        "get-balance, list-transactions, list-wallets, transfer, subtract-money")
    name := flag.String("name", "", "wallet name")
    user := flag.String("user", "", "user name")
    email := flag.String("email", "", "user email")
    password := flag.String("password", "", "user password")
    num := flag.Int("num", 20, "number (wallets or transactions)")
    id := flag.Int("id", 0, "wallet id")
    amount := flag.Int64("amount", 0, "amount")
    from := flag.Int("from", 0, "from wallet id")
    to := flag.Int("to", 0, "to wallet id")
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
            wallets, err := pdb.FindWalletsByField("id", "", *id)
            if len(wallets) == 0 {
                fmt.Println("wallet not found")
                return
            }
            tid, err := pdb.AddMoney(wallets[0], *amount)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Printf("%d: Added %d to wallet %d\n", tid, *amount, *id)
        }
        case "create-wallet": {
            wid, err := pdb.CreateWallet(*name, *user, *email, *password)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(wid)
        }
        case "delete-wallet": {
            nrows, err := pdb.DeleteWallet(*id)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(nrows)
        }
        case "find-wallet": {
            var wallets []*pdb.Wallet
            var err error
            if *id != 0 {
                wallets, err = pdb.FindWalletsByField("id", "", *id)
            } else if *name != "" {
                wallets, err = pdb.FindWalletsByField("name", *name, 0)
            } else if *user != "" {
                wallets, err = pdb.FindWalletsByField("user", *user, 0)
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
            balance, err := pdb.GetBalance(*id)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(balance)
        }
        case "list-wallets": {
            wallets, err := pdb.ListWallets(*num)
            if err != nil {
                fmt.Println(err)
                return
            }
            for _, w := range wallets {
                fmt.Println(w.String())
            }
        }
        case "list-transactions": {
            transactions, err := pdb.ListTransactions(*num)
            if err != nil {
                fmt.Println(err)
                return
            }
            for _, t := range transactions {
                fmt.Println(t.String())
            }
        }
        case "transfer": {
            if *from == 0 {
                fmt.Println("from is required")
                return
            }
            if *to == 0 {
                fmt.Println("to is required")
                return
            }
            if *amount == 0 {
                fmt.Println("amount is required")
                return
            }
            if *password == "" {
                fmt.Println("password is required")
                return
            }
            tid, err := pdb.Transfer(*from , *to, *amount, *password)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Printf("%d: Sent %d from wallet %d to wallet %d\n", tid, *amount, *from, *to)
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
            wallets, err := pdb.FindWalletsByField("id", "", *id)
            if len(wallets) == 0 {
                fmt.Println("wallet not found")
                return
            }
            tid, err := pdb.SubtractMoney(wallets[0], *amount)
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
}
