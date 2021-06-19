package pkg

import (
	"io"
	"os"
	"io/ioutil"
	"path/filepath"
	"archive/zip"
)

type GTFSRT_Table struct {
	TimeStamp				int64		`parquet:"name=timestamp,				type=INT64"`
	Lat							float32	`parquet:"name=lat,							type=FLOAT"`
	Lon							float32	`parquet:"name=lon, 						type=FLOAT"`
	TripId					string	`parquet:"name=trip_id, 				type=BYTE_ARRAY,	convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	RouteId					string	`parquet:"name=route_id, 				type=BYTE_ARRAY,	convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	OccupancyStatus	int32		`parquet:"name=occupancy_status,type=INT32"`
}

type LineStringGeoJSON struct {
	Type     string			`json:"type"`
	Features []Feature	`json:"features"`
}

type Feature struct{
	Type     string						`json:"type"`
	Geometry struct {
		Type        string			`json:"type"`
		Coordinates [][]float64	`json:"coordinates"`
	}													`json:"geometry"`
	Properties Property				`json:"properties"`
}

type Property struct{
	OccupancyStatus int		`json:"OccupancyStatus"`
}

func Dirwalk(dir string) ([]string, []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths, file_names []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
		file_names = append(file_names, file.Name())
	}
	return paths, file_names
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}