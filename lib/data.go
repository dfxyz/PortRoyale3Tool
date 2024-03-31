package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type Data struct {
	Cities map[string]*cityInfo `json:"cities"`
	Groups map[int]*groupInfo   `json:"groups"`
}

type groupInfo struct {
	Cities map[string]struct{} `json:"cities"`
}

type cityInfo struct {
	Buildings map[Good]int `json:"buildings"`
}

func NewData() *Data {
	return &Data{
		Cities: make(map[string]*cityInfo),
		Groups: make(map[int]*groupInfo),
	}
}

func (d *Data) Load(path string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return
	}
	if err := json.Unmarshal(raw, d); err != nil {
		panic("data corrupted")
	}
}

func (d *Data) Save(path string) (err error) {
	raw, err := json.Marshal(d)
	if err != nil {
		return
	}
	err = os.WriteFile(path, raw, 0644)
	return
}

func (d *Data) ListAll() {
	for idx, group := range d.Groups {
		fmt.Printf("Group [%d]:\n", idx)
		for name := range group.Cities {
			fmt.Printf("* [%s]\n", name)
		}
		fmt.Println()
	}

	if len(d.Cities) > 0 {
		fmt.Println("Cities:")
		for name := range d.Cities {
			fmt.Printf("* [%s]\n", name)
		}
	}
}

func (d *Data) ListCity(name string) {
	city, ok := d.Cities[name]
	if !ok {
		fmt.Printf("City [%s] not found\n", name)
		return
	}
	produces, consumes := city.getGoodDetails()
	printCityDetail(name, city.Buildings, produces, consumes)
}

func (city *cityInfo) getGoodDetails() (produces map[Good]int, consumes map[Good]int) {
	produces = make(map[Good]int, len(city.Buildings))
	consumes = make(map[Good]int)

	for produceGood, buildingNum := range city.Buildings {
		if buildingNum <= 0 {
			continue
		}
		info, ok := produceConsumeInfoMap[produceGood]
		if !ok {
			continue
		}
		produceNum := info.produce * buildingNum
		produces[produceGood] = produceNum
		for consumeGood, consumeNum := range info.consumes {
			consumes[consumeGood] += consumeNum * buildingNum
		}
	}
	return
}

func printCityDetail(cityName string, buildings map[Good]int, produces map[Good]int, consumes map[Good]int) {
	fmt.Printf("City [%s]:\n", cityName)
	for produceGood, buildingNum := range buildings {
		produce := produces[produceGood]
		consume := consumes[produceGood]
		printProduceGoodDetail(produceGood, buildingNum, produce, consume)
		fmt.Println()
	}
	for consumeGood, consume := range consumes {
		if _, ok := buildings[consumeGood]; ok {
			continue
		}
		printConsumeGoodDetail(consumeGood, consume)
		fmt.Println()
	}
}

func printProduceGoodDetail(good Good, buildingNum int, produce int, consume int) {
	// <goodName>(<buildingNum>): <balance>
	// <goodName>(<buildingNum>): <balance>(<produceNum>/<consumeNum>)
	fmt.Print("\t")
	printGreenStart()
	fmt.Print(good)
	printColorEnd()
	fmt.Printf("(%d): ", buildingNum)
	if consume <= 0 {
		printSignedNumWithColor(produce)
		return
	}
	printSignedNumWithColor(produce - consume)
	fmt.Print(" (")
	printSignedNumWithColor(produce)
	fmt.Print("/")
	printSignedNumWithColor(-consume)
	fmt.Print(")")
}

func printConsumeGoodDetail(good Good, consume int) {
	// -<goodName>: <consumeNum>
	fmt.Print("\t")
	printRedStart()
	fmt.Print(good)
	printColorEnd()
	fmt.Print(": ")
	printSignedNumWithColor(-consume)
}

func (d *Data) ListGroup(index int) {
	group, ok := d.Groups[index]
	if !ok {
		fmt.Printf("Group [%d] not found\n", index)
		return
	}

	produceBuildings := make(map[Good]int)
	overallProduces := make(map[Good]int)
	overallConsumes := make(map[Good]int)
	for cityName := range group.Cities {
		city, ok := d.Cities[cityName]
		if !ok {
			continue
		}

		produces, consumes := city.getGoodDetails()
		printCityDetail(cityName, city.Buildings, produces, consumes)

		for good, buildingNum := range city.Buildings {
			produceBuildings[good] += buildingNum
		}
		for good, num := range produces {
			overallProduces[good] += num
		}
		for good, num := range consumes {
			overallConsumes[good] += num
		}
	}
	allGoods := make(map[Good]struct{})
	for good := range overallProduces {
		allGoods[good] = struct{}{}
	}
	for good := range overallConsumes {
		allGoods[good] = struct{}{}
	}
	fmt.Printf("Group [%d] balances:\n", index)
	for good := range allGoods {
		buildingNum := produceBuildings[good]
		produce := overallProduces[good]
		consume := overallConsumes[good]
		printOverallGoodDetail(good, buildingNum, produce, consume)
		fmt.Println()
	}
}

func printOverallGoodDetail(good Good, buidlingNum int, produce int, consume int) {
	// <goodName>: <balance>(<produceNum>/<consumeNum>)
	fmt.Printf("\t%s: ", good)
	printSignedNumWithColor(produce - consume)
	if buidlingNum > 0 {
		fmt.Printf("(%d)", buidlingNum)
	}
	fmt.Print(" (")
	printSignedNumWithColor(produce)
	fmt.Print("/")
	printSignedNumWithColor(-consume)
	fmt.Print(")")
}

func printSignedNumWithColor(num int) {
	if num == 0 {
		fmt.Print(0)
		return
	}
	if num > 0 {
		printGreenStart()
		fmt.Printf("+%d", num)
		printColorEnd()
		return
	}
	printRedStart()
	fmt.Printf("%d", num)
	printColorEnd()
}

func printGreenStart() {
	fmt.Print("\033[32m")
}

func printRedStart() {
	fmt.Print("\033[31m")
}

func printColorEnd() {
	fmt.Print("\033[0m")
}

func (d *Data) SetProduceBuilding(cityName string, good Good, buildingNum int) {
	city, ok := d.Cities[cityName]
	if !ok {
		city = &cityInfo{Buildings: make(map[Good]int)}
		d.Cities[cityName] = city
	}
	city.Buildings[good] = buildingNum
	fmt.Printf("City [%s] building [%s] set to %d\n", cityName, good, buildingNum)
}

func (d *Data) GroupAssociate(groupIndex int, cities []string) {
	checkedCities := make([]string, 0, len(cities))
	for _, city := range cities {
		if _, ok := d.Cities[city]; !ok {
			fmt.Printf("City [%s] not found\n", city)
			continue
		}
		checkedCities = append(checkedCities, city)
	}
	if len(checkedCities) <= 0 {
		return
	}

	group, ok := d.Groups[groupIndex]
	if !ok {
		group = &groupInfo{Cities: make(map[string]struct{})}
		d.Groups[groupIndex] = group
	}
	for _, city := range checkedCities {
		group.Cities[city] = struct{}{}
	}
	fmt.Printf("Group [%d] associated with cities: %v\n", groupIndex, checkedCities)
}

func (d *Data) UnsetProduceBuilding(cityName string, good Good) {
	city, ok := d.Cities[cityName]
	if !ok {
		fmt.Printf("City [%s] not found\n", cityName)
		return
	}
	if _, ok := city.Buildings[good]; !ok {
		fmt.Printf("City [%s] building [%s] not found\n", cityName, good)
		return
	}
	delete(city.Buildings, good)
	fmt.Printf("City [%s] building [%s] unset\n", cityName, good)

	if len(city.Buildings) <= 0 {
		d.removeCity(cityName)
	}
}

func (d *Data) removeCity(cityName string) {
	delete(d.Cities, cityName)
	for groupIndex, group := range d.Groups {
		delete(group.Cities, cityName)
		if len(group.Cities) <= 0 {
			d.removeGroup(groupIndex)
		}
	}
	fmt.Printf("City [%s] removed\n", cityName)
}

func (d *Data) UnassociateGroup(groupIndex int, cities []string) {
	group, ok := d.Groups[groupIndex]
	if !ok {
		fmt.Printf("Group [%d] not found\n", groupIndex)
		return
	}
	removedCities := make([]string, 0, len(cities))
	for _, city := range cities {
		if _, ok := group.Cities[city]; !ok {
			fmt.Printf("City [%s] not found in group [%d]\n", city, groupIndex)
			continue
		}
		delete(group.Cities, city)
		removedCities = append(removedCities, city)
	}
	if len(removedCities) > 0 {
		fmt.Printf("Group [%d] unassociated with cities: %v\n", groupIndex, removedCities)
	}
	if len(group.Cities) <= 0 {
		d.removeGroup(groupIndex)
	}
}

func (d *Data) removeGroup(groupIndex int) {
	delete(d.Groups, groupIndex)
	fmt.Printf("Group [%d] removed\n", groupIndex)
}

func (d *Data) RemoveCity(cityName string) {
	if _, ok := d.Cities[cityName]; !ok {
		fmt.Printf("City [%s] not found\n", cityName)
		return
	}
	d.removeCity(cityName)
}

func (d *Data) RemoveGroup(groupIndex int) {
	if _, ok := d.Groups[groupIndex]; !ok {
		fmt.Printf("Group [%d] not found\n", groupIndex)
		return
	}
	d.removeGroup(groupIndex)
}
