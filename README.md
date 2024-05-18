# MINI-BLINK

基于免费版的 [miniblink](https://miniblink.net/) 的 GO 封装，内嵌 miniblink 的 DLL，并根据构建环境自动选择 x86/x64 DLL

## 特点
1. 纯 GO 实现，无须写 C 代码
2. 封装了大部分 miniblink 的 API，面向对象，方便使用。
3. JS 交互（IPC通讯、事件绑定、调用JS函数），具体使用方式，请参考示例
4. 本地目录、BIN资源的加载
5. 内嵌 miniblink 的 DLL，并根据构建环境自动选择 x86/x64 DLL


#### 部分未封装的接口，可以使用以下函数直接调用 `miniblink` 的 DLL
```go
func (mb *Blink) CallFunc(name string, args ...uintptr) (r1 uintptr, r2 uintptr, err error)
```

## 开发环境
- GO 1.20


## 打包

### 打包标签:
- `release` 打包程序，不包含调试信息
- `slim` 不内嵌miniblink的dll，需要手动放入程序根目录或系统默认路径

### 示例
```bash
# 默认打包
go build \
  -tags 'release' \
  -ldflags '-w -s -H=windowsgui' \
  -o miniBlink.exe \
  ./samples/demo-baidu

# 打包32位程序
GOARCH=386 go build \
  -tags 'release' \
  -ldflags '-w -s -H=windowsgui' \
  -o miniBlink.exe \
  ./samples/demo-baidu

```

### 添加程序版本信息、图标。。。 请查看 `demo-baidu` 的 `main.go` 文件。并参阅 [josephspurrier/goversioninfo](https://github.com/josephspurrier/goversioninfo)
