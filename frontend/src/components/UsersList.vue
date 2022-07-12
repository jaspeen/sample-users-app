<script setup lang="ts">
import { provideApolloClient, useMutation, useQuery } from "@vue/apollo-composable"
import gql from "graphql-tag";
import apolloClient from "@/plugins/apollo"
import { ref } from "@vue/reactivity";
import EditUser from "./EditUser.vue";
import { User } from "@/services/models";

provideApolloClient(apolloClient)

const showSideBar = ref(false)

const { result: usersList } = useQuery<{ users: User[] }>(gql`
  query {
    users {
      id
      firstName
      lastName
      email
      phone
      gender
      admin
    }
  }
`, null, { fetchPolicy: "cache-and-network" })

const curUser = ref<User | undefined>(undefined)

const { mutate: removeUserMutation, error, loading } = useMutation(gql`
    mutation removeUser($id: ID!) {
      removeUser (id: $id)
    }
  `);
const removeUserLoading = ref(false)
async function removeUser() {
  removeUserLoading.value = true;
  try {
    const res = await removeUserMutation({ id: curUser.value?.id })
    await apolloClient.reFetchObservableQueries()
    showSideBar.value = false
  } finally {
    removeUserLoading.value = false
  }
}

const showEditUser = ref(false)
const showCreateUser = ref(false)

</script>


<template>
  <v-container>
    <v-dialog v-model="showCreateUser">
      <v-container>
        <EditUser :create="true" v-on:close="showCreateUser = false" />
      </v-container>
    </v-dialog>
    <v-dialog v-model="showEditUser">
      <v-container>
        <EditUser :user="curUser" :create="false" v-on:close="showEditUser = false; showSideBar = false" />
      </v-container>
    </v-dialog>

    <v-col>
      <v-btn variant="outlined" density="compact" color="success" icon @click="showCreateUser = true">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-col>
    <v-navigation-drawer location="right" v-model="showSideBar" temporary>
      <v-card>
        <v-card-title class="d-flex justify-center">{{ curUser?.firstName }} {{ curUser?.lastName }}</v-card-title>
        <v-card-actions>
          <v-row justify="center">
            <v-btn color="primary" variant="tonal" @click="showEditUser = true">Modify</v-btn>
            <v-btn color="error" variant="tonal" :loading="removeUserLoading" @click="removeUser()">Remove</v-btn>
          </v-row>
        </v-card-actions>
      </v-card>
    </v-navigation-drawer>
    <v-table density="compact">
      <thead>
        <tr>
          <th class="text-left">
            Name
          </th>
          <th class="text-left">
            Email
          </th>
          <th class="text-left">
            Phone
          </th>
          <th class="text-left">
            Gender
          </th>
          <th class="text-left">
            Is Admin
          </th>
        </tr>
      </thead>
      <tbody>
        <tr @click="curUser = user; showSideBar = true" v-for="user of usersList?.users" :key="user.id">
          <td>{{ user.firstName + " " + user.lastName }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.phone }}</td>
          <td>{{ user.gender }}</td>
          <td>
            <v-checkbox-btn density="compact" readonly v-model="user.admin" />
          </td>
        </tr>
      </tbody>
    </v-table>
  </v-container>
</template>

