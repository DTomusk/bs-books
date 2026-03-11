import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    meta: { variant: 'main' },
    children: [
      { path: '', component: () => import('pages/LandingPage/_LandingPage.vue') },
      {
        path: 'home',
        component: () => import('pages/Dashboard.vue'),
        meta: { requiresAuth: true },
      },
    ],
  },
  {
    path: '/auth',
    component: () => import('layouts/MainLayout.vue'),
    meta: { variant: 'auth' },
    children: [
      {
        path: 'login',
        component: () => import('pages/AuthPage.vue'),
        props: { mode: 'login', isOwnPage: true },
      },
      {
        path: 'create-account',
        component: () => import('pages/AuthPage.vue'),
        props: { mode: 'register', isOwnPage: true },
      },
    ],
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

declare module 'vue-router' {
  interface RouteMeta {
    variant?: 'main' | 'auth';
    requiresAuth?: boolean;
  }
}

export default routes;
