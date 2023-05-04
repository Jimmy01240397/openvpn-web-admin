import axios from '@/plugins/axios.js';

function getgroup(groupname) {
  return axios.post('/group/get', {
    groupname,
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

function getgroups() {
  return axios.get('/group/getall');
}

export default {
  getgroup,
  getgroups,
}
