// chart génère le graphique Int vs Float pour le rapport TN4.
// Exécuter depuis ce dossier : go run .
// Sortie : ../docs/benchmark-chart.png
package main

import (
	"fmt"
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Données des benchmarks (médianes benchstat, 6 exécutions)
var labels = []string{"1%", "10%", "20%", "30%", "40%", "50%", "60%", "70%", "80%", "90%", "100%"}
var intMS = []float64{0.44, 4.09, 8.11, 11.83, 15.52, 19.28, 22.98, 26.58, 30.94, 34.78, 38.71}
var floatMS = []float64{0.24, 2.11, 4.24, 7.79, 8.99, 11.98, 13.61, 14.69, 16.82, 18.96, 20.93}

// Palette mauve (cohérente avec le rapport)
var (
	colorInt      = color.RGBA{R: 106, G: 62, B: 160, A: 255}  // violet foncé
	colorFloat    = color.RGBA{R: 160, G: 132, B: 210, A: 255} // lavande moyenne
	colorFill     = color.NRGBA{R: 185, G: 165, B: 225, A: 75} // fill translucide (non-premultiplied)
	colorAnnot    = color.RGBA{R: 80, G: 30, B: 140, A: 255}   // violet annotation
	colorEndLabel = color.RGBA{R: 60, G: 60, B: 80, A: 220}    // gris ardoise pour valeurs finales
)

func xyPoints(xs, ys []float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i := range xs {
		pts[i].X = xs[i]
		pts[i].Y = ys[i]
	}
	return pts
}

// fillBetween construit un polygone fermé entre deux courbes.
func fillBetween(xs, top, bottom []float64) plotter.XYs {
	n := len(xs)
	pts := make(plotter.XYs, 2*n+1)
	for i := 0; i < n; i++ {
		pts[i] = plotter.XY{X: xs[i], Y: top[i]}
	}
	for i := 0; i < n; i++ {
		pts[n+i] = plotter.XY{X: xs[n-1-i], Y: bottom[n-1-i]}
	}
	pts[2*n] = pts[0]
	return pts
}

func main() {
	n := len(labels)
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = float64(i)
	}

	p := plot.New()

	// Axes (titre géré par le rapport, pas dans le PNG)
	p.X.Label.Text = "Pourcentage du tableau (%)"
	p.X.Label.TextStyle.Font.Size = vg.Points(10)
	p.Y.Label.Text = "Temps d'exécution (ms)"
	p.Y.Label.TextStyle.Font.Size = vg.Points(10)

	// Plage des axes — marge à droite pour les labels de valeur
	p.Y.Min = 0
	p.Y.Max = 45
	p.X.Min = -0.4
	p.X.Max = float64(n) + 0.6 // espace pour "38.71 ms"

	// Ticks Y explicites : multiples de 10
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 0, Label: "0"}, {Value: 5, Label: "5"},
		{Value: 10, Label: "10"}, {Value: 15, Label: "15"},
		{Value: 20, Label: "20"}, {Value: 25, Label: "25"},
		{Value: 30, Label: "30"}, {Value: 35, Label: "35"},
		{Value: 40, Label: "40"},
	})

	// Ticks X avec les labels du tableau
	p.X.Tick.Marker = plot.ConstantTicks(func() []plot.Tick {
		ticks := make([]plot.Tick, n)
		for i, lbl := range labels {
			ticks[i] = plot.Tick{Value: float64(i), Label: lbl}
		}
		return ticks
	}())

	// Pas de grille — fond propre comme dans le document Word

	// Zone translucide entre les courbes
	poly, err := plotter.NewPolygon(fillBetween(xs, intMS, floatMS))
	if err != nil {
		log.Fatal(err)
	}
	poly.Color = colorFill
	poly.LineStyle.Width = 0
	p.Add(poly)

	// ── Courbe Float (bas) ──
	floatLine, err := plotter.NewLine(xyPoints(xs, floatMS))
	if err != nil {
		log.Fatal(err)
	}
	floatLine.Color = colorFloat
	floatLine.Width = vg.Points(2.5)

	floatScatter, err := plotter.NewScatter(xyPoints(xs, floatMS))
	if err != nil {
		log.Fatal(err)
	}
	floatScatter.GlyphStyle = draw.GlyphStyle{
		Color:  color.White,
		Radius: vg.Points(4),
		Shape:  draw.CircleGlyph{},
	}
	floatScatterFill, err := plotter.NewScatter(xyPoints(xs, floatMS))
	if err != nil {
		log.Fatal(err)
	}
	floatScatterFill.GlyphStyle = draw.GlyphStyle{
		Color:  colorFloat,
		Radius: vg.Points(3),
		Shape:  draw.CircleGlyph{},
	}

	// ── Courbe Int (haut) ──
	intLine, err := plotter.NewLine(xyPoints(xs, intMS))
	if err != nil {
		log.Fatal(err)
	}
	intLine.Color = colorInt
	intLine.Width = vg.Points(2.5)

	intScatter, err := plotter.NewScatter(xyPoints(xs, intMS))
	if err != nil {
		log.Fatal(err)
	}
	intScatter.GlyphStyle = draw.GlyphStyle{
		Color:  color.White,
		Radius: vg.Points(4),
		Shape:  draw.CircleGlyph{},
	}
	intScatterFill, err := plotter.NewScatter(xyPoints(xs, intMS))
	if err != nil {
		log.Fatal(err)
	}
	intScatterFill.GlyphStyle = draw.GlyphStyle{
		Color:  colorInt,
		Radius: vg.Points(3),
		Shape:  draw.CircleGlyph{},
	}

	p.Add(floatLine, floatScatter, floatScatterFill,
		intLine, intScatter, intScatterFill)

	// ── Légende ──
	p.Legend.Add("Int (avec conversion float64)", intLine, intScatterFill)
	p.Legend.Add("Float (sans conversion)", floatLine, floatScatterFill)
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.Padding = vg.Points(4)

	// ── Valeurs finales à 100% ──
	endLabels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: plotter.XYs{
			{X: float64(n-1) + 0.12, Y: intMS[n-1] + 0.8},
			{X: float64(n-1) + 0.12, Y: floatMS[n-1] - 1.5},
		},
		Labels: []string{
			fmt.Sprintf("%.2f ms", intMS[n-1]),
			fmt.Sprintf("%.2f ms", floatMS[n-1]),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	endLabels.TextStyle[0].Color = colorEndLabel
	endLabels.TextStyle[0].Font.Size = vg.Points(8.5)
	endLabels.TextStyle[1].Color = colorEndLabel
	endLabels.TextStyle[1].Font.Size = vg.Points(8.5)
	p.Add(endLabels)

	// ── Annotation ratio au milieu de l'écart à 80% (avant le bout pour éviter le chevauchement) ──
	ratio := intMS[n-1] / floatMS[n-1]
	annotX := float64(n-1) - 2.3 // ~80% du tableau
	annotIdx := n - 3
	midY := (intMS[annotIdx] + floatMS[annotIdx]) / 2
	ratioLabel, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    plotter.XYs{{X: annotX, Y: midY}},
		Labels: []string{fmt.Sprintf("écart ×%.2f", ratio)},
	})
	if err != nil {
		log.Fatal(err)
	}
	ratioLabel.TextStyle[0].Color = colorAnnot
	ratioLabel.TextStyle[0].Font.Size = vg.Points(9)
	p.Add(ratioLabel)

	out := "../docs/benchmark-chart.png"
	if err := p.Save(28*vg.Centimeter, 15*vg.Centimeter, out); err != nil {
		log.Fatal(err)
	}
	log.Printf("Graphique sauvegardé : %s", out)
}
