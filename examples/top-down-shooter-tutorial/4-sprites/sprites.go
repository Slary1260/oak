package main

import (
	"image/color"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/oakmound/oak/v3/render/mod"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/collision/ray"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Collision labels
const (
	Enemy  collision.Label = 1
	Player collision.Label = 2
)

var (
	playerAlive = true
	// Vectors are backed by pointers,
	// so despite this not being a pointer,
	// this does update according to the player's
	// position so long as we don't reset
	// the player's position vector
	playerPos physics.Vector

	sheet [][]*render.Sprite
)

func main() {
	oak.AddScene("tds", scene.Scene{Start: func(ctx *scene.Context) {
		// Initialization
		playerAlive = true
		sprites, err := render.GetSheet(filepath.Join("16x16", "sheet.png"))
		dlog.ErrorCheck(err)
		sheet = sprites.ToSprites()

		// Player setup
		eggplant, err := render.GetSprite(filepath.Join("character", "eggplant-fish.png"))
		playerR := render.NewSwitch("left", map[string]render.Modifiable{
			"left": eggplant,
			// We must copy the sprite before we modify it, or "left"
			// will also be flipped.
			"right": eggplant.Copy().Modify(mod.FlipX),
		})
		if err != nil {
			dlog.Error(err)
		}
		char := entities.NewMoving(100, 100, 32, 32,
			playerR,
			nil, 0, 0)

		char.Speed = physics.NewVector(5, 5)
		playerPos = char.Point.Vector
		render.Draw(char.R, 2)

		char.Bind(event.Enter, func(id event.CID, _ interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			char.Delta.Zero()
			if oak.IsDown(key.W) {
				char.Delta.ShiftY(-char.Speed.Y())
			}
			if oak.IsDown(key.A) {
				char.Delta.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.S) {
				char.Delta.ShiftY(char.Speed.Y())
			}
			if oak.IsDown(key.D) {
				char.Delta.ShiftX(char.Speed.X())
			}
			char.ShiftPos(char.Delta.X(), char.Delta.Y())
			hit := char.HitLabel(Enemy)
			if hit != nil {
				playerAlive = false
			}

			// update animation
			swtch := char.R.(*render.Switch)
			if char.Delta.X() > 0 {
				if swtch.Get() == "left" {
					swtch.Set("right")
				}
			} else if char.Delta.X() < 0 {
				if swtch.Get() == "right" {
					swtch.Set("left")
				}
			}

			return 0
		})

		char.Bind(mouse.Press, func(id event.CID, me interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			mevent := me.(mouse.Event)
			x := char.X() + char.W/2
			y := char.Y() + char.H/2
			ray.DefaultCaster.CastDistance = floatgeom.Point2{x, y}.Sub(floatgeom.Point2{mevent.X(), mevent.Y()}).Magnitude()
			hits := ray.CastTo(floatgeom.Point2{x, y}, floatgeom.Point2{mevent.X(), mevent.Y()})
			for _, hit := range hits {
				hit.Zone.CID.Trigger("Destroy", nil)
			}
			ctx.DrawForTime(
				render.NewLine(x, y, mevent.X(), mevent.Y(), color.RGBA{0, 128, 0, 128}),
				time.Millisecond*50,
				2)
			return 0
		})

		// Create enemies periodically
		event.GlobalBind(event.Enter, func(_ event.CID, frames interface{}) int {
			f := frames.(int)
			if f%EnemyRefresh == 0 {
				NewEnemy(ctx)
			}
			return 0
		})

		// Draw the background
		for x := 0; x < ctx.Window.Width(); x += 16 {
			for y := 0; y < ctx.Window.Height(); y += 16 {
				i := rand.Intn(3) + 1
				// Get a random tile to draw in this position
				sp := sheet[i/2][i%2].Copy()
				sp.SetPos(float64(x), float64(y))
				render.Draw(sp, 1)
			}
		}

	}, Loop: func() bool {
		return playerAlive
	}})

	oak.Init("tds", func(c oak.Config) (oak.Config, error) {
		// This indicates to oak to automatically open and load image and audio
		// files local to the project before starting any scene.
		c.BatchLoad = true
		c.Debug.Level = "Verbose"
		return c, nil
	})
}

// Top down shooter constsv
const (
	EnemyRefresh = 30
	EnemySpeed   = 2
)

// NewEnemy creates an enemy for a top down shooter
func NewEnemy(ctx *scene.Context) {
	x, y := enemyPos(ctx)

	enemyFrame := sheet[0][0].Copy()
	enemyR := render.NewSwitch("left", map[string]render.Modifiable{
		"left":  enemyFrame,
		"right": enemyFrame.Copy().Modify(mod.FlipX),
	})
	enemy := entities.NewSolid(x, y, 16, 16,
		enemyR,
		nil, 0)

	render.Draw(enemy.R, 2)

	enemy.UpdateLabel(Enemy)

	enemy.Bind(event.Enter, func(id event.CID, _ interface{}) int {
		enemy := event.GetEntity(id).(*entities.Solid)
		// move towards the player
		x, y := enemy.GetPos()
		pt := floatgeom.Point2{x, y}
		pt2 := floatgeom.Point2{playerPos.X(), playerPos.Y()}
		delta := pt2.Sub(pt).Normalize().MulConst(EnemySpeed)
		enemy.ShiftPos(delta.X(), delta.Y())

		// update animation
		swtch := enemy.R.(*render.Switch)
		if delta.X() > 0 {
			if swtch.Get() == "left" {
				swtch.Set("right")
			}
		} else if delta.X() < 0 {
			if swtch.Get() == "right" {
				swtch.Set("left")
			}
		}
		return 0
	})

	enemy.Bind("Destroy", func(id event.CID, _ interface{}) int {
		enemy := event.GetEntity(id).(*entities.Solid)
		enemy.Destroy()
		return 0
	})
}

func enemyPos(ctx *scene.Context) (float64, float64) {
	w := ctx.Window.Width()
	h := ctx.Window.Height()
	// Spawn on the edge of the screen
	perimeter := w*2 + h*2
	pos := int(rand.Float64() * float64(perimeter))
	// Top
	if pos < w {
		return float64(pos), 0
	}
	pos -= w
	// Right
	if pos < h {
		return float64(w), float64(pos)
	}
	// Bottom
	pos -= h
	if pos < w {
		return float64(pos), float64(h)
	}
	pos -= w
	// Left
	return 0, float64(pos)
}
