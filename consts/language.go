package consts

const (
	Announcement = "announcement"
)

var languages = map[string]map[string]string{
	"zh-CN": {
		Announcement: `当前集群没有开启鉴权，请参考<a href="https://nacos.io/zh-cn/docs/v2/guide/user/auth.html">文档</a>开启鉴权~`,
	},
	"en-US": {
		Announcement: `Authentication has not been enabled in cluster, please refer to <a href="https://nacos.io/en-us/docs/v2/guide/user/auth.html\">Documentation</a> to enable~`,
	},
}

func GetI18n(language, key string) string {
	if _, ok := languages[language]; !ok {
		language = "zh-CN"
	}
	return languages[language][key]
}
