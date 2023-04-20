package repository

import (
	"database/sql"
	"encoding/xml"
	"net/http"
)

// SDNList структура для представления списка записей в таблице people
type SDNList struct {
	XMLName xml.Name `xml:"sdnList"`
	SDNs    []Person `xml:"sdnEntry"`
}

// Person структура для представления записей в таблице people
type Person struct {
	UID       int    `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	SDNType   string `json:"sdn_type"`
}

// InsertPerson вставляет запись о человеке в таблицу people в БД
func InsertPerson(db *sql.DB, person Person) (int, error) {
	// Выполняем SQL-запрос на вставку записи
	stmt, err := db.Prepare("SELECT insert_person($1, $2, $3)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var uid int
	// Извлекаем UID созданной записи о человеке
	err = stmt.QueryRow(person.FirstName, person.LastName, person.SDNType).Scan(&uid)
	if err != nil {
		return 0, err
	}

	return uid, nil
}

// GetPeople возвращает всех людей из таблицы people в БД, соответствующих заданным параметрам
func GetPeople(db *sql.DB, name string, matchType string) ([]Person, error) {
	var people []Person

	// Выполняем SQL-запрос на получение записей из таблицы people
	stmt, err := db.Prepare("SELECT * FROM get_people($1, $2)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name, matchType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Обходим результаты запроса и сохраняем их в структуры Person
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.UID, &p.FirstName, &p.LastName, &p.SDNType)
		if err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	return people, nil
}

// Update Получает данные из https://www.treasury.gov/ofac/downloads/sdn.xml и сохраняет их в локальную базу данных PostgreSQL
func Update(db *sql.DB) error {
	// Получаем данные из внешнего источника
	data, err := fetchData()
	if err != nil {
		return err
	}

	// Удаляем все записи из таблицы people, чтобы заменить их на новые
	_, err = db.Exec("DELETE FROM people")
	if err != nil {
		return err
	}

	// Сохраняем записи с sdnType=Individual в таблицу people
	for _, item := range data {
		if item.SDNType == "Individual" {
			person := Person{
				FirstName: item.FirstName,
				LastName:  item.LastName,
				SDNType:   item.SDNType,
			}
			if _, err := InsertPerson(db, person); err != nil {
				return err
			}
		}
	}
	return nil
}

func fetchData() ([]Person, error) {
	resp, err := http.Get("https://www.treasury.gov/ofac/downloads/sdn.xml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Преобразуем XML в структуру данных
	var result SDNList
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.SDNs, nil
}
