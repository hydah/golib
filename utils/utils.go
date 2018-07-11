package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"

	math_rand "math/rand"
)

func GenerateRandom16Bytes() []byte {
	data := make([]byte, 16)
	io.ReadFull(rand.Reader, data)
	return data
}

func CurrentTimeMilliseconds() int64 {
	return time.Now().UnixNano() / 1000000
}

func CurrentTimeSeconds() int64 {
	return time.Now().Unix()
}

func TimestampToTimeString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func GetCurrentTime() string {
	now := time.Now()
	month, day, year := now.Month(), now.Day(), now.Year()
	hour, minute, second := now.Hour(), now.Minute(), now.Second()
	millisecond := now.Nanosecond() / 1000000
	longTime := Int64ToStrWithZeroPrefix(int64(hour), 2) + ":" + Int64ToStrWithZeroPrefix(int64(minute), 2) +
		":" + Int64ToStrWithZeroPrefix(int64(second), 2) + "." + Int64ToStrWithZeroPrefix(int64(millisecond), 3)
	longDate := Int64ToStrWithZeroPrefix(int64(year), 4) + "/" + Int64ToStrWithZeroPrefix(int64(month), 2) +
		"/" + Int64ToStrWithZeroPrefix(int64(day), 2)
	return longDate + " " + longTime
}

func GetCurrentDate() string {
	now := time.Now()
	month, day, year := now.Month(), now.Day(), now.Year()
	longDate := fmt.Sprintf("%04d%02d%02d", year, month, day)
	return longDate
}

func StringToMd5(s string) string {
	m := md5.New()
	io.WriteString(m, s)
	return fmt.Sprintf("%x", m.Sum(nil))
}

func MD5(s string) string {
	m := md5.New()
	io.WriteString(m, s)
	return base64.StdEncoding.EncodeToString(m.Sum(nil))
}

func ComputeFileMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	var result []byte
	h := hash.Sum(result)
	return fmt.Sprintf("%x", h), nil
}

func Base64StdEncode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64URLEncode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func Base64StdDecode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func Base64URLDecode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func URLEncode(src string) string {
	return url.QueryEscape(src)
}

func URLDecode(src string) (string, error) {
	return url.QueryUnescape(src)
}

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

func Uint16ToBytesExt(v uint16, buf []byte) {
	binary.BigEndian.PutUint16(buf, v)
}

func BytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

func Uint32ToBytesExt(v uint32, buf []byte) {
	binary.BigEndian.PutUint32(buf, v)
}

// BytesToUint32 BigEndian: []byte --> uint32
func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// Uint64ToBytes BigEndian: uint64 --> []byte
func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// BytesToUint64 BigEndian: []byte --> uint64
func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Int64ToStrWithZeroPrefix(v int64, zeroCnt int) (s string) {
	s = strconv.FormatInt(v, 10)
	needZero := zeroCnt - len(s)
	for i := 0; i < needZero; i++ {
		s = "0" + s
	}
	return
}

func GetInt64FromInterface(data interface{}) (value int64, err error) {
	switch v := data.(type) {
	case int64:
		value = int64(v)
	case uint64:
		value = int64(v)
	case int32:
		value = int64(v)
	case uint32:
		value = int64(v)
	case int16:
		value = int64(v)
	case uint16:
		value = int64(v)
	case int8:
		value = int64(v)
	case uint8:
		value = int64(v)
	case int:
		value = int64(v)
	case float64:
		value = int64(v)
	case float32:
		value = int64(v)
	default:
		err = fmt.Errorf("Can not convert %v to int64.", reflect.TypeOf(data))
	}
	return
}

func GetFloat64FromInterface(data interface{}) (value float64, err error) {
	switch v := data.(type) {
	case int64:
		value = float64(v)
	case uint64:
		value = float64(v)
	case int32:
		value = float64(v)
	case uint32:
		value = float64(v)
	case int16:
		value = float64(v)
	case uint16:
		value = float64(v)
	case int8:
		value = float64(v)
	case uint8:
		value = float64(v)
	case int:
		value = float64(v)
	case float64:
		value = float64(v)
	case float32:
		value = float64(v)
	default:
		err = fmt.Errorf("Can not convert %v to float64.", reflect.TypeOf(data))
	}
	return
}

func ToInt(data interface{}) (value int64) {
	switch v := data.(type) {
	case int64:
		value = int64(v)
	case uint64:
		value = int64(v)
	case int32:
		value = int64(v)
	case uint32:
		value = int64(v)
	case int16:
		value = int64(v)
	case uint16:
		value = int64(v)
	case int8:
		value = int64(v)
	case uint8:
		value = int64(v)
	case int:
		value = int64(v)
	case float64:
		value = int64(v)
	case float32:
		value = int64(v)
	}
	return
}

// String convert to string
func String(any interface{}) string {
	switch val := any.(type) {
	case int, uint, int64, uint64, uint32, int32, uint8, int8, int16, uint16:
		return fmt.Sprintf("%d", val)
	case string, []byte:
		return fmt.Sprintf("%s", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func Krand(size int) string {
	ikind, kinds, result := 3, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := true
	math_rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = math_rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + math_rand.Intn(scope))
	}

	return string(result)
}

func Dec2base64(num int64) (res string) {
	res = ""
	if num < 0 {
		return
	}
	codemap := map[int32]string{
		0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J",
		10: "K", 11: "L", 12: "M", 13: "N", 14: "O", 15: "P", 16: "Q", 17: "R", 18: "S", 19: "T",
		20: "U", 21: "V", 22: "W", 23: "X", 24: "Y", 25: "Z", 26: "a", 27: "b", 28: "c", 29: "d",
		30: "e", 31: "f", 32: "g", 33: "h", 34: "i", 35: "j", 36: "k", 37: "l", 38: "m", 39: "n",
		40: "o", 41: "p", 42: "q", 43: "r", 44: "s", 45: "t", 46: "u", 47: "v", 48: "w", 49: "x",
		50: "y", 51: "z", 52: "0", 53: "1", 54: "2", 55: "3", 56: "4", 57: "5", 58: "6", 59: "7",
		60: "8", 61: "9", 62: "+", 63: "/",
	}
	for _tk := true; _tk || num >= 1; _tk = false {
		res = codemap[int32(num&63)] + res
		num = num >> 6
	}
	return
}
