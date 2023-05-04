<template>
  <div class='layout'>
    <v-card>
      <v-layout full-height=true>
        <v-navigation-drawer
          theme="dark"
          permanent
        >
          <v-list :selected="selectedPage">
            <v-list-item
              v-for="(page, index) in pages"
              :key="index"
              :value="page"
              :prepend-icon="page.icon"
              :title="page.title"
              color="yellow"
              @click="selectPage(page)"
            >
            </v-list-item>
          </v-list>
  
          <template v-slot:append>
            <div class="pa-2">
              <v-btn 
                block
                @click='logout'
              >
                Logout
              </v-btn>
            </div>
          </template>
        </v-navigation-drawer>
        <v-main style="height: 400px"></v-main>
      </v-layout>
    </v-card>
    <router-view class='router'></router-view>
  </div>
</template>

<script>
import userService from '@/services/user';
import groupService from '@/services/group';
export default {
  name: 'LayOut',
  data: () => ({
    userinfo: null,
    admingroup: null,
    selectedPage: null,
  }),
  computed: {
    pages: function() {
      if(this.admingroup != null && this.userinfo != null && this.admingroup.Members.includes(this.userinfo.Username)) {
        return [
          {
              title: 'Dashboard',
              icon: 'mdi-view-dashboard',
              path: '/layout/dashboard',
          },
          {
              title: 'Log',
              icon: 'mdi-book',
              path: '/layout/log',
          },
          {
              title: 'USER',
              icon: 'mdi-account-supervisor',
              path: '/layout/user',
          },
          {
              title: 'GROUP',
              icon: 'mdi-account-group-outline',
              path: '/layout/GROUP',
          },
          {
              title: 'VPNs',
              icon: 'mdi-open-source-initiative',
              path: '/layout/vpns',
          },
        ]
      }
      else {
        return [
          {
              title: 'Dashboard',
              icon: 'mdi-view-dashboard',
              path: '/layout/dashboard',
          },
        ]
      }
    }
  },
  methods: {
    selectPage(page) {
      this.selectedPage = [ page ];
      this.$router.push(page.path);
    },
    async logout() {
      if((await userService.logout()).data) this.$router.push('/');
    }
  },
  async beforeMount() {
    this.userinfo = (await userService.getuser()).data;
    this.admingroup = (await groupService.getgroup("Administrators")).data;
    this.selectPage(this.pages[0]);
  }
}
</script>

<style scoped>
.error-message{
  display: inline;
}
.layout {
  flex-direction: row;
  display: flex;
  width: 100%;
  height: 100%;
}
.router {
  flex-direction: column;
  display: flex;
  width: 100%;
  height: 100%;
}
</style>
