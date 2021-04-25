# Error

## Go Error

* **fatal error 才使用panic**

* 例如索引越界、不可恢复的环境问题、栈溢出，对于其他的错误情况，使用error来进行判定

* **Error are values**

## Error Type

### Sentinel Error
  
* 预定义的特定错误

```go
if err == ErrSomething{...}
//类似 io.EOF
```

* 使用sentinel值是最不灵活的错误处理策略，因为必须使用==来进行比较预先声明的值。提供更多上下文时，返会不同的值将破坏相等性检查
* 上下文会破坏调用者的 ==  

* **不依赖检查 error.Error的输出**

* sentinel errors 成为api公共部分

* sentinel errors 在两个包之间创建了依赖

* **结论：尽可能避免使用sentinel errors**

### Error types

* 实现了error接口的自定义类型。

```go
// /usr/local/go/src/io/fs/fs.go
type PathError struct {
	Op   string
	Path string
	Err  error
}

```

* 调用者要使用类型断言和类型switch，要让自定义的error变为public。会导致和调用者产生强耦合，从而导致API变得脆弱
* **结论：尽量避免使用 error types，虽然错误类型比sentinel error更好，但是error types 共享error values许多相同的问题，至少避免将他们作为公共API的一部分**

### Opaque errors

* 不透明的，最灵活的处理方式
* 因为虽然知道发生了错误，但没有能力看到错误的内部，知道就是起作用了，或者没有起作用
* **不透明错误处理：只需返回错误而不假设其内容**

* **assert errors for behaviour,not type**

## Handing Error

### Indented flow is for errors

* 无错误的正常流程代码，将成为一条直线，而不是缩进的代码

### Wrap errors

* 没有生成错误的 file:line信息。没有导致错误的调用堆栈的堆栈追踪。  
        类似代码处理错误，记录日志并且返回错误

* 日志记录与错误无关且对调试没有帮助的信息应被视为噪音，应予以质疑。**记录日志的原因是因为某些东西失败了，而日志包含了答案**

* **github.com/pkg/errors**

* 使用pkg/errors包，向错误值添加上下文，既可以由人也可以由机器检查

```go
func main(){
    _,err := ReadConfig()
    if(err!=nil){
        fmt.Printf("original error: %T %v\n",errors.Cause(err),errors.Cause(err))
        fmt.Printf("stack trace:\n %+v \n",err)
        os.Exit(1)
    }
}

```

* 在应用代码中，使用 errors.New 或者 errors.Errorf 返回错误

```go
//调用自己方法错误
func parseArgs(args []string) error{
    if len(args > 3){
        return errors.Errorf("not enough arguments......")
    }
}
```

* 调用其他库 包装起来 往上抛

```go
f,err := os.Open(path)
if err!= nil {
    return errors.Wrapf(err,"failed to open %q",path)
}

```

* 直接返回错误，而不是每个错误地方到处打日志

* 在程序的顶部或者是工作的goroutine 顶部 （请求入口），使用 %+v 把堆栈详情记录

```go
func main(){
    err:= app.Run()
    if err!=nil {
        fmt.Printf("FATAL: %+v",err)
        os.Exit(1)
    }
}

```

* 使用errors.Cause 获取root error，再进行和sentinel error判定

### 总结

* 选择warp error是只有 applications 可以选择应用的策略。 具有最高可重用性的包只能返回根错误值。
    此机制与Go标准库中使用的相同（kit库的sql.ErrNoRows）。

* 一旦确定函数/方法将处理错误，错误就不再是错误，需要降级处理。（返回0 或者nil）

## Error Inspection

### Error before Go 1.13

* 最简单的错误检查

```go
if err!=nil{
    //something went wrong
}
```

* 有时需要对sentinel error 进行检查

```go
var ErrNotFound = errors.New("not found")

if err == ErrNotFound{
    //something wasn't found
}

```

* 实现了 error interface的自定义error struct，进行断言使用获取更丰富的上下文

```go
type NotFoundError struct{
    Name string
}
func (e *NotFoundError) Error() string {return e.Name +" : not found"}

if e,ok:=err.(*NotFoundError);ok {
    //e.Name wasn't found
}
```

### Unwrap

* go1.13 errors包包含两个用于检查错误的新函数 ： Is 和 As

```go
if errors.Is(err,ErrNotFound){
    //something wasn't found
}
var e *QueryError
if errors.As(err,&e){
    //err is a *QueryError,and e is set to the error's value
}
```

### Wrapping errors with %w

* 如前所述，使用fmt.Errorf向错误添加附加信息

```go
if err!=nil {
    return fmt.Errorf("decompress %v:%v",name,err)
}
```

* 在Go 1.13中 fmt.Errorf 支持新的 %w 谓词

```go
if err!=nil {
    //return an error which unwraps to err
    return fmt.Errorf("decompress %v:%w",name,err)
}
```

* 用 %w 包装错误可用于 errors.Is 以及 errors.As 

```go
err:= fmt.Errorf("access denied : %w",ErrorPermission)
//...
if errors.Is(err,ErrorPermission){
    //...
}
```
