/* config.js */

const ADDR          = "127.0.0.1:8000";
const WWW           = "http://" + ADDR;
const GAME          = WWW + "/api/games";
const SOCKET        = "ws://" + ADDR;
const GAME_SOCKET   = SOCKET + "/ws/game";

const HOME          = "HOME";
const AWAY          = "AWAY";
const HOME_SCORE    = "HOME_SCORE";
const HOME_FOUL     = "HOME_FOUL";
const HOME_TIMEOUT  = "HOME_TIMEOUT";
const AWAY_SCORE    = "AWAY_SCORE";
const AWAY_FOUL     = "AWAY_FOUL";
const AWAY_TIMEOUT  = "AWAY_TIMEOUT";

const GAME_SETTINGS     = "GAME_SETTINGS";
const CLOCK             = "CLOCK";
const PERIOD            = "PERIOD";
const POSSESSION_HOME   = "POSSESSION_HOME";
const POSSESSION_AWAY   = "POSSESSION_AWAY";

const FOULS             = "Fouls";

const PERIODS = ["1st", "2nd", "3rd", "4th"];
