package blink

type Env struct {
	isSYS64 bool
	isDebug bool
}

func IsSYS64() bool {
	return env.isSYS64
}

func IsDebug() bool {
	return env.isDebug
}

func IsRelease() bool {
	return !env.isDebug
}
