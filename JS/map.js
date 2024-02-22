var map = L.map('map').setView([0, 0], 2);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
  }).addTo(map);
  var index = 0;
  for (const [location] of Object.entries(coordinatesMap)) {
    if (coordinatesMap.hasOwnProperty(location)) {
      var latitude = coordinatesMap[location][0].lat;
      var longitude = coordinatesMap[location][0].lon;
      var marker = L.marker([latitude, longitude]).addTo(map);
      marker.bindPopup('Événement: ' + eventData[index].Dates + ', Ville: ' + formatLocation(location));
      index++;
    }
  }

function formatLocation(location) {
  var words = location.split(/[-_]/).map(function(word) {
    if (word.length <= 3) {
      return word.toUpperCase();
    } else {
      return word.charAt(0).toUpperCase() + word.slice(1);
    }
  });
  var countryIndex = words.length - 1;
  if ((words[countryIndex] === 'Caledonia') || (words[countryIndex] === 'Guinea') ||  (words[countryIndex] === 'Kong') || (words[countryIndex] === 'Guinea') || (words[countryIndex] === 'Guinea') || (words[countryIndex] === 'Guinea')){
    countryIndex -= 1;
  }
  var city = words.slice(0, countryIndex).join(' ');
  var country = words.slice(countryIndex).join(' ');
  return city + ' (' + country + ')';
}
  