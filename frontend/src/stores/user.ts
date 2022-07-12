import { defineStore } from "pinia";

interface User {
  username: string;
  firstName: string;
  lastName: string;
  email: string;
}

const useUserStore = defineStore({
  id: "user",
  state: () => {
    const stored = localStorage.getItem("user");
    const authToken = localStorage.getItem("authToken");
    const refreshToken = localStorage.getItem("refreshToken");
    return {
      user: stored ? JSON.parse(stored) : (null as User | null),
      authToken,
      refreshToken,
    };
  },
});
