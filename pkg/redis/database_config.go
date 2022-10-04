package redis

type DatabaseConfig struct {
	DATABASES struct {
		APPLDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"APPL_DB"`
		ASICDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"ASIC_DB"`
		COUNTERSDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"COUNTERS_DB"`
		LOGLEVELDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"LOGLEVEL_DB"`
		CONFIGDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"CONFIG_DB"`
		PFCWDDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"PFC_WD_DB"`
		FLEXCOUNTERDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"FLEX_COUNTER_DB"`
		STATEDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"STATE_DB"`
		SNMPOVERLAYDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"SNMP_OVERLAY_DB"`
		ERRORDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"ERROR_DB"`
		RESTAPIDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"RESTAPI_DB"`
		GBASICDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"GB_ASIC_DB"`
		GBCOUNTERSDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"GB_COUNTERS_DB"`
		GBFLEXCOUNTERDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"GB_FLEX_COUNTER_DB"`
		EVENTDB struct {
			ID        int    `json:"id"`
			Separator string `json:"separator"`
			Instance  string `json:"instance"`
		} `json:"EVENT_DB"`
	} `json:"DATABASES"`
	VERSION string `json:"VERSION"`
}
