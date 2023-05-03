<template>
  <div class="login-div">
    <div class="error-message">
      <v-alert
        color='red'
        type='error'
        :model-value='loginError'
      >
        {{ message }}
      </v-alert>
    </div>
    <form class='login-form'>
      <v-text-field
        label='Username'
        v-model='username'
        @keyup.enter='login'
      ></v-text-field>
      <v-text-field 
        label='Password'
        v-model='password'
        type='password'
        @keyup.enter='login'
      ></v-text-field>
      <v-btn
        @click='login'
      >
        Login
      </v-btn>
    </form>
  </div>
</template>


<script>
import userService from '@/services/user';
export default {
  name: 'LogIN',
  data: () => ({
    username: '',
    password: '',
    loginError: false,
    message: ''
  }),
  methods: {
    login: async function() {
      const result = await userService.login(this.username, this.password);
      switch (result.status) {
        case 200:
          this.loginError = false;
          this.$router.push('/layout');
          break;
        default:
          this.loginError = true;
          this.message = result.data;
      }
    }
  },
  async beforeMount() {
      var issignin = await userService.issignin();
      if(issignin.data) this.$router.push('/layout');
  }
}
</script>

<style scoped>
.error-message{
  display: inline;
}
.login-div {
  flex-direction: column;
  display: flex;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
}

.login-form {
  margin: 30px;
  width: 40%;
}
</style>
