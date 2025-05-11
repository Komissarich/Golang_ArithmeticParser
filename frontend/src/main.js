//import './assets/main.css'

import { createApp, reactive, watch } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)

const authState =  reactive({
    isLoggedIn: localStorage.getItem('auth') === 'true',
    login() {
        
      this.isLoggedIn = true
      localStorage.setItem('auth', 'true')
    },
    logout() {
      this.isLoggedIn = false
      localStorage.setItem('auth', "false")
      localStorage.setItem('user', "")
      localStorage.setItem('token', "")
    }
  })

export default authState

watch(
    () => authState.isLoggedIn,
    (newVal) => {
      localStorage.setItem('auth', newVal ? 'true' : 'false')
    }
  )
  // Делаем доступным во всех компонентах
  app.provide('authState', authState)

app.use(router)
app.mount('#app')