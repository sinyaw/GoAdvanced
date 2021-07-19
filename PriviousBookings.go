package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var PBData = PriviousBookingsData{} // to other function for reference
var PBFile = &PriviousBookings{}    // to other function for reference
var PBLenght [8]int                 // for printing format purpuse

//data come from after booking date finish
type PriviousBookings struct {
	BookingID     string
	CarID         int
	CarName       string
	StartDate     time.Time //same as vehicle's StartDate
	EndDate       time.Time //same as vehicle's EndDate
	CustomerIC    string
	Price         float64
	DaysOfRenting int
	NextPage      *PriviousBookings
}

type PriviousBookingsData struct {
	PriviousBookings *PriviousBookings
	TotalBookings    int
}

func (PBD *PriviousBookingsData) readData() (PriviousBookingsData, error) {
	f, err := excelize.OpenFile("PriviousBookingsData.xlsx")
	if err == nil {
		rows := f.GetRows("Sheet1")
		for i, row := range rows {
			newDetail := &PriviousBookings{}
			if i > 0 {
				newDetail.BookingID = row[0]
				newDetail.CarID, _ = strconv.Atoi(row[1])
				newDetail.CarName = row[2]
				newDetail.StartDate = stringToTime(row[3])
				newDetail.EndDate = stringToTime(row[4])
				newDetail.CustomerIC = row[5]
				newDetail.Price, _ = strconv.ParseFloat(row[6], 64)
				newDetail.DaysOfRenting, _ = strconv.Atoi(row[7])

				for i, _ := range PBLenght {
					if PBLenght[i] < len(row[i]) {
						PBLenght[i] = len(row[i])
					}
				}

				PBD.addBookings(newDetail) //need to change to stack
			}

			PBD.TotalBookings++
		}
	} else {
		fmt.Println(err)
	}
	return *PBD, nil
}

func (PBD *PriviousBookingsData) StoreData() error {
	currentNode := PBD.PriviousBookings
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("Sheet1")
	if currentNode == nil {
		fmt.Println("No Previous Booking Data to Store!")
	} else {
		i := 2
		for currentNode != nil {

			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i), currentNode.BookingID)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i), currentNode.CarID)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i), currentNode.CarName)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i), timeToString(currentNode.StartDate))
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i), timeToString(currentNode.EndDate))
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i), currentNode.CustomerIC)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(i), currentNode.Price)
			xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(i), currentNode.DaysOfRenting)

			currentNode = currentNode.NextPage
			i++
		}
	}
	xlsx.SetActiveSheet(index)
	err := xlsx.SaveAs("PriviousBookingsData.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (PBD *PriviousBookingsData) updateBookNPriBookData(B *PriviousBookings) {
	if B != nil {
		newNode := &PriviousBookings{
			BookingID:     B.BookingID,
			CarID:         B.CarID,
			CarName:       B.CarName,
			StartDate:     B.StartDate,
			EndDate:       B.EndDate,
			CustomerIC:    B.CustomerIC,
			Price:         B.Price,
			DaysOfRenting: B.DaysOfRenting,
		}
		PBD.addBookings(newNode)
		CarReturn = append(CarReturn, newNode.CarID)
	}
	page = 0
	return
}

func (PBD *PriviousBookingsData) addBookings(completedBooking *PriviousBookings) error {
	if PBD.PriviousBookings == nil {
		PBD.PriviousBookings = completedBooking
	} else {
		completedBooking.NextPage = PBD.PriviousBookings
		PBD.PriviousBookings = completedBooking
	}
	PBD.TotalBookings++
	return nil
}

func (D *PriviousBookingsData) printAllNodes(NRIC string) error {
	cN := D.PriviousBookings
	pi := 0
	q := 0
	s, _ := fmt.Println("Previous Booking")
	fmt.Println(strings.Repeat("-", s-1))
	if cN == nil {
		fmt.Println("no item recorded")
	} else {

		for cN != nil {
			if cN.CustomerIC == NRIC {
				pi++
				q++
				fmt.Printf("%v.  ", pi)
				fmt.Printf("Ref No.%v", cN.BookingID)
				for j := 0; j < (PBLenght[0] - len(cN.BookingID) + 3); j++ {
					fmt.Printf(" ")
				}

				fmt.Printf("%v", cN.CarName)
				for j := 0; j < (PBLenght[2] - len(cN.CarName) + 5); j++ {
					fmt.Printf(" ")
				}

				d := cN.StartDate
				m := d.Month().String()
				fmt.Printf("start: %v %v %v   ", d.Day(), m[:3], d.Year())
				d = cN.EndDate
				m = d.Month().String()
				fmt.Printf("end: %v %v %v   ", d.Day(), m[:3], d.Year())
				fmt.Println()
			} else if NRIC == "*Owner*" {
				pi++
				q++
				fmt.Printf("%v.  ", pi)
				fmt.Printf("Ref No.%v", cN.BookingID)
				for j := 0; j < (PBLenght[0] - len(cN.BookingID) + 3); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("CarID%v", cN.CarID)
				for j := 0; j < (PBLenght[1] - len(strconv.Itoa(cN.CarID)) + 3); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("%v", cN.CarName)
				for j := 0; j < (PBLenght[2] - len(cN.CarName) + 5); j++ {
					fmt.Printf(" ")
				}

				d := cN.StartDate
				m := d.Month().String()
				fmt.Printf("start: %v %v %v   ", d.Day(), m[:3], d.Year())
				d = cN.EndDate
				m = d.Month().String()
				fmt.Printf("end: %v %v %v   ", d.Day(), m[:3], d.Year())

				fmt.Printf("%v ", cN.CustomerIC)
				for j := 0; j < (PBLenght[5] - len(cN.CustomerIC) + 5); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("  $%v", cN.Price)
				for j := 0; j < (MLenght[6] - len(strconv.Itoa(int(cN.Price))) + 5); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("%vdy/s", cN.DaysOfRenting)
				for j := 0; j < (PBLenght[1] - len(strconv.Itoa(cN.DaysOfRenting)) + 3); j++ {
					fmt.Printf(" ")
				}
				fmt.Println()
			}
			cN = cN.NextPage
		}
		if q == 0 {
			fmt.Println("no item recorded")
		}
	}
	page = 0
	return nil
}

func (D *PriviousBookingsData) totalEarning() (sum float64) {
	currentNode := D.PriviousBookings
	for currentNode != nil {
		sum += (currentNode.Price * float64(currentNode.DaysOfRenting))
		currentNode = currentNode.NextPage
	}
	return
}

func getPreviousBookingsData() error {
	defer wg.Done()
	pbd := &PriviousBookingsData{}
	PBD, _ := pbd.readData()
	PBData = PBD
	return nil
}
