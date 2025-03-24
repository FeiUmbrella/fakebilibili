package email

import (
	"fakebilibili/infrastructure/pkg/global"
	"gopkg.in/gomail.v2"
	"strconv"
)

// SendEmail 发送邮箱验证码
func SendEmail(mailTo []string, subject string, body string) error {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": global.Config.EmailConfig.User,
		"pass": global.Config.EmailConfig.Password,
		"host": global.Config.EmailConfig.Host,
		"port": global.Config.EmailConfig.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"])
	m := gomail.NewMessage()
	m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...) // 可以发送多个邮箱
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	if err != nil {
		global.Logger.Errorf("发送邮箱验证码至 %s 失败，内容：%s，错误原因：%s", mailTo, body, err.Error())
	}
	return err
}
