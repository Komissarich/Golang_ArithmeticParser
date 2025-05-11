

<template>
  <div id="app">

  <TheNavbar></TheNavbar>
  <router-view></router-view>
 </div>
</template>



<script>

import { RouterView } from "vue-router";
import TheNavbar from "./components/TheNavbar.vue";


 export default {
   components: {TheNavbar}
 }



</script>