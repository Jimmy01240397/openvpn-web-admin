import axios from '@/plugins/axios.js';

function getuser() {
  return axios.get('/user');
}

function issignin() {
  return axios.get('/user/issignin');
}

function login(username, password) {
  return axios.post('/user/login', {
    username,
    password
  }).then((res) => {
    return {
      status: res.status,
      data: res.data,
    };
  }).catch((err) => {
    return {
      status: err.response.status,
      data: err.response.data.msg,
    };
  });
}

function logout() {
  return axios.get('/user/logout');
}

export default {
  getuser,
  issignin,
  login,
  logout,
}
