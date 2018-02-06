package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
    "github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
    "golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
)

type Enemy struct {
    pos *pixel.Vec
    shape *imdraw.IMDraw
    id int
}

var enemyList []Enemy
var last time.Time
var state int
var pPos pixel.Vec
var pSpeed float64 = 100.0
var eSpeed float64 = 50.0
var playerShape *imdraw.IMDraw

var Tolerance = 2.0

var PlayerShapePoints = [][]float64 {{-5.0, 0.0}, {0.0, -5.0}, {5.0, 0.0}, {0.0, 5.0}}
var EnemyShapePoints = [][]float64 {{-5.0, -5.0}, {5.0, -5.0}, {5.0, 5.0}, {-5.0, 5.0}}

var enemyIDBase = 100

var basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
var basicText = text.New(pixel.V(100, 500), basicAtlas)


func enemyID() int {
    enemyIDBase++
    return enemyIDBase
}

func makeShape(pos *pixel.Vec, shape *imdraw.IMDraw, points [][]float64, colour pixel.RGBA) *imdraw.IMDraw {
    shape.Color = colour
    for _, point := range points {
        shape.Push(pixel.V(pos.X + point[0], pos.Y + point[1]))
    }
    shape.Polygon(1)
    return shape
}

func setup() {
    enemyList = nil
    pPos = pixel.V(200, 200)
    var eX, eY float64
    // When eX, eY are randomised, they always give the same values...
    for i := 0; i < 10; i++ {
        eX = rand.Float64()*1010
        eY = rand.Float64()*752
        fmt.Println(eX, eY)
        v := pixel.V(eX, eY)
        enemyList = append(enemyList, Enemy{&v, imdraw.New(nil), enemyID()})
    }
    playerShape = imdraw.New(nil)

    last = time.Now()

    state = 1 // 1 = start screen, 2 = game
    fmt.Println("Hey it's the start of the game, we're not doing anything yet really.")
}

func processInput(dt float64, win *pixelgl.Window) {
    if win.Pressed(pixelgl.KeyLeft) {
        pPos.X -= pSpeed * dt
    }
    if win.Pressed(pixelgl.KeyRight) {
        pPos.X += pSpeed * dt
    }
    if win.Pressed(pixelgl.KeyDown) {
        pPos.Y -= pSpeed * dt
    }
    if win.Pressed(pixelgl.KeyUp) {
        pPos.Y += pSpeed * dt
    }
}

func moveEnemies(dt float64, win *pixelgl.Window) {
    // Move enemies
    for index, enemy := range enemyList {
        // move enemy towards player position
        if enemy.pos.X > pPos.X {
            enemy.pos.X -= eSpeed * dt
        } else if enemy.pos.X < pPos.X {
            enemy.pos.X += eSpeed * dt
        }
        if enemy.pos.Y > pPos.Y {
            enemy.pos.Y -= eSpeed * dt
        } else if enemy.pos.Y < pPos.Y {
            enemy.pos.Y += eSpeed * dt
        }

        // Game over if player touches an enemy
        if math.Abs(enemy.pos.X - pPos.X) < Tolerance && math.Abs(enemy.pos.Y - pPos.Y) < Tolerance {
            fmt.Println("Game Over!")
            setup()
        }

        // Draw individual enemy
        enemy.shape = imdraw.New(nil)
        makeShape(enemy.pos, enemy.shape, EnemyShapePoints, pixel.RGB(1, 0.6, 0))
        enemy.shape.Draw(win)

        // Structs are copied by value so I need to reassign it???
        // Probably is a neater way...
        enemyList[index] = enemy
    }
    var deadEnemies = make(map[int]bool)

    // Collide/delete enemies
    for index, enemy := range enemyList {
        for innerIndex, innerEnemy := range enemyList {
            if index != innerIndex && math.Abs(enemy.pos.X - innerEnemy.pos.X) < Tolerance && math.Abs(enemy.pos.Y - innerEnemy.pos.Y) < Tolerance {
                deadEnemies[innerEnemy.id] = true
                fmt.Println(enemy.id)
                continue
            }
        }
    }
    if len(deadEnemies) > 0 {
        fmt.Println(deadEnemies)
        for k := 0; k < len(enemyList); k++ {
            fmt.Println("index", k)
            if deadEnemies[enemyList[k].id] {
                fmt.Println("removed ", enemyList[k].id)
                enemyList[k].pos = nil
                enemyList[k].shape = nil
                enemyList[len(enemyList)-1], enemyList[k] = enemyList[k], enemyList[len(enemyList)-1]
                enemyList = enemyList[:len(enemyList)-1]
                k-- // now that we've deleted something, we have to go back
                fmt.Println("index", k)
            }
        }
    }

    if len(enemyList) < 1 {
        fmt.Println("Game Over!")
        setup()
    }
}

func updateLoop(win *pixelgl.Window) {
    if win.Pressed(pixelgl.KeyEscape) {
        win.SetClosed(true)
    }
    if state == 1 {
        win.Clear(colornames.Purple)
        if win.Pressed(pixelgl.KeySpace) {
            state = 2
        }
        basicText.Draw(win, pixel.IM.Scaled(basicText.Orig, 5))
    } else if state == 2 {
        win.Clear(colornames.Aliceblue)
        dt := time.Since(last).Seconds()
        last = time.Now()
        processInput(dt, win)

        moveEnemies(dt, win)

        // Draw player
        playerShape = imdraw.New(nil)
        makeShape(&pPos, playerShape, PlayerShapePoints, pixel.RGB(0.6, 0, 1))
        playerShape.Draw(win)
    }
    // Final window draw
    win.Update()
}

func run() {
    cfg := pixelgl.WindowConfig{
		Title: "JEGO",
		Bounds: pixel.R(0, 0, 1024, 768),
        VSync: true,
	}
    win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

    fmt.Fprintln(basicText, "JUST EVASION")
    setup()

    for !win.Closed() {
        updateLoop(win)
	}
}

func main() {
	pixelgl.Run(run)
    fmt.Println("And we're done here.")
}
