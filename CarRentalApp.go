package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	load bool = true
	page int  = 0
	Exit bool = false
	wg   sync.WaitGroup
)

func pages(p int) {
	fmt.Println()
	var s int
	switch p {
	case 0:
		s, _ = fmt.Println("Car Rental Application")
		fmt.Println(strings.Repeat("=", s-1))
		menu()
	case 1:
		s, _ = fmt.Println("Search Car for Rent")
		fmt.Println(strings.Repeat("=", s-1))
		_, TF := checkCustomerData()
		if TF {
			modelCatogory()
		}
	case 2:
		s, _ = fmt.Println("Current Booking")
		fmt.Println(strings.Repeat("=", s-1))
		NRIC, TF := checkCustomerData()
		if TF {
			BData.printAllNodes(NRIC)
		}
	case 3:
		s, _ = fmt.Println("Previous Booking Record")
		fmt.Println(strings.Repeat("=", s-1))
		NRIC, TF := checkCustomerData()
		if TF {
			priviousBookingRecord(NRIC)
		}
	case 4:
		s, _ = fmt.Println("Add Model Detail")
		fmt.Println(strings.Repeat("=", s-1))
		addModel()
	case 5:
		s, _ = fmt.Println("Model Vehicle's Details")
		fmt.Println(strings.Repeat("=", s-1))
		viewModel()
	case 6:
		s, _ = fmt.Println("Delete of The Model Vehicle")
		fmt.Println(strings.Repeat("=", s-1))
		delModel()
	case 7:
		s, _ = fmt.Println("View Earning From Previous Booking")
		fmt.Println(strings.Repeat("=", s-1))
		viewEarning()
	case 8:
		Exit = true
		fmt.Println("Exit App with Files Saved!!")

	default:
		fmt.Println("Invalid input!")
	}
}

func menu() {
	fmt.Println("1. Search Car for Rent")
	fmt.Println("2. Current Booking")
	fmt.Println("3. Previous Booking Record")
	fmt.Println("4. Add Model Detail")
	fmt.Println("5. View or Modify Model Vehicle's Details")
	fmt.Println("6. Delete of The Model Vehicle")
	fmt.Println("7. View Earning From Previous Booking")
	fmt.Println("8. Exit and Save Files")
	fmt.Println("Select Your Choice:")

	input := userInput()
	v, err := strconv.Atoi(input)
	if len(input) == 0 {
		fmt.Println("No Input Found!")
	} else if err != nil {
		fmt.Println("Invalid input!")
	} else {
		page = v
	}
	return
}

func main() {
	runtime.GOMAXPROCS(50)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println()
			fmt.Println("Sorry for there are some problem occured... Now, we start from the main menu.")
			page = 0
			main()
		}
	}()
	defer storeBackData()

	wg.Add(5)
	go getCustomersData()
	go getModelData()
	go getVehicleData()
	go getBookingsData()
	go getPreviousBookingsData()
	wg.Wait()
	go uMDTVD(&MData, &VData)
	go getCarModel(&MData)

	for {
		pages(page)
		if Exit {
			break
		}
	}

}

func userInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func storeBackData() {
	VData.StoreData()
	MData.StoreData()
	CData.StoreData()
	BData.StoreData()
	PBData.StoreData()
}

//1. Search Car for Rent
func modelCatogory() error {
	var Choice string
	var CarofChoice = map[string]*Vehicles{}
	var TStruct *Vehicles

	for {
		fmt.Println("Selection of car classification:")
		for i, v := range carModel {
			fmt.Printf("%v. %v\n", i+1, v)
		}
		fmt.Println("Select Your Choice: (Enter '*' return to Main Menu)")
		input := userInput()
		v, err := checkChoices(input, len(carModel))
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		} else if v == -1 {
			page = 0
			return nil
		}
		Choice = carModel[v-1]
		fmt.Printf("Classification of '%v' has selected...\n", Choice)
		fmt.Println()
		break
	}

	//Search car Model in the choice catogory
	currentNode := VData.Vehicles
	for currentNode != nil {
		if currentNode.LinkToTypeVehicle.Type == Choice && currentNode.Selecting == false && currentNode.Reserved == false && currentNode.Services == false {
			CarofChoice[currentNode.Name] = currentNode
		}
		currentNode = currentNode.NextPage
	}

	//Selection of Model of Car
	if len(CarofChoice) == 0 {
		fmt.Printf("There are no car in this classification! Please choose another one!\n\n")

		modelCatogory()
		return nil
	}

	for {
		fmt.Println("Selection of car:")
		j := 1
		kSlice := []string{}
		for k, _ := range CarofChoice {
			fmt.Printf("%v. %v\n", j, k)
			kSlice = append(kSlice, k)
			j++
		}
		fmt.Println("Select Your Choice: (Enter '*' return to Main Menu)")
		input := userInput()
		v, err := checkChoices(input, len(CarofChoice))
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		}
		fmt.Printf("'%v' has been choosing for car rental..\n", kSlice[v-1])
		TStruct = CarofChoice[kSlice[v-1]]
		TStruct.Selecting = true // lock for other people to sellect
		fmt.Println()
		CarBooking(TStruct) //pass in Vehicle struct node to function, function from booking files
		break
	}
	return nil
}

func checkChoices(input string, choiceNo int) (int, error) {
	v, err := strconv.Atoi(input) //will also check string
	if input == "*" {
		return -1, nil
	} else if err != nil {
		return 0, err
	} else if v <= 0 || v > choiceNo {
		return 0, errors.New("Invalid input!")
	}
	return v, nil
}

func checkYesNo(input string) (int, error) {
	if len(input) == 0 {
		return 0, errors.New("No Input Found!")
	} else if strings.ToUpper(input) == "YES" {
		return 1, nil
	} else if strings.ToUpper(input) == "NO" {
		return -1, nil
	}
	return 0, errors.New("Invalid input!")
}

func checkInt(input string) (int, error) {
	var inputConvertInt int
	inputString, err := checkString(input)
	if inputString == "*" {
		page = 0
		// pages(0)
		return 0, nil
	} else if err != nil {
		return 0, err
	} else {
		inputConvertInt, err = strconv.Atoi(inputString)
		if err != nil {
			return 0, errors.New("Invalid input!")
		}
	}
	return inputConvertInt, nil
}

func checkString(input string) (string, error) { //(Enter '*' return to Main Menu)
	if len(input) == 0 {
		return "", errors.New("No Input Found!")
	} else if input == "*" {
		page = 0
		// pages(0)
		return "*", nil
	}
	return input, nil
}

func viewModel() {
	MData.printAllNodes()
	for {
		fmt.Println("Do you want to modify any items (Yes or Enter '*' return to Main Menu)? ")
		input := userInput()
		if input == "" {
			fmt.Println("No Input Found!")
		} else if strings.ToUpper(input) == "YES" {
			changesRecord := MData.EditModelData()
			if len(changesRecord) != 0 {
				for k, v := range changesRecord {
					fmt.Printf("%v has change to %v\n", k, v)
				}
			} else {
				fmt.Println("No change have made!")
				userInput()
				page = 0
				// pages(0)
				return
			}
		} else if input == "*" {
			page = 0
			return
		} else {
			fmt.Println("Invalid input!")
		}
	}
}

func addModel() {
	fmt.Println("Existing Model Data")
	MData.printAllNodes()
	MData.AddNewMedal()
}

func addVehicles() {
	fmt.Println("List of Existing Vehicles Inventory")
	VData.addAdditinalVehicles()
	page = 0
	// pages(0)
}

func checkCustomerData() (string, bool) {
	var NRIC, Name string
	fmt.Printf("Please enter your NRIC: ")
	NRIC = userInput()
	if NRIC == "" {
		fmt.Println("No Input Found!")
		page = 0
		return "", false
	}
	fmt.Println()
	TF := CData.existingCustomer(NRIC)
	if !TF {
		for {
			fmt.Println("Do you want to be a member? (Yes or No)")
			input := userInput()
			v, err := checkYesNo(input)
			if v == -1 {
				page = 0
				return "", false
			} else if err != nil {
				fmt.Println(err)
				continue
			} else if v == 1 {
				for {
					fmt.Println("Please enter your name: (Enter '*' return to Main Menu)")
					Name = userInput()
					Name, err := checkString(Name)
					if err != nil {
						fmt.Println(err)
						continue
					} else if input == "*" {
						page = 0
						// pages(0)
						return "", false
					}
					msg := CData.addCustomer(NRIC, Name)
					fmt.Println(msg)
					return NRIC, true
				}
			}
			break
		}
	}
	return NRIC, true
}

func priviousBookingRecord(NRIC string) {
	BData.updateBookings()
	VData.updateVehicle()
	PBData.printAllNodes(NRIC)
}

func viewEarning() {
	PBData.printAllNodes("*Owner*")
	total := PBData.totalEarning()
	fmt.Printf("The total period earning: $%.2f\n", total)
	userInput()
}

func delModel() {
	MData.printAllNodes()
	var input string
	for {
		fmt.Println("Please enter the vehicle name to delete: (Enter '*' return to Main Menu)")
		input = userInput()
		input2, err := checkString(input)
		if err != nil {
			fmt.Println(err)
			continue
		} else if input == "*" {
			page = 0
			return
		}
		datafound, err := MData.deleteVehicle(input2)
		if datafound {
			break
		} else {
			continue
		}
	}
	userInput()
	return
}

// fmt.Printf("%T - %v\n", v, v)
