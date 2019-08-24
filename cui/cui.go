package cui

import (
	"time"

	"github.com/hiroebe/gakki"
	"github.com/nsf/termbox-go"
)

func GetDefaultCUI() gakki.UI {
	return &defaultCUI{
		typedKeys: map[rune]bool{},
	}
}

type defaultCUI struct {
	keyboard       [][]rune
	keyDisplayFunc func(rune) string

	typedKeys map[rune]bool
}

func (c *defaultCUI) SetKeyboard(keyboard [][]rune) {
	c.keyboard = keyboard
}

func (c *defaultCUI) SetKeyDisplayFunc(f func(rune) string) {
	c.keyDisplayFunc = f
}

func (c *defaultCUI) Run(keydownCh, keyupCh chan<- rune) error {
	if err := termbox.Init(); err != nil {
		return err
	}
	c.draw()

MAINLOOP:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				break MAINLOOP
			}

			if c.typedKeys[ev.Ch] {
				break
			}
			c.typedKeys[ev.Ch] = true
			keydownCh <- ev.Ch

			go func() {
				time.Sleep(1 * time.Second)
				c.typedKeys[ev.Ch] = false
				keyupCh <- ev.Ch
				c.draw()
			}()

		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		}

		c.draw()
	}
	return nil
}

func (c *defaultCUI) Close() {
	termbox.Close()
}

func (c *defaultCUI) draw() {
	w, _ := termbox.Size()
	cellsPerKey := w / len(c.keyboard[0])

	for i, keys := range c.keyboard {
		y := i * 2
		for j, r := range keys {
			x := j*cellsPerKey + i

			contents := make([]rune, cellsPerKey-1)
			if c.keyDisplayFunc != nil {
				s := c.keyDisplayFunc(r)
				copy(contents, []rune(s))
			} else {
				contents[len(contents)/2] = r
			}

			fg := termbox.ColorDefault
			bg := termbox.ColorDefault
			if c.typedKeys[r] {
				bg = termbox.ColorRed
			}
			for offset, r := range contents {
				termbox.SetCell(x+offset, y, r, fg, bg)
			}
		}
	}

	termbox.Flush()
}
