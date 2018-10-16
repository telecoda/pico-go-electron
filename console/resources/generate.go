//go:generate file2byteslice -package=fonts -input=./fonts/font.ttf -output=./fonts/font.go -var=Font_ttf

//go:generate file2byteslice -package=images -input=./images/icons.png -output=./images/icons.go -var=Icons_png
//go:generate file2byteslice -package=images -input=./images/logo.png -output=./images/logo.go -var=Logo_png
//go:generate file2byteslice -package=images -input=./images/sprites.png -output=./images/sprites.go -var=Sprites_png

package resources
