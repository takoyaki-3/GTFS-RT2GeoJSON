package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/takoyaki-3/gtfsrt/pkg"
	"github.com/takoyaki-3/goc"
)

func SumFloatTable(df [][]float64)float64{
	sum := 0.0
	for i:=0;i<len(df)-1;i++ {
		r := df[i]
		l := df[i+1]
		p1 := pkg.LatLon{r[0],r[1]}
		p2 := pkg.LatLon{l[0],l[1]}
		
		sum += pkg.HubenyDistance(p1,p2)
	}
	return sum
}

type GTFSandVPOS struct {
	GTFSrecords [][]string
	GTFStitles map[string]int
	Vpos [][]float64
}

// データを保存する出力
type OutLine struct{
	FromStop string
	ToStop string
	OccupancyStatus int
	TimeStamp int64
}

func TripAnalysis(trip GTFSandVPOS)[]OutLine{

	Ans := []OutLine{}

	// 元データの並び替え
	departureTimeCol := trip.GTFStitles["departure_time"]
	sort.Slice(trip.GTFSrecords,func(i,j int) bool { return trip.GTFSrecords[i][departureTimeCol] < trip.GTFSrecords[j][departureTimeCol] })
	sort.Slice(trip.Vpos,func(i,j int) bool { return trip.Vpos[i][2] < trip.Vpos[j][2] })

	// 緯度経度情報のデータフレーム作成
	df1 := [][]float64{}
	df2 := trip.Vpos

	for _,row := range trip.GTFSrecords{
		lat,_ := strconv.ParseFloat(row[trip.GTFStitles["stop_lat"]], 64)
		lon,_ := strconv.ParseFloat(row[trip.GTFStitles["stop_lon"]], 64)
		df1 = append(df1, []float64{lat,lon})
	}

	// 総走行距離を算出
	sum1 := SumFloatTable(df1)
	sum2 := SumFloatTable(df2)

	// VposがないTripを排除
	if len(df2) == 0{
		return Ans
	}

	// マッチング
	progress1 := 0.0

	df1LatLon := pkg.LatLon{-1.0,-1.0}

	lastKn := -1
	for i,r := range df1 {
		minK  := math.MaxFloat64
		minC  := math.MaxFloat64
		minKn := -1

		p1 := pkg.LatLon{r[0],r[1]}
		if 0 < df1LatLon.Lat && 0 < df1LatLon.Lon{
			progress1 += pkg.HubenyDistance(df1LatLon,p1)
		}
		df1LatLon = p1

		df2LatLon := pkg.LatLon{-1.0,-1.0}
		progress2 := 0.0
		for n,l := range df2 {
			p2 := pkg.LatLon{l[0],l[1]}
			if 0 < df2LatLon.Lat && 0 < df2LatLon.Lon {
				progress2 += pkg.HubenyDistance(df2LatLon,p2)
			}
			df2LatLon = p2
			k := pkg.HubenyDistance(p1,p2)
			c := k * 100 + math.Pow((progress1/sum1 - progress2/sum2)*100,4.0)

			if c < minC {
				minK = k
				minC = c
				minKn = n
			}
		}
		if false {
			fmt.Println("bus_stop",i,"is	",int(minK),"m		",minKn)
		}
		if lastKn > minKn {
			fmt.Println("error!!!!")
		}
		if i>0{
			if lastKn<0{
				lastKn = 0
			}
			ind := (lastKn+minKn)/2
			if len(df2) <= ind {
				fmt.Println("???",ind,len(df2))
			}
			Ans = append(Ans, OutLine{
				FromStop				: trip.GTFSrecords[i-1][trip.GTFStitles["stop_id"]],
				ToStop					: trip.GTFSrecords[i][trip.GTFStitles["stop_id"]],
				OccupancyStatus	: int(df2[ind][3]),
				TimeStamp				: int64(df2[ind][2]),
			})
		}
		lastKn = minKn
	}
	return Ans
}

func main(){

	trips := map[string]*GTFSandVPOS{}

	titles,records:=goc.ReadCSV("./df.csv")

	for _,line:=range records{
		tripId := line[titles["trip_id"]]
		if _,ok:=trips[tripId];!ok{
			t := GTFSandVPOS{}
			t.GTFStitles = titles
			t.GTFSrecords = [][]string{}
			trips[tripId] = &t
		}
		trips[tripId].GTFSrecords = append(trips[tripId].GTFSrecords, line)
	}

	// GTFS-RTの読み込み
	titles,records = goc.ReadCSV("./vps.csv")
	
	for _,line:=range records{
		tripId := line[titles["trip_id"]]
		if _,ok:=trips[tripId];!ok{
			fmt.Println(tripId,"not found!")
			continue
		}
		lat,_ := strconv.ParseFloat(line[titles["lat"]], 64)
		lon,_ := strconv.ParseFloat(line[titles["lon"]], 64)
		ostate,_ := strconv.ParseFloat(line[titles["occupancy_status"]], 64)
		timestamp,_ := strconv.ParseFloat(line[titles["timestamp"]], 64)
		trips[tripId].Vpos = append(trips[tripId].Vpos, []float64{lat,lon,timestamp,ostate})
	}

	df := [][]string{}
	df = append(df, []string{"trip_id","from_stop","to_stop","occupancy_status"})
	for tripId,trip:=range trips{
		fmt.Println("now prosessing ",tripId)
		ol := TripAnalysis(*trip)
		for _,v:=range ol{
			df = append(df, []string{
				tripId,
				v.FromStop,
				v.ToStop,
				strconv.Itoa(v.OccupancyStatus),
				strconv.Itoa(int(v.TimeStamp)),
			})
		}
	}
	goc.Write2DStr("edge.csv",df)
}
