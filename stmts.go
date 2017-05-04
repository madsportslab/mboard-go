package main

const (

	GameCreate = "INSERT into games DEFAULT VALUES"

  GameDelete = "DELETE from games WHERE id=?"
	
	GameGet = "SELECT " +
	  "id, data, status, created, updated " +
		"FROM games " +
		"WHERE id=?"

	GamesGet = "SELECT " +
	  "id, data, status, created, updated " + 
		"FROM games " +
		"ORDER BY created DESC"

	GameUpdate = "UPDATE games " +
	  "SET data=?, updated=CURRENT_TIMESTAMP, status=? " +
		"WHERE id=?"

	LogCreate = "INSERT into logs" +
	  "(game_id, data) " + 
		"VALUES ($1, $2)"
	
	LogGet = "SELECT " +
	  "id, data, created, updated " +
		"FROM logs " + 
		"WHERE game_id=?"

)
