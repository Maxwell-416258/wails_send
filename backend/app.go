package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	conn net.Conn
	listener net.Listener
	pendingConn net.Conn // 接收用
}

func NewSender() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.StartServer()
}

func (a *App) Shutdown(ctx context.Context) {
	// 在此处理发送端关闭时的清理操作
}

func (s *App) GetInnerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "无法获取IP"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ip := ipNet.IP.To4()
			if ip == nil {
				continue // 忽略 IPv6
			}

			// 排除 169.254.x.x（APIPA）
			if ip[0] == 169 && ip[1] == 254 {
				continue
			}

			// 只保留私有地址段
			if (ip[0] == 10) ||
				(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) ||
				(ip[0] == 192 && ip[1] == 168) {
				return ip.String()
			}
		}
	}
	return "未找到有效内网IP"
}
func (a *App) CheckIP(ip string) bool {
    // 尝试连接到目标 IP 的 8888 端口
    conn, err := net.DialTimeout("tcp", ip+":8888", 2*time.Second)
    if err != nil {
        return false // 如果连接失败，表示目标 IP 不可达
    }
    defer conn.Close()
    return true // 如果连接成功，表示目标 IP 可达
}

func (a *App) ChooseFile() (string) {
	// 使用 Wails 的文件选择对话框
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
	})
	if err != nil {
		return ""
	}
	fmt.Println("选择的文件路径:", filePath)
	return filePath
}


func (a *App)SendFileMeta(filepath string,ip string) error{
	var port string = "8888"
	file,err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()
	fileMetaData := FileMetadata{
		FileName: fileName,
		FileSize: fileSize,
	}
	// 这里可以将 fileMetaData 发送到指定的 IP 和端口
	// 例如使用 TCP 连接发送数据
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	a.conn = conn
	// 发送文件元数据	
	jsonData, err := json.Marshal(fileMetaData)
	if err != nil {
		return err
	}
	_, err = conn.Write(jsonData)
	if err != nil {
		return err
	}

	go a.waitForConfirmation(filepath,fileSize)
	return nil
}

func (a *App) waitForConfirmation(filePath string,fileSize int64) {
	buf :=make([]byte,1024)
	n,err:=a.conn.Read(buf)
	if err != nil {
		println("Error reading from connection:", err.Error())
		return
	}
	if string(buf[:n]) == "ACCEPT" {
		if err := a.SendFile(filePath,fileSize);err!=nil{
			runtime.EventsEmit(a.ctx,"fileMetaSendError",err.Error())
		}else{
			runtime.EventsEmit(a.ctx,"fileMetaSendSuccess","文件发送成功")
		}
	}else{
		runtime.EventsEmit(a.ctx,"fileMetaSendError","对方拒绝接收文件")
	}
}


func (a *App) StartServer(){
	ln,err:=net.Listen("tcp",":8888")
	if err!=nil{
		println("Error starting server:", err.Error())
		return
	}
	a.listener = ln
	go func() {
		for {
			conn,err := ln.Accept()
			if err != nil {
				println("Error accepting connection:", err.Error())
				continue
			}
			go a.handleConnection(conn)
		}
	}()
}


func (a *App) handleConnection(conn net.Conn) {
	buf := make([]byte,1024)
	n,err:=conn.Read(buf)
	if err != nil {
		conn.Close()
		return
	}
	var fileInfo FileMetadata
	if err:= json.Unmarshal(buf[:n],&fileInfo);err!=nil{
		conn.Close()
		println("Error unmarshalling JSON:", err.Error())
		return
	}
	runtime.EventsEmit(a.ctx,"fileInfo",fileInfo)
	a.pendingConn = conn
}


func (a *App) AcceptFileRequest(fileName string,fileSize int64) error{
	if a.pendingConn == nil{
		return fmt.Errorf("没有可用的连接")
	}
	if _,err := a.pendingConn.Write([]byte("ACCEPT"));err!=nil{
		return err
	}
	//开始接收文件内容
	go func(fileName string,fileSize int64) {
		saveFilePath,err:= a.SelectSaveDirectory()
		if err != nil{
			fmt.Println("选择保存目录失败:", err)
			a.pendingConn.Close()
            a.pendingConn = nil
			return
		}
		saveFilePath = saveFilePath + string(os.PathSeparator) + fileName
		log.Println("选择的保存目录:", saveFilePath)
		err = a.SaveFile(saveFilePath,fileSize,a.pendingConn)
		if err != nil{
			fmt.Println("保存文件失败:", err)
			a.pendingConn.Close()
            a.pendingConn = nil
			return
		}
	}(fileName,fileSize)
	return nil
}

func (a *App) SendFile(filePath string, fileSize int64) error {
    var sent int64
    if a.conn == nil {
        return fmt.Errorf("没有可用的连接")
    }
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("无法打开文件: %v", err)
    }
    defer file.Close()

	defer func() {
		if a.conn != nil {
			a.conn.Close()
			log.Println("发送端连接已关闭")
		}
	}()

    buf := make([]byte, 32*1024)
    for {
        n, err := file.Read(buf)
        if n > 0 {
            _, writeErr := a.conn.Write(buf[:n])
            if writeErr != nil {
                return fmt.Errorf("发送文件时出错: %v", writeErr)
            }
            sent += int64(n)
			progress := float64(sent) / float64(fileSize) * 100
            runtime.EventsEmit(a.ctx, "fileSendProgress", fmt.Sprintf("%.2f", progress))
			log.Println("发送进度:", fmt.Sprintf("%.2f%%", progress))
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("读取文件时出错: %v", err)
        }
    }
    
    runtime.EventsEmit(a.ctx, "fileSendSuccess", "文件发送成功")
    return nil
}


func (a *App) SelectSaveDirectory() (string,error){
	return runtime.OpenDirectoryDialog(a.ctx,runtime.OpenDialogOptions{
		Title: "选择保存目录",
	})
}

func (a *App) SaveFile(filePath string,fileSize int64,conn net.Conn) error {
	var received int64

	file,err := os.Create(filePath)
	if err !=nil{
		return err
	}
	defer file.Close()
	// 从连接中接收数据
	buf := make([]byte,32*1024)
	for {
		n,err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				break
			}
			return fmt.Errorf("接收文件时出错: %v", err)
		}
		_, err = file.Write(buf[:n])
		received += int64(n)
		progress := fmt.Sprintf("%.2f", float64(received)/float64(fileSize)*100)
		runtime.EventsEmit(a.ctx,"fileReceiveProgress",progress)
		if err != nil {
			return fmt.Errorf("写入文件时出错: %v", err)
		}
	}
	runtime.EventsEmit(a.ctx,"fileSaveSuccess","文件保存成功")
	return nil
}