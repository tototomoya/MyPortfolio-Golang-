package main

import (
    "fmt"
    "log"
    "os"
    "encoding/json"
    "strconv"
    "net/http"
    . "main/WalletFacade"
    "main/WalletFacade/Account"
    "main/WalletFacade/Items"
    "main/WalletFacade/Wallet"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/stripe/stripe-go/card"
)

const dir = "./user/"

var publishableKey string

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    godotenv.Load()
	publishableKey = os.Getenv("PUBLISHABLE_KEY")
	stripe.Key = os.Getenv("SECRET_KEY")
}

func get(name string, password string) *WalletFacade {
    var fp *os.File
    filename := dir + name
    fp, err := os.Open(filename)
    if err != nil {
        log.Println(err)
    }
    dec := json.NewDecoder(fp)
    var user *WalletFacade
    dec.Decode(&user)
    return user
}

func exsistFile(filename string) bool {
    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return true
}

func newWalletFacade(name string, password string) *WalletFacade {
    wFacade := &WalletFacade{
        Account:  Account.AccountCons(name, password),
        Wallet:  Wallet.WalletCons(),
        Items:   Items.ItemsCons(),
    }
    cus_params := &stripe.CustomerParams{
        Email:  stripe.String("ktaoao@yahoo.co.jp"),
        Name:   &name,
    }
    cus, _ := customer.New(cus_params)
    wFacade.ID = cus.ID
    card_params := &stripe.CardParams{
        Customer: stripe.String(wFacade.ID),
        Token: stripe.String("tok_visa"),
    }
    card.New(card_params)
    fmt.Printf("%s 様のアカウントを作成しました。\n", wFacade.Account.Name)
    return wFacade
}

func userRegister(c *gin.Context) {
    name := c.Param("name")
    password := c.Param("password")
    user := newWalletFacade(name, password)
    user.Save()
    c.SetCookie("name", user.Account.Name, 60*60*24, "/", "facadestripe.hitabacokyou.repl.co", false, false)
    c.SetCookie("password", user.Account.Password, 60*60*24, "/", "facadestripe.hitabacokyou.repl.co", false, false)
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "User": user,
        "pub": publishableKey,
        "itemList": Items.GetItemList(),
    })
}

func userLogin(c *gin.Context) {
    name := c.Param("name")
    password := c.Param("password")
    filename := dir + name
    if exsistFile(filename) {
        user := get(name, password)
        err := user.Account.CheckAccount(name, password)
        if err != nil {
           c.String(http.StatusCreated, "不正なアクセスです。")
           return
        }
        c.SetCookie("name", user.Account.Name, 60*60*24, "/", "stripe.hitabacokyou.repl.co", false, false)
        c.SetCookie("password", user.Account.Password, 60*60*24, "/", "stripe.hitabacokyou.repl.co", false, false)
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "User": user,
            "pub": publishableKey,
            "itemList": Items.GetItemList(),
        })
        return
    } else {
        c.String(http.StatusCreated, "存在しないユーザです。")
        return 
    }
}

func userLogout(c *gin.Context) {
    c.SetCookie("name", "", 1, "/", "facadestripe.hitabacokyou.repl.co", false, false)
    c.SetCookie("password", "", 1, "/", "facadestripe.hitabacokyou.repl.co", false, false)
    c.String(http.StatusCreated, "ログアウトしました。")
}

func chargeMoney(c *gin.Context) {
    name, err := c.Cookie("name")
    password, err := c.Cookie("password")
    if err != nil {
        c.String(http.StatusCreated, "ログインしてください。")
        return     
    } else {
        user := get(name, password)
        err := user.Account.CheckAccount(name, password)
        if err != nil {
           c.String(http.StatusCreated, "不正なアクセスです。")
           return
        }
        c.Request.ParseMultipartForm(50)
        for _, v := range c.Request.Form["food"] {
            v_int, _ := strconv.Atoi(v)
            user.Items.Add(v_int)
        }
        sum := user.Items.Sum()
        err = user.Buy(user.Account.Name, user.Account.Password, sum)
        if err != nil {
            c.String(http.StatusCreated, err.Error())
            return
        }
        user.Save()
        s := fmt.Sprintf("金額は%v円になります。", sum)
        c.String(http.StatusCreated, s + "\nありがとうございました。")
    }
    return
}

func deposit(c *gin.Context) {
    name := c.Param("name")
    password := c.Param("password")
    var user *WalletFacade
    filename := dir + name
    if exsistFile(filename) {
        user = get(name, password)
    } else {
        c.String(http.StatusCreated, "このユーザは利用できません。")
        return
    }
    err := user.Account.CheckAccount(name, password)
    if err != nil {
        c.String(http.StatusCreated, "不正なアクセスです。")
        return
    }
    amount, _ := strconv.Atoi(c.Param("amount"))
    err = user.AddMoneyToWallet(name, password, amount)
    if err != nil {
        c.String(http.StatusCreated, err.Error())
        return
    }
    s := fmt.Sprintf("%v", amount)
    c.String(http.StatusCreated, s + "円チャージしました。")
}

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("./*.tmpl")
    r.GET("/register/:name/:password", userRegister)
    r.POST("/charge", chargeMoney)
    r.GET("/login/:name/:password", userLogin)
    r.GET("/logout", userLogout)
    r.GET("/deposit/:name/:password/:amount", deposit)
    r.StaticFS("/static", http.Dir("./static"))
    r.Run(":8000")
}