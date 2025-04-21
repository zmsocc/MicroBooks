package main

func main() {
	//server := gin.Default()
	server := InitWebServer()
	server.Run(":8080")
}
