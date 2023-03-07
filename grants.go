package mysqlctl

var grants = map[string]string{
	"ALTER":                   "Alter_priv",
	"ALTER ROUTINE":           "Alter_routine_priv",
	"CREATE":                  "Create_priv",
	"CREATE ROUTINE":          "Create_routine_priv",
	"CREATE TEMPORARY TABLES": "Create_tmp_table_priv",
	"CREATE VIEW":             "Create_view_priv",
	"DELETE":                  "Delete_priv",
	"DROP":                    "Drop_priv",
	"EVENT":                   "Event_priv",
	"EXECUTE":                 "Execute_priv",
	"INDEX":                   "Index_priv",
	"INSERT":                  "Insert_priv",
	"LOCK TABLES":             "Lock_tables_priv",
	"REFERENCES":              "References_priv",
	"SELECT":                  "Select_priv",
	"SHOW VIEW":               "Show_view_priv",
	"TRIGGER":                 "Trigger_priv",
	"UPDATE":                  "Update_priv",
}
