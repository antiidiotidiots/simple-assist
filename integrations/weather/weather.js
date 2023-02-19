// Get the weather from the API
var response = JSON.parse(assist.request(
    "https://api.openweathermap.org/data/2.5/weather?q=Minneapolis&APPID=55f8d4b02cd09c48ee8607bc39b9a844", 
    "get",
    {},
    {}
));

assist.respond("The weather in Minneapolis is " + response.weather[0].description + ".");
var temp = Math.round(response.main.temp - 273.15);
var hasS = "s";
if(temp === 1) {
    hasS = "";
}
assist.respond("The temperature is " + temp + " degree" + hasS + " Celsius.");
assist.respond("The humidity is " + response.main.humidity + "%.");
assist.respond("The wind speed is " + response.wind.speed + " meters per second.");
assist.respond("The wind direction is " + response.wind.deg + " degrees.");
assist.respond("The pressure is " + response.main.pressure + " hPa.");
assist.respond("The cloudiness is " + response.clouds.all + "%.");
assist.respond("The sunrise is at " + new Date(response.sys.sunrise * 1000).toLocaleTimeString() + ".");
assist.respond("The sunset is at " + new Date(response.sys.sunset * 1000).toLocaleTimeString() + ".");
assist.respond("The geo coordinates are [" + response.coord.lon + ", " + response.coord.lat + "].");