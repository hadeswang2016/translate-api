package translate

import (
    "io/ioutil"
    "github.com/bitly/go-simplejson"
    "log"
    "github.com/robertkrimen/otto"
    "strings"
    "net/http"
    "math/rand"
    "crypto/md5"
    "encoding/hex"
    "net/url"
    "time"
    "strconv"
)
func Test_Translate_youdao(t *testing.T){
    words := "test"
    res,err := Translate_youdao(words)
    if err != nil{
        t.Error(err)
    }
}
func Test_Translate_google(t *testing.T){
    words := "test"
    res,err := Translate_google(words)
    if err != nil{
        t.Error(err)
    }
}