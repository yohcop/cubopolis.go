package main

import "sdl"
import "gl"
import "flag"
import (
	"math"
)
const (
  PiOver100 float64 = 0.0174532925199433
  cubeWidth = 0.5
)

var printInfo = flag.Bool("info", false, "print GL implementation information")

var T0 uint32 = 0
var Frames uint32 = 0


func cube(color []float32) {
  gl.Color4f(color[0], color[1], color[2], color[3])
	gl.Begin(gl.QUADS)
	// TOP
	gl.Vertex3f(cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, cubeWidth, cubeWidth)
	// Bottom
	gl.Vertex3f(cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, -cubeWidth)
	// Front
	gl.Vertex3f(cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, cubeWidth)
	// Back
	gl.Vertex3f(cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(cubeWidth, cubeWidth, -cubeWidth)
	// Left
	gl.Vertex3f(-cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, cubeWidth)
	// Right
	gl.Vertex3f(cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, -cubeWidth)
	gl.End()

  gl.Color4f(color[0] / 1.2, color[1] / 1.2, color[2] / 1.2, 1.0)
  gl.Begin(gl.LINE_STRIP);						// Start Drawing Our Player Using Lines
	// TOP
	gl.Vertex3f(cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, cubeWidth, cubeWidth)
	// Front
	gl.Vertex3f(-cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, cubeWidth)
	// Bottom
	gl.Vertex3f(-cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, -cubeWidth)
	// Back
	gl.Vertex3f(-cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(cubeWidth, cubeWidth, -cubeWidth)
	// Right
	gl.Vertex3f(cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, cubeWidth)
	gl.Vertex3f(cubeWidth, -cubeWidth, -cubeWidth)
	// Left
	gl.Vertex3f(-cubeWidth, cubeWidth, cubeWidth)
	gl.Vertex3f(-cubeWidth, cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, -cubeWidth)
	gl.Vertex3f(-cubeWidth, -cubeWidth, cubeWidth)
  gl.End();
  gl.Color4f(1.0, 1.0, 1.0, 1.0)
}

var (
	yrot                    float32 // camera rotation
	xpos, zpos              float32 // camera position
	walkbias, walkbiasangle float32 // head-bobbing....
	lookupdown              float32

	lightAmbient  = []float32{0.5, 0.5, 0.5, 1.0}
	lightDiffuse  = []float32{1.0, 1.0, 1.0, 1.0}
	lightPosition = []float32{0.0, 0.0, 2.0, 1.0}
)

var cubes map[string]uint = make(map[string]uint)
var chunks []*ChunkInfo

func draw() {
  xtrans   := -xpos
  ztrans   := -zpos
  ytrans   := -walkbias - 0.25
  scenroty := 360.0 - yrot

  // Clear the screen and depth buffer
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

  // reset the view
  gl.LoadIdentity()

  // Rotate up and down to look up and down
  gl.Rotatef(lookupdown, 1.0, 0.0, 0.0)
  // Rotate depending on direction player is facing
  gl.Rotatef(scenroty, 0.0, 1.0, 0.0)
  // translate the scene based on player position
  gl.Translatef(xtrans, ytrans - 1.5, ztrans)

	for _, chunk := range chunks {
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z <= 16; z++ {
					if len(chunk.Data[y][x]) > z && chunk.Data[y][x][z] != "" {
						gl.PushMatrix()
						gl.Translated(
							float64(16*chunk.X+x),
							float64(z),
							float64(16*chunk.Y+y))
						gl.CallList(cubes[chunk.Data[y][x][z]])
						gl.PopMatrix()
					}
				}
			}
		}
	}

	gl.PopMatrix()

	sdl.GL_SwapBuffers()

	Frames++
	{
		t := sdl.GetTicks()
		if t-T0 >= 5000 {
			seconds := (t - T0) / 1000.0
			fps := Frames / seconds
			print(Frames, " frames in ", seconds, " seconds = ", fps, " FPS\n")
			T0 = t
			Frames = 0
		}
	}
}


func idle() {
}

/* new window size or exposure */

func reshape(width int, height int) {
	// protect against a divide by zero
	if height == 0 {
		height = 1
	}

	// Setup our viewport
	gl.Viewport(0, 0, width, height)

	// change to the projection matrix and set our viewing volume.
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	// aspect ratio
	aspect := float64(width) / float64(height)

	// Set our perspective.
	// This code is equivalent to using gluPerspective as in the original tutorial.
	var fov, near, far float64
	fov = 45.0
	near = 0.1
	far = 100.0
	top := math.Tan(float64(fov*math.Pi/360.0)) * near
	bottom := -top
	left := aspect * bottom
	right := aspect * top
	gl.Frustum(left, right, bottom, top, near, far)

	// Make sure we're changing the model view and not the projection
	gl.MatrixMode(gl.MODELVIEW)

	// Reset the view
	gl.LoadIdentity()
}

func fetchChunks(chunkY, chunkX int) {
	for y := -chunkY; y < chunkY; y++ {
		for x := -chunkX; x < chunkX; x++ {
			if chunk, err := GetChunk(y, x); err == nil {
				chunks = append(chunks, chunk)
			}
		}
	}
}

func init_() {
	//pos := []float64{5.0, 5.0, 10.0, 0.0}

	colors := map[string][]float32{
		"0": []float32{0.7, 0.7, 0.7, 1.0},
		"1": []float32{1.0, 0.0, 0.0, 1.0},
		"2": []float32{1.0, 0.6, 0.0, 1.0},
		"3": []float32{1.0, 1.0, 0.0, 1.0},
		"4": []float32{0.0, 1.0, 0.0, 1.0},
		"5": []float32{0.0, 0.0, 1.0, 1.0},
		"6": []float32{0.0, 0.6, 1.0, 1.0},
		"7": []float32{1.0, 0.0, 1.0, 1.0},
		"8": []float32{1.0, 1.0, 1.0, 1.0},
		"9": []float32{0.3, 0.3, 0.3, 1.0},
	}

	fetchChunks(2, 2)

  //gl.ShadeModel(gl.SMOOTH)
  gl.ClearColor(0.2, 0.2, 0.6, 0.0)
  gl.ClearDepth(1.0)
  gl.Enable(gl.DEPTH_TEST)
  //gl.DepthFunc(gl.LEQUAL)
  gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

  gl.Lightfv(gl.LIGHT1, gl.AMBIENT,  lightAmbient )
  gl.Lightfv(gl.LIGHT1, gl.DIFFUSE,  lightDiffuse )
  gl.Lightfv(gl.LIGHT1, gl.POSITION, lightPosition)
  gl.Enable(gl.LIGHT1)

  gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
  gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST);

  //gl.Enable(gl.LIGHTING)
  //gl.Enable(gl.BLEND)

	for name, color := range colors {
		/* make a cube */
		cubes[name] = gl.GenLists(1)
		gl.NewList(cubes[name], gl.COMPILE)
		//gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, color)
		cube(color)
		gl.EndList()
	}

	if *printInfo {
		print("GL_RENDERER   = ", gl.GetString(gl.RENDERER), "\n")
		print("GL_VERSION    = ", gl.GetString(gl.VERSION), "\n")
		print("GL_VENDOR     = ", gl.GetString(gl.VENDOR), "\n")
		print("GL_EXTENSIONS = ", gl.GetString(gl.EXTENSIONS), "\n")
	}

}

// handle key press events
func handleKeyPress(keysym []uint8) {
  keys := sdl.GetKeyState()

  if keys[sdl.K_RIGHT] == 1 {
    yrot -= 3.5
  }

  if keys[sdl.K_LEFT] == 1 {
    yrot += 3.5
  }

  if keys[sdl.K_UP] == 1 {
    xpos -= float32(math.Sin(float64(yrot) * PiOver100)) * 0.5
    zpos -= float32(math.Cos(float64(yrot) * PiOver100)) * 0.5
    if walkbiasangle >= 359.0 {
      walkbiasangle = 0.0
    } else {
      walkbiasangle += 10.0
    }
    walkbias = float32(math.Sin(float64(walkbiasangle) * PiOver100)) / 20.0
  }

  if keys[sdl.K_DOWN] == 1 {
    xpos += float32(math.Sin(float64(yrot) * PiOver100)) * 0.5
    zpos += float32(math.Cos(float64(yrot) * PiOver100)) * 0.5
    if walkbiasangle <= 1.0 {
      walkbiasangle = 359.0
    } else {
      walkbiasangle -= 10.0
    }
    walkbias = float32(math.Sin(float64(walkbiasangle) * PiOver100)) / 20.0
  }
}

func main() {

	flag.Parse()

	var done bool
	var keys []uint8

	sdl.Init(sdl.INIT_VIDEO)

	var screen = sdl.SetVideoMode(640, 480, 16, sdl.OPENGL|sdl.RESIZABLE)

	if screen == nil {
		sdl.Quit()
		panic("Couldn't set 300x300 GL video mode: " + sdl.GetError() + "\n")
	}

	sdl.WM_SetCaption("Gears", "gears")

	init_()
	reshape(int(screen.W), int(screen.H))
	done = false
	for !done {
		var event sdl.Event

		idle()
		for event.Poll() {
			switch event.Type {
			case sdl.VIDEORESIZE:
				screen = sdl.SetVideoMode(int(event.Resize().W), int(event.Resize().H), 16,
					sdl.OPENGL|sdl.RESIZABLE)
				if screen != nil {
					reshape(int(screen.W), int(screen.H))
				} else {
					panic("we couldn't set the new video mode??")
				}
				break

			case sdl.QUIT:
				done = true
				break
			}
		}
		keys = sdl.GetKeyState()

    handleKeyPress(keys)
		if keys[sdl.K_ESCAPE] != 0 {
			done = true
		}

		draw()
	}
	sdl.Quit()
	return

}
