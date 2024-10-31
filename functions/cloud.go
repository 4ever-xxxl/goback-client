package functions

import (
	"bytes"
	"encoding/json"
	"goback-client/data"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetRootDir() {
	var rootDirURL = URL + "/api/v1/file/root"
	req, _ := http.NewRequest("GET", rootDirURL, nil)
	req.Header.Set("Authorization", data.Config.JWTTOKEN)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var rootDirResp RootDirResponse
	err = json.NewDecoder(resp.Body).Decode(&rootDirResp)
	if err != nil {
		log.Println("Error decoding root dir response:", err)
		return
	}

	if rootDirResp.Status == 200 {
		data.Config.RootDir = rootDirResp.Data
	}
}

func GetFileList(s string) {
	var fileListURL = URL + "/api/v1/file/list"
	jsonData := []byte(`{"parent_id":"` + s + `"}`)
	req, _ := http.NewRequest("GET", fileListURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", data.Config.JWTTOKEN)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var fileListResp FileListResponse
	err = json.NewDecoder(resp.Body).Decode(&fileListResp)
	if err != nil {
		log.Println("Error decoding file list response:", err)
		return
	}

	if fileListResp.Status == 200 {
		var cloudFiles []data.CloudFile
		for _, item := range fileListResp.Data.Items {
			cloudFiles = append(cloudFiles, data.CloudFile{
				UUID:     item.UUID,
				Filename: item.Filename,
				Filesize: item.Filesize,
				Ext:      item.Ext,
				CreatedAt: func() time.Time {
					t, err := time.Parse(time.RFC3339, item.CreatedAt)
					if err != nil {
						log.Println("Error parsing time:", err)
					}
					return t
				}(),
				UpdatedAt: func() time.Time {
					t, err := time.Parse(time.RFC3339, item.UpdatedAt)
					if err != nil {
						log.Println("Error parsing time:", err)
					}
					return t
				}(),
			})
		}
		data.CloudFileList = cloudFiles
	}
}

func UploadFiles(files []string, parentDir string) error {
	var uploadURL = URL + "/api/v1/file/upload"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for _, file := range files {
		// 获取文件的基本名称
		fileName := filepath.Base(file)

		fileWriter, err := writer.CreateFormFile("files", fileName)
		if err != nil {
			log.Println("Error creating form file:", err)
			return err
		}

		fh, err := os.Open(file)
		if err != nil {
			log.Println("Error opening file:", err)
			return err
		}
		defer fh.Close()

		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			log.Println("Error copying file:", err)
			return err
		}
	}

	err := writer.WriteField("parentDir", parentDir)
	if err != nil {
		log.Println("Error writing parentDir field:", err)
		return err
	}

	contentType := writer.FormDataContentType()
	err = writer.Close()
	if err != nil {
		log.Println("Error closing writer:", err)
		return err
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Authorization", data.Config.JWTTOKEN)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	var uploadResp UploadFileResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadResp)
	if err != nil {
		log.Println("Error decoding upload response:", err)
		return err
	}

	if uploadResp.Status == 200 {
		log.Println("[Cloud] Upload success")
	} else {
		log.Println("[Cloud] Upload failed:", uploadResp.Err)
	}

	return nil
}

func DeleteFile(filename, parentID string) error {
	var deleteURL = URL + "/api/v1/file/delete"

	deleteReq := DeleteFileRequest{
		Filename: filename,
		ParentID: parentID,
	}

	jsonData, err := json.Marshal(deleteReq)
	if err != nil {
		log.Println("Error marshalling delete request:", err)
		return err
	}

	req, err := http.NewRequest("DELETE", deleteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Authorization", data.Config.JWTTOKEN)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	var deleteResp DeleteFileResponse
	err = json.NewDecoder(resp.Body).Decode(&deleteResp)
	if err != nil {
		log.Println("Error decoding delete response:", err)
		return err
	}

	if deleteResp.Status == 200 {
		log.Println("[Cloud] Delete success")
	} else {
		log.Println("[Cloud] Delete failed:", deleteResp.Err)
	}
	return nil
}

func DownloadFile(filename, parentID, saveDir string) error {
	var downloadURL = URL + "/api/v1/file/download"

	downloadReq := DownloadFileRequest{
		Filename: filename,
		ParentID: parentID,
	}

	jsonData, err := json.Marshal(downloadReq)
	if err != nil {
		log.Println("Error marshalling download request:", err)
		return err
	}

	req, err := http.NewRequest("GET", downloadURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Authorization", data.Config.JWTTOKEN)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			log.Println("Error decoding error response:", err)
			return err
		}
		log.Println("Download failed:", errResp.Msg, errResp.Err)
		return nil
	}

	savePath := filepath.Join(saveDir, filename)
	out, err := os.Create(savePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Error saving file:", err)
		return err
	}

	log.Println("File downloaded successfully:", savePath)
	return nil
}
