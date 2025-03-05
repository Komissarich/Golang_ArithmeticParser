const calculator = {
    displayValue: '0',
    firstOperand: null,
    waitingForSecondOperand: false,
    operator: null,
};

  function handleFormSubmit() {
    const expression_str = document.getElementsByClassName("calculator-screen")[0].value
    const response = fetch("http://localhost:8080/api/v1/calculate/", {
        method: "post",
        headers: {
          'Content-Type': 'application/json'
        },
      
        body: JSON.stringify({
            expression: expression_str
        })
      })
      .then( (response) => { 
       
        response.json().then((data) => {
            console.log(response.text, data);
        })
    
      });
  }
  

 function change() {
    if ((document.getElementsByClassName('list')[0].value == 'One expression' || document.getElementsByClassName('list')[0].value == 'One task') && document.getElementsByClassName('ghost').length == 0) {
    
    const inputElement = document.createElement('input');

// Устанавливаем атрибуты элемента (например, тип, placeholder и т.д.)
inputElement.setAttribute('type', 'text');
inputElement.setAttribute('class', 'ghost');
inputElement.setAttribute('placeholder', 'Write id');
inputElement.setAttribute('style', 'height:50px; width:300px');

// Добавляем элемент на страницу
document.body.appendChild(inputElement); // Добавление в тело документа
// Или:
const container = document.getElementsByClassName('magic_input')[0];
container.appendChild(inputElement)

    }

    if ((document.getElementsByClassName('list')[0].value == 'All expressions' || document.getElementsByClassName('list')[0].value == 'All tasks') && document.getElementsByClassName('ghost').length == 1) {
        var elem = document.getElementsByClassName('ghost')[0];
        elem.parentNode.removeChild(elem);
      
    }
 } 

  
function inputDigit(digit) {
    const { displayValue, waitingForSecondOperand } = calculator;
    if (calculator.displayValue === '0') {
        calculator.displayValue = digit;
    } else {
        calculator.displayValue = displayValue + digit;
    }
}

function resetCalculator() {
    calculator.displayValue = '0';
    calculator.firstOperand = null;
    calculator.waitingForSecondOperand = false;
    calculator.operator = null;
}

function updateDisplay() {
    const display = document.querySelector('.calculator-screen');
    display.value = calculator.displayValue;
}

function receiveInfo() {
    const options = document.getElementsByClassName("list")[0]
    var url = ""
   
   if (options.value == 'All expressions')  {
        url = "http://localhost:8080/api/v1/expressions/"
       
   } else if (options.value == 'One expression') {
        id = document.getElementsByClassName("ghost")[0].value
       
        url = "http://localhost:8080/api/v1/expressions/" + id
      
   } else if (options.value == 'All tasks') {
    url = "http://localhost:8080/api/v1/tasks/"
   }
   else if (options.value == 'One task') {
    id = document.getElementsByClassName("ghost")[0].value
    url = "http://localhost:8080/api/v1/tasks/" + id
   }
    const response = fetch(url, {
        method: "get",
       
      })
      .then( (response) => { 
    
        response.json().then((data) => {
            console.log(response.statusText, data);
            const a = document.getElementsByClassName("received-expressions");
            a[0].textContent = JSON.stringify(data, null, 2);
        })
      });
}


const keys = document.querySelector('.calculator-buttons');


keys.addEventListener('click', (event) => {
   
    const { target } = event;
    if (!target.matches('button')) {
        return;
    }

   

    if (target.classList.contains('all-clear')) {
        resetCalculator();
        updateDisplay();
        return;
    }

    if (target.classList.contains("equal-sign")) {
        handleFormSubmit();
        return;
    }
    inputDigit(target.value);
    updateDisplay();
});