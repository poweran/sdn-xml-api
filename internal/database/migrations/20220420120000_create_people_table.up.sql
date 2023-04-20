-- Создание таблицы people:
CREATE TABLE IF NOT EXISTS people (
                                      uid SERIAL PRIMARY KEY,
                                      first_name VARCHAR(255) NOT NULL,
                                      last_name VARCHAR(255) NOT NULL,
                                      sdn_type VARCHAR(255) NOT NULL
);

-- Создание индексов для таблицы people:
CREATE INDEX IF NOT EXISTS people_first_name_idx ON people (first_name);
CREATE INDEX IF NOT EXISTS people_last_name_idx ON people (last_name);
CREATE INDEX IF NOT EXISTS people_sdn_type_idx ON people (sdn_type);

-- Создание функции insert_person для вставки данных в таблицу people:
CREATE OR REPLACE FUNCTION insert_person(p_first_name VARCHAR(255), p_last_name VARCHAR(255), p_sdn_type VARCHAR(255))
    RETURNS VOID AS $$
BEGIN
    INSERT INTO people (first_name, last_name, sdn_type) VALUES (p_first_name, p_last_name, p_sdn_type);
END;
$$ LANGUAGE plpgsql;

-- Создание функции get_people для получения списка людей из таблицы people по заданным параметрам:
CREATE OR REPLACE FUNCTION get_people(p_name VARCHAR(255), p_type VARCHAR(255))
    RETURNS TABLE (uid INT, first_name VARCHAR(255), last_name VARCHAR(255)) AS $$
BEGIN
    IF p_type = 'strong' THEN
        RETURN QUERY SELECT uid, first_name, last_name FROM people WHERE (first_name || ' ' || last_name) = p_name;
    ELSEIF p_type = 'weak' THEN
        RETURN QUERY SELECT uid, first_name, last_name FROM people WHERE first_name ILIKE ('%' || p_name || '%') OR last_name ILIKE ('%' || p_name || '%');
    ELSE
        RETURN QUERY SELECT uid, first_name, last_name FROM people;
    END IF;
END;
$$ LANGUAGE plpgsql;
