package format

import "strings"

var initialisms = map[string]string{
	"Api":   "API",
	"Ascii": "ASCII",
	"Cpu":   "CPU",
	"Css":   "CSS",
	"Dns":   "DNS",
	"Eof":   "EOF",
	"Guid":  "GUID",
	"Html":  "HTML",
	"Http":  "HTTP",
	"Https": "HTTPS",
	"Id":    "ID",
	"Ip":    "IP",
	"Json":  "JSON",
	"Qps":   "QPS",
	"Ram":   "RAM",
	"Rpc":   "RPC",
	"Sla":   "SLA",
	"Smtp":  "SMTP",
	"Sql":   "SQL",
	"Ssh":   "SSH",
	"Tcp":   "TCP",
	"Tls":   "TLS",
	"Ttl":   "TTL",
	"Udp":   "UDP",
	"Ui":    "UI",
	"Uid":   "UID",
	"Uuid":  "UUID",
	"Uri":   "URI",
	"Url":   "URL",
	"Utf8":  "UTF8",
	"Vm":    "VM",
	"Xml":   "XML",
	"Xsrf":  "XSRF",
	"Xss":   "XSS",
}

func NormalizeInitialism(word string) string {
	for k, v := range initialisms {
		if strings.HasPrefix(word, k) {
			return v
		}
	}
	return word
}
