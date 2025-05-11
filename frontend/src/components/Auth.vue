<template>
    <div class="auth-container">
   
      
      <form class="login-form">
        <div class="welcome-header">
        <h1>Добро пожаловать!</h1>
      </div>
      <div class="form-group">
          <label for="login">Почта</label>
          <input 
            type="text" 
            id="login" 
             v-model="email"
            placeholder="Введите почту"
            required
          >
        </div>

    
        
        <div class="form-group">
          <label for="password">Пароль</label>
          <input 
            type="password" 
            id="password" 
            v-model="password"
            placeholder="Введите ваш пароль"
            required
          >
        </div>
        
     
        
        <button type="submit" class="login-btn" @click.prevent="handleLogin">
          Войти
        </button>
      </form>
    </div>

    <div v-if="errorMessage !== ''" class="error-message">
      {{ errorMessage }} 
    </div>


    <div class="reg-container">
        <form class="login-form">
        <button type="submit" class="reg-btn" @click.prevent="handleReg">
          Регистрация
        </button>
        </form>
    </div>
</template>

<script>
import axios from 'axios'
import { inject } from 'vue'

// После успешного входа:



export default {
    name: 'Auth',
    setup() {
      const authState = inject('authState')
      return {
        authState
      }
    },
  data() {
    return {
        email: '',
        password: '',
        errorMessage: ''
    }
  },
  methods: {
    async handleLogin() {
        this.errorMessage = ''
        if (this.username !== ""  && this.password !== "") {
         try { 
          const data = await axios.post(
          '/api/v1/login/',
            {
              email: this.email,
              password: this.password
            },
            {
              headers: {
                'Content-Type': 'application/json', // Важно явно указать!
              },
            }
          );
          
          
          console.log("Succesfully auth")
          localStorage.setItem("auth", "true")
          console.log("new token ", data.token)  
          console.log(data.data)
          localStorage.setItem("user_id", data.data.user_id)
          console.log("new_id", data.user_id)
          localStorage.setItem('token', data.data.token);
          
          this.authState.isLoggedIn = true
        
          this.$router.push("/calculate") 
         } catch (error) {
          console.log(error)
          this.errorMessage = 'Пользователь с таким логином или паролем не найден'
          
         }
        
        }
        else {
            this.errorMessage = 'Пожалуйста, заполните все поля'
        }

        
        
    },
    handleReg() {
      this.$router.push('/register') // Перенаправление после успешного входа
    }
  }
}
</script>



<style scoped>

.error-message {
    max-width: 400px;
  color: #f44336;
  background-color: #ffebee;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 20px;
  margin: 0 auto;
  border-left: 4px solid #f44336;
  text-align: center;
}
.auth-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Arial', sans-serif;
  text-align: center;
}

.reg-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Arial', sans-serif;
  text-align: center;
}
.welcome-header h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 2rem;
  line-height: 1.3;
}

.login-form {
  background: #f8f9fa;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
  text-align: left;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.remember-me {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
}

.remember-me input {
  margin-right: 0.5rem;
}

.login-btn {
  width: 100%;
  padding: 0.75rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.reg-btn {
  width: 100%;
  padding: 0.75rem;
  background-color: #45b9ce;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}
.login-btn:hover {
  background-color: #45a049;
}
</style>