import { sleep, group } from 'k6'
import { LoginTest } from './login.test.js'
import { RegisterTest } from './register.test.js'

export const BASE_URL = 'http://localhost:8080'

export let options = {
    stages: [
        { duration: '30s', target: 50 },
        { duration: '1m30s', target: 100 },
        { duration: '20s', target: 20 },
    ],
};

export default (data) => {
    group('Login', () => { LoginTest() })
    group('Register', () => { RegisterTest() })
    sleep(1)
}