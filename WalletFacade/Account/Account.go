package Account

import "fmt"

type Account struct {
    Name string
    Password string
}

func AccountCons(name, password string) *Account {
    return &Account{Name: name, Password: password}
}

func (a *Account) CheckAccount(name string, password string) error {
    if a.Name == name && a.Password == password {
        return nil
    } else {
        return fmt.Errorf("このアカウントを利用することはできません。")
    }
}