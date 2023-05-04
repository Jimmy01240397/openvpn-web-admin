import axios from '@/plugins/axios.js';

function getvpn() {
  return axios.get('/vpn/get');
}

function getvpns() {
  return axios.get('/vpn/getall');
}

function addvpn(vpnname, enable) {
  return axios.post('/vpn/add', {
    vpnname,
    enable,
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

function deletevpn(vpnname) {
  return axios.post('/vpn/delete', {
    vpnname,
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

function startstopvpn(vpnname, active) {
  return axios.post('/vpn/startstop', {
    vpnname,
    active,
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

export default {
  getvpn,
  getvpns,
  addvpn,
  deletevpn,
  startstopvpn,
}
