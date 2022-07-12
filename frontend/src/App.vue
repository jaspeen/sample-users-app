<script setup lang="ts">
import { computed } from '@vue/reactivity';
import { useErrorsStore } from './stores/errors';

const errors = useErrorsStore()
const hasErrors = computed((): boolean => errors.message != null)

</script>

<template>
  <v-app>
    <!--v-navigation-drawer v-model="drawer" app></v-navigation-drawer-->
    <v-app-bar app elevation="1">
      <!--v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon-->

      <v-toolbar-title>Sample users app</v-toolbar-title>
    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <v-snackbar variant="tonal" v-model="hasErrors" multi-line>
      {{ errors.message }}

      <template v-slot:actions>
        <v-btn color="red" variant="text" @click="errors.message = null">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>



<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /*text-align: center;*/
  color: #2c3e50;
}
</style>
