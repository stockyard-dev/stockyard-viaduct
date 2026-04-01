package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-viaduct/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Route{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var rt store.Route;json.NewDecoder(r.Body).Decode(&rt);if rt.Version==""||rt.Upstream==""{writeError(w,400,"version and upstream required");return};s.db.Create(&rt);writeJSON(w,201,rt)}
func(s *Server)handleDeprecate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{SunsetDate string `json:"sunset_date"`};json.NewDecoder(r.Body).Decode(&req);s.db.Deprecate(id,req.SunsetDate);writeJSON(w,200,map[string]string{"status":"deprecated"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
