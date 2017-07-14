package helper

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)

	RsaPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)

	Version = "4.0.3"

	AesKey       = "jjsdj#$827gfgh38dakljzmknghsj#2k" //32位
	AesPublicKey = "akljzmknm.ahkjkl"                 //16位

	ConfigPath = "./conf/config.conf"

	//网站信息
	Domain      = "http://www.yougam.com"
	SiteName    = "YOUGAM"
	SiteTitle   = "YOUGAM - 记录分享探索世界的有趣轨迹"
	Keywords    = "YOUGAM,游戏有感,yougam,分享有感,探索有感,创造有感,旅游有感,读书有感,美食有感,电影有感,朋友有感"
	Description = "优姬是一个记录生命旅途的每一分精彩和感动，珍藏美好的回忆，分享创造发现的地方.."

	//验证码
	IsCaptcha = false

	//数据库设定
	DataType  = "goleveldb"
	DBConnect = "goleveldb://../yougam/data/tidb/tidb"

	//API通信常量密匙
	AesConstKey = "Ks*x(" //5位

	//七牛云存储
	BUCKET4QINIU = "yougam"
	DOMAIN4QINIU = "7mnnte.com1.z0.glb.clouddn.com"
	AKEY4QINIU   = "7hFd6CFkhjbtr74AgTlXZ7WgY4mpo4pmsTlyK37D"
	SKEY4QINIU   = "F5EmL1X44LA-ypSIX1TUdZrpwxrDwUI2WiMnRVe6"

	//本地存储根路径设定
	FileStorageDir = "../"

	//数据库表前缀设定
	DBTablePrefix = ""
)
