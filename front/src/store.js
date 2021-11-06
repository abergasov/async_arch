import { createStore } from 'vuex'

// Create a new store instance.
const store = createStore({
    state () {
        return {
            auth: 0,
            init_done: false,
            key: localStorage.getItem('key'),
        }
    },
    mutations: {
        increment (state) {
            state.count++
        },
        setAuth (state, payload) {
            state.auth = payload;
        },
        setJWT (state, payload) {
            localStorage.setItem('key', payload);
            state.key = payload;
        },
        initDone (state) {
            state.init_done = true;
        },
    }
})

export default store;