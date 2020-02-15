package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
)
const SQL_INSERT_FILE string = "insert into WorkoutFile(file_name,file_path,status,user_id,version_nb) values (?,?,?,?,?)"
type FileRepository struct {
	connection connections.IMysqlConnection
}

func NewFileRepository(connection connections.IMysqlConnection) *FileRepository {
	return &FileRepository{connection}
}

func (repo *FileRepository) SaveFile(file *models.WorkoutFile) (error) {
	db, err := repo.connection.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}
	
	defer db.Close()

	stmt, err := db.Prepare(SQL_INSERT_FILE);

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(&file.File_Name, &file.File_Path, &file.Status, &file.User_Id, &file.Version_Nb)
	
	if err != nil {
		panic(err.Error())
	}

	rowCnt, err := res.RowsAffected()

	if err != nil {
		panic(err.Error())
	}

	if rowCnt != 1 {
		return errors.New("Failure in sql execution")
	}

	id, err := res.LastInsertId();
	
	if err != nil {
		panic(err.Error())
	}
	file.Workout_File_Id = id;
	return nil

}