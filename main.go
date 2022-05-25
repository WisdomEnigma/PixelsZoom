package main

import (
	"errors"
	"fmt"
	"image/png"
	"os"

	"github.com/WisdomEnigma/PixelsZoom/zoom_pixels"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	app_web := fiber.New(fiber.Config{
		Views: html.New("./views", ".hbs"),
	})

	app_web.Get("/", func(c *fiber.Ctx) error {

		// page rendered
		return c.Render("index", fiber.Map{
			"Title": "PixelsMetrica",
		})

	})

	app_web.Post("/", func(c *fiber.Ctx) error {

		file, err := c.FormFile("image")
		if err != nil {

			code := fiber.StatusInternalServerError

			if err, ok := err.(*fiber.Error); !ok {

				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": errors.New("file might be corrupted"),
			})
		}

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
				"message": errors.New("file might be corrupted"),
			})
		}

		FileInfo, err := os.OpenFile(file.Filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if !ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": errors.New("file is locked"),
			})

		}

		defer FileInfo.Close()

		_, err = os.Stat(FileInfo.Name())
		if os.IsExist(err) {

			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": errors.New("file properties not provided"),
			})
		}

		_content, err := png.Decode(FileInfo)
		if err != nil {
			code := fiber.StatusInternalServerError

			err, ok := err.(*fiber.Error)
			if ok {
				code = err.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(code).JSON(map[string]interface{}{
				"code":    code,
				"message": errors.New("file properties not provided"),
			})
		}

		zoom_pixels.SetImage(_content)
		zoom_pixels.Zoom_KTime(20, FileInfo)

		return c.Render("index", fiber.Map{
			"Title": "PixelsMetrica",
		})
	})

	err := app_web.Listen(":3000")
	if err != nil {
		panic(err)
	}

}