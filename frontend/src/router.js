import { createRouter, createWebHistory } from "vue-router";
import Auth from "./components/Auth.vue";
import Calculator from "./components/Calculator.vue";
import Register from "./components/Register.vue";

import authState from "./main"



const router = createRouter( {
     history: createWebHistory(),
    routes: [
        {path: '/auth', component: Auth, meta: {requireAuth: false}, alias: '/'},
        {path: '/register', component: Register, meta: {requireAuth: false}},
        {path: '/calculate', component: Calculator, meta: {requireAuth: true} },
    ]
})

router.beforeEach((to, from, next) => {
   
    
    if (to.meta.requireAuth === false) {
        next()
    }
    else {
      
        if (to.meta.requireAuth === true && authState.isLoggedIn === false) {
           
          
            return next({ path: "/auth" })
        }
     
        next()
    }
      
 
    
  })

export default router;