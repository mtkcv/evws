module evws

go 1.13

require (
	github.com/gorilla/websocket v1.4.2
	go.uber.org/zap v1.15.0
	golang.org/x/tools v0.0.0-20191112195655-aa38f8e97acc // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/gorilla/websocket v1.4.2 => github.com/mtkcv/websocket v1.4.3-0.20200704063855-5fdbe6fce0b1
