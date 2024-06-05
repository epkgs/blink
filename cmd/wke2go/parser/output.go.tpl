package {{ .PkgName }}

import (
	"log"

	"golang.org/x/sys/windows"
)

type DLL struct {
	*windows.LazyDLL
	{{- range $cname, $gname := .FuncMapping }}
	_{{ $cname }} *Proc
	{{- end}}
}

type Proc struct {
	*windows.LazyProc
}

var globalDLL *DLL

func InitDLL(dll *windows.LazyDLL) *DLL {
	globalDLL := &DLL{
		LazyDLL: dll,
	}

	{{- range $cname, $gname := .FuncMapping }}
	globalDLL._{{ $cname }} = globalDLL.NewProc("{{ $cname }}")
	{{- end}}

	return globalDLL
}

func (dll *DLL) NewProc(name string) *Proc {
	proc := dll.LazyDLL.NewProc(name)
	return &Proc{
		LazyProc: proc,
	}
}

func (proc *Proc) Call(args ...uintptr) (r1 uintptr, r2 uintptr, err error) {
	defer func() {
		if r := recover(); r != nil {

			if r == windows.NOERROR {
				err = nil
				return
			}

			err = r.(error)
			log.Fatalf("Panic by Call Proc: %s, with Err: %s", proc.Name, err.Error())
		}
	}()

	r1, r2, err = proc.LazyProc.Call(args...)

	if err == windows.NOERROR {
		err = nil
	}

	return
}

