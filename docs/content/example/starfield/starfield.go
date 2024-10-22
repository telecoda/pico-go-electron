package main

/*
	This is a project to demo an oldskool starfield

	Copyright 2018 @telecoda

*/

import (
	"math/rand"

	"github.com/telecoda/pico-go-electron/console"
)

const (
	// set console type to one of the predefined consoles
	screenWidth  = 128
	screenHeight = 128
	consoleType  = console.CBM64
)

type cartridge struct {
	*console.BaseCartridge
	s []int
}

/* This is the original tweetcart code
s={}w=128 r=rnd for i=1,w do s[i]={}p=s[i]p[1]=r(w)end::a::cls()for i=1,w do p=s[i]pset(p[1],i,i%3+5)p[1]=(p[1]-i%3)%w end flip()goto a
*/

// Init -  called once
func (c *cartridge) Init() error {

	// init stars
	/*
		s={}
		w=128
		r=rnd
		for i=1,w do
			s[i]={}
			p=s[i]
			p[1]=r(w)
		end
	*/

	w := c.GetWidth()
	//h := c.GetHeight()
	c.s = make([]int, w, w)

	for i := 0; i < w; i++ {
		c.s[i] = rand.Intn(w)
	}

	return nil

}

// Update -  called once every frame
func (c *cartridge) Update() {

}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls(console.PICO8_BLACK)
	for i := 0; i < c.GetHeight(); i++ {
		c.PSet(c.s[i], i, console.ColorID(i%3+5))
		c.s[i] = (c.s[i] - (i % 3)) % c.GetWidth()
		if c.s[i] < 0 {
			c.s[i] += c.GetWidth()
		}
	}
}
