package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"errors"
	"database/sql"
	"fmt"
	"time"
)
const (
	
	SQL_INSERT_FILE string = "insert into WorkoutFile(file_name,file_path,status,user_id,version_nb,cre_ts) values (?,?,?,?,?,?)"
	SQL_UPDATE_FILE_STATUS string = "update WorkoutFile set version_nb=version_nb+1, status=? where version_nb=? and status=? and workout_file_id=?"
	SQL_GET_FILE string = "select * from WorkoutFile where workout_file_id=?"
)
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

	res, err := stmt.Exec(&file.File_Name, &file.File_Path, &file.Status, &file.User_Id, &file.Version_Nb, time.Now())
	
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

func (repo *FileRepository) UpdateFileStatus(file *models.WorkoutFile, newStatus string) (models.WorkoutFile, error) {
	db, err := repo.connection.GetDBObject();
	
	if err != nil {
		panic(err.Error())
	}
	
	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		panic(err.Error())
	}

	res, err := tx.Exec(SQL_UPDATE_FILE_STATUS, newStatus,&file.Version_Nb,&file.Status, &file.Workout_File_Id)

	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	rows, err := res.RowsAffected();

	if err != nil {
		tx.Rollback()
		return models.WorkoutFile{}, err
	}

	if rows != 1 {
		tx.Rollback()
		return models.WorkoutFile{}, errors.New("Version mismatch no row returned")
	}

	var newWorkoutFile models.WorkoutFile

	row := tx.QueryRow(SQL_GET_FILE, file.Workout_File_Id)

	err=row.Scan(&newWorkoutFile.Workout_File_Id, &newWorkoutFile.File_Name, &newWorkoutFile.File_Path, &newWorkoutFile.Status, &newWorkoutFile.User_Id, &newWorkoutFile.Version_Nb, &newWorkoutFile.Cre_Ts, &newWorkoutFile.Mod_Ts);
	fmt.Printf("File being queries: %+v", newWorkoutFile)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return models.WorkoutFile{}, errors.New("File is missing")
		}
		panic(err.Error())
	}
	tx.Commit()
	return newWorkoutFile, nil;


}