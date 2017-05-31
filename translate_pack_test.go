package translate

import (  
    // "log"  
    "strings"
    "fmt"
    // "github.com/tealeg/xlsx"  
    "math/rand"
    "time"
    "crypto/md5"
    "encoding/hex"
    "net/http"
    "io/ioutil"
    "net/url"
    "github.com/bitly/go-simplejson"
     "database/sql"  
    _ "github.com/go-sql-driver/mysql"  
    "github.com/robertkrimen/otto"
    "testing"
)  

func Test_Get_trans_content(t *testing.T){
    word_list := []string{}
    word_list = append(word_list, translate.Translate("test",2))
    word_list = append(word_list, translate.Translate("test",2))
    word_list = append(word_list, "people")
    word_list = append(word_list, translate.Translate("add",2))
    word_list = append(word_list, "name")
    word_list = append(word_list, "mom")
    word_list = append(word_list, "father")
    word_list = append(word_list, "return")
    db_conf := make(map[string]string)
    db_conf["Host"] = "127.0.0.1"
    db_conf["User"] = "root"
    db_conf["Password"] = ""
    db_conf["Dbname"] = "hadtest"
    db_conf["Port"] = "3306"
    ret,err := Get_trans_content(2,word_list,db_conf)
    if err != nil{
        t.Error(err)
    }
}
func Test_Translate(t *testing.T){
    words := "test"
    if Translate(words,2) == ""{
        t.Error("translate error")
    }
}