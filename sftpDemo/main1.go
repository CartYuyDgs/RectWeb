package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

//const ftpDir = "D:\\bonc\\sftp"
//const tomcatDir = "D:\\bonc\\仪表盘读数\\通威\\通威\\apache-tomcat-7.0.106\\webapps\\pipeline_pressure_front"
const (
	Photo  = "photo.jpg"
	Photo1 = "photo1.jpg"
	Photo2 = "photo2.jpg"
	Photo3 = "photo3.jpg"
	Photo4 = "photo4.jpg"
	Photo5 = "photo5.jpg"
	Photo6 = "photo6.jpg"
	Photo7 = "photo7.jpg"
)

func main1() {
	//检查目录变化
	var Photoaddr string
	Photoaddr = filepath.Join(ftpDir, Photo)
	fmt.Println(Photoaddr)

	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer watch.Close()

	err = watch.Add(ftpDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for {
			select {
			case env, ok := <-watch.Events:
				{
					if !ok {
						log.Fatal("aa", env)
						return
					}

					if env.Op == fsnotify.Rename {
						fmt.Println(env.Name, env.Op, env.String())
					}

					if env.Op == fsnotify.Create {
						fmt.Println(env.Name, env.Op, env.String())
						Photoaddr = env.Name

						oldaddr := filepath.Join(tomcatDir, "1.jpg")
						fmt.Println(oldaddr)
						os.Rename(Photoaddr, oldaddr)
					}
					//fmt.Println(env.Name, env.Op, env.String())

					//CopyFile(oldaddr, Photoaddr)
					////CopyShellFile(oldaddr,Photoaddr)
					//log.Println("bb",env)

				}
			case env, ok := <-watch.Errors:
				{
					if !ok {
						log.Fatal(ok)
						return
					}

					log.Println("cc", env)
				}

			}
		}
	}()

	select {}
	//将特定的文件同步到指定文件夹
}

func CopyFile(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func CopyShellFile(dst, src string) error {
	c := fmt.Sprintf("COPY /Y %s %s ", src, dst)
	cmd := exec.Command(c)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("cmd finished")
	return err
}

const ftpDir = "D:\\bonc\\sftp"
const tomcatDir = "D:\\bonc\\仪表盘读数\\通威\\通威\\apache-tomcat-7.0.106\\webapps\\pipeline_pressure_front"
const photoPath = "D:\\bonc\\sftp\\abc.jpg"

type ImageInfo struct {
	Name       string
	FtpAddr    string
	TomcatDir  string
	CreateTime string
	Infochan   chan ChanInfo
}

type ChanInfo struct {
	Name       string
	CreateTime int64
}

func main() {

	image := ImageInfo{
		Name:       "1.jpg",
		FtpAddr:    "D:\\bonc\\sftp\\abc.jpg",
		TomcatDir:  "D:\\bonc\\仪表盘读数\\通威\\通威\\apache-tomcat-7.0.106\\webapps\\pipeline_pressure_front\\1.jpg",
		CreateTime: "",
		Infochan:   nil,
	}

	image.Infochan = make(chan ChanInfo, 5)

	go image.removePhoto()
	image.findPhoto()

}

func (image *ImageInfo) removePhoto() {
	var lasttime int64
	lasttime = 0
	for {
		select {
		case info := <-image.Infochan:
			{
				if info.CreateTime <= lasttime {
					continue
				}
				lasttime = info.CreateTime
				log.Fatal(image.FtpAddr, image.TomcatDir, info.Name, info.CreateTime)
				//fmt.Println(image.FtpAddr, image.TomcatDir,info.Name,info.CreateTime)
				os.Rename("D:\\bonc\\sftp\\abc.jpg", "D:\\bonc\\仪表盘读数\\通威\\通威\\apache-tomcat-7.0.106\\webapps\\pipeline_pressure_front\\1.jpg")
			}
		default:
			continue
		}
	}
}

func (image *ImageInfo) findPhoto() {

	osType := runtime.GOOS
	for {
		fileInfo, err := os.Stat(photoPath)
		if err != nil {
			continue
		}

		fmt.Println(osType, fileInfo)
		if osType == "windows" {
			wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
			tNanSeconds := wFileSys.CreationTime.Nanoseconds() /// 返回的是纳秒
			tSec := tNanSeconds / 1e9                          ///秒
			fmt.Println(tSec)
			var charinfo = ChanInfo{
				Name:       image.Name,
				CreateTime: tSec,
			}
			image.Infochan <- charinfo
		}
	}
}
