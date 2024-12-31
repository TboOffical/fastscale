package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TboOffical/fastscale/commander/utils"
	"github.com/lib/pq"
	"log"
	"os"
)

var (
	AllowedDataTypes  = []string{"string", "int", "bool", "date"}
	KindToSqlMappings = map[string]string{
		"string": "text",
		"int":    "bigint",
		"bool":   "boolean",
		"date":   "date",
	}
)

type ModelSchemaPart struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
	LinkTo   string `json:"link_to"`
	OldName  string `json:"old_name"` //Used if a name needs to be changed, tells the system what the original name was so it can rename instead of drop.
}

// ModelSchema is an individual field in a model.
type ModelSchema struct {
	Name  string            `json:"name"`
	Parts []ModelSchemaPart `json:"parts"`
}

// Model contain the entire model
type Model struct {
	Name    string        `json:"name"`
	Author  string        `json:"author"`
	Version string        `json:"version"`
	Schemas []ModelSchema `json:"schemas"`
}

type ModelManager struct {
	CurrentModel    Model
	CurrentDatabase *sql.DB
}

// schemaToSqlCreate takes a file object and returns the sql code required to create a
// table to reflect it. This does not update the table, that is another function.
func (m *ModelManager) schemaToSqlCreate(f ModelSchema) (string, error) {
	var finalDoc string
	finalDoc += fmt.Sprintln("CREATE TABLE", "fastscale."+f.Name, "(")

	for i, part := range f.Parts {
		if !utils.TableColNameCheck(part.Name) {
			return "", errors.New("part name is invalid, rules are no special characters, spaces, or dashes")
		}

		if !utils.IsInListStrings(part.Kind, AllowedDataTypes) {
			return "", errors.New("part kind is not one of the valid data types")
		}

		finalDoc += fmt.Sprint(part.Name, " ", KindToSqlMappings[part.Kind], utils.If(i == len(f.Parts)-1, "\n", ",\n")) //only put a comma if we are on the last part
	}

	finalDoc += fmt.Sprintln(");")
	return finalDoc, nil
}

// syncTables syncs and updates the tables in the database
func (m *ModelManager) syncTables() error {
	//create/make sure the fastscale schema exists
	_, err := m.CurrentDatabase.Exec("CREATE SCHEMA IF NOT EXISTS fastscale")
	if err != nil {
		return err
	}

	for _, schema := range m.CurrentModel.Schemas {
		//For each schema, get the current structure of the corresponding database table and see if it matches up
		//sometimes we will need to create the table, and sometimes we just need to update the table.

		log.Print("Syncing ", schema.Name, "...")

		create := false //true=create new table, false=update

		rows, err := m.CurrentDatabase.Query(fmt.Sprintf("SELECT * FROM fastscale.%s LIMIT %d", schema.Name, 1))
		var pqE *pq.Error
		if err != nil {
			errors.As(err, &pqE)
			if pqE.Code == "42P01" {
				//Table does not exist
				create = true
			} else {
				return err
			}
		}

		if create {
			//create the table because it does not already exist
			generatedSql, err := m.schemaToSqlCreate(schema)
			if err != nil {
				return err
			}

			_, err = m.CurrentDatabase.Exec(generatedSql)
			if err != nil {
				return err
			}

			log.Print(fmt.Sprint("New table created ", schema.Name))

		} else {
			//update the table

			cols, err := rows.Columns()
			if err != nil {
				return err
			}

			//loop through all parts of the schema to make sure they are represented in the table
			//if not, add them. If a name needs to be changed via OldName, that will be done too.

			for _, part := range schema.Parts {
				//Name needs to be changed
				if len(part.OldName) != 0 && part.OldName != part.Name {
					//change name
					_, err := m.CurrentDatabase.Exec(fmt.Sprintf("ALTER TABLE fastscale.%s RENAME COLUMN %s to %s", schema.Name, part.OldName, part.Name))
					if err != nil {
						return err
					}

					log.Print(fmt.Sprint("*Changed name of ", part.OldName, " to ", part.Name, " in ", schema.Name))
					continue
				}

				//check to see if the part already has a row
				alreadyThere := utils.IsInListStrings(part.Name, cols)

				if !alreadyThere {
					//we are here if the part is newly added to the schema and the table needs to be updated to include it

					//first make the schema name is acceptable
					acceptable := utils.TableColNameCheck(part.Name)
					if !acceptable {
						return errors.New("column name is not acceptable")
					}

					//and that the data type is valid
					acceptable = utils.IsInListStrings(part.Kind, AllowedDataTypes)
					if err != nil {
						return errors.New("data type is not on the list of valid data types")
					}

					_, err = m.CurrentDatabase.Exec(fmt.Sprintf("ALTER TABLE fastscale.%s ADD %s %s", schema.Name, part.Name, KindToSqlMappings[part.Kind]))
					if err != nil {
						return err
					}

					//print a little message to signify that a model was updated
					log.Print(fmt.Sprint("+"+part.Name, " added to schema ", schema.Name))
				}
			}

		}

		//sql, err := m.schemaToSqlCreate(field)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//log.Println(sql)
	}

	return nil
}

// loadModel loads a model from a file into the ModelManager
func (m *ModelManager) loadModel(file string, database *sql.DB) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &m.CurrentModel)
	if err != nil {
		return err
	}

	m.CurrentDatabase = database

	err = m.syncTables()
	if err != nil {
		return err
	}

	log.Print("Model ", m.CurrentModel.Name, " Loaded Successfully")

	return nil
}

// saveToFile saves the configuration loaded into the CurrentModel field to the model.json file in the conf folder
// this can be dangerous as it overwrites the running configuration and the next time the model is reloaded this model will be used instead
// always use pushChanges() instead to sync a release in the database
func (m *ModelManager) saveToFile() error {
	jsonData, err := json.MarshalIndent(m.CurrentModel, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile("./conf/model.json", jsonData, 7777)
	if err != nil {
		return err
	}

	log.Print("Model file updated!")
	return nil
}

// createDefaultModel creates a model if there is not one detected on system statup
func (m *ModelManager) createDefaultModel() error {
	//Create the model in the CurrentModel Object
	cm := &m.CurrentModel

	cm.Name = "Test Model"
	cm.Author = "The Fastscale Team"
	cm.Version = fmt.Sprint("v", FastscaleVersonMajor, ".", FastscaleVersionMinor, ".", FastscaleVersionPatch)

	cm.Schemas = append(cm.Schemas, ModelSchema{
		Name: "book",
		Parts: []ModelSchemaPart{
			{
				Name:     "title",
				Kind:     "string",
				Required: true,
			},
			{
				Name:     "pages",
				Kind:     "int",
				Required: false,
			},
			{
				Name:     "description",
				Kind:     "string",
				Required: false,
			},
			{
				Name:     "author",
				Kind:     "string",
				Required: true,
			},
		},
	})

	err := m.saveToFile()
	if err != nil {
		return err
	}

	log.Print("Default model created")
	return nil
}
