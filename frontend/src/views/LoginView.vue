<script setup lang="ts">
import { reactive, ref } from "@vue/reactivity";
import { useAuthStore } from "@/stores/auth";
import { emailRules } from "@/services/validation"

const form = reactive({
  email: "",
  password: "",
  valid: false
})

const loginError = ref("")
const loading = ref(false)

async function login() {
  loading.value = true
  try {
    const e = await auth.login(form.email, form.password)
    if (e) {
      loginError.value = e
    }
  } finally {
    loading.value = false
  }
}

const auth = useAuthStore()
</script>

<template>
  <v-container>
    <v-form v-model="form.valid" @submit.prevent="login()">
      <v-row justify="center">
        <v-col cols="6">
          <v-text-field label="Email" density="compact" variant="outlined" :rules="emailRules" required
            v-model="form.email"></v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="6">
          <v-text-field label="Password" density="compact" variant="outlined" required type="password"
            v-model="form.password"></v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-btn :disabled="!form.valid" :loading="loading" variant="tonal" color="primary" class="mr-4" type="submit">
          Login
        </v-btn>
      </v-row>
    </v-form>
  </v-container>
</template>