package zoom_pixels

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"reflect"
)

var width, height int = 200, 200

const (
	Red int = iota << 1
	Green
	Blue
	AlphaValue
)

// type RawImage struct {
// 	r, g, b, a uint32
// }

var decodeRawImage image.Image
var deltaPixels Pixel_Diff
var deltaPixelsOps Pixel_Diff
var deltaPixelsAddOps Pixel_Diff

var zoom []Pixel_Diff

func Zoom_KTime(value int, file *os.File) {

	zoom = make([]Pixel_Diff, GetImage().Bounds().Max.X*3)

	for i := 0; i < GetImage().Bounds().Max.X; i++ {
		for j := 0; j < GetImage().Bounds().Max.Y; j++ {

			list := Pixels_Info{
				Value: GetImage().At(i, j),
			}

			another_list := Pixels_Info{
				Value: GetImage().At(i, j+1),
			}

			colorUnit_r, colorUnit_g, colorUnit_b, colorUnit_a := list.Value.RGBA()
			colorUnit_r0, colorUnit_g0, colorUnit_b0, colorUnit_a0 := another_list.Value.RGBA()

			if reflect.DeepEqual(colorUnit_r, uint32(0)) && reflect.DeepEqual(colorUnit_g, uint32(0)) && reflect.DeepEqual(colorUnit_b, uint32(0)) && reflect.DeepEqual(colorUnit_a, uint32(0)) {
				continue
			}

			if reflect.DeepEqual(colorUnit_a0, uint32(0)) && reflect.DeepEqual(colorUnit_r0, uint32(0)) && reflect.DeepEqual(colorUnit_g0, uint32(0)) && reflect.DeepEqual(colorUnit_b0, uint32(0)) {
				continue
			}

			deltaPixels = Substract(colorUnit_r, colorUnit_g, colorUnit_b, colorUnit_a, colorUnit_r0, colorUnit_g0, colorUnit_b0, colorUnit_a0)
			k = value

			deltaPixelsOps = Division(deltaPixels)
			deltaPixelsAddOps = Add(deltaPixels)

			log.Println("Pixels Difference", deltaPixels, "Divison:", deltaPixelsOps, "Add:", deltaPixelsAddOps)

			p, q, r := NewImage(deltaPixels, deltaPixelsOps, deltaPixelsAddOps, i)

			u, v := Is_Sort(p, q, r)

			zoom = append(zoom, u, v)
		}

	}

	newPicture := copy_pixels()
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err := encoder.Encode(file, newPicture)
	if err != nil {
		log.Fatalln("picture encode error:", err)
		return
	}
}

var k int = 0

var img *image.RGBA

type Pixels_Info struct {
	Value color.Color
}

type Pixel_Diff struct {
	r, g, b, a uint32
}

func Substract(r, g, b, a uint32, r0, g0, b0, a0 uint32) Pixel_Diff {

	sub_r := r - r0
	sub_g := g - g0
	sub_b := b - b0
	sub_a := a - a0
	if r > r0 && g > g0 && b > b0 && a > a0 {
		return Pixel_Diff{r: sub_r, g: sub_g, b: sub_b, a: sub_a}
	}

	sub_r = r0 - r
	sub_g = g0 - g
	sub_b = b0 - b
	sub_a = a0 - a
	return Pixel_Diff{r: sub_r, g: sub_g, b: sub_b, a: sub_a}
}

func Division(p Pixel_Diff) Pixel_Diff {

	if k == 0 {
		log.Fatalln("Divison operation is not allowed")
		return Pixel_Diff{}
	}

	dr := int(p.r) / k
	dg := int(p.g) / k
	db := int(p.b) / k
	da := int(p.a) / k

	return Pixel_Diff{r: uint32(dr), g: uint32(dg), b: uint32(db), a: uint32(da)}

}

func Add(p Pixel_Diff) Pixel_Diff {

	da := int(p.r) + k
	dg := int(p.g) + k
	db := int(p.b) + k
	dr := int(p.r) + k

	return Pixel_Diff{r: uint32(dr), g: uint32(dg), b: uint32(db), a: uint32(da)}
}

// var avatar_zoom *image.Paletted

func copy_pixels() *image.Paletted {

	var pictureColor []color.Color

	for i := range zoom {
		pictureColor = []color.Color{
			color.RGBA64{uint16(zoom[i].r), uint16(zoom[i].g), uint16(zoom[i].b), uint16(zoom[i].a)},
		}
	}

	return image.NewPaletted(image.Rect(0, 0, width*k, height*k), pictureColor)
}

func SetImage(im image.Image) { decodeRawImage = im }

func GetImage() image.Image { return decodeRawImage }

func NewImage(p, q, r Pixel_Diff, i int) (Pixel_Diff, Pixel_Diff, Pixel_Diff) {

	return p, q, r
}

func shuffle(s Pixel_Diff, t Pixel_Diff) (Pixel_Diff, Pixel_Diff) {

	temp := Pixel_Diff{}

	temp.r, temp.g, temp.b, temp.a = s.r, s.g, s.b, s.a
	s.r, s.g, s.b, s.a = t.r, t.g, t.b, t.a
	t.r, t.g, t.b, t.a = temp.r, temp.g, temp.b, temp.a

	return t, s
}

func Is_Sort(p, q, r Pixel_Diff) (Pixel_Diff, Pixel_Diff) {

	var u, v Pixel_Diff

	log.Println("Q:", q, " P:", p)
	if q.r > p.r && q.g > p.g && q.b > p.b && q.a > p.a {
		u, v = shuffle(p, q)
		log.Println("Pixels sorts between p & q:", u, v)

		return u, v
	}

	log.Println("R:", r, " P:", p)
	if r.r > p.r && r.g > p.g && r.b > p.b && r.a > p.a {
		u, v = shuffle(p, r)
		log.Println("Pixels sorts between r & p:", u, v)
		return u, v
	}

	log.Println("Q:", q, " R:", r)
	if r.r > q.r && r.g > q.g && r.b > q.b && r.a > q.a {
		u, v = shuffle(q, r)
		log.Println("Pixels sorts between q & r:", u, v)
		return u, v
	}

	return Pixel_Diff{}, Pixel_Diff{}
}
