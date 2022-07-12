<script setup lang="ts">
import { User } from '@/services/models';
import { reactive } from '@vue/reactivity';
import { emailRules, genderVariants } from "@/services/validation"

import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import apolloClient from '@/plugins/apollo';

const props = defineProps<{ user?: User, create: boolean }>()
const emit = defineEmits(['close'])

const form = reactive({
  firstName: props.user?.firstName,
  lastName: props.user?.lastName,
  email: props.user?.email,
  phone: props.user?.phone,
  gender: props.user?.gender,
  admin: props.user?.admin,
  password: "",
  valid: false
})

function emitClose() {
  emit("close")
}

const { mutate: addUserMut, loading: loadingCreate } = useMutation(gql`
  mutation addUser($input: UserInput!){
    addUser(input: $input) {
      id
      firstName
      lastName
      email
      gender
      admin
    }
  }
`)

async function createUser() {
  await addUserMut({
    input: {
      firstName: form.firstName,
      lastName: form.lastName,
      email: form.email,
      gender: form.gender,
      admin: form.admin ?? false,
      password: form.password
    }
  })
  await apolloClient.reFetchObservableQueries()
  emitClose()
}


const { mutate: updateUserMut, loading: loadingUpdate } = useMutation(gql`
  mutation updateUser($id: ID!, $input: UserUpdate!){
    updateUser(id: $id, input: $input) {
      firstName
      lastName
      phone
      gender
      admin
    }
  }
`)

async function updateUser() {
  await updateUserMut({
    id: props.user?.id,
    input: {
      firstName: form.firstName,
      lastName: form.lastName,
      phone: form.phone,
      gender: form.gender,
      admin: form.admin ?? false,
    }
  })
  await apolloClient.reFetchObservableQueries()
  emitClose()
}
</script>

<template>
  <v-card>
    <v-card-title class="d-flex justify-center"><span class="text-h5">{{ props.create ? "Create user" :
        "Modify user"
    }}</span></v-card-title>
    <v-card-text>
      <v-form v-model="form.valid">
        <v-container min-width="400">
          <v-row>
            <v-col cols="6">
              <v-text-field density="compact" variant="outlined" label="Email" :disabled="!props.create"
                :rules="emailRules" required v-model="form.email">
              </v-text-field>
            </v-col>
            <v-col cols="6" v-if="props.create">
              <v-text-field density="compact" variant="outlined" required type="password" label="Password"
                v-model="form.password">
              </v-text-field>
            </v-col>

          </v-row>

          <v-row>
            <v-col cols="6">
              <v-text-field density="compact" variant="outlined" label="First Name" required v-model="form.firstName">
              </v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field density="compact" variant="outlined" label="Last Name" required v-model="form.lastName">
              </v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field density="compact" variant="outlined" label="Phone" v-model="form.phone">
              </v-text-field>
            </v-col>

            <v-col cols="6">
              <v-select density="compact" variant="outlined" :items="genderVariants" v-model="form.gender"
                label="Gender"></v-select>
            </v-col>
            <v-col>
              <v-checkbox density="compact" variant="outlined" v-model="form.admin" label="Is Admin"></v-checkbox>
            </v-col>
          </v-row>
        </v-container>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-row justify="center">
        <v-btn v-if="!props.create" variant="tonal" :loading="loadingUpdate" color="primary" class="mr-4"
          @click="updateUser()">
          Save
        </v-btn>
        <v-btn v-else variant="tonal" :disabled="!form.valid" :loading="loadingCreate" color="primary" class="mr-4"
          @click="createUser()">
          Create
        </v-btn>
        <v-btn type="" class="mr-4" @click="emitClose()">
          Close
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>

</template>