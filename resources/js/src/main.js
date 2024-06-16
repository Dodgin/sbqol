import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { state } from './state'; // Import the state

const app = createApp(App);
app.use(router);
app.mount('#app');
router.replace("/").catch((err) => console.log(err))

// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(rawMessage) {
        const message = JSON.parse(rawMessage);
        switch(message.type) {
            case "throttle":
            {
                let val0 = JSON.parse(message.value);
                val0 = val0[0];
                state.throttleMessage.name = val0.name;
                state.throttleMessage.address = val0.address;
                state.throttleMessage.value = val0.value;
                if(router.currentRoute.value.path != "/hotas") router.replace("/hotas").catch((err) => console.log(err));
                break;
            }
        }
    });
});
