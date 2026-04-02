package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Proxy struct{
	ID string `json:"id"`
	Name string `json:"name"`
	ListenPort int `json:"listen_port"`
	TargetURL string `json:"target_url"`
	Enabled string `json:"enabled"`
	RequestCount int `json:"request_count"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"viaduct.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS proxies(id TEXT PRIMARY KEY,name TEXT NOT NULL,listen_port INTEGER DEFAULT 0,target_url TEXT DEFAULT '',enabled TEXT DEFAULT 'true',request_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Proxy)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO proxies(id,name,listen_port,target_url,enabled,request_count,created_at)VALUES(?,?,?,?,?,?,?)`,e.ID,e.Name,e.ListenPort,e.TargetURL,e.Enabled,e.RequestCount,e.CreatedAt);return err}
func(d *DB)Get(id string)*Proxy{var e Proxy;if d.db.QueryRow(`SELECT id,name,listen_port,target_url,enabled,request_count,created_at FROM proxies WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.ListenPort,&e.TargetURL,&e.Enabled,&e.RequestCount,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Proxy{rows,_:=d.db.Query(`SELECT id,name,listen_port,target_url,enabled,request_count,created_at FROM proxies ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Proxy;for rows.Next(){var e Proxy;rows.Scan(&e.ID,&e.Name,&e.ListenPort,&e.TargetURL,&e.Enabled,&e.RequestCount,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM proxies WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM proxies`).Scan(&n);return n}
