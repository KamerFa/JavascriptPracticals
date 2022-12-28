
// Function that runs on first render and gets data from localStorage //

const body = document.getElementById("body");


function initialTableRender() {
    var keys = Object.keys(localStorage);

    var userKeys = keys.filter(key => key.startsWith("user"));

    if(userKeys.length != 0){
        for (var i = 0; i < userKeys.length; i++) {
            var key = userKeys[i];
            var value = JSON.parse(localStorage.getItem(key));
            addRow(value);
            console.log(value, key);
        }
    }
}

initialTableRender();


// Form handlers //

// get form element
const form = document.getElementById('myForm');

// get form field elements
const formFields = form.elements;

// handle form submission
form.addEventListener('submit', function (event) {
    // prevent form from submitting
    event.preventDefault();

    const formData = new Array;

    // loop through form fields and push values to formData array
    for (let i = 0; i < formFields.length -1; i++) {
        const field = formFields[i];
        if( field.type === "radio"){
            if(field.checked){
                formData.push(field.value);
            }
        } else {
        formData.push(field.value);
        }
    }

    // log formData array
    console.log(formData);
    addRow(formData);
    localStorage.setItem("user"+formData[0], JSON.stringify(formData));

});



function addRow(formData) {

    const table = document.getElementById('myTable');
    const row = document.createElement('tr');
    
    var deleteButton = document.createElement("button");
    deleteButton.innerHTML = "Delete";
    deleteButton.classList = "btn btn-secondary my-1"
    
    var editButton = document.createElement("button");
    editButton.innerHTML = "Edit";
    editButton.classList = "btn btn-secondary mx-1 my-1"

    deleteButton.addEventListener("click", function () {
        // Remove the row from the table
        table.deleteRow(row.rowIndex);

        // Remove the array from local storage
        localStorage.removeItem("user"+formData[0]);
    });
    function editRow() {
        // Get the string from localStorage
        let myString = localStorage.getItem('user' + formData[0]);
        

        // Create a modal with a text field
        let modal = document.createElement('div');
        modal.innerHTML = `
            <div class='modal-edit'>
                <p class="fw-bold font-3">Edit the string:</p>
                <textarea class="shadow" type="text" id="string-input" rows="10" cols="40" value="${myString}">${myString}</textarea>
                <div class='buttons-container mt-4'>
                    <button id="save-button" class="btn btn-dark btn-lg px-4 fs-6 fw-light">Save</button>
                    <button id="cancel-button" class="btn btn-dark btn-lg px-4 fs-6 fw-light">Cancel</button>
                </div>
            </div>
            `;
        document.body.appendChild(modal);

        // Add an event listener to the save button
        let saveButton = document.getElementById('save-button');
        saveButton.addEventListener('click', function () {
            // Get the new value from the text field
            let newString = document.getElementById('string-input').value;

            // Save the new value to localStorage
            localStorage.setItem('user' + formData[0], newString);

            // Close the modal
            modal.remove();
            body.style.position="unset";
        });
        let cancelButton = document.getElementById('cancel-button');
        cancelButton.addEventListener('click', function(){
            modal.remove();
            body.style.position = "unset";
        })
    }

    editButton.addEventListener("click", editRow);

    
    
    
    


    // loop through formData array and create cell for each item
    for (let i = 0; i < formData.length; i++) {
        const cell = document.createElement('td');
        cell.innerHTML = formData[i];

        row.appendChild(cell);
    }

    // append row to table
    table.appendChild(row);
    row.appendChild(deleteButton);
    row.appendChild(editButton);
}
