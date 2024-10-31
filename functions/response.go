package functions

type LoginResponse struct {
	Status int    `json:"Status"`
	Msg    string `json:"Msg"`
	Data   struct {
		UserInfo struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"user_info"`
		AccessToken string `json:"access_token"`
	} `json:"Data"`
	Err string `json:"Err"`
}

type PingResponse struct {
	Msg string `json:"message"`
}

type RootDirResponse struct {
	Status int    `json:"Status"`
	Msg    string `json:"Msg"`
	Data   string `json:"Data"`
	Err    string `json:"Err"`
}

type FileItem struct {
	UUID      string `json:"uuid"`
	Filename  string `json:"filename"`
	Filesize  int    `json:"filesize"`
	Ext       string `json:"ext"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

type FileListData struct {
	Items []FileItem `json:"items"`
	Total int        `json:"total"`
}

type FileListResponse struct {
	Status int          `json:"Status"`
	Msg    string       `json:"Msg"`
	Data   FileListData `json:"Data"`
	Err    string       `json:"Err"`
}

type UploadFileItem struct {
	Filename string `json:"filename"`
	Filesize int    `json:"filesize"`
	Info     string `json:"info"`
}

type UploadFileData struct {
	Items []UploadFileItem `json:"items"`
	Total int              `json:"total"`
}

type UploadFileResponse struct {
	Status int            `json:"Status"`
	Msg    string         `json:"Msg"`
	Data   UploadFileData `json:"Data"`
	Err    string         `json:"Err"`
}

type DeleteFileRequest struct {
	Filename string `json:"filename"`
	ParentID string `json:"parent_id"`
}

type DeleteFileResponse struct {
	Status int         `json:"Status"`
	Msg    string      `json:"Msg"`
	Data   interface{} `json:"Data"`
	Err    string      `json:"Err"`
}

type DownloadFileRequest struct {
	Filename string `json:"filename"`
	ParentID string `json:"parent_id"`
}

type ErrorResponse struct {
	Status int    `json:"Status"`
	Msg    string `json:"Msg"`
	Err    string `json:"Err"`
}
