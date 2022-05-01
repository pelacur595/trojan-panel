package service

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"time"
	"trojan/module/constant"
	"trojan/module/dto"
)

func SendEmail(sendEmailDto *dto.SendEmailDto) error {
	name := constant.SystemName
	systemVo, err := SelectSystemByName(&name)
	if err != nil {
		return err
	}
	ch := make(chan *gomail.Message)
	d := gomail.NewDialer(systemVo.EmailHost, int(systemVo.EmailPort), systemVo.EmailUsername, systemVo.EmailPassword)

	go func() {
		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						logrus.Errorf("mail dail err: %v\n", err)
						return
					}
					open = true
				}

				if err := gomail.Send(s, m); err != nil {
					logrus.Errorf("mail send err: %v\n", err)
					return
				}

			// 30秒没有发送消息则关闭SMTP server连接
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						logrus.Errorf("mail close err: %v\n", err)
						return
					}
					open = false
				}
			}
		}
	}()

	// 发送消息
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(systemVo.EmailUsername, sendEmailDto.FromEmailName)},
		"To":      sendEmailDto.ToEmails,
		"Subject": {sendEmailDto.Subject},
	})
	m.SetBody("text/html", sendEmailDto.Content)
	// 附件选项
	// m.Attach("/home/Alex/lolcat.jpg")

	ch <- m
	// 关闭channel来停止守护进程
	defer close(ch)

	return nil
}
