package WalletFacade

import (
    "fmt"
    "log"
    "io/ioutil"
    "encoding/json"
    . "main/WalletFacade/Account"
    . "main/WalletFacade/Items"
    . "main/WalletFacade/Wallet"
	"github.com/stripe/stripe-go"
    "github.com/stripe/stripe-go/charge"
    "github.com/stripe/stripe-go/customerbalancetransaction"
)

type WalletFacade struct {
    Account      *Account
    Wallet       *Wallet
    Items *Items
    ID string
}

const dir = "./user/"

func (w *WalletFacade) Save() {
    fileName := dir + w.Account.Name
    data, _ := json.MarshalIndent(w, "", " ")
    ioutil.WriteFile(fileName, []byte(data), 0664)
}

func (w *WalletFacade) AddMoneyToWallet(name string, password string, amount int) error {
    defer w.Save()
    log.Println()
    err := w.Account.CheckAccount(name, password)
    if err != nil {
        return fmt.Errorf("このアカウントを使用することはできません。")
    }
    err = w.stripe_deposit(amount)
    if err != nil {
        return err
    }
    w.Wallet.BalanceAdd(amount)
    return nil
}

func (w *WalletFacade) stripe_deposit(amount int) error {
    log.Println()
    params := &stripe.CustomerBalanceTransactionParams{
        Amount: stripe.Int64(int64(amount)),
        Currency: stripe.String(string(stripe.CurrencyJPY)),
        Customer: stripe.String(w.ID),
    }
    _, err := customerbalancetransaction.New(params)
    if err != nil {
        return fmt.Errorf("入金が出来ません。")
    }
    log.Println("stripeの入金処理終了。")
    return nil
}

func (w *WalletFacade) Buy(name string, password string, amount int) error {
    defer w.Save()
    err := w.Account.CheckAccount(name, password)
    if err != nil {
        return err
    }
    err = w.debit(amount)
    if err != nil {
        return err
    }

    w.Items.SumValue += w.Items.Sum()
    w.Items.Done()
    return nil
}

func (w *WalletFacade) debit(amount int) error {
    if w.Wallet.Balance < amount {
        return fmt.Errorf("残高が不足しています。")
    }
    w.Wallet.BalanceDebit(amount)
    _, err := charge.New(&stripe.ChargeParams{
        Amount:       stripe.Int64(int64(amount)),
        Currency:     stripe.String(string(stripe.CurrencyJPY)),
        Customer:   stripe.String(w.ID),
    })
    if err != nil {
        return fmt.Errorf("決済ができません。")
    }
    return nil
}