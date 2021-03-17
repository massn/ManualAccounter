package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type SeriesData struct {
	Name      string
	ChartData []PointData
}

type PointData struct {
	Date  string
	Value float64
}

func Render(gain, valuation SeriesData) error {
	page := components.NewPage()
	page.AddCharts(
		lineMulti(gain),
		lineMulti(valuation),
	)
	f, err := os.Create("line.html")
	if err != nil {
		return err
	}
	page.Render(io.MultiWriter(f))
	return nil
}

func lineMulti(sd SeriesData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: sd.Name,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "shine",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
	)

	xAxis := []string{}
	data := []opts.LineData{}
	for _, pointData := range sd.ChartData {
		data = append(data, opts.LineData{Value: pointData.Value})
		xAxis = append(xAxis, pointData.Date)
	}

	line.SetXAxis(xAxis).
		AddSeries(sd.Name, data)

	return line
}
