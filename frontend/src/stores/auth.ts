import { defineStore } from "pinia";
import gql from "graphql-tag";
import { router } from "@/plugins/router";
import apolloClient from "@/plugins/apollo";
import { User } from "@/services/models";

const baseUrl = `${import.meta.env.VITE_API_URL}/query`;

export const useAuthStore = defineStore({
  id: "auth",
  state: () => {
    const storedUser = localStorage.getItem("user");
    const storedRefreshToken = localStorage.getItem("refreshToken");
    console.log("Stored:" + storedUser);

    return {
      user: storedUser ? JSON.parse(storedUser) : (null as User | null),
      accessToken: null as string | null,
      refreshToken: storedRefreshToken as string | null,
      returnUrl: null as string | null,
    };
  },
  actions: {
    async login(
      username: string,
      password: string
    ): Promise<string | undefined> {
      try {
        const res = await apolloClient.mutate({
          mutation: gql`
            mutation login($username: String!, $password: String!) {
              login(username: $username, password: $password) {
                user {
                  id
                  firstName
                  lastName
                  admin
                  email
                }
                token
                refreshToken
              }
            }
          `,
          variables: { username, password },
        });

        // update pinia state
        this.accessToken = res.data.login.token;
        this.refreshToken = res.data.login.refreshToken;
        const user = res.data.login.user;
        this.user = user;

        // store user details and jwt in local storage to keep user logged in between page refreshes
        localStorage.setItem("user", JSON.stringify(user));
        localStorage.setItem("refreshToken", this.refreshToken as string);

        // redirect to previous url or default to home page
        router.push(this.returnUrl || "/");
      } catch (error: any) {
        return error.message;
      }
    },
    async refresh(): Promise<boolean> {
      if (!this.refreshToken) {
        return false;
      }
      const res = await apolloClient.mutate({
        mutation: gql`
          mutation ($input: String!) {
            renewToken(refreshToken: $input) {
              token
              refreshToken
            }
          }
        `,
        variables: { input: this.refreshToken },
      });
      if (!res.errors) {
        this.accessToken = res.data.renewToken.token;
        console.log("Setting access token: " + this.accessToken);
        return true;
      } else {
        await this.logout();
        return false;
      }
    },
    async fetchUser() {},
    logout() {
      this.accessToken = this.refreshToken = this.user = null;
      localStorage.removeItem("user");
      localStorage.removeItem("refreshToken");
      router.push("/login");
    },
  },
});
