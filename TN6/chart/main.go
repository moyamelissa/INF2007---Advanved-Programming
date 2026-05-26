// chart génère le graphique d'accélération du crawler pour le rapport TN6.
// Exécuter depuis ce dossier : go run .
// Sortie : ../docs/benchmark-chart.png
package main

import (
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Données benchstat (medians, count=1, politenessDelay=0, mai 2026)
var goroutines = []float64{0, 1, 2, 3} // positions X (1, 2, 4, 8 goroutines)
var labels = []string{"1", "2", "4", "8"}

var singleMS = []float64{2.87, 1.67, 5.91, 11.53}
var multiMS = []float64{3.02, 1.52, 0.92, 0.86}

// Palette mauve cohérente avec TN4
var (
	colorSingle = color.RGBA{R: 106, G: 62, B: 160, A: 255}  // violet foncé — serveur unique
	colorMulti  = color.RGBA{R: 160, G: 132, B: 210, A: 255} // lavande — multi-serveurs
	colorLabel  = color.RGBA{R: 60, G: 60, B: 80, A: 220}    // gris ardoise pour annotations
)

func xyPoints(xs, ys []float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i := range xs {
		pts[i].X = xs[i]
		pts[i].Y = ys[i]
	}
	return pts
}

func main() {
	if err := os.MkdirAll("../data", 0755); err != nil {
		log.Fatal(err)
	}

	p := plot.New()

	p.X.Label.Text = "Goroutines"
	p.X.Label.TextStyle.Font.Size = vg.Points(10)
	p.Y.Label.Text = "Temps d'exécution (ms)"
	p.Y.Label.TextStyle.Font.Size = vg.Points(10)

	p.Y.Min = 0
	p.Y.Max = 14
	p.X.Min = -0.3
	p.X.Max = 3.5

	// Ticks X : 1, 2, 4, 8
	p.X.Tick.Marker = plot.ConstantTicks(func() []plot.Tick {
		ticks := make([]plot.Tick, len(labels))
		for i, lbl := range labels {
			ticks[i] = plot.Tick{Value: float64(i), Label: lbl}
		}
		return ticks
	}())

	// Ticks Y
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 0, Label: "0"}, {Value: 2, Label: "2"},
		{Value: 4, Label: "4"}, {Value: 6, Label: "6"},
		{Value: 8, Label: "8"}, {Value: 10, Label: "10"},
		{Value: 12, Label: "12"}, {Value: 14, Label: "14"},
	})

	// ── Courbe serveur unique ──
	singleLine, err := plotter.NewLine(xyPoints(goroutines, singleMS))
	if err != nil {
		log.Fatal(err)
	}
	singleLine.Color = colorSingle
	singleLine.Width = vg.Points(2.5)
	singleLine.Dashes = []vg.Length{vg.Points(6), vg.Points(3)}

	singleScatter, err := plotter.NewScatter(xyPoints(goroutines, singleMS))
	if err != nil {
		log.Fatal(err)
	}
	singleScatter.GlyphStyle = draw.GlyphStyle{
		Color:  color.White,
		Radius: vg.Points(5),
		Shape:  draw.CircleGlyph{},
	}

	singleDot, err := plotter.NewScatter(xyPoints(goroutines, singleMS))
	if err != nil {
		log.Fatal(err)
	}
	singleDot.GlyphStyle = draw.GlyphStyle{
		Color:  colorSingle,
		Radius: vg.Points(3),
		Shape:  draw.CircleGlyph{},
	}

	// ── Courbe multi-serveurs ──
	multiLine, err := plotter.NewLine(xyPoints(goroutines, multiMS))
	if err != nil {
		log.Fatal(err)
	}
	multiLine.Color = colorMulti
	multiLine.Width = vg.Points(2.5)

	multiScatter, err := plotter.NewScatter(xyPoints(goroutines, multiMS))
	if err != nil {
		log.Fatal(err)
	}
	multiScatter.GlyphStyle = draw.GlyphStyle{
		Color:  color.White,
		Radius: vg.Points(5),
		Shape:  draw.SquareGlyph{},
	}

	multiDot, err := plotter.NewScatter(xyPoints(goroutines, multiMS))
	if err != nil {
		log.Fatal(err)
	}
	multiDot.GlyphStyle = draw.GlyphStyle{
		Color:  colorMulti,
		Radius: vg.Points(3),
		Shape:  draw.SquareGlyph{},
	}

	p.Add(singleLine, singleScatter, singleDot)
	p.Add(multiLine, multiScatter, multiDot)

	// Légende
	p.Legend.Add("Serveur unique", singleLine, singleDot)
	p.Legend.Add("Multi-serveurs (8 serveurs)", multiLine, multiDot)
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.TextStyle.Font.Size = vg.Points(9)
	p.Legend.Padding = vg.Points(4)

	if err := p.Save(16*vg.Centimeter, 10*vg.Centimeter, "../data/benchmark-chart.png"); err != nil {
		log.Fatal(err)
	}
	log.Println("Graphique généré : ../data/benchmark-chart.png")
}
