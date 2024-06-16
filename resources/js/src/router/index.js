import { createRouter, createWebHistory } from 'vue-router';
import Needscalibration from '../views/Needscalibration.vue';
import Hotas from '../views/Hotas.vue';

const routes = [
  {
    path: '/',
    name: 'Needscalibration',
    component: Needscalibration
  },
  {
    path: '/hotas',
    name: 'Hotas',
    component: Hotas
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
