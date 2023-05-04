<template>
  <div>
    <v-expansion-panels class="mb-6">
      <v-expansion-panel
        v-for="(userdata, i) in userdatas"
        :key="i"
      >
        <v-expansion-panel-title expand-icon="mdi-plus" collapse-icon="mdi-minus">
          {{ userdata.Username }}
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          test
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
    <h2>Add New User</h2>
    <div class="message">
      <v-alert
        :color='messagecolor'
        :type='messagetype'
        :model-value='showmessage'
      >
        {{ message }}
      </v-alert>
    </div>
    <form class='adduser-form'>
      <v-text-field
        label='Username'
        v-model='username'
        @keyup.enter='adduser'
      ></v-text-field>
      <v-text-field 
        label='Password'
        v-model='password'
        type='password'
        @keyup.enter='adduser'
      ></v-text-field>
      <v-text-field 
        label='Password Confirmation'
        v-model='checkpassword'
        type='password'
        @keyup.enter='adduser'
      ></v-text-field>
      <v-btn
        @click='adduser'
      >
        Add
      </v-btn>
    </form>

  </div>
</template>

<script>
import userService from '@/services/user';
export default {
  name: 'USER',
  data: () => ({
    username: '',
    password: '',
    checkpassword: '',
    messagecolor: 'red',
    messagetype: 'error',
    showmessage: false,
    message: "",
    userdatas: [],
  }),
  methods: {
    adduser: async function() {
      if(this.password !== this.checkpassword) {
        this.messagecolor = 'red';
        this.messagetype = 'error';
        this.showmessage = true;
        this.message = "Password confirmation does not match.";
        return;
      }
      const result = await userService.adduser(this.username, this.password);
      switch (result.status) {
        case 200:
          this.username = '';
          this.password = '';
          this.checkpassword = '';
          this.messagecolor = 'green';
          this.messagetype = 'success'
          this.showmessage = true;
          this.message = "Add user success!";
          this.userdatas = (await userService.getusers()).data;
          break;
        default:
          this.messagecolor = 'red';
          this.messagetype = 'error';
          this.showmessage = true;
          this.message = result.data;
          break;
      }
    }
  },
  async beforeMount() {
    this.userdatas = (await userService.getusers()).data;
  }
}
</script>
