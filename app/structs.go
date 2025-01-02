package main

type HspWeekProgram struct {
	Day         string `json:"day"`
	Begin       string `json:"begin"`
	End         string `json:"end"`
	Temperature int    `json:"temp"`
}
type HspMeta struct {
	SoftwareVersion       string   `json:"sw_version"`
	HardwareVersion       string   `json:"hw_version"`
	BootloaderVersion     string   `json:"bootl_version"`
	WifiSoftwareVersion   string   `json:"wifi_sw_version"`
	WifiBootloaderVersion string   `json:"wifi_bootl_version"`
	SerialNumber          string   `json:"sn"`
	StoveType             string   `json:"typ"`
	Language              string   `json:"language"`
	Nonce                 string   `json:"nonce"`
	EcoEditable           bool     `json:"eco_editable"`
	Ts                    int      `json:"ts"`
	Ean                   string   `json:"ean"`
	Rau                   bool     `json:"rau"`
	WlanFeatures          []string `json:"wlan_features"`
}
type HspStoveError struct {
	Time      string `json:"string"`
	ErrorCode int    `json:"nr"`
}
type HspStove struct {
	Meta               HspMeta          `json:"meta"`
	Start              bool             `json:"prg"`
	StartWeekProgram   bool             `json:"wprg"`
	Mode               string           `json:"mode"`
	TargetTemperature  int              `json:"sp_temp"`
	CurrentTemperature float32          `json:"is_temp"`
	HeatingCurve       int              `json:"ht_char"`
	WeekProgram        []HspWeekProgram `json:"weekprogram"`
	Error              []HspStoveError  `json:"error"`
	EcoMode            bool             `json:"eco_mode"`
	RoomMode           bool             `json:"room_mode"`
	TvlTemperature     int              `json:"tvl_temp"`
	Pgi                bool             `json:"pgi"`
	Ignitions          int              `json:"ignitions"`
	OnTime             int              `json:"on_time"`
	Consumption        int              `json:"consumption"`
	MaintenanceIn      int              `json:"maintenance_in"`
	CleaningIn         int              `json:"cleaning_in"`
}

type HspCommand struct {
	TargetTemperature *int  `json:"sp_temp,omitempty"`
	TvlTemperature    *int  `json:"tvl_temp,omitempty"`
	HeatingCurve      *int  `json:"ht_char,omitempty"`
	Start             *bool `json:"prg,omitempty"`
	StartWeekProgram  *bool `json:"wprg,omitempty"`
	RoomMode          *bool `json:"room_mode,omitempty"`
	EcoMode           *bool `json:"eco_mode,omitempty"`
}

type HspSensorDiscovery struct {
	Device            HspDevice `json:"device"`
	Name              string    `json:"name"`
	UniqueId          string    `json:"uniq_id"`
	UnitOfMeasurement string    `json:"unit_of_meas,omitempty"`
	DeviceClass       string    `json:"dev_cla,omitempty"`
	ForceUpdate       bool      `json:"frc_upd"`
	StateTopic        string    `json:"stat_t"`
	ValueTemplate     string    `json:"val_tpl"`
}
type HspButtonDiscovery struct {
	Device        HspDevice `json:"device"`
	Name          string    `json:"name"`
	UniqueId      string    `json:"uniq_id"`
	CommandTopic  string    `json:"cmd_t"`
	ValueTemplate string    `json:"val_tpl"`
	ForceUpdate   bool      `json:"frc_upd"`
}
type HspSwitchDiscovery struct {
	Device        HspDevice `json:"device"`
	Name          string    `json:"name"`
	UniqueId      string    `json:"uniq_id"`
	CommandTopic  string    `json:"cmd_t"`
	StateTopic    string    `json:"stat_t"`
	StateOff      bool      `json:"pl_off"`
	StateOn       bool      `json:"pl_on"`
	ValueTemplate string    `json:"val_tpl"`
	ForceUpdate   bool      `json:"frc_upd"`
}
type HspInputNumberDiscovery struct {
	Device        HspDevice `json:"device"`
	Name          string    `json:"name"`
	UniqueId      string    `json:"uniq_id"`
	CommandTopic  string    `json:"cmd_t"`
	NumberMin     float	`json:"n_min"`
	NumberMax     float	`json:"n_max"`
	ValueTemplate string    `json:"val_tpl"`
}
type HspDevice struct {
	Ids             []string `json:"ids"`
	Name            string   `json:"name"`
	Manufacturer    string   `json:"mf"`
	Model           string   `json:"mdl"`
	SoftwareVersion string   `json:"sw"`
}
type HspClimateDiscovery struct {
	Device                   HspDevice `json:"device"`
	Name                     string    `json:"name"`
	UniqueId                 string    `json:"uniq_id"`
	ForceUpdate              bool      `json:"frc_upd"`
	ModeStateTopic           string    `json:"mode_stat_t"`
	ModeStateTemplate        string    `json:"mode_stat_tpl"`
	TempCommandTopic         string    `json:"temp_cmd_t"`
	TempStateTopic           string    `json:"temp_stat_t"`
	TempStateTemplate        string    `json:"temp_stat_tpl"`
	CurrentTempStateTopic    string    `json:"curr_temp_t"`
	CurrentTempStateTemplate string    `json:"curr_temp_tpl"`
	MinTemp                  string    `json:"min_temp"`
	MaxTemp                  string    `json:"max_temp"`
	TempStep                 string    `json:"temp_step"`
	Modes                    []string  `json:"modes"`
}
type HspError struct {
	ErrorCode int `json:"nr"`
}
type HspCleanError struct {
	SeenError []HspError `json:"seen_error"`
}

func BoolPointer(b bool) *bool {
	boolVar := b
	return &boolVar
}
func IntPointer(number int) *int {
	intVar := number
	return &intVar
}
func FloatPointer(number float) *float {
	floatVar := number
	return &floatVar
}
