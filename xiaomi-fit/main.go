package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {

	startS := flag.String("start", "", "start time")
	endS := flag.String("end", "", "end time")
	path := flag.String("path", "", "path")
	flag.Parse()

	if *path == "" {
		fmt.Println("no path")
		return
	}
	start, err := time.Parse(time.DateOnly, *startS)
	if err != nil {
		fmt.Println(err)
		return
	}
	end, err := time.Parse(time.DateOnly, *endS)
	if err != nil {
		fmt.Println(err)
		return
	}
	end = end.AddDate(0, 0, 1).Add(-1)

	heartRateDays := []cell[int64]{}
	heartRateHours := []cell[int64]{}
	bloodPreSubDays := []cell[int64]{}
	bloodPreDays := []cell[[]int64]{}
	bloodPreSysDays := []cell[int64]{}
	bloodPreDiaDays := []cell[int64]{}
	singleSpo2Days := []cell[int64]{}
	stressDays := []cell[int64]{}
	sleepDurationDays := []cell[int64]{}
	sleepDeepDurationDays := []cell[int64]{}
	sleepLightDurationDays := []cell[int64]{}
	sleepRemDurationDays := []cell[int64]{}

	for i := start; i.Before(end); i = i.AddDate(0, 0, 1) {
		//for j := 0; j < 48; j++ {
		//	t := i.Add(time.Minute * time.Duration(j) * 30).Unix()
		//	heartRateHours = append(heartRateHours, cell[int64]{key: "heart_rate", filterKey: "heart_rate",
		//		win: thirtyMinutes(t), time: t, getter: getHeartRate})
		//}
		heartRateDays = append(heartRateDays, cell[int64]{key: "heart_rate", filterKey: "heart_rate",
			win: day(i.Unix()), time: i.Unix(), getter: getHeartRate})

		bloodPreSubDays = append(bloodPreSubDays, cell[int64]{key: "blood_pressure_sub", filterKey: "blood_pressure",
			win: day(i.Unix()), time: i.Unix(), getter: getBloodPreSub})
		bloodPreDiaDays = append(bloodPreDiaDays, cell[int64]{key: "blood_pressure_diastolic", filterKey: "blood_pressure",
			win: day(i.Unix()), time: i.Unix(), getter: getBloodPreDia})
		bloodPreSysDays = append(bloodPreSysDays, cell[int64]{key: "blood_pressure_systolic", filterKey: "blood_pressure",
			win: day(i.Unix()), time: i.Unix(), getter: getBloodPreSys})
		bloodPreDays = append(bloodPreDays, cell[[]int64]{key: "blood_pressure", filterKey: "blood_pressure",
			win: day(i.Unix()), time: i.Unix(), getter: getBloodPressure})

		singleSpo2Days = append(singleSpo2Days, cell[int64]{key: "single_spo2", filterKey: "single_spo2",
			win: day(i.Unix()), time: i.Unix(), getter: getSingleSpo2})

		stressDays = append(stressDays, cell[int64]{key: "stress", filterKey: "stress",
			win: day(i.Unix()), time: i.Unix(), getter: getStress})

		sleepDurationDays = append(sleepDurationDays, cell[int64]{key: "sleep_duration", filterKey: "watch_night_sleep",
			win: day(i.Unix()), time: i.Unix(), getter: getSleepDuration})
		sleepDeepDurationDays = append(sleepDeepDurationDays, cell[int64]{key: "sleep_duration_deep", filterKey: "watch_night_sleep",
			win: day(i.Unix()), time: i.Unix(), getter: getSleepDeepDuration})
		sleepLightDurationDays = append(sleepLightDurationDays, cell[int64]{key: "sleep_duration_light", filterKey: "watch_night_sleep",
			win: day(i.Unix()), time: i.Unix(), getter: getSleepLightDuration})
		sleepRemDurationDays = append(sleepRemDurationDays, cell[int64]{key: "sleep_duration_rem", filterKey: "watch_night_sleep",
			win: day(i.Unix()), time: i.Unix(), getter: getSleepRemDuration})

	}

	hearRateMonth := []cell[int64]{}
	bloodPreSubMonth := []cell[int64]{}
	bloodPreDiaMonth := []cell[int64]{}
	bloodPreSysMonth := []cell[int64]{}
	singleSpo2Month := []cell[int64]{}
	stressMonth := []cell[int64]{}
	sleepDurationMonth := []cell[int64]{}
	sleepDeepDurationMonth := []cell[int64]{}
	sleepLightDurationMonth := []cell[int64]{}
	sleepRemDurationMonth := []cell[int64]{}

	for i := start; i.Before(end); i = i.AddDate(0, 1, 0) {
		hearRateMonth = append(hearRateMonth, cell[int64]{key: "心跳", filterKey: "heart_rate",
			win: month(i.Unix()), time: i.Unix(), getter: getHeartRate})

		bloodPreSubMonth = append(bloodPreSubMonth, cell[int64]{key: "血压差", filterKey: "blood_pressure",
			win: month(i.Unix()), time: i.Unix(), getter: getBloodPreSub})
		bloodPreDiaMonth = append(bloodPreDiaMonth, cell[int64]{key: "低压", filterKey: "blood_pressure",
			win: month(i.Unix()), time: i.Unix(), getter: getBloodPreDia})
		bloodPreSysMonth = append(bloodPreSysMonth, cell[int64]{key: "高压", filterKey: "blood_pressure",
			win: month(i.Unix()), time: i.Unix(), getter: getBloodPreSys})

		singleSpo2Month = append(singleSpo2Month, cell[int64]{key: "血氧", filterKey: "single_spo2",
			win: month(i.Unix()), time: i.Unix(), getter: getSingleSpo2})

		stressMonth = append(stressMonth, cell[int64]{key: "压力", filterKey: "stress",
			win: month(i.Unix()), time: i.Unix(), getter: getStress})

		sleepDurationMonth = append(sleepDurationMonth, cell[int64]{key: "睡眠", filterKey: "watch_night_sleep",
			win: month(i.Unix()), time: i.Unix(), getter: getSleepDuration})
		sleepDeepDurationMonth = append(sleepDeepDurationMonth, cell[int64]{key: "深睡眠", filterKey: "watch_night_sleep",
			win: month(i.Unix()), time: i.Unix(), getter: getSleepDeepDuration})
		sleepLightDurationMonth = append(sleepLightDurationMonth, cell[int64]{key: "浅睡眠", filterKey: "watch_night_sleep",
			win: month(i.Unix()), time: i.Unix(), getter: getSleepLightDuration})
		sleepRemDurationMonth = append(sleepRemDurationMonth, cell[int64]{key: "rem睡眠", filterKey: "watch_night_sleep",
			win: month(i.Unix()), time: i.Unix(), getter: getSleepRemDuration})

	}

	f, err := os.Open(*path)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	datas := readCsv(f)
	for {
		data, ok := <-datas
		if !ok {
			break
		}
		if data.Time < start.Unix() || data.Time > end.Unix() {
			continue
		}

		t := typeMapper[data.Key]
		err := json.Unmarshal([]byte(data.Value), t)
		if err == nil {
			collect(data, t, hearRateMonth)
			collect(data, t, bloodPreSubMonth)
			collect(data, t, bloodPreSysMonth)
			collect(data, t, bloodPreDiaMonth)
			collect(data, t, heartRateHours)

			collect(data, t, heartRateDays)
			collect(data, t, bloodPreSysDays)
			collect(data, t, bloodPreDiaDays)
			collect(data, t, bloodPreSubDays)
			collect(data, t, bloodPreDays)

			collect(data, t, singleSpo2Month)
			collect(data, t, singleSpo2Days)

			collect(data, t, stressMonth)
			collect(data, t, stressDays)

			collect(data, t, sleepDurationMonth)
			collect(data, t, sleepDurationDays)

			collect(data, t, sleepDeepDurationMonth)
			collect(data, t, sleepDeepDurationDays)

			collect(data, t, sleepLightDurationMonth)
			collect(data, t, sleepLightDurationDays)

			collect(data, t, sleepRemDurationMonth)
			collect(data, t, sleepRemDurationDays)
		}
	}

	fmt.Printf("# 汇总 %s ~ %s\n\n", *startS, *endS)
	verticalTableWithMd(hearRateMonth, bloodPreSysMonth, bloodPreDiaMonth, bloodPreSubMonth, singleSpo2Month,
		stressMonth, sleepDurationMonth, sleepDeepDurationMonth, sleepLightDurationMonth, sleepRemDurationMonth)
	fmt.Printf("\n# 心跳 %s ~ %s\n\n", *startS, *endS)
	horizontalTableWithMd(heartRateDays)
	fmt.Printf("\n# 血氧 %s ~ %s\n\n", *startS, *endS)
	horizontalTableWithMd(singleSpo2Days)
	fmt.Printf("\n# 压力 %s ~ %s\n\n", *startS, *endS)
	horizontalTableWithMd(stressDays)
	fmt.Printf("\n# 血压 %s ~ %s\n\n", *startS, *endS)
	horizontalTableWithMdFromBlood(bloodPreDiaDays, bloodPreSysDays, bloodPreSubDays)
}

func horizontalTableWithMd(calls []cell[int64]) {
	//fmt.Println("| key | time | max | min | avg | variance | median | maxDiff | count |")
	//fmt.Println("|---|---|---|---|---|---|---|---|---|")
	//for _, v := range calls {
	//	c := staticsCelldata(v.values)
	//	fmt.Printf("| %s | %s | %d | %d | %.2f | %.2f | %d | %d | %d |\n",
	//		v.key, time.Unix(v.time, 0).Format(time.DateOnly), c.max, c.min, c.avg, c.variance, c.median, c.maxDiff, c.count)
	//}

	data := [][]string{}
	data = append(data, []string{"time", "---", "max", "min", "avg", "variance", "median", "maxDiff", "count"})
	for _, v := range calls {
		c := staticsCelldata(v.values)
		row := []string{
			time.Unix(v.time, 0).Format("02"),
			"---",
			fmt.Sprintf("%d", c.max),
			fmt.Sprintf("%d", c.min),
			fmt.Sprintf("%.2f", c.avg),
			fmt.Sprintf("%.2f", c.variance),
			fmt.Sprintf("%d", c.median),
			fmt.Sprintf("%d", c.maxDiff),
			fmt.Sprintf("%d", c.count),
		}
		data = append(data, row)
	}
	for i := 0; i < len(data[0]); i++ {
		for j := 0; j < len(data); j++ {
			fmt.Printf("| %s ", data[j][i])
		}
		fmt.Println("|")
	}
}

func horizontalTableWithMdFromBlood(dia, sys, sub []cell[int64]) {

	// AI 生成的转换方案很别扭, 为啥不直接生成一个 4*n 的数组, 要生成一个 n*4 的数组.
	// 4*n 可以直接 strings.Join(data, "|") 就可以了, 但是 n*4 就不行了, 蛋疼
	// Collect data into a 2D slice
	data := [][]string{}
	data = append(data, []string{"time", "---", "max", "min", "sub"})
	for i := range dia {
		di := staticsCelldata(dia[i].values)
		sy := staticsCelldata(sys[i].values)
		su := staticsCelldata(sub[i].values)
		row := []string{
			time.Unix(dia[i].time, 0).Format("02"),
			"---",
			fmt.Sprintf("%.2f", di.avg),
			fmt.Sprintf("%.2f", sy.avg),
			fmt.Sprintf("%.2f", su.avg),
		}
		data = append(data, row)
	}

	// Transpose and print the data
	for i := 0; i < len(data[0]); i++ {
		for j := 0; j < len(data); j++ {
			fmt.Printf("| %s ", data[j][i])
		}
		fmt.Println("|")
	}
}

func verticalTableWithMd(cells ...[]cell[int64]) {
	fmt.Println("| key | time | max | min | avg | variance | median | maxDiff | count |")
	fmt.Println("|---|---|---|---|---|---|---|---|---|")
	for _, v2 := range cells {
		for _, v := range v2 {
			c := staticsCelldata(v.values)
			fmt.Printf("| %s | %s | %d | %d | %.2f | %.2f | %d | %d | %d |\n",
				v.key, time.Unix(v.time, 0).Format("2006-01"), c.max, c.min, c.avg, c.variance, c.median, c.maxDiff, c.count)
		}
	}
}

func collect[T any](data FitnessData, t any, cells []cell[T]) {
	for i := range cells {
		v := cells[i]
		if v.filterKey == data.Key && v.win(data.Time) {
			if num, ok := v.getter(data, t); ok == nil {
				v.values = append(v.values, num)
				cells[i] = v
				break
			}
		}
	}
}

type windows func(time int64) bool

func day(unixtime int64) windows {
	t := time.Unix(unixtime, 0)
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, time.Local).Unix()
	return func(unixtime int64) bool {
		return unixtime >= start && unixtime <= end
	}
}

func thirtyMinutes(unixtime int64) windows {
	t := time.Unix(unixtime, 0)
	start := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()/30*30, 0, 0, time.Local).Unix()
	end := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()/30*30+29, 59, 999, time.Local).Unix()
	return func(unixtime int64) bool {
		return unixtime >= start && unixtime <= end
	}
}

func month(unixtime int64) windows {
	t := time.Unix(unixtime, 0)
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local).Unix()
	end := time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999, time.Local).Unix()
	return func(unixtime int64) bool {
		return unixtime >= start && unixtime < end
	}
}

type celldata struct {
	min      int64
	max      int64
	avg      float64
	variance float64
	median   int64
	maxDiff  int64
	count    int
}

func staticsCelldata(datas []int64) celldata {
	if len(datas) == 0 {
		return celldata{}
	}
	var sum int64
	var max int64 = math.MinInt64
	var min int64 = math.MaxInt64
	for _, v := range datas {
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	avg := float64(sum) / float64(len(datas))
	var variance float64
	for _, v := range datas {
		variance += math.Pow(float64(v)-avg, 2)
	}
	variance = variance / float64(len(datas))
	median := datas[len(datas)/2]
	maxDiff := max - min
	return celldata{
		min:      min,
		max:      max,
		avg:      avg,
		variance: variance,
		median:   median,
		maxDiff:  maxDiff,
		count:    len(datas),
	}
}

type cell[T any] struct {
	key       string
	filterKey string
	win       windows
	time      int64
	getter    get[T]
	values    []T
}

func readCsv(reader io.Reader) chan FitnessData {
	r := make(chan FitnessData, 20)
	go func() {
		csvr := csv.NewReader(reader)
		for {
			line, err := csvr.Read()
			if err != nil {
				break
			}
			if line[2] == "Key" {
				continue
			}
			ti, err := strconv.ParseInt(line[3], 10, 64)
			if err != nil {
				fmt.Println(line, err)
				continue
			}
			uti, err := strconv.ParseInt(line[5], 10, 64)
			if err != nil {
				fmt.Println(line, err)
			}
			r <- FitnessData{
				Uid:    line[0],
				Sid:    line[1],
				Key:    line[2],
				Time:   ti,
				Value:  line[4],
				UpTime: uti,
			}
		}
		close(r)
	}()
	return r
}

type get[T any] func(data FitnessData, v any) (T, error)

func getHeartRate(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*HeartRate); ok {
		return v.Bpm, nil
	}
	return 0, novalue
}

var novalue = errors.New("no value")

func getBloodPressure(data FitnessData, value any) ([]int64, error) {
	if v, ok := value.(*BloodPressure); ok {
		return []int64{v.DiastolicPressure, v.SystolicPressure}, nil
	}
	return nil, novalue
}

func getBloodPreSub(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*BloodPressure); ok {
		return v.SystolicPressure - v.DiastolicPressure, nil
	}
	return 0, novalue
}

func getBloodPreSys(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*BloodPressure); ok {
		return v.SystolicPressure, nil
	}
	return 0, novalue
}

func getBloodPreDia(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*BloodPressure); ok {
		return v.DiastolicPressure, nil
	}
	return 0, novalue
}

func getSingleSpo2(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*SingleSpo2); ok {
		return v.Spo2, nil
	}
	return 0, novalue
}

func getStress(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*Stress); ok {
		return v.Stress, nil
	}
	return 0, novalue
}

func getSleepAwake(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*WatchNightSleep); ok {
		return v.AwakeCount, nil
	}
	return 0, novalue
}

func getSleepDuration(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*WatchNightSleep); ok {
		return v.Duration, nil
	}
	return 0, novalue
}

func getSleepDeepDuration(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*WatchNightSleep); ok {
		return v.SleepDeepDuration, nil
	}
	return 0, novalue
}

func getSleepLightDuration(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*WatchNightSleep); ok {
		return v.SleepLightDuration, nil
	}
	return 0, novalue
}

func getSleepRemDuration(data FitnessData, value any) (int64, error) {
	if v, ok := value.(*WatchNightSleep); ok {
		return v.SleepRemDuration, nil
	}
	return 0, novalue
}

type FitnessData struct {
	Uid    string
	Sid    string
	Key    string
	Time   int64
	Value  string
	UpTime int64
}

var typeMapper = map[string]any{
	"abnormal_heart_beat": &AbnormalHeartBeat{},
	"blood_pressure":      &BloodPressure{},
	"calories":            &Calories{},
	"dynamic":             &Dynamic{},
	"heart_rate":          &HeartRate{},
	"intensity":           &Intensity{},
	"menstrual_symptoms":  &MenstrualSymptoms{},
	"menstruation":        &Menstruation{},
	"pai":                 &Pai{},
	"resting_heart_rate":  &RestingHeartRate{},
	"single_heart_rate":   &SingleHeartRate{},
	"single_spo2":         &SingleSpo2{},
	"single_stress":       &SingleStress{},
	"steps":               &Steps{},
	"stress":              &Stress{},
	"training_load":       &TrainingLoad{},
	"valid_stand":         &ValidStand{},
	"vo2_max":             &Vo2Max{},
	"watch_daytime_sleep": &WatchDaytimeSleep{},
	"watch_night_sleep":   &WatchNightSleep{},
	"weight":              &Weight{},
}
