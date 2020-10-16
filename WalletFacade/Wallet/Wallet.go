package Wallet

import (
    "fmt"
    "log"
)

type Wallet struct {
    Balance int
}

func WalletCons() *Wallet {
    return &Wallet{Balance: 0}
}

func (w *Wallet) Check() string {
    return fmt.Sprintf("残高は、%v 円です。\n", w.Balance)
}

func (w *Wallet) BalanceAdd(amount int) {
    w.Balance += amount
    s := fmt.Sprintf("Balanceに%v円を追加。Balance: %v円", amount, w.Balance)
    log.Println(s)
}

func (w *Wallet) BalanceDebit(amount int) {
    w.Balance = w.Balance - amount
}