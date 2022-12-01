package metrics

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vokomarov/home-exporter/config"
	"github.com/vokomarov/home-exporter/telegram"
)

const MessageUp = `🟢 Internet back online in home %s`
const MessageDown = `🔴 Internet is offline in home %s`

type InternetStatusMetric struct {
	lastStatus  InternetStatus
	stop        bool
	name        string
	host        string
	port        string
	method      string
	retries     int
	timeout     time.Duration
	interval    time.Duration
	tgChatId    int64
	retryAmount int
}

func NewInternetStatusMetric(config config.Home) *InternetStatusMetric {
	if !config.InternetStatus.Enabled {
		return nil
	}

	metric := InternetStatusMetric{
		name:     config.Name,
		host:     config.InternetStatus.Host,
		port:     config.InternetStatus.Port,
		method:   config.InternetStatus.Method,
		retries:  config.InternetStatus.Retries,
		timeout:  time.Duration(config.InternetStatus.Timeout) * time.Second,
		interval: time.Duration(config.InternetStatus.Interval) * time.Second,
		tgChatId: config.TelegramChatId,
		lastStatus: InternetStatus{
			status: true,
		},
	}

	return &metric
}

func (m *InternetStatusMetric) Run() {
	m.log("Launching for %s://%s:%s [interval: %s] [timeout: %s] [retries: %d]", m.method, m.host, m.port, m.interval, m.timeout, m.retries)

	for {
		if m.stop {
			m.log("Stop signal received")
			break
		}

		actualStatus := m.ping()
		statusChanged := m.lastStatus.String() != actualStatus.String()

		if statusChanged && m.retries == 0 || statusChanged && (actualStatus.IsUp() || m.retryAmount >= m.retries) {
			m.statusChange(actualStatus)
		} else if statusChanged && m.retries > 0 {
			m.statusChangeRetry(actualStatus)
		} else {
			m.retryAmount = 0
		}

		if err := m.trackMetric(); err != nil {
			m.log("Unable to track metric, stopping: %v", err)
			break
		}

		time.Sleep(m.interval)
	}
}

func (m *InternetStatusMetric) statusChange(actualStatus InternetStatus) {
	m.log("Detected status change %s => %s", m.lastStatus.String(), actualStatus.String())

	if err := m.telegramMessage(actualStatus); err != nil {
		m.log("Unable to send telegram bot message: %v", err)
	}

	m.lastStatus = actualStatus
	m.retryAmount = 0

	if !actualStatus.IsUp() {
		m.log("Connection error: %v", actualStatus.Error())
	}
}

func (m *InternetStatusMetric) statusChangeRetry(actualStatus InternetStatus) {
	m.retryAmount++
	m.log("Possible status change %s => %s, retry %d", m.lastStatus.String(), actualStatus.String(), m.retryAmount)
}

func (m *InternetStatusMetric) trackMetric() error {
	metric, err := homeInternetStatus.GetMetricWithLabelValues(
		m.name,
		m.host,
		m.port,
		m.method,
	)
	if err != nil {
		return fmt.Errorf("fetching metric with label values error: %w", err)
	}

	metric.Set(m.lastStatus.MetricValue())

	return nil
}

func (m *InternetStatusMetric) Stop() {
	m.stop = true
	m.log("Sending stop signal")
}

func (m *InternetStatusMetric) log(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	log.Printf("[home=%s] Internet status monitor: %s\n", m.name, message)
}

func (m *InternetStatusMetric) ping() InternetStatus {
	switch m.method {
	case "tcp":
		return NewInternetStatus(m.pingTcp(m.host, m.port, m.timeout))
	case "icmp":
		return NewInternetStatus(m.pingIcmp(m.host, m.timeout))
	default:
		return NewInternetStatus(false, fmt.Errorf("unsupported method"))
	}
}

func (m *InternetStatusMetric) telegramMessage(actualStatus InternetStatus) error {
	if m.tgChatId == 0 {
		return nil
	}

	var message string

	if actualStatus.IsUp() {
		message = fmt.Sprintf(MessageUp, m.name)
	} else {
		message = fmt.Sprintf(MessageDown, m.name)
	}

	if err := telegram.Bot.Send(message, m.tgChatId); err != nil {
		return fmt.Errorf("sending status change telegram message: %w", err)
	}

	return nil
}

func (m *InternetStatusMetric) pingTcp(host, port string, timeout time.Duration) (bool, error) {
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if conn != nil {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				m.log("Unable to close network connection: %v", err)
			}
		}(conn)
	}

	return conn != nil && err == nil, err
}

func (m *InternetStatusMetric) pingIcmp(host string, timeout time.Duration) (bool, error) {
	_, _, err := ICMPPing(host, timeout)
	return err == nil, err
}
