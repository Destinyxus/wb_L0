package main

import "awesomeProject5/publisherOrder/natsPublisher"

func main() {
	pub := natsPublisher.NewPublisher()
	if err := pub.Run(); err != nil {
		panic(err)
	}

}
