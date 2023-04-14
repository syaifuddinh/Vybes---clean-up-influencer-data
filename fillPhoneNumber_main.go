package main

import (
    "database/sql"
    "fmt"
    "strings"
    _ "github.com/go-sql-driver/mysql"
)

func connectToMySQL() (*sql.DB, error) {
    // Replace the connection parameters with your own MySQL configuration
    db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/vybesrawinfluencer")
    if err != nil {
        return nil, err
    }
    return db, nil
}

type InfluencerInterface struct {
	URI string
	PhoneNumber string
}

func executeQuery(queryString string) (bool) {
	db, err := connectToMySQL()
	if err != nil {
        return false
    }
    defer db.Close()
    // Define the query to retrieve the URI and phoneNumber columns from the influencerTable

    // Execute the query and get the result set
    fmt.Println(queryString)
    __, err := db.Exec(queryString)
    if err != nil {
        return false
    }
    fmt.Println(__)

    return true
}

func getDataFromQuery(queryString string) ([]InfluencerInterface, error) {
	db, err := connectToMySQL()
	if err != nil {
        return nil, err
    }
    defer db.Close()
    // Define the query to retrieve the URI and phoneNumber columns from the influencerTable

    // Execute the query and get the result set
    rows, err := db.Query(queryString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Create a slice to hold the results
    influencers := []InfluencerInterface{}

    // Loop through the rows and create an InfluencerInterface struct for each row
    for rows.Next() {
        influencer := InfluencerInterface{}

        // Scan the values from the row into the Influencer struct
        err := rows.Scan(&influencer.URI, &influencer.PhoneNumber)
        if err != nil {
            return nil, err
        }

        // Add the Influencer struct to the results slice
        influencers = append(influencers, influencer)
    }

    // Check for errors that may have occurred during iteration over the rows
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Return the list of Influencer structs
    return influencers, nil
}

func main() {
	var queryString string
	var tableName string
	var uri string
	var phoneNumber string
	var updateCommand string
	queryString = "SELECT rawinfluencer.uri, rawinfluencer.phone_number AS phoneNumber FROM rawinfluencer JOIN ( SELECT uri FROM (SELECT uri, COUNT(id) AS userAmount FROM rawinfluencer GROUP BY uri ) AS influencerGroup WHERE userAmount > 1) AS duplicatedInfluencers ON duplicatedInfluencers.uri = rawinfluencer.uri WHERE phone_number IS NOT NULL"

	influencers, isSuccess := getDataFromQuery(queryString) 
	for _, influencer := range influencers {
		uri = influencer.URI
		phoneNumber = strings.ReplaceAll(influencer.PhoneNumber, "'", "\"")
		tableName = "rawinfluencer"
		updateCommand = "UPDATE " + tableName + " SET phone_number = '" + phoneNumber + "' WHERE uri = '" + uri + "'"
		executeQuery(updateCommand)
	}
	fmt.Println(isSuccess)
}
