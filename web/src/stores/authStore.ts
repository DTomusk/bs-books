import { defineStore } from 'pinia';
import * as authService from 'src/services/authService';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    loading: false,
  }),

  getters: {
    isAuthenticated: (state) => !!state.user,
  },

  actions: {
    async login(email: string, password: string) {
      this.loading = true;

      try {
        const data = await authService.login(email, password);

        localStorage.setItem('token', data.token);
        this.user = data.user;
      } finally {
        this.loading = false;
      }
    },

    logout() {
      localStorage.removeItem('token');
      this.user = null;
    },

    async fetchUser() {
      try {
        this.user = await authService.getMe();
      } catch {
        this.logout();
      }
    },
  },
});
