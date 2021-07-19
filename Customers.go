package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var CData = CustomersData{} // to other function for reference

type Customers struct {
	NRIC     string
	Name     string
	NextPage *Customers
}

type CustomersData struct {
	Customers      *Customers
	TotalCustomers int
}

//Read Customers Data from Excel file by links list formating with Name and ID sorted
func (CD *CustomersData) readData() (CustomersData, error) {
	f, err := excelize.OpenFile("CustomersData.xlsx")
	if err == nil {
		rows := f.GetRows("Sheet1")
		for i, row := range rows {
			newDetail := &Customers{}
			if i > 0 {
				newDetail.NRIC = row[0]
				newDetail.Name = row[1]

				CD.addToDataStructByID(newDetail)
			}

			CD.TotalCustomers++
		}
	} else {
		fmt.Println(err)
	}

	return *CD, nil
}

//Store back Customers Data to Excel file
func (CD *CustomersData) StoreData() error {
	currentNode := CD.Customers
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("Sheet1")
	if currentNode == nil {
		fmt.Println("No Customers Data to Store!")
	} else {
		i := 2
		for currentNode != nil {

			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), currentNode.NRIC)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), currentNode.Name)

			currentNode = currentNode.NextPage
			i++
		}
	}
	xlsx.SetActiveSheet(index)
	err := xlsx.SaveAs("CustomersData.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//Insert Customers to Customers Data by sorting with name
func (CD *CustomersData) addToDataStructByID(newDetail *Customers) {
	if CD.Customers == nil {
		CD.Customers = newDetail
		// fmt.Println("a", newDetail.Name, MD.TypeVehicle.Name)
	} else if strings.ToUpper(CD.Customers.NRIC) > strings.ToUpper(newDetail.NRIC) {
		newDetail.NextPage = CD.Customers
		CD.Customers = newDetail
		// fmt.Println("b", newDetail.Name, MD.TypeVehicle.Name)
	} else {
		currentNode := CD.Customers
		for currentNode.NextPage != nil {
			// fmt.Println(i, newDetail.Name > currentNode.Name)
			if strings.ToUpper(currentNode.NextPage.NRIC) > strings.ToUpper(newDetail.NRIC) {
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

//for Customers data struct to use for printing it's data
func (D *CustomersData) printAllNodes() error {
	currentNode := D.Customers
	if currentNode == nil {
		fmt.Println("no item recorded")
	} else {
		for currentNode != nil {
			fmt.Println(currentNode)
			currentNode = currentNode.NextPage
		}
	}
	page = 0
	return nil
}

func (D *CustomersData) existingCustomer(input string) bool {
	currentNode := D.Customers
	NoUser := "Your detail is not in the data yet!!"
	if currentNode == nil {
		fmt.Println(NoUser)
	} else {
		for currentNode != nil {
			if currentNode.NRIC == input {
				return true
			}
			currentNode = currentNode.NextPage
		}
		fmt.Println(NoUser)
	}
	return false
}

//to set 1st Node for Vehicle Data to 'Vdata'
func getCustomersData() error {
	defer wg.Done()
	cd := &CustomersData{}
	CD, _ := cd.readData()
	CData = CD
	return nil
}

func (CD *CustomersData) addCustomer(NRIC string, Name string) string {
	newDetail := &Customers{NRIC, Name, nil}
	CD.addToDataStructByID(newDetail)
	return "Successfully added!"
}
