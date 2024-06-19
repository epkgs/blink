package resource

import (
	"embed"
	"errors"
	fs "io/fs"
	"net/http"
	netUrl "net/url"
	"os"
)

type Resource map[string]http.FileSystem

func New() *Resource {
	loader := make(Resource)
	return &loader
}

// Bind fileSystem to domain
//
// FileSystem accept type of below
//   - embed.FS
//   - fs.FS
//   - fs.SubFS
//   - http.FileSystem
//   - string of directory (The resource will not embed, you should copy the files to the target build directory)
func (res *Resource) Bind(domain string, fileSystem interface{}) (err error) {
	uri, err := netUrl.Parse(domain)
	if err != nil {
		return
	}
	dm := uri.Host
	if uri.Host == "" {
		dm = uri.Path
	}

	switch v := fileSystem.(type) {
	case http.FileSystem:
		(*res)[dm] = v
	case embed.FS:
		(*res)[dm] = http.FS(v)
	case string:
		(*res)[dm] = http.FS(os.DirFS(v))
	case fs.SubFS:
		(*res)[dm] = http.FS(v)
	case fs.FS:
		(*res)[dm] = http.FS(v)
	default:
		err = errors.New("fs type error, only accept: http.FileSystem, embed.FS, fs.FS, fs.SubFS or string of directory")
	}

	return
}

func (res *Resource) Unbind(domain string) {
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

func (res *Resource) IsExist(domain string) bool {
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

func (res *Resource) GetFile(url string) http.File {

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
