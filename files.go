package b2

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	DELETE_FILE_VERSION_URL string = "/b2api/v1/b2_delete_file_version"
	DOWNLOAD_FILE_BY_ID_URL string = "/b2api/v1/b2_download_file_by_id"
)

type FileVersion struct {
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
}

type FileInfo struct {
	ContentLength int
	ContentType   string
	FileId        string
	FileName      string
	ContentSha1   string
	Info          map[string]string
}

func (c *Client) DeleteFileVersion(fileName string, fileId string) (*FileVersion, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"fileName": "%s", "fileId": "%s"}`, fileName, fileId))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(DELETE_FILE_VERSION_URL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result FileVersion
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) DownloadFileById(fileId string) ([]byte, *FileInfo, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"fileId": "%s"}`, fileId))
	if req, err := http.NewRequest("POST", c.buildFileRequestUrl(DOWNLOAD_FILE_BY_ID_URL), reqBody); err != nil {
		return nil, nil, err
	} else {
		c.setHeaders(req)

		data, header, err := c.requestBytes(req)

		if err != nil {
			return nil, nil, err
		}

		contentLen, err := strconv.Atoi(header.Get("Content-Length"))

		info := &FileInfo{
			ContentLength: contentLen,
			ContentType:   header.Get("Content-Type"),
			FileId:        header.Get("X-Bz-File-Id"),
			FileName:      header.Get("X-Bz-File-Name"),
			ContentSha1:   header.Get("X-Bz-Content-Sha1"),
			Info:          make(map[string]string),
		}

		for key, _ := range header {
			// make the key lowercase as some of the headers B2 returns are lowercase, some uppercase...
			if strings.HasPrefix(strings.ToLower(key), "x-bz-info-") {
				info.Info[key] = header.Get(key)
			}
		}

		return data, info, nil
	}
}

func (c *Client) DownloadFileByName(bucketName string, fileName string) ([]byte, *FileInfo, error) {
	requestPath := fmt.Sprintf("/%s/%s", bucketName, fileName)
	if req, err := http.NewRequest("GET", c.buildFileRequestUrl(requestPath), nil); err != nil {
		return nil, nil, err
	} else {
		req.Header.Set("Authorization", c.AuthToken)

		fmt.Println(req)

		data, header, err := c.requestBytes(req)

		if err != nil {
			return nil, nil, err
		}

		contentLen, err := strconv.Atoi(header.Get("Content-Length"))

		info := &FileInfo{
			ContentLength: contentLen,
			ContentType:   header.Get("Content-Type"),
			FileId:        header.Get("X-Bz-File-Id"),
			FileName:      header.Get("X-Bz-File-Name"),
			ContentSha1:   header.Get("X-Bz-Content-Sha1"),
			Info:          make(map[string]string),
		}

		for key, _ := range header {
			// make the key lowercase as some of the headers B2 returns are lowercase, some uppercase...
			if strings.HasPrefix(strings.ToLower(key), "x-bz-info-") {
				info.Info[key] = header.Get(key)
			}
		}

		return data, info, nil
	}
}
