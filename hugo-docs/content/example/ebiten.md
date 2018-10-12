+++
title = "ebiten"
draft = false
weight = 99

+++

You can use _most_ of the ebiten examples too. 

Just trying opening a file from [$GOPATH/src/github.com/hajimehoshi/ebiten/examples]($GOPATH/src/github.com/hajimehoshi/ebiten/examples)

The only code you _may_ need to change is the scale factor when running ebiten.  Set the value to `1` instead of `2` on some of the demos and we'll handle the scaling for you.

{{< highlight go >}}
func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "(Demo)"); err != nil {
		log.Fatal(err)
	}
}
{{< / highlight >}}

If your code has variables or constants called  **screenWidth** and **screenHeight** these will be used to maintain aspect ratio when the window is resized.