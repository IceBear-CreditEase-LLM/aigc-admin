package main

import (
	"embed"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/cmd/service"
)

var (
	//go:embed web
	webFs embed.FS
	//go:embed data
	dataFs embed.FS
)

func main() {
	service.WebFs = webFs
	service.DataFs = dataFs
	service.Run()
}
