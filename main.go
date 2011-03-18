package main

import "sdl"
import "gl"
import "flag"

var printInfo = flag.Bool("info", false, "print GL implementation information")

var T0 uint32 = 0
var Frames uint32 = 0

const cubeWidth = 0.5

func cube(cube string) {
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
}

var view_rotx float64 = 0.0
var view_roty float64 = 0.0
var view_rotz float64 = 0.0
var view_x float64 = 0.0
var view_y float64 = 0.0
var view_z float64 = 0.0

var cubes map[string]uint = make(map[string]uint)
var chunks []*ChunkInfo

func draw() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.PushMatrix()
	gl.Rotated(view_rotx, 1.0, 0.0, 0.0)
	gl.Rotated(view_roty, 0.0, 1.0, 0.0)
	gl.Rotated(view_rotz, 0.0, 0.0, 1.0)
  gl.Translated(view_x, view_z, view_y)

  for _, chunk := range chunks {
    for y := 0; y < 16; y++ {
      for x := 0; x < 16; x++ {
        for z := 0; z <= 16; z++ {
          if len(chunk.Data[y][x]) > z && chunk.Data[y][x][z] != "" {
            gl.PushMatrix()
            gl.Translated(
                float64(16 * chunk.X + x),
                float64(16 * chunk.Y + y),
                float64(z))
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

	h := float64(height) / float64(width)

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1.0, 1.0, -h, h, 5.0, 1000.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0.0, 0.0, -40.0)
}

func fetchChunks(chunkY, chunkX int) {
  for y := - chunkY; y < chunkY; y++ {
    for x := - chunkX; x < chunkX; x++ {
      if chunk, err := GetChunk(y, x); err == nil {
        chunks = append(chunks, chunk)
      }
    }
  }
}

func init_() {
	pos := []float32{5.0, 5.0, 10.0, 0.0}

  colors := map[string][]float32 {
	  "0": []float32{0.7, 0.7, 0.7, 1.0},
	  "1": []float32{1.0, 0.0, 0.0, 1.0},
    "2": []float32{1,0.66,0,1},
    "3": []float32{1,1,0,1},
    "4": []float32{0,1,0,1},
    "5": []float32{0,0,1,1},
    "6": []float32{0,0.66,1,1},
    "7": []float32{1,0,1,1},
    "8": []float32{1,1,1,1},
    "9": []float32{0.3,0.3,0.3,1},
  }

  fetchChunks(2, 2)
	gl.Lightfv(gl.LIGHT0, gl.POSITION, pos)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)
	gl.Enable(gl.DEPTH_TEST)

  for name, color := range colors {
    /* make a cube */
    cubes[name] = gl.GenLists(1)
    gl.NewList(cubes[name], gl.COMPILE)
    gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, color)
    cube("")
    gl.EndList()
  }

	gl.Enable(gl.NORMALIZE)

	if *printInfo {
		print("GL_RENDERER   = ", gl.GetString(gl.RENDERER), "\n")
		print("GL_VERSION    = ", gl.GetString(gl.VERSION), "\n")
		print("GL_VENDOR     = ", gl.GetString(gl.VENDOR), "\n")
		print("GL_EXTENSIONS = ", gl.GetString(gl.EXTENSIONS), "\n")
	}

}

func main() {

	flag.Parse()

	var done bool
	var keys []uint8

	sdl.Init(sdl.INIT_VIDEO)

	var screen = sdl.SetVideoMode(480, 390, 16, sdl.OPENGL|sdl.RESIZABLE)

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

		if keys[sdl.K_ESCAPE] != 0 {
			done = true
		}
		if keys[sdl.K_UP] != 0 {
			view_rotx += 5.0
		}
		if keys[sdl.K_DOWN] != 0 {
			view_rotx -= 5.0
		}
		if keys[sdl.K_LEFT] != 0 {
			view_rotz += 5.0
		}
		if keys[sdl.K_RIGHT] != 0 {
			view_rotz -= 5.0
		}
		if keys[sdl.K_w] != 0 {
      view_z -= 1.0
    }
		if keys[sdl.K_a] != 0 {
      view_x += 1.0
    }
		if keys[sdl.K_s] != 0 {
      view_z += 1.0
    }
		if keys[sdl.K_d] != 0 {
      view_x -= 1.0
    }

		draw()
	}
	sdl.Quit()
	return

}
