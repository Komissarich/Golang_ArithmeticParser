<template>
    <nav class="navbar">
     <router-link to="/play" class="navbar-brand">Easy-quizy</router-link>
     <div class="navbar-actions">
        <button class = "btn" @click="$router.push('/calculate')">Калькулятор</button>
      
      
       <button 
         class="btn" 
         :class="authState.isLoggedIn ? 'logout' : 'login'"
         @click="authState.isLoggedIn ? handleLogout() : $router.push('/auth')"
       >
         {{ authState.isLoggedIn ? 'Выйти' : 'Войти' }}
       </button>
 
     </div>  
   </nav>

  
 </template>
 
 
 
 <script>
 import { inject, onMounted } from 'vue'
 import { useRouter } from 'vue-router'
 
 export default {
   name: 'TheNavbar',
   methods: {
     goToProfile() {
     
       this.$router.push('/profile/me');
     }
   },
 
   setup() {
     
     const authState = inject('authState')
     const router = useRouter()
 
     const gotoProfile = () => {
       this.$router.push('/profile/me');
     }
     const handleLogout = () => {
       authState.logout()
       router.push('/auth')
     }
     return {
       authState,
       handleLogout,
       router
     }
   },
   
     goToProfile() {
       this.$router.push('/profile/' + localStorage.getItem('username'));
     }
   
   }
 
 </script>
 
 <style scoped>
 .navbar {
   display: flex;
   justify-content: space-between;
   align-items: center;
   padding: 1rem 2rem;
   background-color: #f8f9fa;
   box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
   font-family: Arial, sans-serif;
 }
 
 .navbar-brand {
   font-size: 1.5rem;
   font-weight: bold;
   color: #333;
 }
 
 .navbar-actions {
   display: flex;
   gap: 1rem;
 }
 
 .btn {
   padding: 0.5rem 1rem;
   border: none;
   border-radius: 4px;
   font-size: 1rem;
   cursor: pointer;
   transition: background-color 0.3s;
 }
 
 .create-test {
   background-color: #f1f1f1;
   color: rgb(0, 0, 0);
 }
 
 .create-test:hover {
   background-color:#f1f1f1;
 }
 
 .login {
   background-color: #f1f1f1;
   color: #333;
 }
 
 .login:hover {
   background-color: #ddd;
 }
 </style>