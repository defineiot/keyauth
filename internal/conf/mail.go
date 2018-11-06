package conf

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// MailConf send 验证码
type MailConf struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Email    string `toml:"email"`
	Password string `toml:"password"`

	header map[string]string
	client *smtp.Client
	auth   smtp.Auth
}

// Init 初始化
func (m *MailConf) Init() error {
	auth := smtp.PlainAuth("", m.Email, m.Password, m.Host)
	m.auth = auth

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return err
	}

	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}

	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	m.client = client

	return nil
}

// Send send mail
func (m *MailConf) Send(toEmail, userName string, code int) error {
	msg := m.prepareMessage(toEmail, userName, code)

	return m.sendMailUsingTLS([]string{toEmail}, []byte(msg))
}

// SendMailUsingTLS 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
func (m *MailConf) sendMailUsingTLS(to []string, msg []byte) (err error) {
	if m.auth != nil {
		if ok, _ := m.client.Extension("AUTH"); ok {
			if err = m.client.Auth(m.auth); err != nil {
				return fmt.Errorf("Error during AUTH, %s", err)
			}
		}
	}

	if err := m.client.Mail(m.Email); err != nil {
		return err
	}

	for _, addr := range to {
		if err = m.client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := m.client.Data()
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
	return m.client.Quit()
}

func (m *MailConf) prepareMessage(toEmail, userName string, code int) string {
	// prepare mail's header
	header := make(map[string]string)

	header["From"] = "西牛云开发者平台" + "<" + m.Email + ">"
	header["To"] = toEmail
	header["Subject"] = "西牛云验证码"
	header["Content-Type"] = "text/html;chartset=UTF-8"
	m.header = header

	// prepare mail's body
	body := fmt.Sprintf(`亲爱的 <b>开发者 %s: </b>: <br><br>
			您请求的验证码是： <b>%d</b> ，请您尽快完成验证, 感谢您的支持。<br><br>
			此致<br><br>
			西牛云 Team 敬上<br><br>
			`,
		userName, code)

	message := ""
	for k, v := range m.header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}

	message += "\r\n" + body
	return message
}
