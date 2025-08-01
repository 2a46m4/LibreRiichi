package core

var typeRegistry = map[string]func() interface{}{
	"ActionType":    func() interface{} { return new(ActionType) },
	"Action":        func() interface{} { return new(Action) },
	"ActionHandler": func() interface{} { return new(ActionHandler) },
	"Ron":           func() interface{} { return new(Ron) },
	"Tsumo":         func() interface{} { return new(Tsumo) },
	"Riichi":        func() interface{} { return new(Riichi) },
	"Toss":          func() interface{} { return new(Toss) },
	"Skip":          func() interface{} { return new(Skip) },
	"Pon":           func() interface{} { return new(Pon) },
	"Kan":           func() interface{} { return new(Kan) },
	"Chii":          func() interface{} { return new(Chii) },
	"Draw":          func() interface{} { return new(Draw) },
}
