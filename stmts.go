package main

const (

	GameCreate = "INSERT into games DEFAULT VALUES"

	GameGet = "SELECT " +
	  "data, created " +
		"FROM games " +
		"WHERE id=?"

	GamesGet = "SELECT " +
	  "id, data, status, created, updated " + 
		"FROM games"

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
