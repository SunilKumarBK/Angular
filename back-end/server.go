package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/rs/cors"

    _ "github.com/go-sql-driver/mysql"
)

type Employee struct {
    ID              int    `json:"id"`
    EmpId           int    `json:"empId"`
    FirstName       string `json:"firstName"`
    LastName        string `json:"lastName"`
    Email           string `json:"email"`
    PhoneNo         int    `json:"phoneNo"`
    FatherName      string `json:"fatherName"`
    EmergencyContact int   `json:"emergencyContact"`
    DateOfBirth     string `json:"dateOfBirth"`
    Address         string `json:"address"`
    Experience      bool   `json:"experience"`
    CompanyName string `json:"companyName"`
    Designation string  `json:"designation"`
    JoinDate   string `json:"joinDate"`
    RelievedDate string `json:"relievedDate"`
    TotalDuration string `json:"totalDuration"`
}

type Company struct {
    EmpId       int    `json:"empId"`
   
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    

    rows, err := db.Query(` SELECT 
            e.id,e.empId,e.FirstName, e.LastName, e.Email, e.PhoneNo, 
            e.FatherName, e.EmergencyContact, e.DateOfBirth, e.Address, e.Experience,
            p.CompanyName, p.position, p.startDate, p.endDate, p.duration
        FROM emply e
          JOIN prevcompany p ON e.empId = p.empId`) 
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    defer rows.Close()

    var employees []Employee
    for rows.Next() {
        var emp Employee
        if err := rows.Scan(&emp.ID, &emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.PhoneNo, &emp.FatherName, &emp.EmergencyContact, &emp.DateOfBirth, &emp.Address, &emp.Experience,&emp.CompanyName,&emp.Designation,&emp.JoinDate,&emp.RelievedDate,&emp.TotalDuration); err != nil {
            log.Printf("Error scanning row: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employees = append(employees, emp)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Rows error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Retrieved data: %v", employees)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(employees); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}



func getEmployeeByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    idStr := params["id"]

    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Invalid employee ID: %v", err)
        http.Error(w, "Invalid employee ID", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var emp Employee

    // Check if employee has experience
    err = db.QueryRow("SELECT Experience FROM emply WHERE empId = ?", id).Scan(&emp.Experience)
    if err != nil {
        if err == sql.ErrNoRows {
            // Handle case where employee with given ID doesn't exist
            http.Error(w, "Employee not found", http.StatusNotFound)
            return
        }
        log.Printf("Error checking employee experience: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Fetch employee data based on experience
    if emp.Experience {
        err = db.QueryRow(`SELECT 
                e.ID, 
                e.EmpId, 
                e.FirstName, 
                e.LastName, 
                e.Email, 
                e.PhoneNo, 
                e.FatherName, 
                e.EmergencyContact, 
                e.DateOfBirth, 
                e.Address, 
                e.Experience,
                p.companyName,
                p.position,
                p.startDate,
                p.endDate,
                p.duration
            FROM 
                emply e
            LEFT JOIN
                prevcompany p ON e.empId = p.empId
            WHERE 
                e.empId = ?`, id).Scan(
                    &emp.ID,
                    &emp.EmpId,
                    &emp.FirstName,
                    &emp.LastName,
                    &emp.Email,
                    &emp.PhoneNo,
                    &emp.FatherName,
                    &emp.EmergencyContact,
                    &emp.DateOfBirth,
                    &emp.Address,
                    &emp.Experience,
                    &emp.CompanyName,
                    &emp.Designation,
                    &emp.JoinDate,
                    &emp.RelievedDate,
                    &emp.TotalDuration,
        )
    } else {
        // Fetch basic employee details without company info
        err = db.QueryRow("SELECT * FROM emply WHERE empId = ?", id).Scan(
            &emp.ID,
            &emp.EmpId,
            &emp.FirstName,
            &emp.LastName,
            &emp.Email,
            &emp.PhoneNo,
            &emp.FatherName,
            &emp.EmergencyContact,
            &emp.DateOfBirth,
            &emp.Address,
            &emp.Experience,
        )
    }

    if err != nil {
        log.Printf("Error retrieving employee: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Retrieved employee: %+v", emp)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(emp)
}


// func addEmployeeHandler(w http.ResponseWriter, r *http.Request) {
//     var emp Employee
//     err := json.NewDecoder(r.Body).Decode(&emp)
//     if err != nil {
//         log.Printf("Error decoding request body: %v", err)
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
//     if err != nil {
//         log.Printf("Error opening database: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     defer db.Close()

//     stmt, err := db.Prepare("INSERT INTO emply (empId, firstName, lastName, email, phoneNo, fatherName, emergencyContact, dateOfBirth, address, experience) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
//     if err != nil {
//         log.Printf("Error preparing statement: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     defer stmt.Close()

//     _, err = stmt.Exec(emp.EmpId, emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Experience)
//     if err != nil {
//         log.Printf("Error executing insert statement: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     log.Printf("Employee added successfully: %+v", emp)

//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(emp)
// }

func addEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    var emp Employee
    err := json.NewDecoder(r.Body).Decode(&emp)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Optionally, validate or process emp data here

    db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Printf("Error beginning transaction: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback() // Rollback if any error occurs before commit

    stmt, err := tx.Prepare("INSERT INTO emply (empId, firstName, lastName, email, phoneNo, fatherName, emergencyContact, dateOfBirth, address, experience) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        log.Printf("Error preparing statement: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    // Execute insert into emply table
    _, err = stmt.Exec(emp.EmpId, emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Experience)
    if err != nil {
        log.Printf("Error executing insert statement: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if emp has experience
    // var hasExperience bool
    // err = db.QueryRow("SELECT Experience FROM emply WHERE empId = ?", emp.EmpId).Scan(&hasExperience)
    // // log.Printf(hasExperience,"hasexx")
    // if err != nil {
    //     log.Printf("Error checking employee experience: %v", err)
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    //     return
    // }

    if emp.Experience {
        // Insert into prevcompany table
        companyStmt, err := tx.Prepare("INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)")
        if err != nil {
            log.Printf("Error preparing company statement: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer companyStmt.Close()

        _, err = companyStmt.Exec(emp.CompanyName, emp.Designation, emp.JoinDate, emp.RelievedDate, emp.TotalDuration, emp.EmpId)
        if err != nil {
            log.Printf("Error executing company insert statement: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Employee added successfully: %+v", emp)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(emp)
}




func deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    idStr := params["id"]

    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Invalid employee ID: %v", err)
        http.Error(w, "Invalid employee ID", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

        // Delete from prevcompany table first
        _, err = db.Exec("DELETE FROM prevcompany WHERE empId = ?", id)
        if err != nil {
            log.Printf("Error deleting from prevcompany: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    

    stmt, err := db.Prepare("DELETE FROM emply WHERE empId = ?")
    if err != nil {
        log.Printf("Error preparing delete statement: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    res, err := stmt.Exec(id)
    if err != nil {
        log.Printf("Error executing delete statement: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        log.Printf("No employee found with ID: %d", id)
        http.Error(w, "No employee found with the given ID", http.StatusNotFound)
        return
    }

    log.Printf("Employee deleted successfully: ID %d", id)
    w.WriteHeader(http.StatusNoContent)
}


func updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var emp Employee
    err := json.NewDecoder(r.Body).Decode(&emp)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
    if err != nil {
        log.Printf("Error opening database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Printf("Error beginning transaction: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    _, err = tx.Exec("UPDATE emply SET firstName=?, lastName=?, email=?, phoneNo=?, fatherName=?, emergencyContact=?, dateOfBirth=?, address=?, experience=? WHERE empId=?", emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Experience, id)
    if err != nil {
        log.Printf("Error updating employee: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if emp.Experience {
        _, err = tx.Exec("UPDATE prevcompany SET companyName=?, position=?, startDate=?, endDate=?, duration=? WHERE empId=?", emp.CompanyName, emp.Designation, emp.JoinDate, emp.RelievedDate, emp.TotalDuration, id)
        if err != nil {
            log.Printf("Error updating company: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Employee updated successfully: %+v", emp)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(emp)
}


func main() {
    r := mux.NewRouter()

    r.HandleFunc("/employee", dataHandler).Methods("GET")
    r.HandleFunc("/addemployee", addEmployeeHandler).Methods("POST")
    r.HandleFunc("/getemployee/employee/{id}", getEmployeeByID).Methods("GET")
    r.HandleFunc("/delete/employee/{id}", deleteEmployeeHandler).Methods("DELETE")
    r.HandleFunc("/update/employee/{id}", updateEmployeeHandler).Methods("PUT")


    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:4200"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

    handler := c.Handler(r)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
