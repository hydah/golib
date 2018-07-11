## 解析ini格式配置到结构体中

Want more objective way to play with INI? Cool.

```ini
[note]
content = Hi is a good man!
city = HangZhou

[person]
name = huaiyann
age = 18
```

```go
type NoteSt struct {
	Content string `ini:"content"`
	City    string `ini:"city"`
}

type PersonSt struct {
	Name string `ini:"name"`
	Age  int64  `ini:"age"`
}

type Config struct{
	Note   NoteSt   `ini:"note"`
	Person PersonSt `ini:"person"`
}

func main() {
	cfg:=new(Config)
	_, err := cfgcenter.LoadConfig("path/to/ini",cfg)
	// ...
}
```

### 不支持除结构体指针外的其他指针类型

```ini
[note]
content = Hi is a good man!
city = HangZhou

[person]
name = huaiyann
age = 18
```

```go
type NoteSt struct {
	Content string `ini:"content"`
	City    string `ini:"city"`
}

type PersonSt struct {
	Name string `ini:"name"`
	Age  *int64 `ini:"age"`
}

type Config struct{
	Note   *NoteSt  `ini:"note"`
	Person PersonSt `ini:"person"`
}

func main() {
	cfg:=new(Config)
	_, err := cfgcenter.LoadConfig("path/to/ini",cfg)
	// cfg.Note，结构体指针，被正常解析
	// cfg.Note.City == "HangZhou"
	// cfg.Person.Age 不支持，不会被解析
	// cfg.Person.Age == nil
}
```

### 结构体成员、结构体指针成员，解析时会按照独立的section处理

```ini
[note]
content = Hi is a good man!
city = HangZhou

[person]
name = huaiyann
age = 18
```

```go
type NoteSt struct {
	Content string `ini:"content"`
	City    string `ini:"city"`
}

type PersonSt struct {
	Note NoteSt `ini:"note"`
	Name string `ini:"name"`
	Age  int64 `ini:"age"`
}

type Config struct{
	Person PersonSt `ini:"person"`
}

func main() {
	cfg:=new(Config)
	_, err := cfgcenter.LoadConfig("path/to/ini",cfg)
	// cfg.Person.Note，结构体成员，按照独立的section解析
	// cfg.Person.Name == "huaiyann"
	// cfg.Person.Note.City == "HangZhou"
	// cfg.Person.Note.Content == "Hi is a good man!"
}
```

### 继承时，根据有没有tag，不同处理

有tag的按照tag取section，没有tag的展开

```ini
[note]
content = Hi is a good man!
city = HangZhou

[person]
name = huaiyann
age = 18
```

```go
type NoteSt struct {
	Content string `ini:"content"`
	City    string `ini:"city"`
}

type PersonSt struct {
	Name string `ini:"name"`
	Age  int64 `ini:"age"`
}

type BaseSt struct {
	Note NoteSt `ini:"note"`
}

type Config struct{
	BaseSt
	PersonSt `ini:"person"`
}

func main() {
	cfg:=new(Config)
	_, err := cfgcenter.LoadConfig("path/to/ini",cfg)
	// cfg.Name == "huaiyann" 有tag，取对应tag的section
	// cfg.BaseSt.Note.Content == "Hi is a good man!" 没有tag的展开，展开后分别处理各个成员
}
```