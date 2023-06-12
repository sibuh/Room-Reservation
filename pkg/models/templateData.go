package models

type TemplateData struct {
	MapString map[string]string
	MapInt    map[string]int
	MapFloat  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Error     string
	Warning   string
}
