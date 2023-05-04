import axios from '@/plugins/axios.js';

function getuser() {
  return axios.get('/user/get');
}

function getusers() {
  return axios.get('/user/getall');
}

function issignin() {
  return axios.get('/user/issignin');
}

function adduser(username, password) {
  return axios.post('/user/add', {
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
  getusers,
  issignin,
  login,
  logout,
  adduser,
}
