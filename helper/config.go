package helper

import (
	"strconv"
	"github.com/insionng/yougam/libraries/goconfig"
)

func init() {
	TouchFile(ConfigPath)
	conf, e := goconfig.LoadConfigFile(ConfigPath)
	if e == nil {

		if site, e := conf.GetSection("site"); e == nil {
			Domain = site["Domain"]
			AesConstKey = site["Aes5Keys"]
			SiteName = site["SiteName"]
			SiteTitle = site["SiteTitle"]
			Keywords = site["Keywords"]
			Description = site["Description"]
		} else {

			if v, _ := conf.GetValue("site", "Domain"); v == "" {
				conf.SetValue("site", "Domain", Domain)
			}
			if v, _ := conf.GetValue("site", "SiteName"); v == "" {
				conf.SetValue("site", "SiteName", SiteName)
			}
			if v, _ := conf.GetValue("site", "Aes5Keys"); v == "" {
				conf.SetValue("site", "Aes5Keys", AesConstKey)
			}
			if v, _ := conf.GetValue("site", "SiteTitle"); v == "" {
				conf.SetValue("site", "SiteTitle", SiteTitle)
			}
			if v, _ := conf.GetValue("site", "Keywords"); v == "" {
				conf.SetValue("site", "Keywords", Keywords)
			}
			if v, _ := conf.GetValue("site", "Description"); v == "" {
				conf.SetValue("site", "Description", Description)
			}
		}

		if datebase, e := conf.GetSection("datebase"); e == nil {
			DataType = datebase["DataType"]
			DBConnect = datebase["Connect"]
		} else {

			if v, _ := conf.GetValue("datebase", "DataType"); v == "" {
				conf.SetValue("datebase", "DataType", DataType)
			}
			if v, _ := conf.GetValue("datebase", "Connect"); v == "" {
				conf.SetValue("datebase", "Connect", DBConnect)
			}
		}

		if cloud, e := conf.GetSection("cloud"); e == nil {
			BUCKET4QINIU = cloud["BUCKET4QINIU"]
			DOMAIN4QINIU = cloud["DOMAIN4QINIU"]
			AKEY4QINIU = cloud["AKEY4QINIU"]
			SKEY4QINIU = cloud["DOMAIN4QINIU"]
		} else {
			if v, _ := conf.GetValue("cloud", "BUCKET4QINIU"); v == "" {
				conf.SetValue("cloud", "BUCKET4QINIU", "")
			}
			if v, _ := conf.GetValue("cloud", "DOMAIN4QINIU"); v == "" {
				conf.SetValue("cloud", "DOMAIN4QINIU", "")
			}
			if v, _ := conf.GetValue("cloud", "AKEY4QINIU"); v == "" {
				conf.SetValue("cloud", "AKEY4QINIU", "")
			}
			if v, _ := conf.GetValue("cloud", "SKEY4QINIU"); v == "" {
				conf.SetValue("cloud", "SKEY4QINIU", "")
			}
		}

		if mail, e := conf.GetSection("mail"); e == nil {
			SmtpHost = mail["SmtpHost"]
			SmtpPort = mail["SmtpPort"]
			MailUser = mail["MailUser"]
			MailPassword = mail["MailPassword"]
			MailAdline = mail["MailAdline"]
		} else {
			if v, _ := conf.GetValue("mail", "SmtpHost"); v == "" {
				conf.SetValue("mail", "SmtpHost", SmtpHost)
			}
			if v, _ := conf.GetValue("mail", "SmtpPort"); v == "" {
				conf.SetValue("mail", "SmtpPort", SmtpPort)
			}
			if v, _ := conf.GetValue("mail", "MailUser"); v == "" {
				conf.SetValue("mail", "MailUser", MailUser)
			}
			if v, _ := conf.GetValue("mail", "MailPassword"); v == "" {
				conf.SetValue("mail", "MailPassword", MailPassword)
			}
			if v, _ := conf.GetValue("mail", "MailAdline"); v == "" {
				conf.SetValue("mail", "MailAdline", MailAdline)
			}

			goconfig.SaveConfigFile(conf, ConfigPath)

		}
		if mail, e := conf.GetSection("signup"); e != nil {
			IsCaptcha = true
		} else {
			if isCaptcha, err := strconv.ParseBool(mail["captcha"]); err == nil {
				IsCaptcha = isCaptcha
			}
		}
	}
}

func Theme() (theme string) {

	TouchFile(ConfigPath)

	if conf, e := goconfig.LoadConfigFile(ConfigPath); e == nil {
		if style, e := conf.GetValue("themes", "style"); (style == "") || (e != nil) {
			conf.SetValue("themes", "style", "yougam")
			goconfig.SaveConfigFile(conf, ConfigPath)
			theme = "yougam"
		} else {
			theme = style
		}
	}

	return theme

}

func IsSendMail() (isSendMail bool) {

	TouchFile(ConfigPath)

	if conf, e := goconfig.LoadConfigFile(ConfigPath); e == nil {
		if i, e := conf.GetValue("signup", "sendmail"); (i == "") || (e != nil) {
			conf.SetValue("signup", "sendmail", "false")
			goconfig.SaveConfigFile(conf, ConfigPath)
			isSendMail = false
		} else {
			if i == "true" {
				isSendMail = true
			} else {
				isSendMail = false
			}
		}
	}

	return isSendMail

}
