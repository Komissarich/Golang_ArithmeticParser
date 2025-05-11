<template>
    <div class="calculator-container">
      <div class="calculator">
        <!-- Экран калькулятора -->
        <div class="calculator-screen">
          <input 
            type="text" 
            v-model="displayValue" 
            readonly
          >
        </div>
  
        <!-- Клавиши калькулятора -->
        <div class="calculator-keys">
          <!-- Верхний ряд: AC и скобки -->
          <div class="top-row">
            <button class="all-clear" @click="resetCalculator">AC</button>
            <button class="bracket" value="(">(</button>
            <button class="bracket" value=")">)</button>
          </div>
  
          <!-- Основные цифры -->
          <div class="digits-grid">
            <button value="7">7</button>
            <button value="8">8</button>
            <button value="9">9</button>
            <button value="4">4</button>
            <button value="5">5</button>
            <button value="6">6</button>
            <button value="1">1</button>
            <button value="2">2</button>
            <button value="3">3</button>
            <button value="0" class="zero-btn">0</button>
            <button value=".">.</button>
          
          </div>
  
          <!-- Операторы справа -->
          <div class="operators-column">
            <button class="operator" value="+">+</button>
            <button class="operator" value="-">-</button>
            <button class="operator" value="*">*</button>  
            <button class="operator" value="/">/</button>
            <button class="equal-sign" @click="handleFormSubmit">=</button>
          </div>
        </div>
      </div>
  
      <!-- Управляющие элементы -->
      <div class="controls">
        <select v-model="selectedOption" @change="handleOptionChange">
          <option>All expressions</option>
          <option>One expression</option>
          <option>All tasks</option>
          <option>One task</option>
        </select>
        <textarea 
        v-if="showIdInput" 
        v-model="idInput" 
        placeholder="Enter ID here..."
        rows="3"
        style="margin-top: 10px; width: 200px;"
    ></textarea>
        <button class="receive-btn" @click="receiveInfo">Receive</button>
        
        <div class="output">
          <pre>{{ receivedData }}</pre>
        </div>
      </div>
    </div>
  </template>
  <script>
  import axios from 'axios'
  export default {
    data() {
      return {
        displayValue: '0',
        firstOperand: null,
        waitingForSecondOperand: false,
        operator: null,
        selectedOption: 'All expressions',
        showIdInput: false,
        idInput: '',
        receivedData: null
      }
    },
    methods: {
        handleOptionChange() {
      this.showIdInput = this.selectedOption === 'One expression' || this.selectedOption === 'One task';
      if (!this.showIdInput) {
        this.idInput = ''; // Очищаем поле, если оно скрыто
      }
    },
      inputDigit(digit) {
        if (this.displayValue === '0') {
          this.displayValue = digit;
        } else {
          this.displayValue += digit;
        }
      },
      
      resetCalculator() {
        this.displayValue = '0';
        this.firstOperand = null;
        this.waitingForSecondOperand = false;
        this.operator = null;
      },
      async handleFormSubmit() {
        try {
            const token = localStorage.getItem('token');
            console.log("token",token)
            const response = await axios.post(
            "/api/v1/calculate/",
            {
                expression: this.displayValue
            },
            {
                headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
                }
            }
            );

        console.log('Expression id', response.data.id);
    
    } catch (error) {
        console.error('Calculation failed:', error);
        this.displayValue = 'Error';
    }
    },
      handleOptionChange() {
        this.showIdInput = 
          this.selectedOption === 'One expression' || 
          this.selectedOption === 'One task';
      },
      async receiveInfo() {
 
    const token = localStorage.getItem('token');
    const user_id = localStorage.getItem('user_id'); // Предполагаем, что username сохраняется при логине
    

    let url;
    let method = 'get';
    let data = null;

    const baseConfig = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    };

    if (this.selectedOption === 'All expressions') {
      url = "/api/v1/expressions/";
      try { 
            const response = await axios.get(
          '/api/v1/expressions/',
            {
              headers: {
                'Content-Type': 'application/json', // Важно явно указать!
                 'Authorization': `Bearer ${token}`
              },
            }
        )
        console.log(response.data)
            this.receivedData = response.data;
        } catch (error) {
            console.error('Request failed:', error);
          //  this.receivedData = { error: error.response?.data?.message || error.message };
        }
    } 
    else if (this.selectedOption === 'One expression') {
        try { 
            const response = await axios.post(
          '/api/v1/get_expression/',
            {
              expression_id: this.idInput,
              user_id: user_id
            },
            {
              headers: {
                'Content-Type': 'application/json', // Важно явно указать!
              },
            }
        )
            this.receivedData = response.data;
        } catch (error) {
            console.error('Request failed:', error);
            this.receivedData = { error: error.response?.data?.message || error.message };
        }
        
   
    } 
    else if (this.selectedOption === 'All tasks') {
    
      const response = await axios.get(
          '/api/v1/tasks/',
            {
              task_id: this.idInput
            },
            {
              headers: {
                'Content-Type': 'application/json', // Важно явно указать!
              },
            }
        )
    } 
    else if (this.selectedOption === 'One task') {
        try {
            const data = await axios.post(
          '/api/v1/get_task/',
            {
              expression_id: this.idInput,
              user_id: user_id
            },
            {
              headers: {
                'Content-Type': 'application/json', // Важно явно указать!
              },
            }
        )
    
            this.receivedData = response.data;
         } catch (error) {
    console.error('Request failed:', error);
   
  }
        }
    
     
   
},
    },
    mounted() {
      // Add click handlers for all buttons except special ones
      const buttons = document.querySelectorAll('.calculator-keys button:not(.all-clear):not(.equal-sign)');
      buttons.forEach(button => {
        button.addEventListener('click', () => {
          this.inputDigit(button.value);
        });
      });
    }
  }
  </script>
    
    <style scoped>
    .calculator-container {
      max-width: 320px;
      margin: 0 auto;
      font-family: Arial, sans-serif;
    }
    
    .calculator {
      background: #f0f2f5;
      border-radius: 12px;
      padding: 20px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    }
    
    .calculator-screen input {
      width: 90%;
      height: 60px;
      padding: 0 15px;
      margin-bottom: 20px;
      font-size: 28px;
      text-align: right;
      border: none;
      border-radius: 8px;
      background: #fff;
      box-shadow: inset 0 2px 4px rgba(0,0,0,0.1);
    }
    
    .calculator-keys {
      display: grid;
      grid-template-columns: 1fr 60px;
      gap: 10px;
    }
    
    .top-row {
      grid-column: 1 / -1;
      display: grid;
      grid-template-columns: 2fr 1fr 1fr;
      gap: 8px;
      margin-bottom: 8px;
    }
    
    .digits-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 8px;
    }
    
    .operators-column {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }
    
    button {
      height: 50px;
      border: none;
      border-radius: 8px;
      font-size: 18px;
      cursor: pointer;
      transition: all 0.2s;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    
    button:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 8px rgba(0,0,0,0.15);
    }
    
    button:active {
      transform: translateY(0);
    }
    
    .all-clear {
      background: #ff4444;
      color: white;
      font-weight: bold;
    }
    
    .bracket {
      background: #e0e3e7;
      font-weight: bold;
    }
    
    .zero-btn {
      grid-column: span 2;
    }
    
    .operator {
      background: #e0e3e7;
      font-weight: bold;
    }
    
    .equal-sign {
      background: #4CAF50;
      color: white;
      font-weight: bold;
      flex-grow: 1;
    }
    
    .controls {
      margin-top: 20px;
      background: #f0f2f5;
      padding: 20px;
      border-radius: 12px;
    }
    
    select, .receive-btn {
      width: 100%;
      padding: 12px;
      margin-bottom: 10px;
      border: none;
      border-radius: 8px;
      font-size: 16px;
    }
    
    .receive-btn {
      background: #2196F3;
      color: white;
    }
    
    .output {
      margin-top: 15px;
      padding: 15px;
      background: #fff;
      border-radius: 8px;
      min-height: 100px;
      max-height: 200px;
      overflow: auto;
    }
    
    .output pre {
      margin: 0;
      white-space: pre-wrap;
    }
    </style>