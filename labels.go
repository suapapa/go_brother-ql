package brother_ql

// FormFactor specifies the type of label media.
type FormFactor int

const (
	// DieCut represents die-cut labels with fixed size.
	DieCut FormFactor = iota + 1
	// Endless represents continuous tape labels.
	Endless
	// RoundDieCut represents round die-cut labels.
	RoundDieCut
	// PtouchEndless represents P-touch continuous tape labels.
	PtouchEndless
)

// Color specifies the color capability of the label.
type Color int

const (
	// BlackWhite is standard black on white printing.
	BlackWhite Color = iota
	// BlackRedWhite is black and red on white printing (two-color).
	BlackRedWhite
)

// Label represents the specifications of a label type.
type Label struct {
	Identifier         string // e.g., "62", "29x90"
	TapeSize           [2]int // width, length in mm
	FormFactor         FormFactor
	DotsTotal          [2]int   // total dots [width, length]
	DotsPrintable      [2]int   // printable dots [width, length]
	OffsetR            int      // right margin alignment offset
	FeedMargin         int      // feed margin in dots
	RestrictedToModels []string // allowed models, if any
	Color              Color
}


func initLabels() []Label {
	return []Label{
		{Identifier: "12", TapeSize: [2]int{12, 0}, FormFactor: Endless, DotsTotal: [2]int{142, 0}, DotsPrintable: [2]int{106, 0}, OffsetR: 29, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "29", TapeSize: [2]int{29, 0}, FormFactor: Endless, DotsTotal: [2]int{342, 0}, DotsPrintable: [2]int{306, 0}, OffsetR: 6, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "38", TapeSize: [2]int{38, 0}, FormFactor: Endless, DotsTotal: [2]int{449, 0}, DotsPrintable: [2]int{413, 0}, OffsetR: 12, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "50", TapeSize: [2]int{50, 0}, FormFactor: Endless, DotsTotal: [2]int{590, 0}, DotsPrintable: [2]int{554, 0}, OffsetR: 12, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "54", TapeSize: [2]int{54, 0}, FormFactor: Endless, DotsTotal: [2]int{636, 0}, DotsPrintable: [2]int{590, 0}, OffsetR: 0, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "62", TapeSize: [2]int{62, 0}, FormFactor: Endless, DotsTotal: [2]int{732, 0}, DotsPrintable: [2]int{696, 0}, OffsetR: 12, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "62red", TapeSize: [2]int{62, 0}, FormFactor: Endless, DotsTotal: [2]int{732, 0}, DotsPrintable: [2]int{696, 0}, OffsetR: 12, FeedMargin: 35, Color: BlackRedWhite},
		{Identifier: "102", TapeSize: [2]int{102, 0}, FormFactor: Endless, DotsTotal: [2]int{1200, 0}, DotsPrintable: [2]int{1164, 0}, OffsetR: 12, FeedMargin: 35, RestrictedToModels: []string{"QL-1050", "QL-1060N", "QL-1100", "QL-1110NWB", "QL-1115NWB"}, Color: BlackWhite},
		{Identifier: "103", TapeSize: [2]int{104, 0}, FormFactor: Endless, DotsTotal: [2]int{1224, 0}, DotsPrintable: [2]int{1200, 0}, OffsetR: 12, FeedMargin: 35, RestrictedToModels: []string{"QL-1050", "QL-1060N", "QL-1100", "QL-1110NWB", "QL-1115NWB"}, Color: BlackWhite},
		{Identifier: "17x54", TapeSize: [2]int{17, 54}, FormFactor: DieCut, DotsTotal: [2]int{201, 636}, DotsPrintable: [2]int{165, 566}, OffsetR: 0, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "17x87", TapeSize: [2]int{17, 87}, FormFactor: DieCut, DotsTotal: [2]int{201, 1026}, DotsPrintable: [2]int{165, 956}, OffsetR: 0, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "23x23", TapeSize: [2]int{23, 23}, FormFactor: DieCut, DotsTotal: [2]int{272, 272}, DotsPrintable: [2]int{202, 202}, OffsetR: 42, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "29x42", TapeSize: [2]int{29, 42}, FormFactor: DieCut, DotsTotal: [2]int{342, 495}, DotsPrintable: [2]int{306, 425}, OffsetR: 6, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "29x90", TapeSize: [2]int{29, 90}, FormFactor: DieCut, DotsTotal: [2]int{342, 1061}, DotsPrintable: [2]int{306, 991}, OffsetR: 6, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "39x90", TapeSize: [2]int{38, 90}, FormFactor: DieCut, DotsTotal: [2]int{449, 1061}, DotsPrintable: [2]int{413, 991}, OffsetR: 12, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "39x48", TapeSize: [2]int{39, 48}, FormFactor: DieCut, DotsTotal: [2]int{461, 565}, DotsPrintable: [2]int{425, 495}, OffsetR: 6, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "52x29", TapeSize: [2]int{52, 29}, FormFactor: DieCut, DotsTotal: [2]int{614, 341}, DotsPrintable: [2]int{578, 271}, OffsetR: 0, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "60x86", TapeSize: [2]int{60, 87}, FormFactor: DieCut, DotsTotal: [2]int{708, 1024}, DotsPrintable: [2]int{672, 954}, OffsetR: 18, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "62x29", TapeSize: [2]int{62, 29}, FormFactor: DieCut, DotsTotal: [2]int{732, 341}, DotsPrintable: [2]int{696, 271}, OffsetR: 12, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "62x100", TapeSize: [2]int{62, 100}, FormFactor: DieCut, DotsTotal: [2]int{732, 1179}, DotsPrintable: [2]int{696, 1109}, OffsetR: 12, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "102x51", TapeSize: [2]int{102, 51}, FormFactor: DieCut, DotsTotal: [2]int{1200, 596}, DotsPrintable: [2]int{1164, 526}, OffsetR: 12, FeedMargin: 0, RestrictedToModels: []string{"QL-1050", "QL-1060N", "QL-1100", "QL-1110NWB", "QL-1115NWB"}, Color: BlackWhite},
		{Identifier: "102x152", TapeSize: [2]int{102, 153}, FormFactor: DieCut, DotsTotal: [2]int{1200, 1804}, DotsPrintable: [2]int{1164, 1660}, OffsetR: 12, FeedMargin: 0, RestrictedToModels: []string{"QL-1050", "QL-1060N", "QL-1100", "QL-1110NWB", "QL-1115NWB"}, Color: BlackWhite},
		{Identifier: "103x164", TapeSize: [2]int{104, 164}, FormFactor: DieCut, DotsTotal: [2]int{1224, 1941}, DotsPrintable: [2]int{1200, 1822}, OffsetR: 12, FeedMargin: 0, RestrictedToModels: []string{"QL-1100", "QL-1110NWB"}, Color: BlackWhite},
		{Identifier: "d12", TapeSize: [2]int{12, 12}, FormFactor: RoundDieCut, DotsTotal: [2]int{142, 142}, DotsPrintable: [2]int{94, 94}, OffsetR: 113, FeedMargin: 35, Color: BlackWhite},
		{Identifier: "d24", TapeSize: [2]int{24, 24}, FormFactor: RoundDieCut, DotsTotal: [2]int{284, 284}, DotsPrintable: [2]int{236, 236}, OffsetR: 42, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "d58", TapeSize: [2]int{58, 58}, FormFactor: RoundDieCut, DotsTotal: [2]int{688, 688}, DotsPrintable: [2]int{618, 618}, OffsetR: 51, FeedMargin: 0, Color: BlackWhite},
		{Identifier: "pt24", TapeSize: [2]int{24, 0}, FormFactor: PtouchEndless, DotsTotal: [2]int{128, 0}, DotsPrintable: [2]int{128, 0}, OffsetR: 0, FeedMargin: 14, Color: BlackWhite},
	}
}

// AllLabels contains all predefined label configurations.
var AllLabels []Label

func init() {
	AllLabels = initLabels()
}

func getLabel(id string) (Label, bool) {
	for _, l := range AllLabels {
		if l.Identifier == id {
			return l, true
		}
	}
	return Label{}, false
}
