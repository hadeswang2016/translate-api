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
//English is translated into Chinese with youdao api
func Translate_youdao(words string)(string,error) {
    res,err := youdao_tanslate(words)
    return res,err 
}
//English is translated into Chinese with google api
func Translate_google(words string)(string,error) {
    res,err := google_translate(words)
    return res,err 
}
//youdao api handle
func youdao_tanslate(words string)(string,error) {
    youdao_api_conf,err :=  read_conf("apiconf.json","youdao")
    if err !=nil {
        log.Panicf("get api conf fail : %s", err.Error())
    } 
    //words string,max_limit int,trunk int,split_str string,api_name string,api_conf_data map[string]interface{}
    translate_max_limit,_:=strconv.Atoi(youdao_api_conf["translate_max_limit"].(string))
    translate_trunk_v,_:=strconv.Atoi(youdao_api_conf["translate_trunk"].(string))
    res,err := translate_trunk(words,translate_max_limit,translate_trunk_v,".","youdao",youdao_api_conf)
    return res,err
}
//google api handle
func google_translate(words string) (string,error) {
    google_api_conf,err :=  read_conf("apiconf.json","google")
    if err !=nil {
        log.Panicf("get api conf fail : %s", err.Error())
    }
    translate_max_limit,_:=strconv.Atoi(google_api_conf["translate_max_limit"].(string))
    translate_trunk_v,_:=strconv.Atoi(google_api_conf["translate_trunk"].(string))
    //words string,max_limit int,trunk int,split_str string,api_name string,api_conf_data map[string]interface{}
    res,err:= translate_trunk(words,translate_max_limit,translate_trunk_v,".","google",nil)
    return res,err
}

//Segmented translation
func translate_trunk(words string,max_limit int,trunk int,split_str string,api_name string,api_conf_data map[string]interface{})(res string,err error) {
    words_str_len := len(words)
    if words_str_len > max_limit{
        words_arr := strings.Split(words, split_str)
        words_arr_len := len(words_arr)
        words_loop_mod := words_arr_len%trunk;
        words_loop := words_arr_len/trunk;
        if words_loop_mod != 0{
            words_loop++
        }
        var final_text string
        for i:=0;i<words_loop;i++{
            startindex := i*trunk
            endindex := i*trunk+trunk
            if endindex > words_arr_len{
                endindex = words_arr_len
            }
            words_trunk := words_arr[startindex:endindex]
            for _,w:=range words_trunk{
                trans_res,_ := switch_translate_api(api_name,api_conf_data,w)
                final_text += trans_res
            }
        }
        return final_text,nil
    }else{
        trans_res,err:=switch_translate_api(api_name,api_conf_data,words)
        return trans_res,err
    }
}
//choose api to translate
func switch_translate_api(api_name string,api_conf_data map[string]interface{},words string)(string,error) {
    if api_name == "youdao"{
        res,err := youdao_trans_trunk(api_conf_data["base_url"].(string),api_conf_data["key"].(string),words,api_conf_data["scret"].(string))
        return res,err
    }else{
        res,err := google_translate_trunk(words)
        return res,err
    }
}
func get_google_tkk(words string) (string,error)  {
    vm := otto.New()
    vm.Set("words",words)
    vm.Run(`
        var tkk = eval('((function(){var a\x3d1745533258;var b\x3d184024323;return 415516+\x27.\x27+(a+b)})())');
        var b = function (a, b) {
            for (var d = 0; d < b.length - 2; d += 3) {
                var c = b.charAt(d + 2),
                    c = "a" <= c ? c.charCodeAt(0) - 87 : Number(c),
                    c = "+" == b.charAt(d + 1) ? a >>> c : a << c;
                a = "+" == b.charAt(d) ? a + c & 4294967295 : a ^ c
            }
            return a
        }
        var tk =  function (a,TKK) {
            //console.log(a,TKK);
            for (var e = TKK.split("."), h = Number(e[0]) || 0, g = [], d = 0, f = 0; f < a.length; f++) {
                var c = a.charCodeAt(f);
                128 > c ? g[d++] = c : (2048 > c ? g[d++] = c >> 6 | 192 : (55296 == (c & 64512) && f + 1 < a.length && 56320 == (a.charCodeAt(f + 1) & 64512) ? (c = 65536 + ((c & 1023) << 10) + (a.charCodeAt(++f) & 1023), g[d++] = c >> 18 | 240, g[d++] = c >> 12 & 63 | 128) : g[d++] = c >> 12 | 224, g[d++] = c >> 6 & 63 | 128), g[d++] = c & 63 | 128)
            }
            a = h;
            for (d = 0; d < g.length; d++) a += g[d], a = b(a, "+-a^+6");
            a = b(a, "+-3^+b+-f");
            a ^= Number(e[1]) || 0;
            0 > a && (a = (a & 2147483647) + 2147483648);
            a %= 1E6;
            return a.toString() + "." + (a ^ h)
        }
        ret_tkk = tk(words,tkk) 
    `)
    value,err := vm.Get("ret_tkk")
    if err != nil{
        return "",err
    }else{
        tk,_:=value.ToString()
        return tk,nil
    }
}
func google_translate_trunk(words_trunk string) (string,error) {
    //replace special string
    words_trunk = strings.Replace(words_trunk, "&", "＆", -1)
    words_trunk = strings.Replace(words_trunk, ";", "；", -1)
    words_trunk = strings.Replace(words_trunk, "+", "＋", -1)
    google_api_conf,err :=  read_conf("apiconf.json","google")
    if err !=nil {
        log.Panicf("get api conf fail : %s", err.Error())
    }
    tkk,err:= get_google_tkk(words_trunk)
    if err !=nil {
        log.Panicf("get tkk fail : %s", err.Error())
    }
    base_url := google_api_conf["base_url"].(string) + "tk="+tkk
    req, err := http.NewRequest("POST", base_url, strings.NewReader("q="+words_trunk))
    req.Header.Add("content-type", "application/x-www-form-urlencoded")
    req.Header.Add("cache-control", "no-cache")
    res, err := http.DefaultClient.Do(req)
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)

    cn_json,err := simplejson.NewJson(body)
    cn_data,err:=cn_json.GetIndex(0).Array()
    var translate_str string
    for i,val:= range cn_data{
        if i < len(cn_data) - 1{
            for j,vv:=range val.([]interface{}){
                if j == 0{
                    translate_str += vv.(string)
                }
            }
        }else{
            break
        }
    }
    return translate_str,err
}


//youdao api area funcs
func youdao_trans_trunk(api_base_url string,app_key string,words string,app_scret string)(res string,err error) {
    salt := make_salt(4)
    generate_sign := generate_sign(app_key,words,salt,app_scret)
    words_ := url.QueryEscape(words)
    api_url := api_base_url + "?q=" + words_ + "&from=en&to=zh_CHS&appKey="+app_key+"&salt="+salt+"&sign="+generate_sign
    resp, err := http.Get(api_url)
    if err != nil {
        return words,err
    }
    defer resp.Body.Close()
    response, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return words,err
    }
    cn_json, err := simplejson.NewJson(response)
    cn_data,err := cn_json.Get("translation").Array() 
    if len(cn_data) == 0{
        return words,err
    }else{
        return cn_data[0].(string),err
    }
}

//creaet random strings
func make_salt(len int)(res string) {
    up_case_string := "ABCDEFGHIJKLMNPQRSTUVWXYZ"
    lowwer_case_string := "abcdefghijklmnpqrstuvwxyz"
    nums_string := "0123456789"
    spec_string := "!@#$%_-"
    strings_ := up_case_string+lowwer_case_string+nums_string+spec_string
    str_len := strings.Count(strings_, "") - 1
    // 根据时间设置随机数种子
    rand.Seed(int64(time.Now().Nanosecond()))
    var return_str string
    for i:=0;i<len;i++{
        randindex := rand.Intn(str_len)
        return_str += string([]byte(strings_)[randindex:randindex+1])
    }
    return return_str
}
func get_md5_str(s string) string {
    h := md5.New()
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}
//create youdao sign
func generate_sign(app_key string,q string,salt string,app_scret string)(res string) {
    new_string := get_md5_str(app_key+q+salt+app_scret)
    return new_string
}
////common funcs
//parse json data
func par_json(body []byte,field string) (res map[string]interface{},err error) {
    cn_json,err := simplejson.NewJson(body)
    cn_data,err := cn_json.Get(field).Map()
    return cn_data,err
}
//read conf from file
func read_conf(path string,api_name string)(res map[string]interface{},err error) {
    api_conf, err := ioutil.ReadFile(path)
    if err != nil {
        log.Panicf("readFile: %s", err.Error())
    }
    api_conf_data,err:= par_json(api_conf,api_name)
    return api_conf_data,err
}
