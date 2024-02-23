var map = L.map('map').setView([0, 0], 2); // Center of the map and zoom level 
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
  }).addTo(map);
  var index = 0;
  for (const [location] of Object.entries(coordinatesMap)) { // Loop through the coordinatesMap object
    if (coordinatesMap.hasOwnProperty(location)) { // Check if the property exists
      var latitude = coordinatesMap[location][0].lat; 
      var longitude = coordinatesMap[location][0].lon;
      var marker = L.marker([latitude, longitude]).addTo(map); // Add a marker to the map
      marker.bindPopup('Événement: ' + eventData[index].Dates + ', Ville: ' + formatLocation(location));
      index++; // Increment the index
    }
  }

function formatLocation(location) {
  var words = location.split(/[-_]/).map(function(word) { // Split the location string and map through the words
    if (word.length <= 3) { // Check if the length of the word is less than or equal to 3
      return word.toUpperCase(); // Return the word in uppercase
    } else {
      return word.charAt(0).toUpperCase() + word.slice(1);
    }
  });
  var countryIndex = words.length - 1; // Get the index of the last word
  if ((words[countryIndex] === 'Caledonia') || (words[countryIndex] === 'Guinea') ||  (words[countryIndex] === 'Kong') || (words[countryIndex] === 'Guinea') || (words[countryIndex] === 'Guinea') || (words[countryIndex] === 'Guinea')){
    countryIndex -= 1;
  }
  /// Check if the last word is a country
  /// If it is put the country in parenthesis
  var city = words.slice(0, countryIndex).join(' ');
  var country = words.slice(countryIndex).join(' ');
  return city + ' (' + country + ')';
}
  