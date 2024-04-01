# MINI-BLINK

## 介绍
基于免费版的 [miniblink](https://miniblink.net/) 的 GO 封装

1. 不使用 CGO
2. 面向对象
3. JS 交互以及事件绑定
4. 本地目录、BIN资源的加载
5. 内嵌 miniblink 的 DLL，并根据构建环境自动选择 x86/x64 DLL


#### 部分未封装的接口，可以使用以下函数直接调用 `miniblink` 的 DLL
```go
func (mb *Blink) CallFunc(name string, args ...uintptr) (r1 uintptr, r2 uintptr, err error)
```

## 开发环境
- GO 1.22


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
