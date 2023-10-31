package controllers

import (
	"fmt"
	"go-crud/configs"
	"go-crud/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

func GetCel(c *gin.Context) {
	// Establish a database connection
	db := configs.ConnectDB()

	// offset := 0     // Start at row 10,001
	// limit := 100000 // Retrieve 10,000 rows

	// Query the "cel" table to retrieve data
	query := "SELECT * FROM cel "
	rows, err := db.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	defer rows.Close()

	var results []models.Cel

	for rows.Next() {
		var each = models.Cel{}

		if err := rows.Scan(&each.ID, &each.Eventtype, &each.Eventtime, &each.Cid_name, &each.Cid_num, &each.Cid_ani, &each.Cid_rdnis, &each.Cid_dnid, &each.Exten, &each.Context, &each.Channame, &each.Appname, &each.Appdata, &each.Amaflags, &each.Accountcode, &each.Uniqueid, &each.Linkedid, &each.Peer, &each.Userdeftype, &each.Extra); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		results = append(results, each)
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}

	headerRow := sheet.AddRow()
	headerRow.AddCell().Value = "ID"
	headerRow.AddCell().Value = "Eventtype"
	headerRow.AddCell().Value = "Eventtime"
	headerRow.AddCell().Value = "Cid_name"
	headerRow.AddCell().Value = "Cid_num"
	headerRow.AddCell().Value = "Cid_ani"
	headerRow.AddCell().Value = "Cid_rdnis"
	headerRow.AddCell().Value = "Cid_dnid"
	headerRow.AddCell().Value = "Exten"
	headerRow.AddCell().Value = "Context"
	headerRow.AddCell().Value = "Channame"
	headerRow.AddCell().Value = "Appname"
	headerRow.AddCell().Value = "Appdata"
	headerRow.AddCell().Value = "Amaflags"
	headerRow.AddCell().Value = "Accountcode"
	headerRow.AddCell().Value = "Uniqueid"
	headerRow.AddCell().Value = "Linkedid"
	headerRow.AddCell().Value = "Peer"
	headerRow.AddCell().Value = "Userdeftype"
	headerRow.AddCell().Value = "Extra"

	// var row_count int = 0

	for _, each := range results {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(each.ID)
		row.AddCell().Value = each.Eventtype
		row.AddCell().Value = each.Eventtime
		row.AddCell().Value = each.Cid_name
		row.AddCell().Value = each.Cid_num
		row.AddCell().Value = each.Cid_ani
		row.AddCell().Value = each.Cid_rdnis
		row.AddCell().Value = each.Cid_dnid
		row.AddCell().Value = each.Exten
		row.AddCell().Value = each.Context
		row.AddCell().Value = each.Channame
		row.AddCell().Value = each.Appname
		row.AddCell().Value = each.Appdata
		row.AddCell().Value = strconv.Itoa(each.Amaflags)
		row.AddCell().Value = each.Accountcode
		row.AddCell().Value = each.Uniqueid
		row.AddCell().Value = each.Linkedid
		row.AddCell().Value = each.Peer
		row.AddCell().Value = each.Userdeftype
		row.AddCell().Value = each.Extra
		// row_count++
		// //print .. "code at "" row
		// fmt.Println("code at ", row_count, " row")

	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=GO_data.xlsx")

	// Write the Excel file data to the response
	err = file.Write(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	//send file
}

func ExcelHandler(c *gin.Context) {
	// totalRowsStr := c.Param("total_rows")

	// fmt.Println(totalRowsStr)
	// totalRows, err := strconv.Atoi(totalRowsStr)
	// if err != nil {
	// 	c.String(http.StatusBadRequest, "Invalid total_rows parameter")
	// 	return
	// }

	// if totalRows <= 0 {
	// 	c.String(http.StatusBadRequest, "total_rows must be a positive integer")
	// 	return
	// }

	totalRows := 100000

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("RandomStrings")
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Error creating Excel sheet")
		return
	}

	charSet := "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < totalRows; i++ {
		row := sheet.AddRow()
		for j := 0; j < 20; j++ {
			cell := row.AddCell()
			cell.SetString(charSet)
		}
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=RandomStrings.xlsx")

	err = file.Write(c.Writer)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Error writing Excel file")
		return
	}
}
