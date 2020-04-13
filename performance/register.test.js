import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL } from "./test-script.js";
import faker from 'cdnjs.com/libraries/Faker'

export const RegisterTest = () => {
    const headers = {
        'Content-Type': 'application/json'
    }
    const user = {
        username: faker.internet.userName(),
        password: faker.internet.password(),
    }

    let res = http.post(`${BASE_URL}/register`, JSON.stringify(user), { headers: headers });
    check(res, {
        'status was 200': r => r.status == 200,
        'transaction time OK': r => r.timings.duration < 200,
    });
    sleep(1);
}