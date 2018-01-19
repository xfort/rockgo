package rockgo

import (
	"testing"
	"log"
)

func TestSubstr(t *testing.T) {
	title := "王者荣耀-别花18888金币买这三个英雄, 英雄降价时肯定后悔"
	log.Println(len([]rune(title)))

	title = Substr(title, 0, 29)

	log.Println(len([]rune(title)))


}
