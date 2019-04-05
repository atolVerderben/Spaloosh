package isledef

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	assets "github.com/atolVerderben/spaloosh/isledef/internal"
	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

var (
	//ZoomLevel is the overall zoom of the game for ease of use
	ZoomLevel                 float64
	SpalooshSheet             *SpriteSheet
	AIBroke                   = false
	TimeRanOut                = false
	rightAngle                = 1.5708
	Game                      *tentsuyu.Game
	ScreenWidth, ScreenHeight float64
	GamePlayer                *Player
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

//NewGame begins a new game of spaloosh!
func NewGame(w, h float64) (game *tentsuyu.Game, err error) {
	rand.Seed(time.Now().UnixNano())
	ScreenWidth, ScreenHeight = w, h

	game, _ = tentsuyu.NewGame(w, h)
	ebiten.SetRunnableInBackground(true) //The timer gets messed up if this isn't enabled

	game.Input.RegisterButton("Forward", ebiten.KeyW, ebiten.KeyUp)
	game.Input.RegisterButton("Up", ebiten.KeyW, ebiten.KeyUp)
	game.Input.RegisterButton("Down", ebiten.KeyS, ebiten.KeyDown)
	game.Input.RegisterButton("Backward", ebiten.KeyS, ebiten.KeyDown)
	game.Input.RegisterButton("Left", ebiten.KeyA, ebiten.KeyLeft)
	game.Input.RegisterButton("Right", ebiten.KeyD, ebiten.KeyRight)
	game.Input.RegisterButton("Escape", ebiten.KeyEscape)
	game.Input.RegisterButton("ToggleFullscreen", ebiten.KeyF11, ebiten.KeyF)
	game.Input.RegisterButton("ChangeScreenScale", ebiten.KeyF10)
	game.Input.RegisterButton("ToggleSound", ebiten.KeyF9, ebiten.KeyM)
	game.Input.RegisterButton("RotateRight", ebiten.KeyRight)
	game.Input.RegisterButton("RotateLeft", ebiten.KeyLeft)

	game.UIController.AddFont(FntSmallPixel, assets.ReturnPixelFont())
	game.LoadImages(func() *tentsuyu.ImageManager {
		return loadImages()
	})
	game.LoadAudio(func() *tentsuyu.AudioPlayer {
		return loadAudio()
	})

	game.SetGameStateLoop(func() error {
		switch game.GetGameState().Msg() {
		case GameStateMsgReqTitle:
			game.SetGameState(CreateTitleMain())
		case GameStateMsgReqMainMenu:
			game.SetGameState(CreateMainMenu(game))
		case GameStateMsgReqMain:
			game.SetGameState(NewGameMain(game))
		case GameStateMsgReqBattle:
			game.SetGameState(NewGameBattle(game))
		case GameStateMsgReqBattleCharacterSelect:
			game.SetGameState(CreateBattleCharSelect(game))
		case GameStateMsgReqMPMainMenu:
			game.SetGameState(CreateMPMainMenu(game))
		case GameStateMsgReqSetIP:
			game.SetGameState(CreateMPSetIP(game))
		case GameStateMsgReqMPHelp:
			game.SetGameState(CreateMPHelp(game))
		case GameStateMsgReqMPStage:
			game.SetGameState(CreateMPStage(game))
		case GameStateMsgReqMPMain:
			game.SetGameState(NewMPBattle(game))
		case GameStateMsgReqHostingRooms:
			game.SetGameState(CreateMPRooms(game))
		case GameStateGameOver:
			game.SetPauseState(CreateGameOver(game, false))
		case GameStateGameWin:
			game.SetPauseState(CreateGameOver(game, true))
		case GameStateMsgReqMPGameOverLose:
			game.SetPauseState(CreateMPGameOver(game, false))
		case GameStateMsgReqMPGameOverWin:
			game.SetPauseState(CreateMPGameOver(game, true))
		case GameStateMsgPause:
			game.SetPauseState(CreatePaused())
		case GameStateMsgUnPause:
			game.UnPause()
		case GameStateMsgReqLostConnection:
			game.SetPauseState(CreateLostConnection())
		case tentsuyu.GameStateMsgNotStarted:
			game.SetGameState(CreateTitleMain())
		}

		if game.Input.Button("ToggleSound").JustPressed() {
			ToggleSound()
		}
		if game.Input.Button("ChangeScreenScale").JustPressed() {
			switch game.GameData.Settings["Scale"].ValueInt {
			case 1:
				game.GameData.Settings["Scale"].ValueInt = 2
				ebiten.SetScreenScale(2)
			case 2:
				game.GameData.Settings["Scale"].ValueInt = 1
				ebiten.SetScreenScale(1)
			case 3:
				game.GameData.Settings["Scale"].ValueInt = 1
				ebiten.SetScreenScale(1)
			}
		}

		return nil
	})
	InitGameData(game.GameData, GameModeNormalTimed)
	//ZoomLevel = 3.0
	GamePlayer = CreatePlayer(game)

	Game = game

	return
}

func DrawBackground(game *tentsuyu.Game) error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameBackground].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameBackground].Frame["y"],
		Width:  SpalooshSheet.Frames[frameBackground].Frame["w"],
		Height: SpalooshSheet.Frames[frameBackground].Frame["h"],
	}

	if err := game.Screen.DrawImage(game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	return nil
}

func loadImages() *tentsuyu.ImageManager {
	imageManager := tentsuyu.NewImageManager()

	//tentsuyu.ImageManager.LoadImageFromFile("spaloosh-sheet", "assets/spaloosh-sheet.png")
	sImg, err := assets.LoadSpalooshSheet()
	if err != nil {
		log.Fatal(err)
	}
	imageManager.AddImage("spaloosh-sheet", sImg)

	drawImage, err := ebiten.NewImage(64, 64, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	drawImage.Fill(color.RGBA{R: 0, G: 0, B: 252, A: 255})
	imageManager.AddImage("blue", drawImage)

	textImg, err := ebiten.NewImage(1, 1, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	textImg.Fill(color.RGBA{R: 238, G: 252, B: 255, A: 255})
	imageManager.AddImage("textBubble", textImg)

	SpalooshSheet = ReadSpriteSheetJSON(assets.LoadSpriteSheetJSON()) //ReadSpriteSheet("assets/spaloosh-sheet.json")

	//ebiten.SetFullscreen(true)
	return imageManager

}

//ToggleFullscreen toggles the game in or out of full screen
func ToggleFullscreen() {
	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	} else {
		ebiten.SetFullscreen(true)
	}
}

//ToggleSound turns the sound on and off
func ToggleSound() {
	Game.AudioPlayer.MuteAll(true)
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
