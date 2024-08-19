package utils

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseNetscapeCookieFile(filePath string) ([]*http.Cookie, error) {

	cookies := make([]*http.Cookie, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return cookies, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "# ") { // 忽略以#开头的注释行，#后加个空格，否则会跳过token的注释行
			continue
		}
		parts := strings.Fields(line) // 用空格分割行
		if len(parts) >= 7 {
			expires, _ := strconv.ParseInt(parts[4], 10, 64)
			// 创建http.Cookie实例
			cookie := &http.Cookie{
				Name:     parts[5],
				Value:    parts[6],
				Path:     parts[2],
				Domain:   strings.TrimPrefix(parts[0], "#HttpOnly_"), // 必须去掉前缀#HttpOnly_
				HttpOnly: parts[1] == "TRUE",
				Secure:   parts[3] == "TRUE",
				Expires:  time.Unix(expires, 0),
			}
			cookies = append(cookies, cookie)
			//fmt.Println(cookie.String())
		}
	}

	if err := scanner.Err(); err != nil {
		return cookies, fmt.Errorf("reading file: %v", err)
	}

	return cookies, nil
}
