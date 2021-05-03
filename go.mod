module github.com/oakmound/oak/v2

require (
	github.com/200sc/klangsynthese v0.2.2-0.20201022002431-a0e14a8c862b
	github.com/BurntSushi/xgbutil v0.0.0-20190907113008-ad855c713046 // indirect
	github.com/disintegration/gift v1.2.0
	github.com/flopp/go-findfont v0.0.0-20201114153133-e7393a00c15b // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/go-mp3 v0.3.1 // indirect
	github.com/oakmound/libudev v0.2.1
	github.com/oakmound/shiny v0.4.2
	github.com/oakmound/w32 v2.1.0+incompatible
	github.com/oov/directsound-go v0.0.0-20141101201356-e53e59c700bf // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/yobert/alsa v0.0.0-20200618200352-d079056f5370 // indirect
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6
	golang.org/x/mobile v0.0.0-20190415191353-3e0bab5405d6
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6
)

go 1.16

replace github.com/oakmound/shiny => ../shiny

replace github.com/oakmound/w32 => ../w32
