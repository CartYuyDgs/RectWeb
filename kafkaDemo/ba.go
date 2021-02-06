//package main
//
//import (
//	"bytes"
//	"encoding/base64"
//	"fmt"
//	"github.com/Shopify/sarama"
//	"io/ioutil"
//	"os"
//	"path/filepath"
//	"time"
//	"flag"
//)
//
//
//var (
//	producer sarama.AsyncProducer
//	message  sarama.ProducerMessage
//)
//
//var (
//	localDir  string = "/root/yuy/images/"
//	localAddr string = "172.16.96.57:9090"
//	topic     string = "test1"
//)
//
//// 已经发送成功的文件名称
//var fileSent string
//
//// KafkaMsgData 承载Kafka消息的载体
//type KafkaMsgData struct {
//	Photobase64 []byte
//}
//
//// MarshalJSON 自定义json序列化
//func (d KafkaMsgData) MarshalJSON() ([]byte, error) {
//	buffer := bytes.NewBufferString(`{"time":"`)
//	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05"))
//	buffer.WriteString(`","imgtype":"jpg","data":"`)
//	buffer.Write(d.Photobase64)
//	buffer.WriteString(`"}`)
//	return buffer.Bytes(), nil
//}
//
//func init() {
//	flag.StringVar(&localDir, "local_address", localDir, "input your localdir")
//	flag.StringVar(&localAddr, "localDir", localAddr, "input your local address")
//	flag.StringVar(&topic, "topic", topic, "input your kafka topic")
//	flag.Parse()
//
//	message = sarama.ProducerMessage{
//		Topic: topic,
//		Key:   sarama.StringEncoder("key"),
//	}
//}
//
//func main() {
//	fmt.Println("start..................")
//	//producer = createProducer([]string{localAddr})
//	//defer producer.Close()
//
//	ticker := time.NewTicker(1 * time.Second)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//			fmt.Println("----ok")
//			glice()
//		}
//	}
//}
//
//func createProducer(address []string) sarama.AsyncProducer {
//	config := sarama.NewConfig()
//	// 等待服务器所有副本都保存成功后，再返回响应
//	config.Producer.RequiredAcks = sarama.WaitForAll
//	// 随机向partition发送消息
//	config.Producer.Partitioner = sarama.NewRandomPartitioner
//	// 是否等待成功和失败后的响应，只有上面的RequireAcks设置不是NoResponse，这里才有用。
//	config.Producer.Return.Successes = true
//	config.Producer.Return.Errors = true
//	// 设置读写超时时间为2秒，默认为10秒
//	config.Producer.Timeout = 2 * time.Second
//	// 尝试发送消息最大次数
//	config.Producer.Retry.Max = 3
//
//	// 使用配置，新建一个异步生产者
//	producer, err := sarama.NewAsyncProducer(address, config)
//	if err != nil {
//		panic(err)
//	}
//
//	go func() {
//		for {
//			select {
//			case msg := <-producer.Successes():
//				if msg != nil {
//					fmt.Printf("succeed, offset=%d, timestamp=%s, partitions=%d\n", msg.Offset, msg.Timestamp.String(), msg.Partition)
//				}
//			case perr := <-producer.Errors():
//				if perr != nil {
//					fmt.Printf("error=%v\n", perr.Err)
//				}
//			}
//		}
//	}()
//
//	return producer
//}
//
//func glice() error {
//	dir, err := os.Open(localDir)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	filenames, _ := dir.Readdirnames(2)
//	dir.Close()
//	fmt.Println("glice.........")
//	for _, filename := range filenames {
//		if filename <= fileSent {
//			os.Remove(filepath.Join(localDir, filename))
//		} else {
//			send2kafka(filename)
//			os.Remove(filepath.Join(localDir, filename))
//		}
//	}
//
//	return nil
//}
//
//func send2kafka(photoWO string) error {
//	content, err := ioutil.ReadFile(filepath.Join(localDir, photoWO))
//	if err != nil {
//		return err
//	}
//	fmt.Println("send .........",photoWO)
//	var data = KafkaMsgData{
//		Photobase64: make([]byte, base64.StdEncoding.EncodedLen(len(content))),
//	}
//	base64.StdEncoding.Encode(data.Photobase64, content)
//
//	encoded, _ := data.MarshalJSON()
//	encoded = encoded
//	//
//	//msg := message
//	//msg.Value = sarama.ByteEncoder(encoded)
//	//producer.Input() <- &msg
//	//
//	//fileSent = photoWO
//	return nil
//}
//
