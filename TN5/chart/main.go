// chart génère le graphique de benchmarks pour le rapport TN5.
// Exécuter depuis ce dossier : go run .
// Sortie : ../data/worker-count-chart.png
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

// Données des benchmarks (médianes benchstat, 6 exécutions — mai 2026)

// BenchmarkWorkerCount
var workerLabels = []string{"1", "2", "4", "8", "16", "32"}
var workerActualMS = []float64{10.39, 9.021, 8.473, 6.926, 6.307, 6.208}

// Palette mauve/violet (cohérente avec TN4)
var (
	colorActual   = color.RGBA{R: 106, G: 62, B: 160, A: 255}  // violet foncé (courbe principale)
	colorIdeal    = color.RGBA{R: 160, G: 132, B: 210, A: 255} // lavande moyenne (courbe idéale)
	colorFill     = color.NRGBA{R: 185, G: 165, B: 225, A: 75} // fill translucide
	colorEndLabel = color.RGBA{R: 60, G: 60, B: 80, A: 220}    // gris ardoise
	colorAnnot    = color.RGBA{R: 80, G: 30, B: 140, A: 255}   // violet annotation
	colorWhiteBg  = color.RGBA{R: 255, G: 255, B: 255, A: 255} // fond blanc pur
)

// xyPoints convertit deux tranches parallèles (abscisses, ordonnées) en points gonum.
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

// buildXS génère une tranche d'entiers [0, 1, …, n-1] en float64 pour l'axe X.
func buildXS(n int) []float64 {
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = float64(i)
	}
	return xs
}

// makeWorkerCountChart trace le temps réel vs le scaling linéaire idéal
// pour BenchmarkWorkerCount (1–32 goroutines, texte 100k mots).
func makeWorkerCountChart() {
	n := len(workerLabels)
	xs := buildXS(n)

	// Scaling idéal : temps[0] / nb_workers
	divisors := []float64{1, 2, 4, 8, 16, 32}
	idealMS := make([]float64, n)
	for i, d := range divisors {
		idealMS[i] = workerActualMS[0] / d
	}

	p := plot.New()
	p.BackgroundColor = colorWhiteBg
	p.X.Label.Text = "Nombre de goroutines"
	p.X.Label.TextStyle.Font.Size = vg.Points(10)
	p.Y.Label.Text = "Temps d'exécution (ms)"
	p.Y.Label.TextStyle.Font.Size = vg.Points(10)

	p.Y.Min = 0
	p.Y.Max = 12
	p.X.Min = -0.4
	p.X.Max = float64(n-1) + 0.6 // espace pour le label final

	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 0, Label: "0"}, {Value: 2, Label: "2"},
		{Value: 4, Label: "4"}, {Value: 6, Label: "6"},
		{Value: 8, Label: "8"}, {Value: 10, Label: "10"},
		{Value: 12, Label: "12"},
	})
	p.X.Tick.Marker = plot.ConstantTicks(func() []plot.Tick {
		ticks := make([]plot.Tick, n)
		for i, lbl := range workerLabels {
			ticks[i] = plot.Tick{Value: float64(i), Label: lbl}
		}
		return ticks
	}())

	// Zone translucide entre courbe réelle (haut) et idéale (bas)
	poly, err := plotter.NewPolygon(fillBetween(xs, workerActualMS, idealMS))
	if err != nil {
		log.Fatal(err)
	}
	poly.Color = colorFill
	poly.LineStyle.Width = 0
	p.Add(poly)

	// ── Courbe idéale (pointillé) ──
	idealLine, err := plotter.NewLine(xyPoints(xs, idealMS))
	if err != nil {
		log.Fatal(err)
	}
	idealLine.Color = colorIdeal
	idealLine.Width = vg.Points(2)
	idealLine.Dashes = []vg.Length{vg.Points(6), vg.Points(3)}

	idealScatter, err := plotter.NewScatter(xyPoints(xs, idealMS))
	if err != nil {
		log.Fatal(err)
	}
	idealScatter.GlyphStyle = draw.GlyphStyle{
		Color:  colorIdeal,
		Radius: vg.Points(3.5),
		Shape:  draw.SquareGlyph{},
	}

	// ── Courbe réelle ──
	actualLine, err := plotter.NewLine(xyPoints(xs, workerActualMS))
	if err != nil {
		log.Fatal(err)
	}
	actualLine.Color = colorActual
	actualLine.Width = vg.Points(2.5)

	actualScatter, err := plotter.NewScatter(xyPoints(xs, workerActualMS))
	if err != nil {
		log.Fatal(err)
	}
	actualScatter.GlyphStyle = draw.GlyphStyle{
		Color:  color.White,
		Radius: vg.Points(4),
		Shape:  draw.CircleGlyph{},
	}
	actualScatterFill, err := plotter.NewScatter(xyPoints(xs, workerActualMS))
	if err != nil {
		log.Fatal(err)
	}
	actualScatterFill.GlyphStyle = draw.GlyphStyle{
		Color:  colorActual,
		Radius: vg.Points(3),
		Shape:  draw.CircleGlyph{},
	}

	p.Add(idealLine, idealScatter,
		actualLine, actualScatter, actualScatterFill)

	p.Legend.Add("Réel", actualLine, actualScatterFill)
	p.Legend.Add("Linéaire idéal", idealLine, idealScatter)
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.Padding = vg.Points(4)

	// Valeur finale (courbe réelle)
	endLabel, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    plotter.XYs{{X: float64(n-1) + 0.12, Y: workerActualMS[n-1] + 0.45}},
		Labels: []string{fmt.Sprintf("%.2f ms", workerActualMS[n-1])},
	})
	if err != nil {
		log.Fatal(err)
	}
	endLabel.TextStyle[0].Color = colorActual
	endLabel.TextStyle[0].Font.Size = vg.Points(8.5)
	p.Add(endLabel)

	// Annotation "plateau ≈ 6 ms" sur la zone de plateau
	plateauLabel, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    plotter.XYs{{X: float64(n-1) - 1.7, Y: workerActualMS[n-1] + 0.85}},
		Labels: []string{"plateau ~ 6 ms"},
	})
	if err != nil {
		log.Fatal(err)
	}
	plateauLabel.TextStyle[0].Color = colorAnnot
	plateauLabel.TextStyle[0].Font.Size = vg.Points(8.5)
	p.Add(plateauLabel)

	out := "../data/worker-count-chart.png"
	if err := p.Save(28*vg.Centimeter, 15*vg.Centimeter, out); err != nil {
		log.Fatal(err)
	}
	log.Printf("Graphique sauvegardé : %s", out)
}

func main() {
	makeWorkerCountChart()
}
