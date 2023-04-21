package repository

import (
	"database/sql"
	"encoding/xml"
	"github.com/lib/pq"
	"net/http"
)

// sdnList структура для представления списка записей в таблице people
type sdnList struct {
	XMLName xml.Name `xml:"sdnList"`
	SDNs    []Person `xml:"sdnEntry"`
}

// Person структура для представления записей в таблице people
type Person struct {
	UID       int    `json:"uid"       xml:"uid"`
	FirstName string `json:"firstName" xml:"firstName"`
	LastName  string `json:"lastName"  xml:"lastName"`
	SDNType   string `json:"sdnType"   xml:"sdnType"`
}

// InsertPerson вставляет запись о человеке в таблицу people в БД
func InsertPerson(db *sql.DB, person Person) (int, error) {
	// Выполняем SQL-запрос на вставку записи
	stmt, err := db.Prepare("SELECT insert_person($1, $2, $3, $4)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	stmt.QueryRow(person.UID, person.FirstName, person.LastName, person.SDNType)

	return person.UID, nil
}

func InsertPeople(db *sql.DB, people []Person) error {
	if len(people) == 0 {
		return nil
	}

	// Настройка COPY FROM запроса
	schemaName := "public"
	tableName := "people"
	columns := []string{"uid", "first_name", "last_name", "sdn_type"}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyInSchema(schemaName, tableName, columns...))
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			panic(err)
		}
	}()

	//loop through an array of struct filled with data, or read from a file
	for _, row := range people {
		_, err := stmt.Exec(row.UID, row.FirstName, row.LastName, row.SDNType)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
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

// Update Получает данные из https://www.treasury.gov/ofac/downloads/sdn.xml и сохраняет их в локальную базу данных
// PostgreSQL. Возвращает список полученных персон.
func Update(db *sql.DB) (people []Person, err error) {
	// Получаем данные из внешнего источника
	data, err := fetchData()
	if err != nil {
		return people, err
	}

	// Удаляем все записи из таблицы people, чтобы заменить их на новые
	_, err = db.Exec("DELETE FROM people")
	if err != nil {
		return people, err
	}

	// Сохраняем записи с sdnType=Individual в таблицу people
	for _, item := range data {
		if item.SDNType == "Individual" {
			person := Person{
				UID:       item.UID,
				FirstName: item.FirstName,
				LastName:  item.LastName,
				SDNType:   item.SDNType,
			}
			people = append(people, person)
		}
	}

	if err := InsertPeople(db, people); err != nil {
		return people, err
	}
	return people, nil
}

func fetchData() ([]Person, error) {
	resp, err := http.Get("https://www.treasury.gov/ofac/downloads/sdn.xml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Преобразуем XML в структуру данных
	var result sdnList
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.SDNs, nil
}
