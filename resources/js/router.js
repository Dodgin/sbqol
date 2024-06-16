const routes = [
    { path: '/', component: Home },
    { path: '/calibration', component: Calibration },
    { path: '/keybinds', component: Keybinds },
];

const router = new VueRouter({
    routes,
    mode: 'hash'
});
