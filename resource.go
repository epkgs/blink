package blink

import (
	"net/http"
	netUrl "net/url"
	"os"
)

type ResourceLoader map[string]http.FileSystem

var Resource = NewResourceLoader()

func NewResourceLoader() *ResourceLoader {
	loader := make(ResourceLoader)
	return &loader
}

// bin-data
func (res *ResourceLoader) Bind(domain string, fs http.FileSystem) {
	uri, err := netUrl.Parse(domain)
	if err != nil {
		return
	}
	dm := uri.Host
	if uri.Host == "" {
		dm = uri.Path
	}
	(*res)[dm] = fs
}

// 本地资源
func (res *ResourceLoader) BindDir(domain string, dir string) {
	uri, err := netUrl.Parse(domain)
	if err != nil {
		return
	}
	dm := uri.Host
	if uri.Host == "" {
		dm = uri.Path
	}
	(*res)[dm] = http.FS(os.DirFS(dir))
}

func (res *ResourceLoader) Unbind(domain string) {
	uri, err := netUrl.Parse(domain)
	if err != nil {
		return
	}
	dm := uri.Host
	if uri.Host == "" {
		dm = uri.Path
	}
	delete(*res, dm)
}

func (res *ResourceLoader) IsExist(domain string) bool {
	uri, err := netUrl.Parse(domain)
	if err != nil {
		return false
	}
	dm := uri.Host
	if uri.Host == "" {
		dm = uri.Path
	}
	_, exist := (*res)[dm]
	return exist
}

func (res *ResourceLoader) GetFile(url string) http.File {

	uri, err := netUrl.Parse(url)
	if err != nil {
		return nil
	}

	//只响应http
	if uri.Scheme != "http" && uri.Scheme != "https" {
		return nil
	}

	dm := uri.Host
	if uri.Host == "" {
		dm = uri.Path
	}
	fs, exist := (*res)[dm]
	if !exist {
		return nil
	}

	f, err := fs.Open(uri.Path)
	if err != nil {
		return nil
	}

	return f
}
