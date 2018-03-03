package models

import (
	"archive/zip"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	mb                  = 1024 * 1024
	maxCaseFileSize     = 10
	maxCaseZipTotalSize = 50
	maxCaseFileCount    = 500
)

var (
	errCaseFileCountLimitExceeded     = fmt.Errorf("ファイルの数は%v個以下にしてください。", maxCaseFileCount)
	errCaseFileSizeLimitExceeded      = fmt.Errorf("展開後の各ファイルのサイズは%vMiB以下にしてください。", maxCaseFileSize)
	errTotalCaseFileSizeLimitExceeded = fmt.Errorf("展開後の合計ファイルサイズは%vMiB以下にしてください。", maxCaseZipTotalSize)
)

func GetBcryptCost() int {
	//適当に調整する
	return 14
}

// length bytesのランダムなBase64エンコードされた文字列を返す
func GenerateRandomBase64String(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

func MaxLong(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func MaxDuration(a, b time.Duration) time.Duration {
	if a < b {
		return b
	}
	return a
}

func DefaultString(value, def string) string {
	if value == "" {
		return def
	}
	return value
}

func EqualTime(t1, t2 time.Time) bool {
	if t1.Unix() < 0 {
		t1 = time.Unix(-t1.Unix(), 0)
	}
	if t2.Unix() < 0 {
		t2 = time.Unix(-t2.Unix(), 0)
	}

	u1 := t1.Add(500 * time.Millisecond).Unix()
	u2 := t2.Add(500 * time.Millisecond).Unix()
	return u1 == u2
}

func checkValidZip(reader *zip.Reader) error {
	var total uint64

	if maxCaseFileCount < len(reader.File) {
		return errCaseFileCountLimitExceeded
	}

	for _, f := range reader.File {
		if f.UncompressedSize64 <= maxCaseFileSize*mb {
			total += f.UncompressedSize64 / mb
			continue
		}

		return errCaseFileSizeLimitExceeded
	}

	if maxCaseZipTotalSize < total {
		return errTotalCaseFileSizeLimitExceeded
	}

	return nil
}

func UniqueUsers(users []User) []User {
	m := make(map[User]bool, len(users))
	res := make([]User, 0, len(users))
	for _, u := range users {
		if m[u] {
			continue
		}
		m[u] = true
		res = append(res, u)
	}

	return res
}
