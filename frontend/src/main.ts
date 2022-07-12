import { createApp, provide, h } from "vue";
import App from "./App.vue";
import vuetify from "@/plugins/vuetify";
import { loadFonts } from "@/plugins/webfontloader";
loadFonts();

import { createPinia } from "pinia";
const pinia = createPinia();

// renew access token using refresh token
const auth = useAuthStore(pinia);
await auth.refresh();

import { ApolloClients } from "@vue/apollo-composable";
import apolloClient from "@/plugins/apollo";

import { router } from "@/plugins/router";
import { useAuthStore } from "@/stores/auth";
createApp({
  setup() {
    provide(ApolloClients, apolloClient);
  },

  render: () => h(App),
})
  .use(pinia)
  .use(router)
  .use(vuetify)
  .mount("#app");
