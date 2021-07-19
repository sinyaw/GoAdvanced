package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var carno = map[string]int{} //to use for update TypeVehicle's Inventory
var carModel = []string{}    // to store the Vehicle Model Type
var MData = ModelData{}      // to other function for reference
var VData = VehiclesData{}   // to other function for reference
var MLenght [7]int           // for printing format purpuse
var VLenght [8]int           // for printing format purpuse
var CarReturn []int          // for update vehicle status

// var VLenght [8]int

type TypeVehicle struct {
	Type          string
	Name          string
	NumberOfSeats int
	Inventory     int
	DailyCost     float64
	WeeklyCost    float64
	MonthlyCost   float64
	NextPage      *TypeVehicle
}

type ModelData struct {
	TypeVehicle *TypeVehicle
	TotalType   int
}

type Vehicles struct {
	ID                int
	Name              string
	Selecting         bool
	Reserved          bool
	Customer          string
	StartDate         time.Time
	EndDate           time.Time
	Services          bool
	LinkToTypeVehicle *TypeVehicle // pointer to the model detail
	NextPage          *Vehicles
}

type VehiclesData struct {
	Vehicles      *Vehicles
	TotalVehicles int
}

//Read Model Data from Excel file by links list formating with name sorted
func (MD *ModelData) readData() (ModelData, error) {
	f, err := excelize.OpenFile("ModelData.xlsx")
	if err == nil {
		rows := f.GetRows("Sheet1")
		for i, row := range rows {
			newDetail := &TypeVehicle{}
			if i > 0 {
				newDetail.Type = row[0]
				newDetail.Name = row[1]
				newDetail.NumberOfSeats, _ = strconv.Atoi(row[2])
				newDetail.Inventory, _ = strconv.Atoi(row[3])
				newDetail.DailyCost, _ = strconv.ParseFloat(row[4], 64)
				newDetail.WeeklyCost, _ = strconv.ParseFloat(row[5], 64)
				newDetail.MonthlyCost, _ = strconv.ParseFloat(row[6], 64)

				//for printing purpuse
				for i, _ := range MLenght {
					if MLenght[i] < len(row[i]) {
						MLenght[i] = len(row[i])
					}
				}

				//for update vehicle QTY
				if len(carno) == 0 {
					carno[newDetail.Name] = 0

				} else {
					i := len(carno)
					for k, _ := range carno {
						if newDetail.Name == k {
							break
						}
						i--
					}
					if i == 0 {
						carno[newDetail.Name] = 0
					}
				}

				MD.addToDataStructByName(newDetail)
			}

			MD.TotalType++
		}
	} else {
		fmt.Println(err)
	}
	return *MD, nil
}

//Store back Model Data to Excel file
func (MD *ModelData) StoreData() error {
	currentNode := MD.TypeVehicle
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("Sheet1")
	if currentNode == nil {
		fmt.Println("No Model Data to Store!")
	} else {
		i := 2
		for currentNode != nil {
			// fmt.Println(currentNode, currentNode.Name) // testing for getting the data from struct
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), currentNode.Type)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), currentNode.Name)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i), currentNode.NumberOfSeats)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i), carno[currentNode.Name])
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i), currentNode.DailyCost)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i), currentNode.WeeklyCost)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(i), currentNode.MonthlyCost)
			currentNode = currentNode.NextPage
			i++
		}
	}
	xlsx.SetActiveSheet(index)
	err := xlsx.SaveAs("ModelData.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//Add Model Detial
func (MD *ModelData) AddNewMedal() error {
	var addedRecord = map[string]string{}
	var err error
	var input string
	newDetail := &TypeVehicle{}

	if MD.TypeVehicle != nil {

	Outer:
		for {
			fmt.Println("Enter New Vehicle Name: (Name cannot repeat or Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
				continue Outer
			} else if input == "*" {
				page = 0
				return nil
			}

			currentNode := MD.TypeVehicle
			for currentNode != nil {
				if input == currentNode.Name {
					fmt.Println("The input name is already exist in the data!")
					continue Outer
				}
				currentNode = currentNode.NextPage
			}
			break
		}
		newDetail.Name = input
		addedRecord["Name"] = input
		for {
			fmt.Println("Type Of Vehicle (Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
			} else if input == "*" {
				page = 0
				return nil
			} else {
				newDetail.Type = input
				addedRecord["Type"] = input
				break
			}
		}
		for {
			fmt.Println("Number Of Seats (Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
			} else if input == "*" {
				page = 0
				return nil
			} else {
				newDetail.NumberOfSeats, _ = strconv.Atoi(input)
				addedRecord["Number Of Seats"] = input
				break
			}
		}
		for {
			fmt.Println("Inventory (Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
			} else if input == "*" {
				page = 0
				return nil
			} else {
				newDetail.Inventory, _ = strconv.Atoi(input)
				addedRecord["Inventory"] = input
				break
			}
		}
		for {
			fmt.Println("Daily Rental Cost (Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
			} else if input == "*" {
				page = 0
				return nil
			} else {
				newDetail.DailyCost, _ = strconv.ParseFloat(input, 64)
				addedRecord["DailyCost"] = input
				break
			}
		}
		for {
			fmt.Println("Weekly Rental Cost (Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
			} else if input == "*" {
				page = 0
				return nil
			} else {
				newDetail.WeeklyCost, _ = strconv.ParseFloat(input, 64)
				addedRecord["WeeklyCost"] = input
				break
			}
		}
		for {
			fmt.Println("Monthly Rental Cost (Enter '*' return to Main Menu)")
			input = userInput()
			input, err = checkString(input)
			if err != nil {
				fmt.Println(err)
			} else if input == "*" {
				page = 0
				return nil
			} else {
				newDetail.MonthlyCost, _ = strconv.ParseFloat(input, 64)
				addedRecord["MonthlyCost"] = input
				break
			}
		}
		for {
			fmt.Println("Conform to add the Model detial? (Yes or No)")
			input := userInput()
			if strings.ToUpper(input) == "YES" {
				MD.addToDataStructByName(newDetail)
				MD.TotalType++
				break
			} else if strings.ToUpper(input) == "NO" {
				break
			} else {
				fmt.Println("Invalid Input!")
				continue
			}
		}
	}
	return nil
}

//Add Vehicles Quantity
func (VD *VehiclesData) addAdditinalVehicles() {
	currentNode := MData.TypeVehicle
	var namelist []string
	var inventorylist []int
	var i int
	var index int
	var qty int
	var input string

	for currentNode != nil {
		i++
		fmt.Printf("%v.\t%v qty\t %v\n", i, currentNode.Inventory, currentNode.Name)
		inventorylist = append(inventorylist, currentNode.Inventory)
		namelist = append(namelist, currentNode.Name)
		currentNode = currentNode.NextPage
	}
	for {
		fmt.Println("Choose item number for the vehicle to add: (Enter '*' return to Main Menu)")
		input = userInput()
		v, err := checkChoices(input, i)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		} else if v == -1 {
			page = 0
			return
		}
		index = v
		break
	}

	for {
		fmt.Printf("Currently '%v' inventory is %v. How many qty to add?\n", namelist[index-1], inventorylist[index-1])
		input2 := userInput()
		v, err := checkInt(input2)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		} else if v == -1 {
			page = 0
			return
		}
		qty = v
		break
	}

	for i := 1; i <= qty; i++ {
		fmt.Println("go throught")
		newDetail := &Vehicles{
			ID:   VData.TotalVehicles + 1,
			Name: namelist[index],
		}
		VD.addToDataStructByName(newDetail)
	}
	VData.StoreData()
	page = 0
	pages(0)
}

//Insert Model Detial to Model Data by sorting with name
func (MD *ModelData) addToDataStructByName(newDetail *TypeVehicle) {
	if MD.TypeVehicle == nil {
		MD.TypeVehicle = newDetail
		// fmt.Println("a", newDetail.Name, MD.TypeVehicle.Name)
	} else if strings.ToUpper(MD.TypeVehicle.Name) > strings.ToUpper(newDetail.Name) {
		newDetail.NextPage = MD.TypeVehicle
		MD.TypeVehicle = newDetail
		// fmt.Println("b", newDetail.Name, MD.TypeVehicle.Name)
	} else {
		currentNode := MD.TypeVehicle
		for currentNode.NextPage != nil {
			// fmt.Println(i, newDetail.Name > currentNode.Name)
			if strings.ToUpper(currentNode.NextPage.Name) > strings.ToUpper(newDetail.Name) {
				newDetail.NextPage = currentNode.NextPage
				currentNode.NextPage = newDetail
				break
			}
			currentNode = currentNode.NextPage
		}
		if currentNode.NextPage == nil {
			currentNode.NextPage = newDetail
		}
	}
}

//Edit Current Model Data
func (MD *ModelData) EditModelData() map[string]string {
	input := ""
	var changesRecord = map[string]string{}
loop:
	for {
		fmt.Println("Please enter name for the vehicle type for data modification: (Enter '*' to main menu)")
		input = userInput()
		if input == "*" {
			page = 0
			return changesRecord
		} else if len(input) == 0 {
			fmt.Println("No Input Found!")
			continue
		} else {

			currentNode := MD.TypeVehicle
			newDetail := &TypeVehicle{}
			for currentNode.NextPage != nil {
				if MD.TypeVehicle.Name == input {
					newDetail, changesRecord = pullOutEditModelData(currentNode)
					newDetail.NextPage = MD.TypeVehicle.NextPage
					MD.TypeVehicle = newDetail
					break loop
				}

				if currentNode == nil {
					break
				} else if currentNode.NextPage.Name == input {
					newDetail, changesRecord = pullOutEditModelData(currentNode)
					newDetail.NextPage = currentNode.NextPage.NextPage
					currentNode.NextPage = newDetail
					break loop
				}
				currentNode = currentNode.NextPage
			}
			fmt.Println("No data name for the input!")
			continue loop
		}
	}

	return changesRecord
}

// Pull out function from Edit Current Medel Data
func pullOutEditModelData(currentNode *TypeVehicle) (*TypeVehicle, map[string]string) {
	newDetail := &TypeVehicle{}
	var changesRecord = map[string]string{}
	newDetail.Name = currentNode.Name

	fmt.Printf("Type Of Vehicle: %v (Enter for no change)\n", currentNode.Type)
	input := userInput()
	if input == "" {
		newDetail.Type = currentNode.Type
	} else {
		newDetail.Type = input
		changesRecord["Type"] = input
	}

	fmt.Printf("Number Of Seats: %v (Enter for no change)\n", currentNode.NumberOfSeats)
	input = userInput()
	if input == "" {
		newDetail.NumberOfSeats = currentNode.NumberOfSeats
	} else {
		newDetail.NumberOfSeats, _ = strconv.Atoi(input)
		changesRecord["Number Of Seats"] = input
	}

	fmt.Printf("Inventory: %v (Enter for no change)\n", currentNode.Inventory)
	input = userInput()
	if input == "" {
		newDetail.Inventory = currentNode.Inventory
	} else {
		newDetail.Inventory, _ = strconv.Atoi(input)
		changesRecord["Inventory"] = input
	}

	fmt.Printf("Daily Rental Cost: %v (Enter for no change)\n", currentNode.DailyCost)
	input = userInput()
	if input == "" {
		newDetail.DailyCost = currentNode.DailyCost
	} else {
		newDetail.DailyCost, _ = strconv.ParseFloat(input, 64)
		changesRecord["DailyCost"] = input
	}

	fmt.Printf("Weekly Rental Cost: %v (Enter for no change)\n", currentNode.WeeklyCost)
	input = userInput()
	if input == "" {
		newDetail.WeeklyCost = currentNode.WeeklyCost
	} else {
		newDetail.WeeklyCost, _ = strconv.ParseFloat(input, 64)
		changesRecord["WeeklyCost"] = input
	}

	fmt.Printf("Monthly Rental Cost: %v (Enter for no change)\n", currentNode.MonthlyCost)
	input = userInput()
	if input == "" {
		newDetail.MonthlyCost = currentNode.MonthlyCost
	} else {
		newDetail.MonthlyCost, _ = strconv.ParseFloat(input, 64)
		changesRecord["MonthlyCost"] = input
	}
	return newDetail, changesRecord
}

//Read Vehicles Data from Excel file by links list formating with Name and ID sorted
func (VD *VehiclesData) readData() (VehiclesData, error) {
	f, err := excelize.OpenFile("VehiclesData.xlsx")
	if err == nil {
		rows := f.GetRows("Sheet1")
		for i, row := range rows {
			newDetail := &Vehicles{}
			if i > 0 {
				newDetail.ID, _ = strconv.Atoi(row[0])
				newDetail.Name = row[1]
				newDetail.Selecting, _ = strconv.ParseBool(row[2])
				newDetail.Reserved, _ = strconv.ParseBool(row[3])
				newDetail.Customer = row[4]
				newDetail.StartDate = stringToTime(row[5])
				newDetail.EndDate = stringToTime(row[6])
				newDetail.Services, _ = strconv.ParseBool(row[7])

				// to updated the model inventory
				for k, v := range carno {
					if newDetail.Name == k {
						v++
						carno[k] = v
					}
				}

				// for printing purpuse it is not good for the items is huge
				for i, _ := range VLenght {
					if VLenght[i] < len(row[i]) {
						VLenght[i] = len(row[i])
					}
				}

				VD.addToDataStructByName(newDetail)
			}

			VD.TotalVehicles++
		}
	} else {
		fmt.Println(err)
	}

	return *VD, nil
}

//Store back Vehicles Data to Excel file
func (VD *VehiclesData) StoreData() error {
	currentNode := VD.Vehicles
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("Sheet1")
	if currentNode == nil {
		fmt.Println("No Vehicles Data to Store!")
	} else {
		i := 2
		for currentNode != nil {

			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), currentNode.ID)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), currentNode.Name)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i), "FALSE")
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i), currentNode.Reserved)
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i), currentNode.Customer)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i), timeToString(currentNode.StartDate))
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(i), timeToString(currentNode.EndDate))
			xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(i), currentNode.Services)

			currentNode = currentNode.NextPage
			i++
		}
	}
	xlsx.SetActiveSheet(index)
	err := xlsx.SaveAs("VehiclesData.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//Insert Vehicles to Vehicles Data by sorting with name
func (VD *VehiclesData) addToDataStructByName(newDetail *Vehicles) {
	if VD.Vehicles == nil {
		VD.Vehicles = newDetail
	} else if strings.ToUpper(VD.Vehicles.Name) >= strings.ToUpper(newDetail.Name) && VD.Vehicles.ID > newDetail.ID {
		newDetail.NextPage = VD.Vehicles
		VD.Vehicles = newDetail
	} else {
		currentNode := VD.Vehicles
		for currentNode.NextPage != nil {
			if strings.ToUpper(currentNode.NextPage.Name) >= strings.ToUpper(newDetail.Name) && currentNode.NextPage.ID > newDetail.ID {
				newDetail.NextPage = currentNode.NextPage
				currentNode.NextPage = newDetail
				break
			}
			currentNode = currentNode.NextPage
		}
		if currentNode.NextPage == nil {
			currentNode.NextPage = newDetail
		}
	}
}

//update Model Data To Vihicle Data
func uMDTVD(MD *ModelData, VD *VehiclesData) error {
	MDcurrentNode := MD.TypeVehicle
	for MDcurrentNode != nil {
		VDcurrentNode := VD.Vehicles
		for VDcurrentNode != nil {
			if VDcurrentNode.Name == MDcurrentNode.Name {
				VDcurrentNode.LinkToTypeVehicle = MDcurrentNode
			}
			VDcurrentNode = VDcurrentNode.NextPage
		}
		MDcurrentNode = MDcurrentNode.NextPage
	}
	return nil
}

func (VD *VehiclesData) updateVehicle() {
	currentNode := VD.Vehicles
	for currentNode != nil {
		for _, v := range CarReturn {
			if VD.Vehicles.ID == v {
				VD.Vehicles.Reserved = false
			}

			if currentNode == nil {
				break
			} else if currentNode.NextPage.ID == v {
				currentNode.NextPage.Reserved = false
			}
		}
		currentNode = currentNode.NextPage
	}

}

func (MD *ModelData) deleteVehicle(input string) (bool, error) {
	currentNode := MD.TypeVehicle
	if currentNode == nil {
		fmt.Println("Existing No Data for delect!")
		return true, nil
	}
	for currentNode.NextPage != nil {
		if MD.TypeVehicle.Name == input {
			MD.TypeVehicle = currentNode.NextPage
			fmt.Printf("The '%v' has deleted from the data!", input)
			return true, nil
		}
		if currentNode == nil {
			break
		} else if currentNode.NextPage.Name == input {
			currentNode.NextPage = currentNode.NextPage.NextPage
			fmt.Printf("The '%v' has deleted from the data!", input)
			return true, nil
		}
		currentNode = currentNode.NextPage
	}
	fmt.Println("No data name for the input!")
	return false, nil
}

//for Vehicles data struct to use for printing it's data
func (D *VehiclesData) printAllNodes() error {
	cN := D.Vehicles
	if cN == nil {
		fmt.Println("no item recorded")
	} else {
		for cN != nil {
			fmt.Printf("%v\t", cN.ID)
			fmt.Printf("%v", cN.Name)
			for j := 0; j < (MLenght[1] - len(cN.Name) + 3); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("%v\t", cN.Selecting)
			fmt.Printf("%v\t", cN.Reserved)
			fmt.Printf("%v\t", cN.Services)
			fmt.Println()
			cN = cN.NextPage
		}
	}
	return nil
}

//for Model data struct to use for printing it's data
func (D *ModelData) printAllNodes() error {
	cN := D.TypeVehicle
	pi := 0
	if cN == nil {
		fmt.Println("no item recorded")
	} else {
		for cN != nil {
			pi++
			fmt.Printf("%v.\t", pi)
			fmt.Printf("%v", cN.Name)
			for j := 0; j < (MLenght[1] - len(cN.Name) + 3); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("%v", cN.Type)
			for j := 0; j < (MLenght[0] - len(cN.Type) + 3); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("%v seats", cN.NumberOfSeats)
			for j := 0; j < (MLenght[2] - len(strconv.Itoa(cN.NumberOfSeats)) + 5); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("%v Qty", cN.Inventory)
			for j := 0; j < (MLenght[3] - len(strconv.Itoa(cN.Inventory)) + 6); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("$%v", cN.DailyCost)
			for j := 0; j < (MLenght[4] - len(strconv.Itoa(int(cN.DailyCost))) + 5); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("$%v", cN.WeeklyCost)
			for j := 0; j < (MLenght[5] - len(strconv.Itoa(int(cN.WeeklyCost))) + 5); j++ {
				fmt.Printf(" ")
			}
			fmt.Printf("$%v", cN.MonthlyCost)
			for j := 0; j < (MLenght[6] - len(strconv.Itoa(int(cN.MonthlyCost))) + 5); j++ {
				fmt.Printf(" ")
			}
			fmt.Println()
			cN = cN.NextPage
		}
	}
	page = 0
	return nil
}

//View Vehicles by Catogory
func getCarModel(MD *ModelData) ([]string, error) {
	// var TType string
	var add bool = false
	currentNode := MD.TypeVehicle
	for currentNode != nil {
		if len(carModel) == 0 {
			carModel = append(carModel, currentNode.Type)
		} else {
			for _, v := range carModel {
				if currentNode.Type == v {
					add = false
					break
				}
				add = true
			}
			if add {
				carModel = append(carModel, currentNode.Type)
			}

		}
		currentNode = currentNode.NextPage
	}
	return carModel, nil
}

//to set 1st Node for Model Data to 'Mdata'
func getModelData() error {
	defer wg.Done()
	md := &ModelData{}
	MD, _ := md.readData()
	MData = MD
	return nil
}

//to set 1st Node for Vehicle Data to 'Vdata'
func getVehicleData() error {
	defer wg.Done()
	vd := &VehiclesData{}
	VD, _ := vd.readData()
	VData = VD
	return nil
}
