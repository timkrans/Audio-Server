package utils

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "time"
)

type LocalStorage struct {
    BasePath string
}

func (ls LocalStorage) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
    os.MkdirAll(ls.BasePath, os.ModePerm)
    filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
    fullPath := filepath.Join(ls.BasePath, filename)
    out, err := os.Create(fullPath)
    if err != nil {
        return "", err
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    if err != nil {
        return "", err
    }
    return fullPath, nil
}