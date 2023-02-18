// Get weather from a public API
// Author:
//
// Created: 2019-10-10
// Updated: 2019-10-10
// Version: 1.0.0
// License: MIT
//
// Description:
// This script will get the weather from a public API and display it on the screen.
//
// Dependencies:
// None
//
// Configuration:
// None
//
// Commands:
// None

// Get the weather from the API
function getWeather() {
    // Create a new XMLHttpRequest object
    var request = new XMLHttpRequest();

    // Open a new connection, using the GET request on the URL endpoint
    request.open('GET', 'https://api.openweathermap.org/data/2.5/weather?q=Minneapolis&APPID=55f8d4b02cd09c48ee8607bc39b9a844    ', true);

    request.onload = function () {
        // Begin accessing JSON data here
        var data = JSON.parse(this.response);

        // If the request was successful
        if (request.status >= 200 && request.status < 400) {
            // Log the data
            console.log(data);

            // Get the weather description and convert it to lowercase
            var weatherDescription = data.weather[0].description.toLowerCase();

            // Get the temperature and convert it to fahrenheit
            var tempF = Math.round((data.main.temp - 273.15) * 1.8 + 32);

            // Get the wind speed
            var windSpeed = data.wind.speed;

            // Get the wind direction
            var windDirection = data.wind.deg;

            // Get the humidity
            var humidity = data.main.humidity;

            // Create the message
            var message = "The weather is " + weatherDescription + ". The temperature is " + tempF + " degrees fahrenheit. The wind speed is " + windSpeed + " and the wind direction is " + windDirection + ". The humidity is " + humidity + "%.";

            // Display the message
            document.getElementById("weather").innerHTML = message;
        } else {
            // Log the error
            console.log('error');
        }
    }

    // Send request
    request.send();
}

// Run the function when the page loads
getWeather();