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

function updateURL() {
        // Get the min and max year values
        var minYear = document.querySelector('.input-min').value;
        var maxYear = document.querySelector('.input-max').value;
    
        // Update the URL with the new parameters
        window.location.href = '/result?minYear=' + minYear + '&maxYear=' + maxYear;
    }
