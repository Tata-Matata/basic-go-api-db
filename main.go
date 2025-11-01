package main

func main() {
	app := App{}

	app.Initialize()

	app.Run(":10000")

}
