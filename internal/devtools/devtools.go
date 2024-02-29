package devtools

import "io"

func Asset(file string) ([]byte, error) {
	f, e := AssetFile().Open(file)

	if e != nil {
		return nil, e
	}

	defer f.Close()

	return io.ReadAll(f)
}
