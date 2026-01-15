package models

// ThaiWaterAPIResponse represents the complete API response wrapper
type ThaiWaterAPIResponse struct {
	Result string              `json:"result"`
	Data   []ThaiWaterResponse `json:"data"`
	Scale  ScaleInfo           `json:"scale"`
}

// ThaiWaterResponse represents individual water level data from Thai Water API
type ThaiWaterResponse struct {
	ID                    int64    `json:"id"`
	WaterlevelDatetime    string   `json:"waterlevel_datetime"`
	WaterlevelM           *float64 `json:"waterlevel_m"`
	WaterlevelMSL         string   `json:"waterlevel_msl"`
	WaterlevelMSLPrevious string   `json:"waterlevel_msl_previous"`
	FlowRate              *float64 `json:"flow_rate"`
	Discharge             *float64 `json:"discharge"`
	StoragePercent        string   `json:"storage_percent"`
	SortOrder             *int     `json:"sort_order"`
	StationType           string   `json:"station_type"`
	SituationLevel        int      `json:"situation_level"`
	Agency                Agency   `json:"agency"`
	Basin                 Basin    `json:"basin"`
	Station               Station  `json:"station"`
	Geocode               Geocode  `json:"geocode"`
	DiffWLBank            string   `json:"diff_wl_bank"`
	DiffWLBankText        string   `json:"diff_wl_bank_text"`
	RiverGID              int      `json:"river_gid"`
	RiverName             string   `json:"river_name"`
}

// Agency represents the agency information
type Agency struct {
	ID              int           `json:"id"`
	AgencyName      MultiLangText `json:"agency_name"`
	AgencyShortname MultiLangText `json:"agency_shortname"`
}

// Basin represents the basin information
type Basin struct {
	ID        int           `json:"id"`
	BasinCode int           `json:"basin_code"`
	BasinName MultiLangText `json:"basin_name"`
}

// Station represents the telemetry station information
type Station struct {
	ID                 int           `json:"id"`
	TeleStationName    MultiLangText `json:"tele_station_name"`
	TeleStationLat     float64       `json:"tele_station_lat"`
	TeleStationLong    float64       `json:"tele_station_long"`
	TeleStationOldcode string        `json:"tele_station_oldcode"`
	TeleStationType    string        `json:"tele_station_type"`
	LeftBank           float64       `json:"left_bank"`
	RightBank          float64       `json:"right_bank"`
	MinBank            float64       `json:"min_bank"`
	GroundLevel        float64       `json:"ground_level"`
	Offset             *float64      `json:"offset"`
	SubBasinID         int           `json:"sub_basin_id"`
	AgencyID           int           `json:"agency_id"`
	GeocodeID          int           `json:"geocode_id"`
	HydroID            *int          `json:"hydro_id"`
	Qmax               *float64      `json:"qmax"`
	IsKeyStation       bool          `json:"is_key_station"`
	WarningLevelM      *float64      `json:"warning_level_m"`
	CriticalLevelM     *float64      `json:"critical_level_m"`
	CriticalLevelMSL   *float64      `json:"critical_level_msl"`
}

// Geocode represents the geographic location information
type Geocode struct {
	AreaCode     string        `json:"area_code"`
	AreaName     MultiLangText `json:"area_name"`
	AmphoeCode   string        `json:"amphoe_code"`
	AmphoeName   MultiLangText `json:"amphoe_name"`
	TumbonCode   string        `json:"tumbon_code"`
	TumbonName   MultiLangText `json:"tumbon_name"`
	ProvinceCode string        `json:"province_code"`
	ProvinceName MultiLangText `json:"province_name"`
}

// MultiLangText represents text in multiple languages
type MultiLangText struct {
	TH string `json:"th"`
	EN string `json:"en"`
	JP string `json:"jp"`
}

// ScaleInfo represents water level scale and classification information
type ScaleInfo struct {
	Scale    []ScaleLevel     `json:"scale"`
	Rule     []ScaleRule      `json:"rule"`
	Level    map[string]Level `json:"level"`
	NotToday NotTodayInfo     `json:"not_today"`
	RuleWeb  []ScaleRule      `json:"rule_web"`
}

// ScaleLevel represents a water level classification scale
type ScaleLevel struct {
	Operator  string `json:"operator"`
	Term      string `json:"term"`
	Color     string `json:"color"`
	ColorName string `json:"colorname"`
	Situation string `json:"situation"`
	Trans     string `json:"trans"`
	Text      string `json:"text"`
}

// ScaleRule represents a rule for water level classification
type ScaleRule struct {
	Operator string `json:"operator"`
	Term     string `json:"term"`
	Level    int    `json:"level"`
}

// Level represents level color and translation information
type Level struct {
	Color     string `json:"color"`
	ColorName string `json:"colorname"`
	Trans     string `json:"trans"`
}

// NotTodayInfo represents information for non-current data
type NotTodayInfo struct {
	Color     string `json:"color"`
	ColorName string `json:"colorname"`
	Text      string `json:"text"`
}
