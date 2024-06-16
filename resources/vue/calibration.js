const Calibration = {
    data() {
        return {
            items: [
                { name: 'Item1', value: 'Value1', address: 'Address1' },
                { name: 'Item2', value: 'Value2', address: 'Address2' },
                // Add more items as needed
            ],
        };
    },
    template: `
        <div>
            <h1>Calibration Complete</h1>
            <ul>
                <li v-for="(item, index) in items" :key="index">
                    {{ item.name }}: {{ item.value }} ({{ item.address }})
                </li>
            </ul>
        </div>
    `,
    name: 'Calibration'
};
