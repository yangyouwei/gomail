package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"net"
	"net/smtp"
	"os"
)

var (
	host string
	port int
	email  string
	pwd string

)

var Usage = func() {
	fmt.Println("Usage: COMMAND args1 args2 args3")
	fmt.Println("args1 is email address")
	fmt.Println("args2 is the mesages's title")
	fmt.Println("args3 is messages's content")
}

func init()  {
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
		return
	}

	host, err = cfg.GetValue("main","smtp_server")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "smtp_server", err)
	}

	port, err = cfg.Int("main","port")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "port", err)
	}

	email, err = cfg.GetValue("main","email")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "email", err)
	}

	pwd, err = cfg.GetValue("main","passworld")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "passworld", err)
	}
}

func main(){

	//fmt.Println(host)
	//fmt.Println(port)
	//fmt.Println(email)
	//fmt.Println(pwd)

	args := os.Args

	if args == nil || len(args) < 2 {
		Usage() //如果用户没有输入,或参数个数不够,则调用该函数提示用户
		return
	}

	toEmail := &args[1]  // 目标地址 ，这个是程序运行是的参数。

	header   :=  make(map[string]string)

	header["From"] = "test"+"<" +email+">"
	header["To"] = *toEmail
	header["Subject"] = args[2]
	header["Content-Type"] = "text/html;chartset=UTF-8"

	body  := args[3]

	message := ""

	for k,v :=range header{
		message  += fmt.Sprintf("%s:%s\r\n",k,v)
	}

	message +="\r\n"+body


	auth :=smtp.PlainAuth(
		"",
		email,
		pwd,
		host,
	)

	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{*toEmail},
		[]byte(message),
	)

	if err  !=nil{
		panic(err)
	}

}


//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Panicln("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}