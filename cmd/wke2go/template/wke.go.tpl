package miniblink

import (
	"syscall"
)

{{range $i, $fc := .Funcs }}

{{- if $fc.HasComment}}
{{$fc.Comment}}
{{- end}}
func {{$fc.Name}}({{$fc.Args}}){{$fc.Result}} {
	syscall.SyscallN("{{$fc.Name}}", {{$fc.Args}})
}
{{end}}
