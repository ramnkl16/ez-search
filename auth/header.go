package auth

type Header struct {
	HeaderXPublic string `header:"x-public"`
	HeaderXauth   string `header:"x-auth"`
	HeaderXDebug  string `header:"x-debug"`
	HeaderNS      string `header:"x-ns"`
}
