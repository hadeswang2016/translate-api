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
)  




func Get_trans_content(translator int,word_list []string,db_conf map[string]string)(int64,error) {
    timestamp := time.Now().Unix()
    db, err := sql.Open("mysql", db_conf["User"]+":"+db_conf["Password"]+"@tcp("+db_conf["Host"]+":"+db_conf["Port"]+")/"+db_conf["Dbname"]+"?charset=utf8") 
    checkErr(err)  
    //插入数据  
    var table string
    if translator == 1{
        table = "translate_youdao"
    }else if translator == 2{
        table = "translate_google"
    }
    stmt, err := db.Prepare("INSERT INTO "+table+" SET title=?,intro=?,cover_image=?,content=?,url=?,ref_url=?,tag=?,source=?,translate_time=?")  
    checkErr(err)  
    res, err := stmt.Exec(word_list[0], word_list[1], word_list[2],word_list[3],word_list[4],word_list[5],word_list[6],word_list[7],timestamp)  
    checkErr(err)  
    id, err := res.LastInsertId()  
    checkErr(err)  
    fmt.Printf("INSERT INTO "+table+" SUCCESS id:%d\r\n",id)
    return id,err
}


func Translate(words string,translator int)(res string) {
	var result string
	if translator == 1{
		result = youdao_tanslate(words)
	}else{
		result = google_translate(words)
	}
	return result
}
/*
	有道翻译接口 字数限制 
*/
func youdao_tanslate(words string)(res string) {
	api_base_url := "http://openapi.youdao.com/api"
	app_key := "5ee079a739e1c466"
	app_scret := "O3ZmaVtRc1qlHaAQeLMinw6slKWHdyKc"
	words_str_len := len(words)
	//文本超过3000 分段翻译
	if words_str_len > 3000{
		words_arr := strings.Split(words, ".")//按句子分割
		words_len := len(words_arr)
		words_step := 25 //每次翻译25个句子
		words_loop_mod := words_len%words_step;
		words_loop := words_len/words_step;
		if words_loop_mod != 0{
			words_loop++
		}
		var final_text string
		for i:=0;i<words_loop;i++{
			startindex := i*words_step
			endindex := i*words_step+words_step
			if endindex > words_len{
				endindex = words_len
			}
			words_trunk := words_arr[startindex:endindex]
			for _,w:=range(words_trunk){
				final_text += youdao_trans_batch(api_base_url,app_key,w,app_scret)
			}
		}
		return final_text
	}else{
		//未超过文本翻译限制 直接翻译
		return youdao_trans_batch(api_base_url,app_key,words,app_scret)
	}

}


func youdao_trans_batch(api_base_url string,app_key string,words string,app_scret string)(res string) {
	salt := make_salt(4)
	generate_sign := generate_sign(app_key,words,salt,app_scret)
	words_ := url_en_code(words)
	api_url := api_base_url + "?" + words_ + "&from=en&to=zh_CHS&appKey="+app_key+"&salt="+salt+"&sign="+generate_sign
	response,err := httpGet(string(api_url))
    if err != nil{
		return "error"    	
    }
    cn_json, _ := simplejson.NewJson(response)
	cn_data,_ := cn_json.Get("translation").Array() 
	 // fmt.Println(words,"//",cn_data,"//")
	if len(cn_data) == 0{
		return words
	}else{
		return cn_data[0].(string)
	}
}


func generate_sign(app_key string,q string,salt string,app_scret string)(res string) {
	new_string := get_md5_str(app_key+q+salt+app_scret)
	return new_string
}

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



func httpGet(url string)(body_r []byte,err_ error) {
    resp, err := http.Get(url)
    if err != nil {
        return 
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }
    return body,err
}


func url_en_code(url_r string)(res string) {
	l3, _ := url.Parse("http://www.baidu.com?q="+url_r)
	return l3.Query().Encode()
}


func checkErr(err error) {  
    if err != nil {  
        panic(err)  
    }  
}  

///////////////google 翻译区


func parse_js(words string) string  {
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
		return "error"
	}else{
		tk,_:=value.ToString()
		return tk
	}
}



func google_translate(words string) string {
	words_str_len := len(words)
	//文本超过4000 分段翻译
    if words_str_len > 5000{
        words_arr := strings.Split(words, ".")//按句子分割
        words_len := len(words_arr)
        words_step := 20 //每次翻译20个句子
        words_loop_mod := words_len%words_step;
        words_loop := words_len/words_step;
        if words_loop_mod != 0{
            words_loop++
        }
        var final_text string
        for i:=0;i<words_loop;i++{
            startindex := i*words_step
            endindex := i*words_step+words_step
            if endindex > words_len{
                endindex = words_len
            }
            words_trunk := words_arr[startindex:endindex]
            for _,v:=range words_trunk{
            	final_text += google_translate_batch(v)
            }
        }
        return final_text
    }else{
        //未超过文本翻译限制 直接翻译
        return google_translate_batch(words)
    }
}

func google_translate_batch(words_trunk string)string {
	words_trunk = strings.Replace(words_trunk, "&", "＆", -1)
	tk := parse_js(words_trunk)
	to_trans_words := url_en_code(words_trunk)
	url := "http://192.168.1.81/Ad/Lmb/translate_for_go?tk="+tk+"&words="+to_trans_words
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	return string(body)
}