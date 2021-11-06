import { createStore } from 'vuex'

// Create a new store instance.
const store = createStore({
    state () {
        return {
            auth: 0,
        }
    },
    mutations: {
        increment (state) {
            state.count++
        },
        setAuth (state, payload) {
            state.auth = payload;
        },
    }
})

export default store;