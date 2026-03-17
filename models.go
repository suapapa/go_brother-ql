package brother_ql

type Model struct {
	Identifier          string
	MinMaxLengthDots    [2]int
	MinMaxFeed          [2]int
	NumberBytesPerRow   int
	AdditionalOffsetR   int
	ModeSetting         bool
	Cutting             bool
	ExpandedMode      bool
	Compression         bool
	TwoColor           bool
	NumInvalidateBytes int
}

func defaultModel(id string, minMax [2]int) Model {
	return Model{
		Identifier:          id,
		MinMaxLengthDots:    minMax,
		MinMaxFeed:          [2]int{35, 1500},
		NumberBytesPerRow:   90,
		AdditionalOffsetR:   0,
		ModeSetting:         true,
		Cutting:             true,
		ExpandedMode:      true,
		Compression:         true,
		TwoColor:           false,
		NumInvalidateBytes: 200,
	}
}

func initModels() []Model {
	var m []Model

	m500 := defaultModel("QL-500", [2]int{295, 11811})
	m500.Compression = false; m500.ModeSetting = false; m500.ExpandedMode = false; m500.Cutting = false
	m = append(m, m500)

	m550 := defaultModel("QL-550", [2]int{295, 11811})
	m550.Compression = false; m550.ModeSetting = false
	m = append(m, m550)

	m560 := defaultModel("QL-560", [2]int{295, 11811})
	m560.Compression = false; m560.ModeSetting = false
	m = append(m, m560)

	m570 := defaultModel("QL-570", [2]int{150, 11811})
	m570.Compression = false; m570.ModeSetting = false
	m = append(m, m570)

	m = append(m, defaultModel("QL-580N", [2]int{150, 11811}))
	m = append(m, defaultModel("QL-650TD", [2]int{295, 11811}))

	m700 := defaultModel("QL-700", [2]int{150, 11811})
	m700.Compression = false; m700.ModeSetting = false
	m = append(m, m700)

	m = append(m, defaultModel("QL-710W", [2]int{150, 11811}))
	m = append(m, defaultModel("QL-720NW", [2]int{150, 11811}))

	m800 := defaultModel("QL-800", [2]int{150, 11811})
	m800.TwoColor = true; m800.Compression = false; m800.NumInvalidateBytes = 400
	m = append(m, m800)

	m810 := defaultModel("QL-810W", [2]int{150, 11811})
	m810.TwoColor = true; m810.NumInvalidateBytes = 400
	m = append(m, m810)

	m820 := defaultModel("QL-820NWB", [2]int{150, 11811})
	m820.TwoColor = true; m820.NumInvalidateBytes = 400
	m = append(m, m820)

	m1050 := defaultModel("QL-1050", [2]int{295, 35433})
	m1050.NumberBytesPerRow = 162; m1050.AdditionalOffsetR = 44
	m = append(m, m1050)

	m1060 := defaultModel("QL-1060N", [2]int{295, 35433})
	m1060.NumberBytesPerRow = 162; m1060.AdditionalOffsetR = 44
	m = append(m, m1060)

	m1100 := defaultModel("QL-1100", [2]int{301, 35434})
	m1100.NumberBytesPerRow = 162; m1100.AdditionalOffsetR = 44
	m = append(m, m1100)

	m1110 := defaultModel("QL-1110NWB", [2]int{301, 35434})
	m1110.NumberBytesPerRow = 162; m1110.AdditionalOffsetR = 44
	m = append(m, m1110)

	m1115 := defaultModel("QL-1115NWB", [2]int{301, 35434})
	m1115.NumberBytesPerRow = 162; m1115.AdditionalOffsetR = 44
	m = append(m, m1115)

	p750 := defaultModel("PT-P750W", [2]int{31, 14172})
	p750.NumberBytesPerRow = 16
	m = append(m, p750)

	p900 := defaultModel("PT-P900W", [2]int{57, 28346})
	p900.NumberBytesPerRow = 70
	m = append(m, p900)

	return m
}

var AllModels []Model

func init() {
	AllModels = initModels()
}

func GetModel(id string) (Model, bool) {
	for _, m := range AllModels {
		if m.Identifier == id {
			return m, true
		}
	}
	return Model{}, false
}
