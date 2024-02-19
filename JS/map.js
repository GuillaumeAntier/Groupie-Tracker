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
