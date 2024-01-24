package main

type (
	Dynamic struct {
		Calories  int64 `json:"calories"`
		Distance  int64 `json:"distance"`
		EndTime   int64 `json:"end_time"`
		StartTime int64 `json:"start_time"`
		Steps     int64 `json:"steps"`
		Type      int64 `json:"type"`
	}
	RestingHeartRate struct {
		Bpm      int64 `json:"bpm"`
		DateTime int64 `json:"date_time"`
	}
	Weight struct {
		Bmi    float64 `json:"bmi"`
		Time   int64   `json:"time"`
		Weight float64 `json:"weight"`
	}
	Stress struct {
		Time   int64 `json:"time"`
		Stress int64 `json:"stress"`
	}
	Menstruation struct {
		Status     int64 `json:"status"`
		DateTime   int64 `json:"date_time"`
		UpdateTime int64 `json:"update_time"`
	}
	AbnormalHeartBeat struct {
		EndTime   int64 `json:"end_time"`
		StartTime int64 `json:"start_time"`
	}
	SingleStress struct {
		Stress int64 `json:"stress"`
		Time   int64 `json:"time"`
	}
	Intensity struct {
		Time int64 `json:"time"`
	}
	ValidStand struct {
		EndTime   int64 `json:"end_time"`
		StartTime int64 `json:"start_time"`
	}
	BloodPressure struct {
		SystolicPressure  int64 `json:"systolic_pressure"`
		DiastolicPressure int64 `json:"diastolic_pressure"`
		Time              int64 `json:"time"`
	}
	TrainingLoad struct {
		CurrentDayTrainLoad int64 `json:"current_day_train_load"`
		DateTime            int64 `json:"date_time"`
		WtlSum              int64 `json:"wtl_sum"`
		WtlSumOptimalMax    int64 `json:"wtl_sum_optimal_max"`
		WtlSumOptimalMin    int64 `json:"wtl_sum_optimal_min"`
		WtlSumOverreaching  int64 `json:"wtl_sum_overreaching"`
	}
	HeartRate struct {
		Time int64 `json:"time"`
		Bpm  int64 `json:"bpm"`
	}
	Steps struct {
		Time     int64 `json:"time"`
		Steps    int64 `json:"steps"`
		Distance int64 `json:"distance"`
		Calories int64 `json:"calories"`
	}
	SingleSpo2 struct {
		Spo2 int64 `json:"spo2"`
		Time int64 `json:"time"`
	}
	Pai struct {
		HighZonePai   float64 `json:"high_zone_pai"`
		LowZonePai    float64 `json:"low_zone_pai"`
		MediumZonePai float64 `json:"medium_zone_pai"`
		DailyPai      float64 `json:"daily_pai"`
		DateTime      int64   `json:"date_time"`
		TotalPai      float64 `json:"total_pai"`
	}
	SingleHeartRate struct {
		Bpm  int64 `json:"bpm"`
		Time int64 `json:"time"`
	}
	Vo2Max struct {
		Time   int64 `json:"time"`
		Vo2Max int64 `json:"vo2_max"`
	}
	MenstrualSymptoms struct {
		Pain     int64 `json:"pain"`
		DateTime int64 `json:"date_time"`
	}
	Calories struct {
		Time     int64 `json:"time"`
		Calories int64 `json:"calories"`
	}
	WatchNightSleep struct {
		AwakeCount         int64 `json:"awake_count"`
		SleepAwakeDuration int64 `json:"sleep_awake_duration"`
		Bedtime            int64 `json:"bedtime"`
		BreathQuality      int64 `json:"breath_quality"`
		SleepDeepDuration  int64 `json:"sleep_deep_duration"`
		SleepLightDuration int64 `json:"sleep_light_duration"`
		SleepRemDuration   int64 `json:"sleep_rem_duration"`
		Duration           int64 `json:"duration"`
		Items              []struct {
			EndTime   int64 `json:"end_time"`
			State     int64 `json:"state"`
			StartTime int64 `json:"start_time"`
		}
		DateTime   int64 `json:"date_time"`
		Timezone   int64 `json:"timezone"`
		WakeUpTime int64 `json:"wake_up_time"`
	}
	WatchDaytimeSleep struct {
		Bedtime  int64 `json:"bedtime"`
		Duration int64 `json:"duration"`
		Items    []struct {
			EndTime   int64 `json:"end_time"`
			State     int64 `json:"state"`
			StartTime int64 `json:"start_time"`
		}
		DateTime   int64 `json:"date_time"`
		Timezone   int64 `json:"timezone"`
		WakeUpTime int64 `json:"wake_up_time"`
	}
)
