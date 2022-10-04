package redis

type DatabaseConfig struct {
	DATABASES struct {
		APPLDB          Database `json:"APPL_DB"`
		ASICDB          Database `json:"ASIC_DB"`
		COUNTERSDB      Database `json:"COUNTERS_DB"`
		LOGLEVELDB      Database `json:"LOGLEVEL_DB"`
		CONFIGDB        Database `json:"CONFIG_DB"`
		PFCWDDB         Database `json:"PFC_WD_DB"`
		FLEXCOUNTERDB   Database `json:"FLEX_COUNTER_DB"`
		STATEDB         Database `json:"STATE_DB"`
		SNMPOVERLAYDB   Database `json:"SNMP_OVERLAY_DB"`
		ERRORDB         Database `json:"ERROR_DB"`
		RESTAPIDB       Database `json:"RESTAPI_DB"`
		GBASICDB        Database `json:"GB_ASIC_DB"`
		GBCOUNTERSDB    Database `json:"GB_COUNTERS_DB"`
		GBFLEXCOUNTERDB Database `json:"GB_FLEX_COUNTER_DB"`
		EVENTDB         Database `json:"EVENT_DB"`
	} `json:"DATABASES"`
	VERSION string `json:"VERSION"`
}

type Database struct {
	ID        int    `json:"id"`
	Separator string `json:"separator"`
	Instance  string `json:"instance"`
}
