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

  var delta   = mins * 60 - c.seconds;
  var ndelta  = delta - 1;
  var seconds = delta % 60;
  var minutes = Math.floor(delta/60);
  var tenths  = 10 - c.tenths;
  console.log("delta: " + delta + " seconds: " + seconds + " minutes: " +
    minutes + " tenths: " + tenths);

  if(delta == 60) {

    if(minutes == 1) {

      if(tenths == 10) {
        return minutes + ":00";
      } else {
        return ndelta + "." + tenths;
      }
      
    } else {
      return minutes + ":59." + tenths;
    }
    
    return str;

  } else if(minutes == 0) {

    if(ndelta == -1) {
      return "0.0";
    } else if(tenths == 10) {
      return ndelta + ".0";
    } else {
      return ndelta + "." + tenths;
    }
    
  } else if(seconds < 10) {
    return minutes + ":0" + ndelta;
  } else {
    return minutes + ":" + ndelta;
  }

} // gameClockToString


function shotClockToString(c, secs) {

  return secs - c.seconds;
  
} // shotClockToString


function updateScore(team, val) {

  if(team == HOME) {
    document.getElementById("homeScore").innerHTML = val;
  } else {
    document.getElementById("awayScore").innerHTML = val;
  }
  
} // updateScore


function updateTimeouts(team, val) {
  
  var t   = null;
  var h5  = document.createElement("h5");
  var s   = document.createElement("span");

  if(team == HOME) {

    t  = document.getElementById("homeTimeouts");
    s.className = "fa fa-5x fa-circle light-blue pull-right";

  } else {

    t  = document.getElementById("awayTimeouts");
    s.className = "fa fa-5x fa-circle light-blue pull-left";

  }

  var count = parseInt(val);

  while(t.firstChild) {
    t.removeChild(t.firstChild);
  }

  for(var j = 0; j < count; j++) {
    
    var x = s.cloneNode(true);
    h5.appendChild(x);

  }

  t.appendChild(h5);
  
} // updateTimeouts
  

function updateFouls(team, val) {

  var f = null;

  var span = document.createElement("span");
  var text  = document.createTextNode(FOULS);
  
  span.className = "light-blue";
  span.appendChild(text);
  
  if(team == HOME) {
    f = document.getElementById("homeFoul");
  } else {
    f = document.getElementById("awayFoul");
  }

  f.innerHTML = val + " ";
  f.appendChild(span);

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


function updateClock(game, shot, mins, secs) {

  document.getElementById("shotClock").innerHTML = shotClockToString(shot, secs);
  document.getElementById("gameClock").innerHTML = gameClockToString(game, mins);

} // updateClock


function updateTeam(team, data) {

  if(data == null || data == undefined ||
    data == "") {
    return;
  }

  updateFouls(team, data.fouls);
  updateTimeouts(team, data.timeouts);
  updateScore(team, calcScore(data.points));

} // updateTeam


function updateDisplay(data) {

  if(data == null || data == undefined ||
    data == "") {
    return;
  }

  var j = JSON.parse(data);

  if(j instanceof Array) {
    return;
  }

  updatePeriod(j.period);

  if(j.possession) {
    updatePossession(HOME);
  } else {
    updatePossession(AWAY);
  }

  updateTeam(HOME, j.home);
  updateTeam(AWAY, j.away);
  
  updateClock(j.game, j.shot, j.settings.minutes, j.settings.shot);

  updateScore(HOME, calcScore(j.home.points));
  updateScore(AWAY, calcScore(j.away.points));

  updateFouls(HOME, j.home.fouls);
  updateFouls(AWAY, j.away.fouls);

  updateTimeouts(HOME, j.home.timeouts, j.settings.timeouts);
  updateTimeouts(AWAY, j.away.timeouts, j.settings.timeouts);

} // updateDisplay
