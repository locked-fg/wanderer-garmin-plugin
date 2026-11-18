package garmintrack

import (
	"encoding/xml"
)

const gpxNamespace = "http://www.topografix.com/GPX/1/1"

type gpxDocument struct {
	XMLName xml.Name `xml:"gpx"`
	Version string   `xml:"version,attr"`
	Creator string   `xml:"creator,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Track   gpxTrack `xml:"trk"`
}

type gpxTrack struct {
	Name string     `xml:"name"`
	Seg  gpxSegment `xml:"trkseg"`
}

type gpxSegment struct {
	Points []gpxPoint `xml:"trkpt"`
}

type gpxPoint struct {
	Lat  float64  `xml:"lat,attr"`
	Lon  float64  `xml:"lon,attr"`
	Ele  *float64 `xml:"ele,omitempty"`
	Time string   `xml:"time,omitempty"`
}

// BuildGPX renders track points as a GPX 1.1 document. Element content is
// escaped by encoding/xml, so malformed timestamps cannot break the XML.
func BuildGPX(points []Point) []byte {
	if len(points) == 0 {
		return nil
	}

	doc := gpxDocument{
		Version: "1.1",
		Creator: "Garmin Personal API",
		Xmlns:   gpxNamespace,
		Track: gpxTrack{
			Name: "Activity Track",
			Seg:  gpxSegment{Points: make([]gpxPoint, 0, len(points))},
		},
	}
	for _, point := range points {
		gp := gpxPoint{Lat: point.Latitude, Lon: point.Longitude, Time: point.Timestamp}
		if point.Elevation != 0 {
			ele := point.Elevation
			gp.Ele = &ele
		}
		doc.Track.Seg.Points = append(doc.Track.Seg.Points, gp)
	}

	out, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil
	}
	return append([]byte(xml.Header), append(out, '\n')...)
}
