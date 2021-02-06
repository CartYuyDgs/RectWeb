package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"time"
)

const addr = "172.16.96.56:32731"
const user = "bonc"
const passwd = "bonc@123"
const port = 32731

var timChan = make(chan int, 2)

const (
	TIME_START = 1
	//TIME_END = 2
)

func SftpConnect() (*sftp.Client, error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passwd))
	config := ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return nil, err
	}

	c, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	fmt.Println("create sftp client ok.")
	return c, nil
}

func main() {
	path := "D:\\bonc\\sftp\\b.png"
	c, err := SftpConnect()
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	go func() {
		for {
			select {
			case env := <-timChan:
				if env == TIME_START {
					fmt.Println("------------------------3")
					time_start := time.Now()
					sftpPath := "upload/photo.jpg"
					fp, err := c.Open(sftpPath)
					if err != nil {
						log.Fatal(err)
					}
					fp1, _ := os.Create(path)
					if err != nil {
						log.Fatal(err)
					}
					io.Copy(fp1, fp)
					fp1.Close()
					fp.Close()
					fmt.Printf("消耗时间为：%v \n", time_start.Sub(time.Now()))
				}
			default:
				break
			}
		}

	}()

	TimeFunc(timChan)

	//select {}
}

func TimeFunc(timChan chan int) {
	for {
		timChan <- TIME_START
		time.Sleep(50 * time.Millisecond)
		fmt.Println("------------------------2")
	}
}

//func ftpConnect() {
//	c, err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	err = c.Login(user, passwd)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	defer c.Quit()
//	defer c.Logout()
//
//	re, err := c.CurrentDir()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	println(re)
//
//}

//wd,_ := c.Getwd()
//fmt.Println(wd)
//defer c.Close()
//time_start := time.Now()
////fp, err := c.Open("upload/")
////if err != nil {
////	log.Fatal(err)
////}
////
////defer fp.Close()
////fileInfo, err := c.ReadDir("upload/")
////stat_t := fileInfo[0].Sys().(*sftp.FileStat)
////fmt.Println(stat_t.Atime)
//sftpPath := "upload/渤.png"
////fmt.Println(fileInfo[0].Name())
////fmt.Println(sftpPath)
////fmt.Println(fileInfo[0].ModTime())
//fp, err := c.Open(sftpPath)
//if err != nil {
//log.Fatal(err)
//}
//fp1 , _ := os.Create(path)
//
//io.Copy(fp1, fp)
//defer fp1.Close()
//
//fmt.Printf("消耗时间为：%v",time_start.Sub(time.Now()))
//defer fp.Close()
