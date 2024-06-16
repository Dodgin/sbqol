<template>
    <div>
      <h1>HOTAS</h1>
      <div v-if="message">
        <p>Name: {{ message.name }}</p>
        <p>Address: {{ message.address }}</p>
        <p>Value: {{ message.value }}</p>
      </div>
    </div>
  </template>
  
  <script>
  import { computed } from 'vue';
  import { state } from '../state'; // Import the state
  
  export default {
    name: 'Hotas',
    setup() {
      const message = computed(() => state.throttleMessage);
      return { message };
    }
  }
  </script>
  
  <style scoped>
  /* Add your styles here */
  </style>
  