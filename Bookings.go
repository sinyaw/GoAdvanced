package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var BData = BookingsData{} // to other function for reference
var BLenght [8]int         // for printing format purpuse

type Bookings struct {
	BookingID     string
	CarID         int
	CarName       string
	StartDate     time.Time //same as vehicle's StartDate
	EndDate       time.Time //same as vehicle's EndDate
	CustomerIC    string
	Price         float64
	DaysOfRenting int
	NextPage      *Bookings
}

type BookingsData struct {
	FirstBookings *Bookings
	LastBookings  *Bookings
	TotalBookings int
}

//Read Bookings Data from Excel file by links list formating with Name and ID sorted
func (BD *BookingsData) readData() (BookingsData, error) {
	f, err := excelize.OpenFile("BookingsData.xlsx")
	if err == nil {
		rows := f.GetRows("Sheet1")
		for i, row := range rows {
			newDetail := &Bookings{}
			if i > 0 {
				newDetail.BookingID = row[0]
				newDetail.CarID, _ = strconv.Atoi(row[1])
				newDetail.CarName = row[2]
				newDetail.StartDate = stringToTime(row[3])
				newDetail.EndDate = stringToTime(row[4])
				newDetail.CustomerIC = row[5]
				newDetail.Price, _ = strconv.ParseFloat(row[6], 64)
				newDetail.DaysOfRenting, _ = strconv.Atoi(row[7])

				for i, _ := range BLenght {
					if BLenght[i] < len(row[i]) {
						BLenght[i] = len(row[i])
					}
				}

				BD.addBookings(newDetail)
			}

			BD.TotalBookings++
		}
	} else {
		fmt.Println(err)
	}

	return *BD, nil
}

//Store back Bookings Data to Excel file
func (BD *BookingsData) StoreData() error {
	currentNode := BD.FirstBookings
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("Sheet1")
	if currentNode == nil {
		fmt.Println("No Bookings Data to Store!")
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
	err := xlsx.SaveAs("BookingsData.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//Insert Bookings to Bookings Data by Queues with ID /** realist need to find out issue
func (BD *BookingsData) addBookings(newBooking *Bookings) error {
	if BD.FirstBookings == nil {
		BD.FirstBookings = newBooking
	} else if BD.LastBookings != nil {
		newBooking.NextPage = BD.LastBookings
		BD.LastBookings = newBooking
	} /* else {
	BED := BD.FirstBookings.EndDate
	NBED := newBooking.EndDate
	dy := NBED.Sub(BED).Hours() / 24
	if dy <= 0 {
		newBooking.NextPage = BD.FirstBookings
		BD.FirstBookings = newBooking
	} else {
		currentNode := BD.FirstBookings
		BED := currentNode.NextPage.EndDate
		NBED := newBooking.EndDate
		dy := NBED.Sub(BED).Hours() / 24
		for currentNode.NextPage != nil && dy <= 0 {
			currentNode = currentNode.NextPage
		}
		newBooking.NextPage = currentNode.NextPage
		currentNode.NextPage = newBooking
	}*/
	fmt.Println(newBooking) //** need to delete
	BD.LastBookings = newBooking

	BD.TotalBookings++
	return nil
}

func (BD *BookingsData) deBookings() (string, error) {
	if BD.FirstBookings == nil {
		return "", errors.New("No data!")
	}
	if BD.TotalBookings == 1 {
		BD.FirstBookings = nil
		BD.LastBookings = nil
	} else {
		BD.FirstBookings = BD.FirstBookings.NextPage
	}
	BD.TotalBookings--
	return "", nil
}

//Transfer booking record to previous booking record
func (BD *BookingsData) updateBookings() {
	currentNode := BD.FirstBookings
	fmt.Println(currentNode)
	for currentNode.NextPage != nil {
		fmt.Println("Test Print0")
		if currentNode != nil {
			BED := BD.FirstBookings.EndDate
			fmt.Println("Test Print 1", BED.Sub(time.Now()).Hours()/24)
			fmt.Println(BED, time.Now())
			if (BED.Sub(time.Now()).Hours() / 24) < 0 {
				fmt.Println("Test Print 2", BED.Sub(time.Now()).Hours()/24)
				PBFile.BookingID = currentNode.BookingID
				PBFile.CarID = currentNode.CarID
				PBFile.CarName = currentNode.CarName
				PBFile.StartDate = currentNode.StartDate
				PBFile.EndDate = currentNode.EndDate
				PBFile.CustomerIC = currentNode.CustomerIC
				PBFile.Price = currentNode.Price
				PBFile.DaysOfRenting = currentNode.DaysOfRenting
				PBData.updateBookNPriBookData(PBFile)
				BD.FirstBookings = currentNode.NextPage
			}

			if currentNode.NextPage != nil {
				BED := currentNode.NextPage.EndDate
				fmt.Println(BED.Sub(time.Now()).Hours() / 24)
				if (BED.Sub(time.Now()).Hours() / 24) < 0 {
					PBFile.BookingID = currentNode.NextPage.BookingID
					PBFile.CarID = currentNode.NextPage.CarID
					PBFile.CarName = currentNode.NextPage.CarName
					PBFile.StartDate = currentNode.NextPage.StartDate
					PBFile.EndDate = currentNode.NextPage.EndDate
					PBFile.CustomerIC = currentNode.NextPage.CustomerIC
					PBFile.Price = currentNode.NextPage.Price
					PBFile.DaysOfRenting = currentNode.NextPage.DaysOfRenting
					PBData.updateBookNPriBookData(PBFile)
					currentNode.NextPage = currentNode.NextPage.NextPage
				}
			}
		}
		currentNode = currentNode.NextPage
	}
}

//for Bookings data struct to use for printing it's data
func (D *BookingsData) printAllNodes(NRIC string) error {
	cN := D.FirstBookings
	pi := 0
	q := 0
	s, _ := fmt.Println("Current Booking")
	fmt.Println(strings.Repeat("-", s-1))
	if cN == nil {
		fmt.Println("no item recorded")
		page = 0
		return nil
	} else {
		for cN != nil {

			if cN.CustomerIC == NRIC {
				q++
				pi++
				fmt.Printf("%v.  ", pi)
				fmt.Printf("Ref No.%v", cN.BookingID)
				for j := 0; j < (BLenght[0] - len(cN.BookingID) + 3); j++ {
					fmt.Printf(" ")
				}

				fmt.Printf("%v", cN.CarName)
				for j := 0; j < (BLenght[2] - len(cN.CarName) + 5); j++ {
					fmt.Printf(" ")
				}

				d := cN.StartDate
				m := d.Month().String()
				fmt.Printf("start: %v %v %v\t", d.Day(), m[:3], d.Year())
				d = cN.EndDate
				m = d.Month().String()
				fmt.Printf("end: %v %v %v\t", d.Day(), m[:3], d.Year())
				fmt.Println()
			} else if NRIC == "*Owner*" {
				q++
				pi++
				fmt.Printf("%v.  ", pi)
				fmt.Printf("Ref No.%v", cN.BookingID)
				for j := 0; j < (BLenght[0] - len(cN.BookingID) + 3); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("CarID%v", cN.CarID)
				for j := 0; j < (BLenght[1] - len(strconv.Itoa(cN.CarID)) + 3); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("%v", cN.CarName)
				for j := 0; j < (BLenght[2] - len(cN.CarName) + 5); j++ {
					fmt.Printf(" ")
				}

				d := cN.StartDate
				m := d.Month().String()
				fmt.Printf("start: %v %v %v\t", d.Day(), m[:3], d.Year())
				d = cN.EndDate
				m = d.Month().String()
				fmt.Printf("end: %v %v %v\t", d.Day(), m[:3], d.Year())

				fmt.Printf("%v", cN.CustomerIC)
				for j := 0; j < (BLenght[5] - len(cN.CustomerIC) + 5); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("$%v", cN.Price)
				for j := 0; j < (MLenght[6] - len(strconv.Itoa(int(cN.Price))) + 5); j++ {
					fmt.Printf(" ")
				}
				fmt.Printf("%vdy/s", cN.DaysOfRenting)
				for j := 0; j < (BLenght[1] - len(strconv.Itoa(cN.DaysOfRenting)) + 3); j++ {
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
	userInput()
	page = 0
	return nil
}

func getBookingsData() error {
	defer wg.Done()
	bd := &BookingsData{}
	BD, _ := bd.readData()
	BData = BD
	return nil
}

func CarBooking(v *Vehicles) {
	now := time.Now().AddDate(0, 0, 1)
	Y, M, D := now.Date()
	after := now.AddDate(1, 0, -1)
	Y1, M1, D1 := after.Date()
	BID := ""

	fmt.Println("Please enter period of date for car rental. (Enter '*' return to Main Menu)")
	fmt.Printf("Open period (%v/%v/%v - %v/%v/%v)\n", D, int(M), Y, D1, int(M1), Y1)
	msg := "Start Date (dd/mm/yyyy): "
	sd, sm, sy := checkDateInput(msg, now, after)
	sDate := ReturnADate(sd, sm, sy)
	if sd+sm+sy == 0 {
		return
	}
	s1 := strconv.Itoa(Y)
	s2 := strconv.Itoa(int(M))
	s3 := strconv.Itoa(D - 1)
	s4 := strconv.Itoa(BData.TotalBookings + 1)
	BID = s1 + s2 + s3 + "-" + s4

	//Define variables for checking
	var ed, em, ey int
	var days float64
	var eDate time.Time
	for {
		msg = "End Date (dd/mm/yyyy): "
		ed, em, ey = checkDateInput(msg, now, after)
		eDate = ReturnADate(ed, em, ey)
		days = eDate.Sub(sDate).Hours() / 24
		if days <= 0 {
			fmt.Println("Invalid Date! The end date you choosed is same or early then your rental start date!")
			continue
		}
		break
	}
	fmt.Printf("The period of rental for the car from %v/%v/%v to %v/%v/%v\n", sd, sm, sy, ed, em, ey)
	fmt.Println("(Enter '*' return to Main Menu)")
	fmt.Printf("Submit your NRIC: ")
	input := userInput()
	input, err := checkString(input)
	if err != nil {
		fmt.Println(err)
	} else if input == "*" {
		page = 0
		pages(0)
		return
	}
	NRIC := input

	var p float64
	if days >= 30 {
		p = v.LinkToTypeVehicle.MonthlyCost / 30
	} else if days >= 7 && days < 30 {
		p = v.LinkToTypeVehicle.WeeklyCost / 7
	} else {
		p = v.LinkToTypeVehicle.DailyCost
	}

	//update booking and vehicle data
	newNode := &Bookings{}
	newNode.BookingID = BID
	newNode.CarID = v.ID
	newNode.CarName = v.Name
	newNode.StartDate = sDate
	newNode.EndDate = eDate
	newNode.CustomerIC = NRIC
	newNode.Price = p
	newNode.DaysOfRenting = int(days)

	// dy := sDate.Sub(time.Now()).Hours() / 24
	// if dy <= 3 {
	v.Reserved = true
	v.StartDate = sDate
	v.EndDate = eDate
	v.Selecting = false
	v.Customer = NRIC
	// }

	BData.addBookings(newNode)
	fmt.Printf("Booking for '%v' successfully submitted.", v.Name)
	//return back to main page
	page = 0

}
