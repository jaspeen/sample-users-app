import { defineStore } from "pinia";

export const useErrorsStore = defineStore("errors", {
  state: () => ({
    message: null as string | null,
    category: null as string | null,
    fields: { input: {} as any },
  }),
});
