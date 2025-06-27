package main

import (
	_ "bondi-push-notification/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// privateKey, publicKey, _ := webpush.GenerateVAPIDKeys()
	// fmt.Println(privateKey)
	// fmt.Println(publicKey)
	beego.Run()

}
