### translate-api 使用go翻译
-调用google和有道api翻译


### 使用
>安装
```
go get github.com/hadeswang2016/translate-api
```

>引入包
```
import github.com/hadeswang2016/translate-api
```

>调用
```
trans_title,_ := translate.Translate_google(title)

trans_words,_ := translate.Translate_youdao(words)

```
### Documentation

Visit the docs on [gopkgdoc](https://godoc.org/github.com/hadeswang2016/translate-api)