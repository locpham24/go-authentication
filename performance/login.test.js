import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL } from "./test-script.js";

const USERNAME = 'tindt'
const PASSWORD = 'password'

export const LoginTest = () => {
    const headers = {
        'Content-Type': 'application/json'
    }
    const user = {
        username: USERNAME,
        password: PASSWORD,
    }

    let res = http.post(`${BASE_URL}/login`, JSON.stringify(user), { headers: headers });
    check(res, {
        'status was 200': r => r.status == 200,
        'transaction time OK': r => r.timings.duration < 200,
    });
    sleep(1);
}