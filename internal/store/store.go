package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Route struct{ID int64 `json:"id"`;Version string `json:"version"`;Upstream string `json:"upstream"`;Description string `json:"description"`;Deprecated bool `json:"deprecated"`;SunsetDate string `json:"sunset_date"`;RequestCount int64 `json:"request_count"`;CreatedAt time.Time `json:"created_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"viaduct.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS routes(id INTEGER PRIMARY KEY AUTOINCREMENT,version TEXT NOT NULL UNIQUE,upstream TEXT NOT NULL,description TEXT DEFAULT '',deprecated INTEGER DEFAULT 0,sunset_date TEXT DEFAULT '',request_count INTEGER DEFAULT 0,created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)}
func(db *DB)Create(r *Route)error{dep:=0;if r.Deprecated{dep=1};res,err:=db.Exec(`INSERT INTO routes(version,upstream,description,deprecated,sunset_date)VALUES(?,?,?,?,?)`,r.Version,r.Upstream,r.Description,dep,r.SunsetDate);if err!=nil{return err};r.ID,_=res.LastInsertId();return nil}
func(db *DB)List()([]Route,error){rows,_:=db.Query(`SELECT id,version,upstream,description,deprecated,sunset_date,request_count,created_at FROM routes ORDER BY version DESC`);defer rows.Close();var out[]Route;for rows.Next(){var r Route;var dep int;rows.Scan(&r.ID,&r.Version,&r.Upstream,&r.Description,&dep,&r.SunsetDate,&r.RequestCount,&r.CreatedAt);r.Deprecated=dep==1;out=append(out,r)};return out,nil}
func(db *DB)Deprecate(id int64,sunsetDate string){db.Exec(`UPDATE routes SET deprecated=1,sunset_date=? WHERE id=?`,sunsetDate,id)}
func(db *DB)IncrementCount(version string){db.Exec(`UPDATE routes SET request_count=request_count+1 WHERE version=?`,version)}
func(db *DB)Delete(id int64){db.Exec(`DELETE FROM routes WHERE id=?`,id)}
func(db *DB)Stats()(map[string]interface{},error){var total,deprecated int;db.QueryRow(`SELECT COUNT(*) FROM routes`).Scan(&total);db.QueryRow(`SELECT COUNT(*) FROM routes WHERE deprecated=1`).Scan(&deprecated);return map[string]interface{}{"routes":total,"deprecated":deprecated},nil}
