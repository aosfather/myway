package core

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"
	"strconv"
	"time"
)

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b

}

func Int2BytesTo(v int, ret []byte) {
	ret[0] = byte(v >> 24)
	ret[1] = byte(v >> 16)
	ret[2] = byte(v >> 8)
	ret[3] = byte(v)
}

func Byte2Int(data []byte) int {
	return int((int(data[0])&0xff)<<24 | (int(data[1])&0xff)<<16 | (int(data[2])&0xff)<<8 | (int(data[3]) & 0xff))
}

const (
	timeLayout = "2006-01-02 15:04:05"
)

func GetNowDate() string {
	return time.Now().Format(timeLayout)
}

func ToDate(t time.Time) string {
	return t.Format(timeLayout)
}

//字符生成md5
func GetMd5str(value string) string {
	data := []byte(value)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1

}

func ToInt(value string) int64 {
	result, err := strconv.Atoi(value)
	if err == nil {
		return int64(result)
	}
	return 0
}

func CreateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(uuid)
	return uuid
}
