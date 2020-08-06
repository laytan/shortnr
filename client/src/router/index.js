import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import { isLoggedIn, doRefresh } from '../auth';

const loggedInToDashboard = (_, __, next) => {
  if (isLoggedIn.value) {
    next({ name: 'Dashboard', query: { 'from-login': 'true' } });
    return;
  }
  next();
};

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import(/* webpackChunkName: "dashboard" */ '../views/Dashboard.vue'),
    beforeEnter: (_, __, next) => {
      if (!isLoggedIn.value) {
        next({ name: 'Login' });
        return;
      }
      next();
    },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import(/* webpackChunkName: "login" */ '../views/Login.vue'),
    beforeEnter: loggedInToDashboard,
  },
  {
    path: '/register',
    name: 'Register',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "register" */ '../views/Register.vue'),
    beforeEnter: loggedInToDashboard,
  },
  {
    path: '/:id',
    name: 'Redirect',
    component: () => import(/* webpackChunkName: "redirect" */ '../views/Redirect.vue'),
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import(/* webpackChunkName: "not-found" */ '../views/NotFound.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Only try to refresh when we actually left the app
let hasCheckedThisLoad = false;

// Check login
router.beforeEach(async (to, __, next) => {
  if (!hasCheckedThisLoad && !isLoggedIn.value) {
    hasCheckedThisLoad = true;

    // Homepage does not need to wait for the refresh to finish
    if (to.name === 'Home' || to.name === 'Redirect') {
      doRefresh();
    } else {
      await doRefresh();
    }
  }

  next();
});

export default router;
