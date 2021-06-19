package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"sort"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	pb "github.com/takoyaki-3/GTFS-RT2GeoJSON/pb"

	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"

	"github.com/takoyaki-3/GTFS-RT2GeoJSON/pkg"
)

func main(){
	fmt.Println("start")

	rtTable := []pkg.GTFSRT_Table{}
	tripVpos := map[string][]pb.VehiclePosition{}

	// GTFS-RT ファイル一覧を取得
	paths ,_:= pkg.Dirwalk("GTFS-RTs")

	for k,v:=range paths{

		// .gitignoreファイル除外
		if strings.HasSuffix(v, ".gitignore"){ continue }

		// 途中経過を1000ファイルごとに出力
		if k%1000==0{ fmt.Println(v) }

		rawBytes, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Println("error")
			continue
		}

		feed := pb.FeedMessage{}
		err = proto.Unmarshal(rawBytes,&feed)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 位置情報取得
		for _,e:=range feed.Entity{
			vc := e.Vehicle
			if vc.Position == nil || vc.Trip == nil {
				continue
			}

			// TripID 毎に分割して保存
			tripVpos[*vc.Trip.TripId] = append(tripVpos[*vc.Trip.TripId],*vc)
			
			// parquet 出力用
			p := vc.Position
			rtTable = append(rtTable, pkg.GTFSRT_Table{
				TimeStamp				: int64(*vc.Timestamp),
				TripId					: *vc.Trip.TripId,
				Lat							:	*p.Latitude,
				Lon							: *p.Longitude,
				OccupancyStatus	: int32(*vc.OccupancyStatus),
			})
		}
	}

	// 並列実行用
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// parquet 出力処理
	go func(){
		defer wg.Done()

		w, err := os.Create("GTFS-RT.parquet")
		if err != nil {
			log.Println("Can't create local file", err)
			return
		}
		pw, err := writer.NewParquetWriterFromWriter(w, new(pkg.GTFSRT_Table), 4)
		if err != nil {
			log.Println("Can't create parquet writer", err)
			return
		}
		pw.RowGroupSize = 128 * 1024 * 1024
		pw.CompressionType = parquet.CompressionCodec_SNAPPY
		for _,v := range rtTable{
			if err = pw.Write(v); err != nil {
				log.Println("Write error", err)
			}
		}
		if err = pw.WriteStop(); err != nil {
			log.Println("WriteStop error", err)
			return
		}
		log.Println("Write Finished")
		w.Close()
	}()

	//////////////////////////////////////////////////////////////////////////////////

	go func(){
		defer wg.Done()

		// 混雑度置換
		OccupancyStatus := map[string]int{}
		OccupancyStatus["EMPTY"] 											= 0
		OccupancyStatus["MANY_SEATS_AVAILABLE"] 			= 1
		OccupancyStatus["FEW_SEATS_AVAILABLE"] 				= 2
		OccupancyStatus["STANDING_ROOM_ONLY"] 				= 3
		OccupancyStatus["CRUSHED_STANDING_ROOM_ONLY"] = 4
		OccupancyStatus["FULL"] 											= 5
		OccupancyStatus["NOT_ACCEPTING_PASSENGERS"] 	= 6

		// GeoJSON出力用
		fc := pkg.LineStringGeoJSON{}
		fc.Type = "FeatureCollection"

		// 時刻順に並び替え
		for _,v:=range tripVpos{
			f := pkg.Feature{}
			f.Type = "Feature"
			f.Geometry.Type = "LineString"
		
			sort.SliceStable(v, func(i,j int) bool { return *v[i].Timestamp < *v[j].Timestamp })

			last := ""
			for _,v:=range v{
				p := v.Position

				f.Geometry.Coordinates = append(f.Geometry.Coordinates, []float64{float64(*p.Longitude),float64(*p.Latitude),float64(0.0),float64(*v.Timestamp)})
				if last != v.OccupancyStatus.String() && last != ""{
					f.Properties.OccupancyStatus = OccupancyStatus[last]
					fc.Features = append(fc.Features, f)
					f = pkg.Feature{}
					f.Type = "Feature"
					f.Geometry.Type = "LineString"
					f.Geometry.Coordinates = append(f.Geometry.Coordinates, []float64{float64(*p.Longitude),float64(*p.Latitude),float64(0.0),float64(*v.Timestamp)})
				}
				last = v.OccupancyStatus.String()
			}

			f.Properties.OccupancyStatus = OccupancyStatus[last]
			fc.Features = append(fc.Features, f)
		}

		// JSON 出力
		file, err := os.Create("GTFS-RT.json")
		defer file.Close()
		if err != nil{
			fmt.Println(err)
		}
		rawJSON, _ := json.Marshal(&fc)
		file.Write(rawJSON)
	}()

	wg.Wait()
}