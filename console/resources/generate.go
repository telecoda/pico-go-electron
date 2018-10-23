//go:generate file2byteslice -package=fonts -input=./fonts/font.ttf -output=./fonts/font.go -var=Font_ttf

//go:generate file2byteslice -package=images -input=./images/sprites.gif -output=./images/sprites.go -var=Sprites_png

package resources
