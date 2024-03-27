package lib

type Good uint32

const (
	Wood Good = iota
	Brick
	Wheat
	Fruit
	Corn
	Sugar
	Hemp
	Cotton
	Dye
	Tobacco
	Coffee
	Cacao
	Metal
	Bread
	Meat
	Rum
	Rope
	Textile
	Tool
	Clothes
)

type produceConsumeInfo struct {
	produce  int
	consumes map[Good]int
}

var produceConsumeInfoMap = map[Good]produceConsumeInfo{
	Wood: {
		produce: 60,
	},
	Brick: {
		produce: 60,
	},
	Wheat: {
		produce: 60,
	},
	Fruit: {
		produce: 40,
	},
	Corn: {
		produce: 40,
	},
	Sugar: {
		produce: 40,
	},
	Hemp: {
		produce: 40,
	},
	Cotton: {
		produce: 40,
	},
	Dye: {
		produce: 20,
	},
	Tobacco: {
		produce: 20,
	},
	Coffee: {
		produce: 20,
		consumes: map[Good]int{
			Tool: 5,
		},
	},
	Cacao: {
		produce: 20,
		consumes: map[Good]int{
			Tool: 5,
		},
	},
	Metal: {
		produce: 30,
		consumes: map[Good]int{
			Wood: 15,
		},
	},
	Bread: {
		produce: 20,
		consumes: map[Good]int{
			Wheat: 10,
			Sugar: 10,
		},
	},
	Meat: {
		produce: 10,
		consumes: map[Good]int{
			Corn: 20,
		},
	},
	Rum: {
		produce: 10,
		consumes: map[Good]int{
			Wood:  5,
			Sugar: 10,
		},
	},
	Rope: {
		produce: 20,
		consumes: map[Good]int{
			Hemp: 20,
		},
	},
	Textile: {
		produce: 20,
		consumes: map[Good]int{
			Cotton: 20,
		},
	},
	Tool: {
		produce: 20,
		consumes: map[Good]int{
			Wood:  10,
			Metal: 20,
		},
	},
	Clothes: {
		produce: 10,
		consumes: map[Good]int{
			Dye:     10,
			Textile: 10,
		},
	},
}

func GoodFromStr(s string) (Good, bool) {
	switch s {
	case "wood":
		return Wood, true
	case "brick":
		return Brick, true
	case "wheat":
		return Wheat, true
	case "fruit":
		return Fruit, true
	case "corn":
		return Corn, true
	case "sugar":
		return Sugar, true
	case "hemp":
		return Hemp, true
	case "cotton":
		return Cotton, true
	case "dye":
		return Dye, true
	case "tobacco":
		return Tobacco, true
	case "coffee":
		return Coffee, true
	case "cacao":
		return Cacao, true
	case "metal":
		return Metal, true
	case "bread":
		return Bread, true
	case "meat":
		return Meat, true
	case "rum":
		return Rum, true
	case "rope":
		return Rope, true
	case "textile":
		return Textile, true
	case "tool":
		return Tool, true
	case "clothes":
		return Clothes, true
	default:
		return 0, false
	}
}

func (g Good) String() string {
	switch g {
	case Wood:
		return "Wood"
	case Brick:
		return "Brick"
	case Wheat:
		return "Wheat"
	case Fruit:
		return "Fruit"
	case Corn:
		return "Corn"
	case Sugar:
		return "Sugar"
	case Hemp:
		return "Hemp"
	case Cotton:
		return "Cotton"
	case Dye:
		return "Dye"
	case Tobacco:
		return "Tobacco"
	case Coffee:
		return "Coffee"
	case Cacao:
		return "Cacao"
	case Metal:
		return "Metal"
	case Bread:
		return "Bread"
	case Meat:
		return "Meat"
	case Rum:
		return "Rum"
	case Rope:
		return "Rope"
	case Textile:
		return "Textile"
	case Tool:
		return "Tool"
	case Clothes:
		return "Clothes"
	default:
		return "<unknown>"
	}
}
