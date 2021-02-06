package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"path/filepath"
	"time"
)

var fileSent []string
var LocalDir string = "/home/nvidia/boncAI/picture/"
var LocalAddr string = "172.16.96.57:9090"
var Topic string = "test1"

type KafkaData struct {
	Time        string `json:"time"`
	Msg         string `json:"imgtype"`
	Photobase64 string `json:"data"`
}

func init() {
	flag.StringVar(&LocalAddr, "local_address", LocalAddr, "input your localdir")
	flag.StringVar(&LocalDir, "localDir", LocalDir, "input your local address")
	flag.StringVar(&Topic, "topic", Topic, "input your kafka topic")
	flag.Parse()

	fmt.Printf("--*--*--*--&-imagelocal: %s, localaddress: %s, topic: %s -&--*--*--*--\n", LocalDir, LocalAddr, Topic)
}

func main() {
	var adds []string
	adds = append(adds, LocalAddr)

	log.Println("Info: connect kafka and create product start.....")
	product, err := createProducer(adds)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer product.Close()
	log.Println("Info: connect kafka and create product success.....")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := glice(product)
			log.Printf("Info: glice result: %v", err)
		}
	}
}

func glice(product sarama.AsyncProducer) error {
	var noOkImage = "out.bmp"
	var isMove = false
	dir, err := os.Open(LocalDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	filenames, _ := dir.Readdirnames(3)
	dir.Close()
	log.Println("Info: glice start......")
	for _, filename := range filenames {

		log.Printf("Info: find filename: %s", filename)
		if noOkImage[:2] == filename[:2] {
			log.Printf("Info: filename %s No ok ", filename)
			continue
		}

		for _, file := range fileSent {
			log.Printf("Info: find filename and lastfile: %s, %s", filename, file)
			if filename == file {
				if err := os.Remove(filepath.Join(LocalDir, filename)); err != nil {
					log.Printf("Info: remove filename: %s failed", filename)
					isMove = true
				}
			}
		}

		if isMove {
			isMove = false
			continue
		}

		err := send2kafka(product, filename)
		if err != nil {
			log.Printf("Info: send err:", err)
			return err
		}
		fileaddr := filepath.Join(LocalDir, filename)
		if err = os.Remove(fileaddr); err != nil {
			fileSent = append(fileSent, filename)
		}

	}

	return nil
}

func send2kafka(product sarama.AsyncProducer, photo string) error {

	log.Printf("Info: Image %s are send to kafka.", photo)
	data, err := configData(photo)
	if err != nil {
		fmt.Println("Info: create config data err: ", err)
		return nil
	}

	//发送data
	producerSend(product, Topic, *data)

	return nil
}

func createProducer(address []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后，再返回响应
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.MaxMessageBytes = 8000000
	// 随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应，只有上面的RequireAcks设置不是NoResponse，这里才有用。
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 设置读写超时时间为2秒，默认为10秒
	config.Producer.Timeout = 2 * time.Second
	// 尝试发送消息最大次数
	config.Producer.Retry.Max = 3

	//fmt.Println("start to make a producer")
	// 使用配置，新建一个异步生产者
	producer, e := sarama.NewAsyncProducer(address, config)
	if e != nil {
		fmt.Println("fail to make a producer, error:", e)
		return nil, e
	}

	return producer, nil
}

func producerSend(producer sarama.AsyncProducer, topic string, kafkadata KafkaData) {

	// 循环判断哪个通道发送过来数据。
	fmt.Println("Info: start goroutine to get response")
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case suc := <-p.Successes():
				if suc != nil {
					fmt.Printf("succeed, offset=%d, timestamp=%s, partitions=%d\n", suc.Offset, suc.Timestamp.String(), suc.Partition)
					//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				}
			case fail := <-p.Errors():
				if fail != nil {
					fmt.Printf("error:  %v\n", fail.Err)
				}
			}
		}
	}(producer)

	// 发送消息
	values, err := json.Marshal(kafkadata)
	if err != nil {
		log.Fatalln(err)
	}

	strKey := "key"
	srcValue := string(values)

	//time.Sleep(50 * time.Millisecond)

	// 发送的消息对应的主题。
	// 注意：这里的msg必须是新构建的变量。不然，发送过去的消息内容都是一样的，因为批次发送消息的关系。
	msg := &sarama.ProducerMessage{
		Topic: topic,
	}

	msg.Timestamp = time.Now()

	// 设置消息的key
	msg.Key = sarama.StringEncoder(strKey)
	// 设置消息的value，将字符串转化为字节数组
	msg.Value = sarama.ByteEncoder(srcValue)
	//fmt.Println(value)

	// 使用通道发送
	producer.Input() <- msg

}

func CopyImage(baseLocal string, name string, im []byte) {

	timeName := baseLocal + time.Now().Format("15:04:05") + name
	fmt.Println(timeName)
	file, err := os.Create(timeName)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := bufio.NewWriter(file)

	writer.Write(im[:])

	defer file.Close()
}

func configData(photo string) (*KafkaData, error) {

	//baseLocal := "/home/nvidia/boncAI/picture_save/"
	//baseLocal := "/root/yuy/image_save/"
	var info = KafkaData{
		Msg:  "jpg",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	time1 := time.Now()

	filePath := LocalDir + photo
	fmt.Println(filePath)
	Objfile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return &info, errors.New("file not find!")
	}

	//base64
	buffer := make([]byte, 30000000)
	n, _ := Objfile.Read(buffer)
	fmt.Println("photo len:", n)

	//CopyImage(baseLocal,photo,buffer[:n])

	encodestring := base64.StdEncoding.EncodeToString(buffer[:n])

	if len(encodestring) > 0 {
		//info.Photobase64 = append(info.Photobase64, encodestring)
		info.Photobase64 = encodestring
	} else {
		//fmt.Println("lenthon: ",len(encodestring))
		return &info, errors.New("file error!")
	}
	fmt.Println(encodestring[:20])
	fmt.Println("open file time:", time.Now().Sub(time1))
	//读取文件之后删除
	//Objfile.Close()
	//err = os.Remove(photo)
	//if err != nil {
	//	fmt.Println(err)
	//}

	return &info, nil
}
