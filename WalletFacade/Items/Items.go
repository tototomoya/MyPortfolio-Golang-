package Items

import (
    . "main/WalletFacade/Items/Item"
    "os"
    "strconv"
    "strings"
)

type Items struct {
    List []*Item
    SumValue int
}

var ItemList map[int]map[string]int

func init() {
    itemName := strings.Split(os.Getenv("name"), ",")
    itemValue := strings.Split(os.Getenv("value"), ",")
    itemNum, _ := strconv.Atoi(os.Getenv("num"))
    
    ItemList = make(map[int]map[string]int, itemNum)
    for i := 0; i < itemNum; i++ {
        val, _ := strconv.Atoi(itemValue[i])
        ItemList[i] = map[string]int{itemName[i]: val}
    }
}

func ItemsCons() *Items {
    return &Items{}
}

func GetItemList() map[int]map[string]int {
    return ItemList
}

func (i *Items) Add(index int) {
    for key, val := range ItemList[index] {
        i.List = append(i.List, &Item{Name: key, Price: val, Del: false})
    }
}

func (i *Items) Done() {
    N := len(i.List)
    for n := 0; n < N; n++ {
        if !i.List[n].Del {
            i.List[n].Del = true
        }
    }
}

func (i *Items) Sum() int {
    N := len(i.List)
    var sum int
    for n := 0; n < N; n++ {
        if i.List[n].Del {
            continue
        }
        sum += i.List[n].Price
    }
    return sum
}
