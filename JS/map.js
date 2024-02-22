
var map = L.map('map').setView([0, 0], 13);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
  }).addTo(map);
  var index = 0;
  for (const [location] of Object.entries(coordinatesMap)) {
    if (coordinatesMap.hasOwnProperty(location)) {
      var latitude = coordinatesMap[location][0].lat;
      var longitude = coordinatesMap[location][0].lon;
      var marker = L.marker([latitude, longitude]).addTo(map);
      marker.bindPopup('Événement: ' + eventData[index].Dates + ', Ville: ' + location);
      index++;
    }
  }

var nextEventIndex = 0;
var nextDateIndex = 0;
  
function showNextEvent(coordinatesMap) {
    var eventElement = document.querySelector('#event');
    var events = Array.from(document.querySelector('#events').children);
  
    if (nextEventIndex < events.length) {
        var location = events[nextEventIndex].children[0].textContent;
        var dates = Array.from(events[nextEventIndex].children[1].children).map(li => li.textContent);
  
        if (nextDateIndex < dates.length) {
            eventElement.textContent = 'Location: ' + location + ', Date: ' + dates[nextDateIndex];
            nextDateIndex++;
            
            console.log(location);
            console.log(coordinatesMap);
            if (coordinatesMap.hasOwnProperty(location)) {
              var locationCoordinates = coordinatesMap[location];
              console.log(locationCoordinates);
              if (locationCoordinates && locationCoordinates.length > 0) {
                var coordinates = locationCoordinates[0];
                if (coordinates) {
                  map.setView([coordinates.lat, coordinates.lon], 13);
                }
              }
            } else {
              console.error('Invalid location: ' + location);
          }
            
        } else {
            nextEventIndex++;
            nextDateIndex = 0;
            showNextEvent();
        }
    }
    if (nextEventIndex >= events.length) {
        nextEventIndex = events.length - 1;
        nextDateIndex = events[nextEventIndex].children[1].children.length - 1;
        eventElement.textContent = 'No more events';
    }
}


function showPreviousEvent(coordinatesMap, map) {
    var eventElement = document.querySelector('#event');
    var events = Array.from(document.querySelector('#events').children);

    if (nextEventIndex < 0) {
        nextEventIndex = 0;
        nextDateIndex = 0;
        eventElement.textContent = 'No more previous events';
        return;
    }

    var location = events[nextEventIndex].children[0].textContent;
    var dates = Array.from(events[nextEventIndex].children[1].children).map(li => li.textContent);

    if (nextDateIndex < 0) {
        nextEventIndex--;
        if (nextEventIndex >= 0) {
            nextDateIndex = events[nextEventIndex].children[1].children.length - 1;
            showPreviousEvent();
        } else {
            nextEventIndex = 0;
            nextDateIndex = 0;
            eventElement.textContent = 'No more previous events';
        }
    } else {
        if (nextDateIndex >= dates.length) {
            nextDateIndex = dates.length - 1;
            console.log(dates.length, nextDateIndex, dates[nextDateIndex]);
        }
        eventElement.textContent = 'Location: ' + location + ', Date: ' + dates[nextDateIndex];
        nextDateIndex--;

        var coordinates = coordinatesMap[location][0];
        if (coordinates) {
            map.setView([coordinates.lat, coordinates.lon], 13);
        }
    }
}

