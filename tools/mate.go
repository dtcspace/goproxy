package main

import (
	"bufio"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	//app := kingpin.New("mate", "proxy mate")
	//auth := app.Command("auth", "mate for proxy authenticate")
	//reset := auth.Command("reset", "reset proxy authenticate")

	//authFile := auth.Flag("auth-file", "proxy authenticate file,\"username:password\" each line in file").Short('F').String()

	//if _, err := os.Stat(authFile); os.IsNotExist(err) {
	//	return err
	//}

	users := []string{"ronghui.wang"}

	rand.Seed(time.Now().UnixNano())
	authUsers, err := listAuthUser("auth")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range users {
		authUsers[v] = randStr(16)
	}

	err = addAuthUser("auth", authUsers)
	if err != nil {
		fmt.Println(err)
		return
	}

	message := `
	Your Proxy Information : <br>
		username -> %s <br>
		password -> %s <br>
	`
	for _, v := range users {
		sendMail(v, fmt.Sprintf(message, v, authUsers[v]))
	}

}

func listAuthUser(authFile string) (map[string]string, error) {
	fd, err := os.Open(authFile)
	defer func() {
		silenceErr(fd.Close())
	}()
	if err != nil {
		return nil, err
	}

	users := make(map[string]string)
	buff := bufio.NewReader(fd)

	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		ud := strings.Split(string(data), ":")
		if len(ud) != 2 {
			fmt.Printf("unidentifiable format:%s\n", string(data))
			continue
		}
		users[ud[0]] = ud[1]
	}
	return users, err
}

func addAuthUser(authFile string, authUsers map[string]string) error {

	fd, err := os.OpenFile(authFile, os.O_RDWR, 0666)
	defer func() {
		silenceErr(fd.Close())
	}()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fd)

	for k, v := range authUsers {
		_, err = w.WriteString(k + ":" + v + "\n")
		if err != nil {
			return err
		}
	}

	silenceErr(w.Flush())
	silenceErr(fd.Sync())
	return nil
}

func sendMail(user string, message string) {
	host := "smtp.feishu.cn"
	port := 587
	userName := "notifications@opendhc.io"
	password := "XorxhyxtQhI1RUvu"

	m := gomail.NewMessage()
	m.SetHeader("From", "dtc-proxy<"+userName+">")
	m.SetHeader("To", user+"@dhc.com.cn")
	m.SetHeader("Subject", "DTC Proxy Notification")
	m.SetBody("text/html", fmt.Sprintf(message))

	d := gomail.NewDialer(
		host,
		port,
		userName,
		password,
	)
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	silenceErr(d.DialAndSend(m))
	return
}

func silenceErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func randStr(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}
