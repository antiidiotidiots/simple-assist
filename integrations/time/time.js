var date = new Date(Date.now());
var dateValues = [
   date.getFullYear(),
   date.getMonth() + 1,
   date.getDate(),
   date.getHours() % 12,
   date.getMinutes(),
   date.getSeconds(),
   date.getHours() > 12
];

message = "It is " + dateValues[3] + ":" + dateValues[4] + " ";

if(dateValues[6]) {
    message += "PM"
} else {
    message += "AM"
}

message += "."