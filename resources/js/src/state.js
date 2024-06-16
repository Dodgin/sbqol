// src/state.js
import { reactive } from 'vue';

export const state = reactive({
  throttleMessage: {
    name: '',
    address: '',
    value: 0
  }
});
