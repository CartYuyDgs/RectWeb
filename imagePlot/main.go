package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Point_areas struct {
	point map[string][]int
}

type date struct {
	Token       string   `json:"token"`
	Imgtype     string   `json:"imgtype"`
	Photobase64 string   `json:"data"`
	Point_area  []string `json:"point_area"`
	Time        string   `json:"time"`
}

type resultInfo struct {
	rename  string
	code    string
	voltage string
	time    string
}

type resultList struct {
	date []resultInfo
}

const url = "http://172.16.96.57:10321/pressure_reading"
const token = "bonccvlab"
const imgtype = ".jpg"

//var point_area = make(map[string][]int)
//var points  = []int{882,187,1179,187,1179,481,882,481, 6}
var point_area = fmt.Sprintf("{'%s': [882,187,1179,187,1179,481,882,481, 6]}", "voltage")

const image_save = "D:\\code\\src\\github.com\\CartYuyDgs\\RectWeb\\imagePlot\\picture_save\\"

func main() {

	infoList := resultList{}

	names, err := readImage(image_save)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, name := range names {
		CreateInfo(name, &infoList)
		if i == 10 {
			break
		}
	}
	fmt.Println(infoList.date)
	plotInfo(infoList)
}

func readImage(addr string) ([]string, error) {
	namelist := []string{}

	files, err := ioutil.ReadDir(addr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, f := range files {
		//fmt.Println(f.Name())
		namelist = append(namelist, f.Name())
	}

	return namelist, nil
}

func CreateInfo(photo string, list *resultList) {

	image_addr := image_save + photo

	Objfile, err := os.Open(image_addr)
	defer Objfile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//base64
	buffer := make([]byte, 500000)
	n, _ := Objfile.Read(buffer)
	encodestring := base64.StdEncoding.EncodeToString(buffer[:n])

	//fmt.Println(encodestring)

	var info = date{
		Imgtype: imgtype,
		Token:   token,
	}

	info.Point_area = append(info.Point_area, point_area)
	info.Photobase64 = encodestring

	info.Time = time.Now().Format("15:04:05")

	//post请求
	contentType := "application/json"

	res, err := postSend(url, info, contentType)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("result: ", res, reflect.TypeOf(res))
	resinfo, err := getInfoFromRes(res, photo)

	list.date = append(list.date, *resinfo)
}

func postSend(url string, info date, contentType string) (string, error) {
	//client := &http.Client{}
	jsonStr, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(url)
	//fmt.Println(info)
	fmt.Println(contentType)

	resp, err := http.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), err
}

func getInfoFromRes(res string, name string) (resultinfo *resultInfo, err error) {
	//{"success":200,"message":"检测成功","data":[{"voltage":4.75}],"time":"10:54:20"}

	reg := regexp.MustCompile(`"success":(.*?),"message":"检测成功","data":\[\{"voltage":(.*?)\}\],"time":"(.*?)"`)
	re := reg.FindAllStringSubmatch(res, -1)

	result := resultInfo{
		rename:  name,
		code:    re[0][1],
		voltage: re[0][2],
		time:    re[0][3],
	}

	fmt.Println(result)
	return &result, nil
}

func plotInfo(list resultList) {
	points := plotter.XYs{}
	for i := 0; i < len(list.date); i++ {
		num, _ := strconv.ParseFloat(list.date[i].voltage, 64)
		points = append(points, plotter.XY{
			X: float64(i),
			Y: num,
		})
	}

	plt, err := plot.New()

	if err != nil {
		panic(err)
	}
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, 10, float64(len(list.date))

	if err := plotutil.AddLines(plt,
		"line1", points,
	); err != nil {
		panic(err)
	}

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "01-draw-line.png"); err != nil {
		panic(err)
	}

}
