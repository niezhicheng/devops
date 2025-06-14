package controllers

import (
	"archive/zip"
	"devops/global"
	"devops/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpFileInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	Type        string    `json:"type"`
	ModifyTime  time.Time `json:"modifyTime"`
	Permissions string    `json:"permissions"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

// 获取SFTP文件列表
func GetSftpFiles(c *gin.Context) {
	hostID := c.Param("id")
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	// 获取主机信息
	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接SSH
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer sshClient.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP连接失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 获取文件列表
	files, err := sftpClient.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("读取目录失败: %v", err)})
		return
	}

	// 构建响应数据
	var fileList []SftpFileInfo
	for _, file := range files {
		fileInfo := SftpFileInfo{
			Name:        file.Name(),
			Path:        filepath.Join(path, file.Name()),
			Size:        file.Size(),
			ModifyTime:  file.ModTime(),
			Permissions: file.Mode().String(),
		}

		if file.IsDir() {
			fileInfo.Type = "directory"
		} else {
			fileInfo.Type = "file"
		}

		fileList = append(fileList, fileInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"list": fileList,
	})
}

// 上传文件到SFTP
func UploadSftpFile(c *gin.Context) {
	hostID := c.Param("id")
	path := c.PostForm("path")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传的文件"})
		return
	}

	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径不能为空"})
		return
	}

	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH连接
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP客户端创建失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("打开上传文件失败: %v", err)})
		return
	}
	defer src.Close()

	// 创建目标文件
	remotePath := filepath.Join(path, file.Filename)
	dst, err := sftpClient.Create(remotePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建远程文件失败: %v", err)})
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传文件失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "上传成功"})
}

// 从SFTP下载文件
func DownloadSftpFile(c *gin.Context) {
	hostID := c.Param("id")
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径不能为空"})
		return
	}

	// 获取主机信息
	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接SSH
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer sshClient.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP连接失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 打开远程文件
	srcFile, err := sftpClient.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("打开文件失败: %v", err)})
		return
	}
	defer srcFile.Close()

	// 获取文件信息
	fileInfo, err := srcFile.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取文件信息失败: %v", err)})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 发送文件内容
	_, err = io.Copy(c.Writer, srcFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("下载文件失败: %v", err)})
		return
	}
}

// 删除SFTP文件
func DeleteSftpFile(c *gin.Context) {
	hostID := c.Param("id")
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径不能为空"})
		return
	}

	// 获取主机信息
	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接SSH
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer sshClient.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP连接失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 获取文件信息
	fileInfo, err := sftpClient.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取文件信息失败: %v", err)})
		return
	}

	// 删除文件或目录
	if fileInfo.IsDir() {
		err = sftpClient.RemoveDirectory(filePath)
	} else {
		err = sftpClient.Remove(filePath)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// 重命名SFTP文件
func RenameSftpFile(c *gin.Context) {
	hostID := c.Param("id")
	oldPath := c.Query("oldPath")
	newPath := c.Query("newPath")
	if oldPath == "" || newPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径不能为空"})
		return
	}

	// 获取主机信息
	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接SSH
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer sshClient.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP连接失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 重命名文件
	err = sftpClient.Rename(oldPath, newPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("重命名失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "重命名成功",
	})
}

// 下载SFTP目录（压缩）
func DownloadSftpDir(c *gin.Context) {
	hostID := c.Param("id")
	dirPath := c.Query("path")
	if dirPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目录路径不能为空"})
		return
	}

	// 获取主机信息
	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接SSH
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer sshClient.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP连接失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "sftp-download-*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建临时目录失败: %v", err)})
		return
	}
	defer os.RemoveAll(tempDir)

	// 下载目录内容到临时目录
	err = downloadDir(sftpClient, dirPath, tempDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("下载目录失败: %v", err)})
		return
	}

	// 创建ZIP文件
	zipPath := filepath.Join(tempDir, "download.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建ZIP文件失败: %v", err)})
		return
	}
	defer zipFile.Close()

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 将目录内容添加到ZIP文件
	err = addDirToZip(zipWriter, tempDir, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建ZIP文件失败: %v", err)})
		return
	}

	// 关闭ZIP写入器
	err = zipWriter.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("关闭ZIP文件失败: %v", err)})
		return
	}

	// 读取ZIP文件内容
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("读取ZIP文件失败: %v", err)})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(dirPath)))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Length", fmt.Sprintf("%d", len(zipData)))

	// 发送ZIP文件内容
	c.Data(http.StatusOK, "application/zip", zipData)
}

// 下载目录内容到本地临时目录
func downloadDir(sftpClient *sftp.Client, remotePath, localPath string) error {
	// 获取远程目录内容
	files, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return err
	}

	// 遍历目录内容
	for _, file := range files {
		remoteFilePath := filepath.Join(remotePath, file.Name())
		localFilePath := filepath.Join(localPath, file.Name())

		if file.IsDir() {
			// 创建本地目录
			err = os.MkdirAll(localFilePath, 0755)
			if err != nil {
				return err
			}
			// 递归下载子目录
			err = downloadDir(sftpClient, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		} else {
			// 下载文件
			remoteFile, err := sftpClient.Open(remoteFilePath)
			if err != nil {
				return err
			}
			defer remoteFile.Close()

			localFile, err := os.Create(localFilePath)
			if err != nil {
				return err
			}
			defer localFile.Close()

			_, err = io.Copy(localFile, remoteFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// addDirToZip 将目录内容添加到zip文件
func addDirToZip(zipWriter *zip.Writer, basePath string, prefix string) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(basePath, file.Name())
		zipPath := filepath.Join(prefix, file.Name())

		if file.IsDir() {
			// 递归处理子目录
			if err := addDirToZip(zipWriter, filePath, zipPath); err != nil {
				return err
			}
		} else {
			// 添加文件到zip
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			writer, err := zipWriter.Create(zipPath)
			if err != nil {
				return err
			}

			if _, err := writer.Write(fileContent); err != nil {
				return err
			}
		}
	}
	return nil
}

// CompressSftpDir 压缩SFTP目录
func CompressSftpDir(c *gin.Context) {
	hostID := c.Param("id")
	dirPath := c.Query("path")

	if dirPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目录路径不能为空"})
		return
	}

	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH连接
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP客户端创建失败: %v", err)})
		return
	}
	defer sftpClient.Close()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "sftp_compress_*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建临时目录失败: %v", err)})
		return
	}
	defer os.RemoveAll(tempDir)

	// 下载目录内容到临时目录
	if err := downloadDir(sftpClient, dirPath, tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("下载目录失败: %v", err)})
		return
	}

	// 创建zip文件
	zipPath := filepath.Join(tempDir, filepath.Base(dirPath)+".zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建zip文件失败: %v", err)})
		return
	}
	defer zipFile.Close()

	// 创建zip写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 将目录内容添加到zip文件
	if err := addDirToZip(zipWriter, tempDir, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建zip文件失败: %v", err)})
		return
	}

	// 关闭zip写入器
	if err := zipWriter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("关闭zip文件失败: %v", err)})
		return
	}

	// 将zip文件上传到SFTP服务器
	zipFileName := filepath.Base(dirPath) + ".zip"
	remoteZipPath := filepath.Join(filepath.Dir(dirPath), zipFileName)
	remoteFile, err := sftpClient.Create(remoteZipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建远程zip文件失败: %v", err)})
		return
	}
	defer remoteFile.Close()

	// 读取本地zip文件并上传
	localZipFile, err := os.Open(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("打开本地zip文件失败: %v", err)})
		return
	}
	defer localZipFile.Close()

	if _, err := io.Copy(remoteFile, localZipFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传zip文件失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "压缩成功"})
}

// WebShell 处理WebShell连接
func WebShell(c *gin.Context) {
	hostID := c.Param("id")

	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主机不存在"})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()

	// 创建SSH连接
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		log.Printf("SSH连接失败: %v", err)
		return
	}
	defer client.Close()

	// 创建新的会话
	session, err := client.NewSession()
	if err != nil {
		log.Printf("创建SSH会话失败: %v", err)
		return
	}
	defer session.Close()

	// 设置伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// 获取终端尺寸
	termWidth := 200
	termHeight := 40

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		log.Printf("请求伪终端失败: %v", err)
		return
	}

	// 获取标准输入输出
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Printf("获取标准输入失败: %v", err)
		return
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("获取标准输出失败: %v", err)
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Printf("获取标准错误失败: %v", err)
		return
	}

	// 启动shell
	if err := session.Shell(); err != nil {
		log.Printf("启动shell失败: %v", err)
		return
	}

	// 处理WebSocket消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("读取WebSocket消息失败: %v", err)
				return
			}

			// 检查是否是调整终端大小的消息
			if len(message) > 0 && message[0] == '\x1b' {
				// 解析终端大小调整消息
				parts := strings.Split(string(message), ";")
				if len(parts) == 3 && strings.HasPrefix(parts[0], "\x1b[8") {
					height, _ := strconv.Atoi(parts[1])
					width, _ := strconv.Atoi(parts[2])
					if height > 0 && width > 0 {
						termHeight = height
						termWidth = width
						log.Printf("调整终端大小: %dx%d", width, height)
						// 调整终端大小
						if err := session.WindowChange(termHeight, termWidth); err != nil {
							log.Printf("调整终端大小失败: %v", err)
						}
						continue
					}
				}
			}

			if _, err := stdin.Write(message); err != nil {
				log.Printf("写入标准输入失败: %v", err)
				return
			}
		}
	}()

	// 处理标准输出
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if err != nil {
				log.Printf("读取标准输出失败: %v", err)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				log.Printf("发送WebSocket消息失败: %v", err)
				return
			}
		}
	}()

	// 处理标准错误
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stderr.Read(buffer)
			if err != nil {
				log.Printf("读取标准错误失败: %v", err)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				log.Printf("发送WebSocket消息失败: %v", err)
				return
			}
		}
	}()

	// 等待会话结束
	session.Wait()
}

// UploadFile 处理文件上传
func UploadFile(c *gin.Context) {
	hostID := c.Param("id")

	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主机不存在"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传的文件"})
		return
	}

	// 创建SSH连接
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SSH连接失败"})
		return
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建SFTP客户端失败"})
		return
	}
	defer sftpClient.Close()

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开上传文件失败"})
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := sftpClient.Create(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目标文件失败"})
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "复制文件内容失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
}

// DownloadFile 处理文件下载
func DownloadFile(c *gin.Context) {
	hostID := c.Param("id")
	filename := c.Query("file")

	var host models.Host
	if err := global.DB.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主机不存在"})
		return
	}

	// 创建SSH连接
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), sshConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SSH连接失败"})
		return
	}
	defer client.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建SFTP客户端失败"})
		return
	}
	defer sftpClient.Close()

	// 打开源文件
	src, err := sftpClient.Open(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开源文件失败"})
		return
	}
	defer src.Close()

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")

	// 复制文件内容到响应
	if _, err := io.Copy(c.Writer, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "复制文件内容失败"})
		return
	}
}
