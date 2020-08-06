import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import './style.scss';
import 'bootstrap/dist/js/bootstrap.bundle';
import Icon from './components/Icon.vue';
import './registerServiceWorker';

createApp(App).use(router).component('Icon', Icon).mount('#app');
