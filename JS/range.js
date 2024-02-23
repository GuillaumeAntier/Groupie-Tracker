const rangeInput = document.querySelectorAll('.range-input input'); // Get the range inputs
priceInput = document.querySelectorAll('.price-input input'); 
progress = document.querySelector('.slider .progress'); // Get the progress bar

let yeargap = 0; // The minimum year gap

let totalRange = parseInt(rangeInput[1].max) - parseInt(rangeInput[0].min); // The total range of the slider 

rangeInput.forEach(input => { // Loop through the range inputs
    input.addEventListener('input', e => { // Add an event listener to the input
        let minVal = parseInt(rangeInput[0].value) - parseInt(rangeInput[0].min), // Get the current value of the min input
            maxVal = parseInt(rangeInput[1].value) - parseInt(rangeInput[1].min); // Get the current value of the max input

        // Update the progress bar
        if (maxVal - minVal <= yeargap) {
            if (e.target.className === "input-min") { // Check if the input is the min input
                rangeInput[0].value = maxVal - yeargap;
                progress.style.left = (minVal / totalRange) * 100 + '%';
            } else { // If the input is the max input
                rangeInput[1].value = minVal + yeargap;
                progress.style.right = 100 - (maxVal / totalRange) * 100 + '%';
            }
        } 
    });
});

rangeInput.forEach(input => { // Loop through the range inputs
    input.addEventListener('input', e => { // Add an event listener to the input
        let minVal = parseInt(rangeInput[0].value),
            maxVal = parseInt(rangeInput[1].value);
        
        document.querySelector('#minYear').textContent = minVal;
        document.querySelector('#maxYear').textContent = maxVal;


        if (maxVal - minVal < yeargap) { // Check if the difference between the max and min values is less than the year gap
            if (e.target.className === "range-min") { // Check if the input is the min input
                rangeInput[0].value = maxVal - yeargap;
            } else { // If the input is the max input
                rangeInput[1].value = minVal + yeargap;
            }
        } else { // If the difference between the max and min values is greater than the year gap
            priceInput[0].value = minVal;
            priceInput[1].value = maxVal;
            progress.style.left = ((minVal - 1958) / totalRange) * 100 + '%';
            progress.style.right = 100 - ((maxVal - 1958) / totalRange) * 100 + '%';
        }
    });
});


let rangeMin = document.querySelector('.range-min');
let rangeMax = document.querySelector('.range-max');


rangeMin.addEventListener('input', filterGroups);
rangeMax.addEventListener('input', filterGroups);

function filterGroups() { 
    // Get the current values of the sliders
    let minYear = rangeMin.value;
    let maxYear = rangeMax.value;

    // Filter the groups based on the slider values
    let filteredGroups = Groups.filter(group => group.creationDate >= minYear && group.creationDate <= maxYear);

    // Generate the HTML for the filtered groups
    let html = filteredGroups.map(group => `
    
    <input type="checkbox" id="popup-${group.name}" class="popup-trigger">
    <label for="popup-${group.name}">
        <img src="${group.image}" alt="${group.name}" class="img-artiste">
    </label>
    <div class="popup">
        <label for="popup-${group.name}" class="close-button">x</label>
        <h2>${group.name}</h2>
        <div class="popup-content">
            <div class="Members">
                <h2>Members : </h2>
                ${group.members.map(member => `<p>${member}</p>`).join('')}
            </div>
            <div class = "date-of-creation">
                <h2>Date of creation : </h2>
                <p>${group.creationDate}</p>
            </div>
            <div class = "first-album">
                <h2>First album : </h2>
                <p>${group.firstAlbum}</p>
            </div>
            <form class ="event-btn" action = "/event" method="get">
              <input type="hidden" name="id" value="${group.id}">
              <button class = "btn" type = "submit" >Events</button>
            </form>
        </div>
    </div>
    `).join('');

    // Update the HTML of the .cover-artiste-grid section
    document.querySelector('.cover-artiste-grid').innerHTML = html;
}

