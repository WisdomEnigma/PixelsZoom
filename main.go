package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	direc "github.com/WisdomEnigma/PixelsZoom/dir"
	"github.com/WisdomEnigma/PixelsZoom/zoom_pixels"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/nfnt/resize"
)

var mountDir = os.Getenv("Mounted_Workspace")
var port = os.Getenv("PORT")

func main() {

	if !strings.Contains(mountDir, "") && !reflect.DeepEqual(mountDir, "3000") {

		log.Fatalln("No working directory specified")
		panic("No working directory specified")
	}

	app_web := fiber.New(fiber.Config{
		Views: html.New("./views", ".hbs"),
	})

	app_web.Get("/", func(c *fiber.Ctx) error {

		// page rendered
		return c.Render("index", fiber.Map{
			"Title": "PixelsMetrica",
		})

	})

	// Load Image from hbs Form
	app_web.Post("/", func(c *fiber.Ctx) error {

		// if file have some issue then throw an exception called StatusInternalServerError; along with error message and error code
		file, err := c.FormFile("image")
		if err != nil {

			code := fiber.StatusInternalServerError

			if err, ok := err.(*fiber.Error); !ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": "file might be corrupted",
			})
		}

		// Save image in local.. again if any issue the n throw exception

		err = c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
		if err != nil {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if !ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": "file might be corrupted",
			})
		}

		// // Create new file or open image file
		// FileInfo, err := os.OpenFile(file.Filename, os.O_RDWR|os.O_CREATE, 0755)
		mdir, err := direc.Chdir(file.Filename)
		if err != nil {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if !ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": "file don't have credentials",
			})

		}

		log.Println("File:", mdir.Name())
		// // close the file descriptor
		// defer mdir.Close()

		// check whelther file exists
		_, err = os.Stat(mdir.Name())
		if os.IsExist(err) {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": "file credentials are not valid",
			})
		}

		read, err := ioutil.ReadFile(file.Filename)
		if err != nil {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": " Read operation failed",
			})

		}

		err = ioutil.WriteFile(mdir.Name(), read, 0644)

		if err != nil {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": " Write operation failed",
			})

		}

		decode, err := png.Decode(mdir)

		if err != nil {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": " decode operation failed",
			})

		}

		// Set Image allow to you to set image pixel values
		zoom_pixels.SetImage(decode)

		picture := resize.Resize(50, 50, decode, resize.Lanczos3)

		png.Encode(mdir, picture)

		// Zoom K Times function have two params to scale up shared content; Zoom Ktime return error if scaling is not supported
		// _newImage, err := zoom_pixels.Zoom_KTime(2, mdir)
		// if err != nil {

		// 	code := fiber.StatusInternalServerError

		// 	err, ok := err.(*fiber.Error)
		// 	if ok {
		// 		code = err.Code
		// 	}

		// 	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		// 	return c.Status(code).JSON(map[string]interface{}{
		// 		"code":    code,
		// 		"message": "scale up resolution corrupted",
		// 	})
		// }

		// err = png.Encode(mdir, _newImage)
		// if err != nil {

		// 	code := fiber.StatusInternalServerError

		// 	err, ok := err.(*fiber.Error)
		// 	if ok {
		// 		code = err.Code
		// 	}

		// 	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		// 	return c.Status(code).JSON(map[string]interface{}{
		// 		"code":    code,
		// 		"message": "scale up resolution failed",
		// 	})
		// }
		// Zoom out pixels reverse process of the image
		//zoom_pixels.ZoomOutPixels(FileInfo, 5)

		return c.Render("index", fiber.Map{
			"Title": "PixelsMetrica",
		})

	})

	err := app_web.Listen(":" + port)
	if err != nil {
		panic(err)
	}

}
