/* madpi.js */

function calcScore(data) {

  var keys = Object.keys(data);
  
  var total = 0;

  for(var i = 0; i < keys.length; i++) {
    total = total + data[keys[i]];
  }

  return total;

} // calcScore


function gameClockToString(c, mins) {

  return "shit";

  
} // gameClockToString


function shotClockToString(c, seconds) {

  console.log(c.seconds);
  console.log(c.tenths);
  
  return "fag";

} // shotClockToString


function updateScore(team, val) {

  if(team == HOME) {
    document.getElementById("homeScore").innerHTML = val;
  } else {
    document.getElementById("awayScore").innerHTML = val;
  }
  
} // updateScore


function updateTimeouts(team, val) {
  
  var t = null;
  var s = document.createElement("span");

  if(team == HOME) {

    t  = document.getElementById("homeTimeouts");
    s.className = "fa fa-4x fa-circle light-blue";

  } else {

    t  = document.getElementById("awayTimeouts");
    s.className = "fa fa-4x fa-circle light-blue pull-right";

  }

  var count = parseInt(val);

  while(t.firstChild) {
    t.removeChild(t.firstChild);
  }

  for(var j = 0; j < count; j++) {
    
    var x = s.cloneNode(true);
    t.appendChild(x);

  }
  
} // updateTimeouts
  

function updateFouls(team, val) {

  var f = null;

  var small = document.createElement("small");
  var text  = document.createTextNode(FOULS);
  
  small.className = "light-blue";
  small.appendChild(text);
  
  if(team == HOME) {
    f = document.getElementById("homeFoul");
  } else {
    f = document.getElementById("awayFoul");
  }

  f.innerHTML = val + " ";
  f.appendChild(small);

} // updateFouls


function updatePossession(team) {
  
  if(team == HOME) {

    document.getElementById("homePos").className = "col-lg-5 col-md-5 col-sm-5 col-xs-5 shade border2";
    document.getElementById("awayPos").className = "col-lg-5 col-md-5 col-sm-5 col-xs-5 shade";

  } else {

    document.getElementById("awayPos").className = "col-lg-5 col-md-5 col-sm-5 col-xs-5 shade border2";
    document.getElementById("homePos").className = "col-lg-5 col-md-5 col-sm-5 col-xs-5 shade";

  }

} // updatePossession


function updatePeriod(val) {

  var p   = parseInt(val);
  var str = PERIODS[0];

  if(p > 3) {
    str = "OT" + (p - 3);
  } else {
    str = PERIODS[p];
  }

  document.getElementById("period").innerHTML = str;

} // updatePeriod


function updateClock(obj) {

  document.getElementById("shotClock").innerHTML = shotClockToString(obj.shot);
  document.getElementById("gameClock").innerHTML = gameClockToString(obj.game);

} // updateClock


function updateTeam(team, data) {

  updateFouls(team, data.fouls);
  updateTimeouts(team, data.timeouts);
  updateScore(team, calcScore(data.points));

} // team


function updateDisplay(data) {

  var j = JSON.parse(data);

  console.log(j);

  updatePeriod(j.period);

  if(j.possession) {
    updatePossession(HOME);
  } else {
    updatePossession(AWAY);
  }

  updateTeam(HOME, j.home);
  updateTeam(AWAY, j.away);
  
  updateClock(j);

  updateScore(HOME, calcScore(j.home.points));
  updateScore(AWAY, calcScore(j.away.points));

  updateFouls(HOME, j.home.fouls);
  updateFouls(AWAY, j.away.fouls);

  updateTimeouts(HOME, j.home.timeouts, j.settings.timeouts);
  updateTimeouts(AWAY, j.away.timeouts, j.settings.timeouts);

} // updateDisplay
