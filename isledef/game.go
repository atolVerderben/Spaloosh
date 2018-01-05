package isledef

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/atolVerderben/tentsuyu"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	//ZoomLevel is the overall zoom of the game for ease of use
	ZoomLevel        float64
	SpalooshSheet    *SpriteSheet
	PlaySoundEffects = true
	AIBroke          = false
	TimeRanOut       = false
	rightAngle       = 1.5708
)

const (
	frameBake int = iota
	frameBomb
	frameGhostGirlSad
	frameGhostGirl
	frameGridBig
	frameGridSmall
	frameKaboom1
	frameKaboom2
	frameKaboom3
	frameKraken
	frameMishipeshu
	frameNureLost
	frameNureWin
	frameNureOnna
	frameBackground
	frameSpaloosh1
	frameSpaloosh2
	frameSpaloosh3
	frameSpalooshLogo
	frameUmibozu
	frameVampHunterSad
	frameVampHunter
	frameVampSad
	frameVamp
)

//Game represents the game itself
type Game struct {
	imageLoadedCh    chan error
	audioLoadedCh    chan error
	gameState        GameState
	pausedState      GameState
	gameData         *GameData
	img              map[string]*ebiten.Image
	orientation      int
	screen           *ebiten.Image
	player           *Player
	background       *backgroundImageParts
	mainState        GameState
	toggleScreenText *tentsuyu.TextElement
}

//NewGame begins a new game of spaloosh!
func NewGame(w, h int) (game *Game, err error) {
	rand.Seed(time.Now().UnixNano())
	defer func() {
		if ferr := finalizeAudio(); ferr != nil && err == nil {
			err = ferr
		}
		if err != nil {
			game = nil
		}
	}()
	game = &Game{
		img:           map[string]*ebiten.Image{},
		imageLoadedCh: make(chan error),
		audioLoadedCh: make(chan error),
		gameData:      NewGameData(GameModeNormalTimed),
	}
	ebiten.SetRunnableInBackground(true) //The timer gets messed up if this isn't enabled

	tentsuyu.BootUp(float64(w), float64(h))

	//GameWorldManager = CreateWorldManager()
	go func() {
		if err := game.loadImages(); err != nil {
			game.imageLoadedCh <- err
		}
		close(game.imageLoadedCh)
	}()
	go func() {
		if err := loadAudio(); err != nil {
			game.audioLoadedCh <- err
		}
		close(game.audioLoadedCh)
	}()

	tentsuyu.Input.RegisterButton("Forward", ebiten.KeyW, ebiten.KeyUp)
	tentsuyu.Input.RegisterButton("Up", ebiten.KeyW, ebiten.KeyUp)
	tentsuyu.Input.RegisterButton("Down", ebiten.KeyS, ebiten.KeyDown)
	tentsuyu.Input.RegisterButton("Backward", ebiten.KeyS, ebiten.KeyDown)
	tentsuyu.Input.RegisterButton("Left", ebiten.KeyA, ebiten.KeyLeft)
	tentsuyu.Input.RegisterButton("Right", ebiten.KeyD, ebiten.KeyRight)
	tentsuyu.Input.RegisterButton("Escape", ebiten.KeyEscape)
	tentsuyu.Input.RegisterButton("ToggleFullscreen", ebiten.KeyF11, ebiten.KeyF)
	tentsuyu.Input.RegisterButton("ChangeScreenScale", ebiten.KeyF10)
	tentsuyu.Input.RegisterButton("ToggleSound", ebiten.KeyF9, ebiten.KeyM)
	tentsuyu.Input.RegisterButton("RotateRight", ebiten.KeyRight)
	tentsuyu.Input.RegisterButton("RotateLeft", ebiten.KeyLeft)
	//ZoomLevel = 3.0
	game.player = CreatePlayer(game)

	return
}

func (g *Game) DrawBackground() error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameBackground].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameBackground].Frame["y"],
		Width:  SpalooshSheet.Frames[frameBackground].Frame["w"],
		Height: SpalooshSheet.Frames[frameBackground].Frame["h"],
	}

	if err := g.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	return nil
}

func (g *Game) loadImages() error {

	font, _ := truetype.Parse(goregular.TTF)
	tentsuyu.Components.UIController.AddFontFile(FntGoRegular, font)
	mono, _ := truetype.Parse(gomono.TTF)
	tentsuyu.Components.UIController.AddFontFile(FntGoMono, mono)
	//tentsuyu.Components.UIController.AddFontFile("font1", tt)
	tentsuyu.Components.UIController.AddFont(FntSmallPixel, "assets/font/small_pixel.ttf")

	/*tentsuyu.ImageManager.LoadImageFromFile(string(ShipTypeBattleShip), "assets/bake-kujira.png")
	tentsuyu.ImageManager.LoadImageFromFile(string(ShipTypeCruiser), "assets/Kraken.png")
	tentsuyu.ImageManager.LoadImageFromFile(string(ShipTypePatrol), "assets/umibozu-horiz.png") //TODO: adjust this to use normal png file
	tentsuyu.ImageManager.LoadImageFromFile("grid", "assets/GRID.png")
	tentsuyu.ImageManager.LoadImageFromFile("explosion", "assets/explosion.png")
	tentsuyu.ImageManager.LoadImageFromFile("bullet", "assets/bomb.png")
	tentsuyu.ImageManager.LoadImageFromFile("girl", "assets/nure-onna.png")
	tentsuyu.ImageManager.LoadImageFromFile("nure-lose", "assets/nure-onna-lose.png")
	tentsuyu.ImageManager.LoadImageFromFile("nure-win", "assets/nure-onna-win.png")
	tentsuyu.ImageManager.LoadImageFromFile("BG", "assets/shoreline.png")
	tentsuyu.ImageManager.LoadImageFromFile("title", "assets/spalooshfolio.png")
	tentsuyu.ImageManager.LoadImageFromFile("spaloosh", "assets/spaloosh3.png")*/
	tentsuyu.ImageManager.LoadImageFromFile("spaloosh-sheet", "assets/spaloosh-sheet.png")

	drawImage, err := ebiten.NewImage(64, 64, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	drawImage.Fill(color.RGBA{R: 0, G: 0, B: 252, A: 255})
	tentsuyu.ImageManager.AddImage("blue", drawImage)
	g.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("blue"), count: 20}
	g.background.SetSize(1600, 1600)

	textImg, err := ebiten.NewImage(1, 1, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	textImg.Fill(color.RGBA{R: 238, G: 252, B: 255, A: 255})
	tentsuyu.ImageManager.AddImage("textBubble", textImg)

	SpalooshSheet = ReadSpriteSheet("assets/spaloosh-sheet.json")

	//ebiten.SetFullscreen(true)
	return nil

}

//ToggleFullscreen toggles the game in or out of full screen
func (g *Game) ToggleFullscreen() {
	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	} else {
		ebiten.SetFullscreen(true)
	}
}

func (g *Game) ToggleSound() {
	PlaySoundEffects = !PlaySoundEffects
}

//Loop controls everything that goes on in the game
func (g *Game) Loop(screen *ebiten.Image) error {
	if g.imageLoadedCh != nil || g.audioLoadedCh != nil {
		select {
		case err := <-g.imageLoadedCh:
			if err != nil {
				return err
			}
			g.imageLoadedCh = nil
		case err := <-g.audioLoadedCh:
			if err != nil {
				return err
			}
			g.audioLoadedCh = nil
		default:
		}
	}
	if g.imageLoadedCh != nil || g.audioLoadedCh != nil {
		return ebitenutil.DebugPrint(screen, "Now Loading...")
	}
	if g.toggleScreenText == nil {
		g.toggleScreenText = tentsuyu.NewTextElementStationary(742, 0, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"F11: Toggle Fullscreen", "F9: Mute Sound Effects"}, color.White, 8)
	}

	/*if err := audioContext.Update(); err != nil {
		return err
	}*/

	tentsuyu.Input.Update()
	g.screen = screen

	if g.gameState == nil {
		//g.mainState = NewGameMain(g)
		g.gameState = CreateTitleMain() //g.mainState
	} else {
		switch g.gameState.Msg() {
		case GameStateMsgReqTitle:
			g.gameState = CreateTitleMain()
		case GameStateMsgReqMainMenu:
			g.gameState = CreateMainMenu(g)
		case GameStateMsgReqMain:
			g.mainState = NewGameMain(g)
			g.gameState = g.mainState
		case GameStateMsgReqBattle:
			g.mainState = NewGameBattle(g)
			g.gameState = g.mainState
		case GameStateMsgReqBattleCharacterSelect:
			g.mainState = CreateBattleCharSelect(g)
			g.gameState = g.mainState
		case GameStateMsgReqMPMainMenu:
			g.mainState = CreateMPMainMenu(g)
			g.gameState = g.mainState
		case GameStateMsgReqSetIP:
			g.mainState = CreateMPSetIP(g)
			g.gameState = g.mainState
		case GameStateMsgReqMPHelp:
			g.mainState = CreateMPHelp(g)
			g.gameState = g.mainState
		case GameStateMsgReqMPStage:
			g.mainState = CreateMPStage(g)
			g.gameState = g.mainState
		case GameStateMsgReqMPMain:
			g.mainState = NewMPBattle(g)
			g.gameState = g.mainState
		case GameStateMsgReqHostingRooms:
			g.mainState = CreateMPRooms(g)
			g.gameState = g.mainState
		case GameStateGameOver:
			g.pausedState = g.gameState
			g.gameState = CreateGameOver(g, false)
		case GameStateGameWin:
			g.pausedState = g.gameState
			g.gameState = CreateGameOver(g, true)
		case GameStateMsgReqMPGameOverLose:
			g.pausedState = g.gameState
			g.gameState = CreateMPGameOver(g, false)
		case GameStateMsgReqMPGameOverWin:
			g.pausedState = g.gameState
			g.gameState = CreateMPGameOver(g, true)
		case GameStateMsgPause:
			g.pausedState = g.gameState
			g.gameState = CreatePaused()
		case GameStateMsgUnPause:
			g.gameState = g.pausedState
			g.gameState.SetMsg(GameStateMsgNone)
		case GameStateMsgReqLostConnection:
			g.pausedState = g.gameState
			g.gameState = CreateLostConnection()
		}
	}
	if g.gameState != nil {
		g.gameState.Update(g)
	}
	if !ebiten.IsRunningSlowly() {
		if err := g.gameState.Draw(g); err != nil {
			return err
		}
		g.toggleScreenText.Draw(g.screen)
	}
	g.gameData.Update()

	if tentsuyu.Input.Button("ToggleFullscreen").JustPressed() {
		g.ToggleFullscreen()
	}
	if tentsuyu.Input.Button("ToggleSound").JustPressed() {
		g.ToggleSound()
	}
	if tentsuyu.Input.Button("ChangeScreenScale").JustPressed() {
		switch g.gameData.currentScale {
		case 1:
			g.gameData.currentScale = 2
			ebiten.SetScreenScale(2)
		case 2:
			g.gameData.currentScale = 1
			ebiten.SetScreenScale(1)
		case 3:
			g.gameData.currentScale = 1
			ebiten.SetScreenScale(1)
		}
	}

	return nil
	//return ebitenutil.DebugPrint(screen, fmt.Sprintf("\nFPS: %.2f", ebiten.CurrentFPS()))
}

//All possible fonts
const (
	FntGoRegular    string = "goregular"
	FntGoBold              = "gobold"
	FntGoItalic            = "goitalic"
	FntGoMono              = "gomono"
	FntGoBoldItalic        = "gobolditalic"
	FntKorean              = "korean"
	FntKoreanBold          = "koreanbold"
	FntEmoji               = "notoemoji"
	FntSymbols             = "notosymbols"
	FntSmallPixel          = "smallpixel"
)
