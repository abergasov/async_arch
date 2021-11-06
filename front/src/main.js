import { createApp } from 'vue'
import App from './App.vue'
import axios from 'axios';
import store from "./store";
import Vant from 'vant';
import 'vant/lib/index.css';

const app = createApp(App);
app.use(store);
app.use(Vant);

app.config.globalProperties.askBackend = (url, param) => {
    let config = {
        headers: {
            m: window.userId,
        }
    }
    return new Promise((resolve, reject) => {
        axios.post(`/api/v1/${url}`, param, config)
            .then(resp => {
                resolve(resp.data)
            })
            .catch(error => {
                let code = +error.response.status;
                let message = ''
                switch (code) {
                    case 401:
                        this.$store.commit('setAuth', 0);
                        message = 'Unauthorized';
                        break;
                    case 409:
                        message = 'Already exist';
                        break;
                    case 400:
                        message = 'Bad request';
                        break;
                }
                reject(error)
            })
    });
};

app.mount('#app')