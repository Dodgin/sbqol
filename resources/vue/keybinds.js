const Keybinds = {
    data() {
        return {
            keybinds: [
                { name: 'FcuForward', defaultKey: 'w', key: '' },
                { name: 'FcuBackward', defaultKey: 's', key: '' },
                { name: 'FcuLeft', defaultKey: 'a', key: '' },
                // Add more keybinds as needed
            ],
        };
    },
    template: `
        <div>
            <h1>Keybind Input Boxes</h1>
            <ul>
                <li v-for="(keybind, index) in keybinds" :key="index">
                    {{ keybind.name }} [ {{ keybind.defaultKey }} ]: <input v-model="keybind.key" />
                </li>
            </ul>
        </div>
    `,
    name: 'Keybinds'
};
